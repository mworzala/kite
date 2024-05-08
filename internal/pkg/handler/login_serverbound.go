package handler

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mworzala/kite"
	"github.com/mworzala/kite/internal/pkg/crypto"
	"github.com/mworzala/kite/pkg/proto"
	"github.com/mworzala/kite/pkg/proto/packet"
)

var _ proto.Handler = (*ServerboundLoginHandler)(nil)

const doAuth = false

var privateKey *rsa.PrivateKey

func init() {
	var err error
	privateKey, err = rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}

	// Validate Private Key
	err = privateKey.Validate()
	if err != nil {
		panic(err)
	}
}

type ServerboundLoginHandler struct {
	Player *kite.Player
}

func NewServerboundLoginHandler(p *kite.Player) proto.Handler {
	return &ServerboundLoginHandler{p}
}

func (h *ServerboundLoginHandler) HandlePacket(pp proto.Packet) (err error) {
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
	}
	return proto.UnknownPacket
}

var username string
var vt []byte

func (h *ServerboundLoginHandler) handleLoginStart(p *packet.ClientLoginStart) error {
	username = p.Name

	if !doAuth {
		resp := &packet.ServerLoginSuccess{
			UUID:     "3bc51b9d-52be-4c4a-a3d6-7cc0bd6e6ea8",
			Username: "notmattw",
		}
		return h.Player.SendPacket(resp)
	}

	publicKey, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		panic(err)
	}

	verifyToken := make([]byte, 16)
	_, err = rand.Read(verifyToken)
	if err != nil {
		panic(err)
	}
	vt = verifyToken

	resp := &packet.ServerEncryptionRequest{
		ServerID:    "",
		PublicKey:   publicKey,
		VerifyToken: verifyToken,
	}
	return h.Player.SendPacket(resp)
}

func (h *ServerboundLoginHandler) handleEncryptionResponse(p *packet.ClientEncryptionResponse) error {
	if !doAuth {
		panic("unexpected encryption response")
	}

	decryptedVerifyToken, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, p.VerifyToken)
	if err != nil {
		panic(err)
	} else if !bytes.Equal(vt, decryptedVerifyToken) {
		panic(errors.New("verifyToken not match"))
	}

	// get sharedSecret
	sharedSecret, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, p.SharedSecret)
	if err != nil {
		panic(err)
	}

	// Read and write encrypted data
	h.Player.Conn.EnableEncryption(sharedSecret)

	// Do mojang auth with session server
	publicKey, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		panic(err)
	}
	authHash := crypto.Sha1([]byte(""), sharedSecret, publicKey)

	resp, err := http.Get("https://sessionserver.mojang.com/session/minecraft/hasJoined?username=" + username + "&serverId=" + authHash)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var authres Resp
	err = json.Unmarshal(body, &authres)

	return h.Player.SendPacket(&packet.ServerLoginSuccess{
		UUID:     "3bc51b9d-52be-4c4a-a3d6-7cc0bd6e6ea8",
		Username: "notmattw",
	})
}

// Resp is the response of authentication
type Resp struct {
	Name string
	ID   uuid.UUID
	//Properties []user.Property
}

func (h *ServerboundLoginHandler) handleLoginAcknowledged(p *packet.ClientLoginAcknowledged) error {
	serverConn, err := net.Dial("tcp", "localhost:25565")
	if err != nil {
		panic(err)
	}

	println(fmt.Sprintf("login ack: %dms", time.Since(start).Milliseconds()))

	doneCh := make(chan bool)
	remote, readLoop := proto.NewConn(packet.Clientbound, serverConn)
	remote.SetState(packet.Handshake, nil)
	err = remote.SendPacket(&packet.ClientHandshake{
		ProtocolVersion: 766,
		ServerAddress:   "localhost:25577",
		ServerPort:      25577,
		NextState:       packet.Login,
	})
	if err != nil {
		panic(err)
	}
	remote.SetState(packet.Login, NewClientboundLoginHandler(h.Player, remote, doneCh))
	err = remote.SendPacket(&packet.ClientLoginStart{
		Name: "notmattw",
		UUID: "3bc51b9d-52be-4c4a-a3d6-7cc0bd6e6ea8",
	})
	if err != nil {
		panic(err)
	}
	go readLoop()

	<-doneCh

	println(fmt.Sprintf("enter config: %dms", time.Since(start).Milliseconds()))

	h.Player.SetState(packet.Config, NewServerboundConfigurationHandler(h.Player))
	return nil
}
