package main

import (
	"bytes"
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/mworzala/kite"
	"github.com/mworzala/kite/pkg/mojang"
	packet2 "github.com/mworzala/kite/pkg/packet"
	"github.com/mworzala/kite/pkg/proto"
	"github.com/mworzala/kite/pkg/velocity"
)

func (p *Player) handleClientLoginPacket(pp proto.Packet) (err error) {
	switch pp.Id {
	case packet2.ClientLoginLoginStartID:
		pkt := new(packet2.ClientLoginStart)
		if err = pp.Read(pkt); err != nil {
			return err
		}
		return p.handleClientLoginStart(pkt)
	case packet2.ClientLoginLoginAcknowledgedID:
		pkt := new(packet2.ClientLoginAcknowledged)
		if err = pp.Read(pkt); err != nil {
			return err
		}
		return p.handleClientLoginAcknowledged(pkt)
	case packet2.ClientLoginEncryptionResponseID:
		pkt := new(packet2.ClientEncryptionResponse)
		if err = pp.Read(pkt); err != nil {
			return err
		}
		return p.handleClientEncryptionResponse(pkt)
	default:
		return proto.UnknownPacket
	}
}

func (p *Player) handleClientLoginStart(pkt *packet2.ClientLoginStart) (err error) {
	p.Username = pkt.Name

	return p.conn.SendPacket(&packet2.ServerEncryptionRequest{
		ServerID:           "",
		PublicKey:          p.proxy.PublicKey,
		VerifyToken:        p.proxy.VerifyToken,
		ShouldAuthenticate: true,
	})
}

func (p *Player) handleClientEncryptionResponse(pkt *packet2.ClientEncryptionResponse) (err error) {
	// TODO: verify token should be per-player, also it is only ever needed during login so we need some temporary state
	// 	     somewhere for that.
	decryptedVerifyToken, err := rsa.DecryptPKCS1v15(nil, p.proxy.PrivateKey, pkt.VerifyToken)
	if err != nil {
		panic(err)
	} else if !bytes.Equal(p.proxy.VerifyToken, decryptedVerifyToken) {
		panic(errors.New("verifyToken not match"))
	}

	sharedSecret, err := rsa.DecryptPKCS1v15(nil, p.proxy.PrivateKey, pkt.SharedSecret)
	if err != nil {
		panic(err)
	}

	// Read and write encrypted data
	if err = p.conn.EnableEncryption(sharedSecret); err != nil {
		return err
	}

	// Do serverside auth with session server
	profile, err := mojang.HasJoined(context.Background(), p.Username, "", sharedSecret, p.proxy.PublicKey)
	if err != nil {
		return err
	} else if profile == nil {
		return errors.New("client did not do self auth")
	}

	p.UUID = profile.ID
	p.Username = profile.Name
	properties := make([]packet2.ProfileProperty, len(profile.Properties))
	for i, prop := range profile.Properties {
		p := packet2.ProfileProperty{Name: prop.Name, Value: prop.Value}
		if prop.Signature != "" {
			p.Signature = &prop.Signature
		}
		properties[i] = p
	}

	p.Profile = &packet2.GameProfile{
		UUID:       p.UUID.String(),
		Username:   p.Username,
		Properties: properties,
	}

	return p.conn.SendPacket(&packet2.ServerLoginSuccess{
		GameProfile: *p.Profile,
	})
}

func (p *Player) handleClientLoginAcknowledged(pkt *packet2.ClientLoginAcknowledged) (err error) {
	// This should never happen in normal operation, but a client could just send a login ack
	// immediately in an attempt to bypass auth. So don't let that happen :)
	if p.Profile == nil {
		// The player didn't do encryption. Don't let them through.
		return fmt.Errorf("missing player profile")
	}

	p.remote, err = p.connectServerSync("localhost", 25565)
	if err != nil {
		panic(err) // TODO
	}

	// Server connection is already in config
	p.conn.SetState(packet2.Config)

	return nil
}

