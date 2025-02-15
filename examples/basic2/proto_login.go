package main

import (
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/mworzala/kite"
	"github.com/mworzala/kite/pkg/mojangutil"
	"github.com/mworzala/kite/pkg/packet"
	"github.com/mworzala/kite/pkg/velocity"
)

func (p *Player) handleClientLoginPacket(pb kite.PacketBuffer) (err error) {
	switch pb.Id {
	case packet.ClientLoginLoginStartID:
		pkt := new(packet.ClientLoginStart)
		if err = pb.Read(pkt); err != nil {
			return err
		}
		return p.handleClientLoginStart(pkt)
	case packet.ClientLoginLoginAcknowledgedID:
		pkt := new(packet.ClientLoginAcknowledged)
		if err = pb.Read(pkt); err != nil {
			return err
		}
		return p.handleClientLoginAcknowledged(pkt)
	case packet.ClientLoginEncryptionResponseID:
		pkt := new(packet.ClientEncryptionResponse)
		if err = pb.Read(pkt); err != nil {
			return err
		}
		return p.handleClientEncryptionResponse(pkt)
	default:
		return errors.New("unexpected login packet")
	}
}

func (p *Player) handleClientLoginStart(pkt *packet.ClientLoginStart) (err error) {
	p.Username = pkt.Name

	return p.conn.SendPacket(&packet.ServerEncryptionRequest{
		ServerID:           "",
		PublicKey:          p.proxy.MojKeyPair.PublicKey(),
		VerifyToken:        p.conn.GetNonce(),
		ShouldAuthenticate: true,
	})
}

func (p *Player) handleClientEncryptionResponse(pkt *packet.ClientEncryptionResponse) (err error) {
	// Complete server side auth and respond with their profile.
	profile, err := mojangutil.HandleEncryptionResponse(p.conn, p.proxy.MojKeyPair, p.Username, pkt.VerifyToken, pkt.SharedSecret)

	p.UUID = profile.ID
	p.Username = profile.Name
	p.Profile = &profile

	return p.conn.SendPacket(&packet.ServerLoginSuccess{
		GameProfile: profile,
	})
}

func (p *Player) handleClientLoginAcknowledged(_ *packet.ClientLoginAcknowledged) (err error) {
	// This should never happen in normal operation, but a client could just send a login ack
	// immediately in an attempt to bypass auth. So don't let that happen :)
	if p.Profile == nil {
		return fmt.Errorf("missing player profile")
	}

	p.remote, err = p.connectServerSync("localhost", 25565)
	if err != nil {
		panic(err) // TODO
	}

	// Server connection is already in config
	p.conn.SetState(packet.Config)

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

	remote := kite.NewConn(packet.Clientbound, serverConn, p.handleServerPacket)
	p.remote = remote // TODO: this whole function is bad
	go remote.ReadLoop()

	// Handshake immediately, then we are in login.
	handshake := &packet.ClientHandshake{
		ProtocolVersion: 768,
		ServerAddress:   address,
		ServerPort:      port,
		Intent:          packet.IntentLogin,
	}
	if err = remote.SendPacket(handshake); err != nil {
		return nil, err
	}
	remote.SetState(packet.Login)

	// Begin login
	err = remote.SendPacket(&packet.ClientLoginStart{
		Name: p.Username,
		UUID: p.UUID,
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

func (p *Player) handleServerLoginPacket(pb kite.PacketBuffer) (err error) {
	//todo we should handle encryption request here to create an error that the backend server is in online-mode
	switch pb.Id {
	case packet.ServerLoginDisconnectID:
		pkt := new(packet.ServerLoginDisconnect)
		if err = pb.Read(pkt); err != nil {
			return err
		}
		return p.handleServerLoginDisconnect(pkt)
	case packet.ServerLoginEncryptionRequestID:
		pb.Consume()
		println("Server requested authorization, is it in offline mode?")
		p.Disconnect("An error has occurred")
		return nil
	case packet.ServerLoginPluginRequestID:
		pkt := new(packet.ServerLoginPluginRequest)
		if err = pb.Read(pkt); err != nil {
			return err
		}
		return p.handleServerLoginPluginRequest(pkt)
	case packet.ServerLoginLoginSuccessID:
		pkt := new(packet.ServerLoginSuccess)
		if err = pb.Read(pkt); err != nil {
			return err
		}
		return p.handleServerLoginSuccess(pkt)
	default:
		return errors.New("unexpected login packet")
	}
}

func (p *Player) handleServerLoginDisconnect(pkt *packet.ServerLoginDisconnect) error {
	p.pendingLoginChan <- fmt.Errorf("disconnect: %s", pkt.Reason)
	return nil
}

func (p *Player) handleServerLoginPluginRequest(pkt *packet.ServerLoginPluginRequest) error {
	if pkt.Channel != "velocity:player_info" {
		println("unhandled plugin request", pkt.Channel)
		return p.remote.SendPacket(&packet.ClientLoginPluginResponse{
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
	return p.remote.SendPacket(&packet.ClientLoginPluginResponse{
		MessageID: pkt.MessageID,
		Data:      forward,
	})
}

func (p *Player) handleServerLoginSuccess(pkt *packet.ServerLoginSuccess) error {
	err := p.remote.SendPacket(&packet.ClientLoginAcknowledged{})
	if err != nil {
		return err
	}

	// Yay! We are connected to the remote server
	p.remote.SetState(packet.Config)
	println("login success")
	p.pendingLoginChan <- nil
	return nil
}
