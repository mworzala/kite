package kite

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"errors"
	"fmt"
	"time"

	"github.com/mworzala/kite/internal/pkg/sessionserver"
	"github.com/mworzala/kite/pkg/proto"
	"github.com/mworzala/kite/pkg/proto/packet"
)

type ClientMojangLoginHandler[T any] struct {
	ClientMojangLoginHandlerOpts[T]

	Player *Player[T]

	publicKey   []byte
	verifyToken []byte

	remoteCtx    context.Context
	remoteCancel context.CancelFunc
	remote       *proto.Conn
}

type ClientMojangLoginHandlerOpts[T any] struct {
	PrivateKey        *rsa.PrivateKey
	InitialServerFunc func() ServerInfo

	ClientConfigHandlerFunc func(*Player[T]) proto.Handler
	ServerConfigHandlerFunc func(*Player[T]) proto.Handler
}

func MakeClientMojangLoginHandler[T any](opts ClientMojangLoginHandlerOpts[T]) func(*proto.Conn) proto.Handler {
	return func(conn *proto.Conn) proto.Handler {
		if opts.InitialServerFunc == nil {
			panic("InitialServerFunc is needed") //todo should not be a panic
		}

		publicKey, err := x509.MarshalPKIXPublicKey(&opts.PrivateKey.PublicKey)
		if err != nil {
			panic(err)
		}

		verifyToken := make([]byte, 16)
		_, err = rand.Read(verifyToken)
		if err != nil {
			panic(fmt.Errorf("failed to generate random verify token: %w", err))
		}

		return &ClientMojangLoginHandler[T]{
			ClientMojangLoginHandlerOpts: opts,
			Player:                       &Player[T]{Conn: conn},
			publicKey:                    publicKey,
			verifyToken:                  verifyToken,
		}
	}
}

func (h *ClientMojangLoginHandler[T]) HandlePacket(pp proto.Packet) (err error) {
	switch pp.Id {
	case packet.ClientLoginLoginStartID:
		p := new(packet.ClientLoginStart)
		if err = pp.Read(p); err != nil {
			return err
		}
		return h.HandleLoginStart(p)
	case packet.ClientLoginLoginAcknowledgedID:
		p := new(packet.ClientLoginAcknowledged)
		if err = pp.Read(p); err != nil {
			return err
		}
		return h.HandleLoginAcknowledged(p)
	case packet.ClientLoginEncryptionResponseID:
		p := new(packet.ClientEncryptionResponse)
		if err = pp.Read(p); err != nil {
			return err
		}
		return h.HandleEncryptionResponse(p)
	default:
		return proto.UnknownPacket
	}
}

func (h *ClientMojangLoginHandler[T]) HandleLoginStart(p *packet.ClientLoginStart) (err error) {
	h.Player.Username = p.Name

	return h.Player.SendPacket(&packet.ServerEncryptionRequest{
		ServerID:           "",
		PublicKey:          h.publicKey,
		VerifyToken:        h.verifyToken,
		ShouldAuthenticate: true,
	})
}

func (h *ClientMojangLoginHandler[T]) HandleEncryptionResponse(p *packet.ClientEncryptionResponse) error {
	decryptedVerifyToken, err := rsa.DecryptPKCS1v15(rand.Reader, h.PrivateKey, p.VerifyToken)
	if err != nil {
		panic(err)
	} else if !bytes.Equal(h.verifyToken, decryptedVerifyToken) {
		panic(errors.New("verifyToken not match"))
	}

	// get sharedSecret
	sharedSecret, err := rsa.DecryptPKCS1v15(rand.Reader, h.PrivateKey, p.SharedSecret)
	if err != nil {
		panic(err)
	}

	// Read and write encrypted data
	if err = h.Player.Conn.EnableEncryption(sharedSecret); err != nil {
		return err
	}

	// Do serverside auth with session server
	profile, err := sessionserver.HasJoined(context.Background(), h.Player.Username, "", sharedSecret, h.publicKey)
	if err != nil {
		return err
	} else if profile == nil {
		return errors.New("client did not do self auth")
	}

	h.Player.UUID = profile.ID
	h.Player.Username = profile.Name
	properties := make([]packet.ProfileProperty, len(profile.Properties))
	for i, prop := range profile.Properties {
		p := packet.ProfileProperty{Name: prop.Name, Value: prop.Value}
		if prop.Signature != "" {
			p.Signature = &prop.Signature
		}
		properties[i] = p
	}

	h.Player.Profile = &packet.GameProfile{
		UUID:       h.Player.UUID.String(),
		Username:   h.Player.Username,
		Properties: properties,
	}

	// At this point we can start connecting to the target server.
	// We will need to wait for it before processing any config state packets.
	targetServer := h.InitialServerFunc()
	var ctx context.Context
	ctx, h.remoteCancel = context.WithTimeout(context.Background(), 30*time.Second)
	println("connecting to", targetServer.Address)
	h.remoteCtx, h.remote, err = proto.CreateServerConn(ctx, targetServer.Address, uint16(targetServer.Port),
		h.Player.Profile, targetServer.Secret, NewServerVelocityLoginHandler)
	if err != nil {
		return err
	}

	return h.Player.SendPacket(&packet.ServerLoginSuccess{
		GameProfile: *h.Player.Profile,
	})
}

func (h *ClientMojangLoginHandler[T]) HandleLoginAcknowledged(_ *packet.ClientLoginAcknowledged) error {
	// This should never happen in normal operation, but a client could just send a login ack
	// immediately in an attempt to bypass auth. So don't let that happen :)
	if h.Player.Profile == nil {
		// The player didn't do encryption. Don't let them through.
		return fmt.Errorf("missing player profile")
	}

	defer h.remoteCancel()

	// Wait for the remote server connection (in fail or success)
	println("waiting for remote login")
	<-h.remoteCtx.Done()
	cause := context.Cause(h.remoteCtx)
	if errors.Is(cause, context.DeadlineExceeded) {
		// This would trigger from the timeout above. This is a complete fail.
		return fmt.Errorf("timed out connecting to remote server: %w", h.remoteCtx.Err())
	} else if !errors.Is(cause, context.Canceled) {
		// Otherwise, we failed with the given cause. This should cause a disconnect failure.
		return fmt.Errorf("failed to connect to remote server: %w", cause)
	}

	// A clean cancel is the success case, so we don't need to do anything.
	// We are now connected to the remote.

	h.remote.SetRemote(h.Player.Conn)
	h.Player.Conn.SetRemote(h.remote)

	if h.ClientConfigHandlerFunc != nil {
		h.Player.SetState(packet.Config, h.ClientConfigHandlerFunc(h.Player))
	} else {
		h.Player.SetState(packet.Config, NewClientConfigHandler(h.Player))
	}
	if h.ServerConfigHandlerFunc != nil {
		h.remote.SetState(packet.Config, h.ServerConfigHandlerFunc(h.Player))
	} else {
		h.remote.SetState(packet.Config, NewServerConfigHandler(h.Player))
	}

	return nil
}

var _ proto.Handler = (*ClientMojangLoginHandler[any])(nil)