func (p *Player) connectServerSync(address string, port uint16) (*kite.Conn, error) {
	if p.pendingLoginChan != nil {
		panic("already connecting to a server")
	}
	p.pendingLoginChan = make(chan error)
	defer func() {
		close(p.pendingLoginChan)
		p.pendingLoginChan = nil
	}()

	serverConn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", address, port))
	if err != nil {
		return nil, fmt.Errorf("failed to dial remote: %w", err)
	}

	remote := kite.NewConn(packet2.Clientbound, serverConn, p.handleServerPacket)
	p.remote = remote // TODO: this whole function is bad
	go remote.ReadLoop()

	// Handshake immediately, then we are in login.
	handshake := &packet2.ClientHandshake{
		ProtocolVersion: 768,
		ServerAddress:   address,
		ServerPort:      port,
		Intent:          packet2.IntentLogin,
	}
	if err = remote.SendPacket(handshake); err != nil {
		return nil, err
	}
	remote.SetState(packet2.Login)

	// Begin login
	err = remote.SendPacket(&packet2.ClientLoginStart{
		Name: p.Username,
		UUID: p.UUID.String(),
	})
	if err != nil {
		panic(err)
	}

	// todo this is a weird solution, pendingloginchan being closed would also trigger this.
	select {
	case err = <-p.pendingLoginChan:
	case <-time.After(30 * time.Second):
		err = errors.New("login timed out")
	}
	if err != nil {
		remote.Close()
		return nil, err
	}

	return remote, nil
}

func (p *Player) handleServerLoginPacket(pp proto.Packet) (err error) {
	//todo we should handle encryption request here to create an error that the backend server is in online-mode
	switch pp.Id {
	case packet2.ServerLoginDisconnectID:
		pkt := new(packet2.ServerLoginDisconnect)
		if err = pp.Read(pkt); err != nil {
			return err
		}
		return p.handleServerLoginDisconnect(pkt)
	case packet2.ServerLoginPluginRequestID:
		pkt := new(packet2.ServerLoginPluginRequest)
		if err = pp.Read(pkt); err != nil {
			return err
		}
		return p.handleServerLoginPluginRequest(pkt)
	case packet2.ServerLoginLoginSuccessID:
		pkt := new(packet2.ServerLoginSuccess)
		if err = pp.Read(pkt); err != nil {
			return err
		}
		return p.handleServerLoginSuccess(pkt)
	default:
		return proto.UnknownPacket
	}
}

func (p *Player) handleServerLoginDisconnect(pkt *packet2.ServerLoginDisconnect) error {
	p.pendingLoginChan <- fmt.Errorf("disconnect: %s", pkt.Reason)
	return nil
}

func (p *Player) handleServerLoginPluginRequest(pkt *packet2.ServerLoginPluginRequest) error {
	if pkt.Channel != "velocity:player_info" {
		println("unhandled plugin request", pkt.Channel)
		return p.remote.SendPacket(&packet2.ClientLoginPluginResponse{
			MessageID: pkt.MessageID,
			Data:      nil, // Unhandled message
		})
	}

	requestVersion := velocity.DefaultForwardingVersion
	if len(pkt.Data) > 0 {
		requestVersion = int(pkt.Data[0])
	}
	forward, err := velocity.CreateSignedForwardingData(requestVersion, []byte(p.proxy.VelocitySecret), "127.0.0.1", p.Profile)
	if err != nil {
		return err
	}

	println("responding to velocity forwarding request")
	return p.remote.SendPacket(&packet2.ClientLoginPluginResponse{
		MessageID: pkt.MessageID,
		Data:      forward,
	})
}

func (p *Player) handleServerLoginSuccess(pkt *packet2.ServerLoginSuccess) error {
	err := p.remote.SendPacket(&packet2.ClientLoginAcknowledged{})
	if err != nil {
		return err
	}

	// Yay! We are connected to the remote server
	p.remote.SetState(packet2.Config)
	println("login success")
	p.pendingLoginChan <- nil
	return nil
}
