package handler

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/mworzala/kite"
	"github.com/mworzala/kite/internal/pkg/sessionserver"
	"github.com/mworzala/kite/pkg/proto"
	"github.com/mworzala/kite/pkg/proto/packet"
)

type ClientMojangLoginHandler struct {
	ClientMojangLoginHandlerOpts

	player *kite.Player

	publicKey   []byte
	verifyToken []byte

	remoteCtx    context.Context
	remoteCancel context.CancelFunc
	remote       *proto.Conn
}

type ClientMojangLoginHandlerOpts struct {
	PrivateKey *rsa.PrivateKey
}

func MakeClientMojangLoginHandler(opts ClientMojangLoginHandlerOpts) func(*proto.Conn) proto.Handler {
	return func(conn *proto.Conn) proto.Handler {

		publicKey, err := x509.MarshalPKIXPublicKey(&opts.PrivateKey.PublicKey)
		if err != nil {
			panic(err)
		}

		verifyToken := make([]byte, 16)
		_, err = rand.Read(verifyToken)
		if err != nil {
			panic(fmt.Errorf("failed to generate random verify token: %w", err))
		}

		return &ClientMojangLoginHandler{
			ClientMojangLoginHandlerOpts: opts,
			player:                       &kite.Player{Conn: conn},
			publicKey:                    publicKey,
			verifyToken:                  verifyToken,
		}
	}
}

func (h *ClientMojangLoginHandler) HandlePacket(pp proto.Packet) (err error) {
	switch pp.Id {
	case packet.ClientLoginLoginStartID:
		p := new(packet.ClientLoginStart)
		if err = pp.Read(p); err != nil {
			return err
		}
		return h.handleLoginStart(p)
	case packet.ClientLoginLoginAcknowledgedID:
		p := new(packet.ClientLoginAcknowledged)
		if err = pp.Read(p); err != nil {
			return err
		}
		return h.handleLoginAcknowledged(p)
	case packet.ClientLoginEncryptionResponseID:
		p := new(packet.ClientEncryptionResponse)
		if err = pp.Read(p); err != nil {
			return err
		}
		return h.handleEncryptionResponse(p)
	default:
		return proto.UnknownPacket
	}
}

func (h *ClientMojangLoginHandler) handleLoginStart(p *packet.ClientLoginStart) (err error) {
	h.player.Username = p.Name

	return h.player.SendPacket(&packet.ServerEncryptionRequest{
		ServerID:           "",
		PublicKey:          h.publicKey,
		VerifyToken:        h.verifyToken,
		ShouldAuthenticate: true,
	})
}

func (h *ClientMojangLoginHandler) handleEncryptionResponse(p *packet.ClientEncryptionResponse) error {
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
	if err = h.player.Conn.EnableEncryption(sharedSecret); err != nil {
		return err
	}

	// Do serverside auth with session server
	profile, err := sessionserver.HasJoined(context.Background(), h.player.Username, "", sharedSecret, h.publicKey)
	if err != nil {
		return err
	} else if profile == nil {
		return errors.New("client did not do self auth")
	}

	h.player.UUID = profile.ID
	h.player.Username = profile.Name
	properties := make([]packet.ProfileProperty, len(profile.Properties))
	for i, prop := range profile.Properties {
		p := packet.ProfileProperty{Name: prop.Name, Value: prop.Value}
		if prop.Signature != "" {
			p.Signature = &prop.Signature
		}
		properties[i] = p
	}

	h.player.Profile = &packet.GameProfile{
		UUID:       h.player.UUID.String(),
		Username:   h.player.Username,
		Properties: properties,
	}

	// At this point we can start connecting to the target server.
	// We will need to wait for it before processing any config state packets.
	h.remoteCtx, h.remoteCancel = context.WithTimeout(context.Background(), 30*time.Second)
	var cancel context.CancelCauseFunc
	h.remoteCtx, cancel = context.WithCancelCause(h.remoteCtx)
	h.remote, err = h.createRemoteConn(cancel, "localhost", 25565)
	if err != nil {
		return err
	}

	return h.player.SendPacket(&packet.ServerLoginSuccess{
		GameProfile:         *h.player.Profile,
		StrictErrorHandling: true,
	})
}

func (h *ClientMojangLoginHandler) handleLoginAcknowledged(_ *packet.ClientLoginAcknowledged) error {
	// This should never happen in normal operation, but a client could just send a login ack
	// immediately in an attempt to bypass auth. So don't let that happen :)
	if h.player.Profile == nil {
		// The player didn't do encryption. Don't let them through.
		return fmt.Errorf("missing player profile")
	}

	defer h.remoteCancel()

	// Wait for the remote server connection (in fail or success)
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

	h.remote.SetRemote(h.player.Conn)
	h.player.Conn.SetRemote(h.remote)

	h.remote.SetState(packet.Config, NewServerConfigHandler(h.player, h.remote))
	h.player.SetState(packet.Config, NewClientConfigHandler(h.player))
	return nil
}

func (h *ClientMojangLoginHandler) createRemoteConn(cancel context.CancelCauseFunc, address string, port uint16) (*proto.Conn, error) {
	serverConn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", address, port))
	if err != nil {
		return nil, fmt.Errorf("failed to dial remote: %w", err)
	}
	remote, readLoop := proto.NewConn(packet.Clientbound, serverConn)

	// Handshake immediately, then we are in login.
	handshake := &packet.ClientHandshake{
		ProtocolVersion: 766,
		ServerAddress:   address,
		ServerPort:      port,
		Intent:          packet.IntentLogin,
	}
	if err = remote.SendPacket(handshake); err != nil {
		return nil, err
	}

	// Setup velocity forwarding handler & begin login.
	remote.SetState(packet.Login, NewServerVelocityLoginHandler(remote, cancel, h.player.Profile))
	err = remote.SendPacket(&packet.ClientLoginStart{
		Name: h.player.Username,
		UUID: h.player.UUID.String(),
	})
	if err != nil {
		panic(err)
	}

	// Start reading from the remote connection
	go readLoop()

	return remote, nil
}

var _ proto.Handler = (*ClientMojangLoginHandler)(nil)
