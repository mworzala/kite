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

	//todo at this point we should start connecting to the target server.

	return h.player.SendPacket(&packet.ServerLoginSuccess{
		GameProfile:         *h.player.Profile,
		StrictErrorHandling: true,
	})
}

func (h *ClientMojangLoginHandler) handleLoginAcknowledged(p *packet.ClientLoginAcknowledged) error {
	serverConn, err := net.Dial("tcp", "localhost:25565")
	if err != nil {
		panic(err)
	}

	doneCh := make(chan bool)
	remote, readLoop := proto.NewConn(packet.Clientbound, serverConn)
	remote.SetState(packet.Handshake, nil)
	err = remote.SendPacket(&packet.ClientHandshake{
		ProtocolVersion: 766,
		ServerAddress:   "localhost:25577",
		ServerPort:      25577,
		Intent:          packet.IntentLogin,
	})
	if err != nil {
		panic(err)
	}
	remote.SetState(packet.Login, NewClientboundLoginHandler(h.player, remote, doneCh))
	err = remote.SendPacket(&packet.ClientLoginStart{
		Name: h.player.Username,
		UUID: h.player.UUID.String(),
	})
	if err != nil {
		panic(err)
	}
	go readLoop()

	<-doneCh

	h.player.SetState(packet.Config, NewServerboundConfigurationHandler(h.player))
	return nil
}

var _ proto.Handler = (*ClientMojangLoginHandler)(nil)
