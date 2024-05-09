package handler

import (
	"fmt"

	"github.com/mworzala/kite"
	"github.com/mworzala/kite/internal/pkg/velocity"
	"github.com/mworzala/kite/pkg/proto"
	"github.com/mworzala/kite/pkg/proto/packet"
)

var _ proto.Handler = (*ClientboundLoginHandler)(nil)

type ClientboundLoginHandler struct {
	Player *kite.Player
	Remote *proto.Conn
	DoneCh chan bool
}

func NewClientboundLoginHandler(p *kite.Player, remote *proto.Conn, doneCh chan bool) proto.Handler {
	return &ClientboundLoginHandler{p, remote, doneCh}
}

func (h *ClientboundLoginHandler) HandlePacket(pp proto.Packet) (err error) {
	//todo we should handle encryption request here to create an error that the backend server is in online-mode
	switch pp.Id {
	case packet.ServerLoginDisconnectID:
		p := new(packet.ServerLoginDisconnect)
		if err = pp.Read(p); err != nil {
			return err
		}
		return h.handleDisconnect(p)
	case packet.ServerLoginPluginRequestID:
		p := new(packet.ServerLoginPluginRequest)
		if err = pp.Read(p); err != nil {
			return err
		}
		return h.handlePluginRequest(p)
	case packet.ServerLoginLoginSuccessID:
		p := new(packet.ServerLoginSuccess)
		if err = pp.Read(p); err != nil {
			return err
		}
		return h.handleLoginSuccess(p)
	}

	return proto.UnknownPacket
}

func (h *ClientboundLoginHandler) handleDisconnect(p *packet.ServerLoginDisconnect) error {
	println("disconnect for", p.Reason)
	h.Player.Close()
	return nil
}

func (h *ClientboundLoginHandler) handlePluginRequest(p *packet.ServerLoginPluginRequest) error {
	if p.Channel != "velocity:player_info" {
		return h.Remote.SendPacket(&packet.ClientLoginPluginResponse{
			MessageID: p.MessageID,
			Data:      nil, // Unhandled message
		})
	}

	// If we don't have a player profile yet for some reason this is a problem
	// This is a sanity check, there isn't a reason (i know of) not to have a profile yet.
	if h.Player.Profile == nil {
		return fmt.Errorf("missing player profile")
	}

	requestVersion := velocity.DefaultForwardingVersion
	if len(p.Data) > 0 {
		requestVersion = int(p.Data[0])
	}
	forward, err := velocity.CreateSignedForwardingData([]byte("test12345"), h.Player.Profile, requestVersion)
	if err != nil {
		return err
	}

	return h.Remote.SendPacket(&packet.ClientLoginPluginResponse{
		MessageID: p.MessageID,
		Data:      forward,
	})
}

func (h *ClientboundLoginHandler) handleLoginSuccess(p *packet.ServerLoginSuccess) error {
	err := h.Remote.SendPacket(&packet.ClientLoginAcknowledged{})
	if err != nil {
		return err
	}

	h.Remote.SetRemote(h.Player.Conn)
	h.Player.Conn.SetRemote(h.Remote)
	h.Remote.SetState(packet.Config, NewClientboundConfigurationHandler(h.Player, h.Remote))
	h.DoneCh <- true

	return nil
}

type ClientboundLoginHandler2 struct {
	Player *kite.Player
	Remote *proto.Conn
}

func NewClientboundLoginHandler2(p *kite.Player, remote *proto.Conn) proto.Handler {
	return &ClientboundLoginHandler2{p, remote}
}

func (h *ClientboundLoginHandler2) HandlePacket(pp proto.Packet) (err error) {
	switch pp.Id {
	case packet.ServerLoginLoginSuccessID:
		p := new(packet.ServerLoginSuccess)
		if err = pp.Read(p); err != nil {
			return err
		}
		return h.handleLoginSuccess(p)
	}

	return proto.UnknownPacket
}

func (h *ClientboundLoginHandler2) handleLoginSuccess(p *packet.ServerLoginSuccess) error {
	err := h.Remote.SendPacket(&packet.ClientLoginAcknowledged{})
	if err != nil {
		return err
	}

	// Disconnect them from their old server
	oldRemote := h.Player.Conn.GetRemote()
	h.Player.Conn.SetRemote(nil)
	oldRemote.SetRemote(nil)
	oldRemote.Close()

	doneCh := make(chan bool)
	h.Player.SetState(packet.Play, &WaitForStartConfigHandler{
		Player: h.Player,
		Remote: h.Remote,
		DoneCh: doneCh,
	})
	if err = h.Player.SendPacket(&packet.ServerStartConfiguration{}); err != nil {
		return err
	}

	<-doneCh

	return nil
}

type WaitForStartConfigHandler struct {
	Player *kite.Player
	Remote *proto.Conn
	DoneCh chan bool
}

func (h *WaitForStartConfigHandler) HandlePacket(pp proto.Packet) (err error) {
	switch pp.Id {
	case packet.ClientPlayConfigurationAckID:
		h.Remote.SetRemote(h.Player.Conn)
		h.Player.Conn.SetRemote(h.Remote)

		h.Player.SetState(packet.Config, NewServerboundConfigurationHandler(h.Player))
		h.Remote.SetState(packet.Config, NewClientboundConfigurationHandler(h.Player, h.Remote))

		h.DoneCh <- true

		return nil
	}
	return nil // Eat any other packet for now
}
