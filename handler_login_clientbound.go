package kite

import (
	"context"

	"github.com/mworzala/kite/pkg/proto"
	"github.com/mworzala/kite/pkg/proto/packet"
)

type ClientboundLoginHandler2[T any] struct {
	Player *Player[T]
	Remote *proto.Conn
}

func NewClientboundLoginHandler2[T any](p *Player[T], remote *proto.Conn) proto.Handler {
	return &ClientboundLoginHandler2[T]{p, remote}
}

func (h *ClientboundLoginHandler2[T]) HandlePacket(pp proto.Packet) (err error) {
	switch pp.Id {
	case packet.ServerLoginLoginSuccessID:
		p := new(packet.ServerLoginSuccess)
		if err = pp.Read(p); err != nil {
			return err
		}
		return h.handleLoginSuccess(p)
	default:
		return proto.UnknownPacket
	}
}

func (h *ClientboundLoginHandler2[T]) handleLoginSuccess(p *packet.ServerLoginSuccess) error {
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
	h.Player.SetState(packet.Play, &WaitForStartConfigHandler[T]{
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

type WaitForStartConfigHandler2 struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func (h *WaitForStartConfigHandler2) HandlePacket(pp proto.Packet) (err error) {
	switch pp.Id {
	case packet.ClientPlayConfigurationAckID:
		p := new(packet.ClientConfigurationAck)
		if err = pp.Read(p); err != nil {
			return
		}
		return h.handleConfigAck(p)
	default:
		return nil // Eat any other packet for now
	}
}

func (h *WaitForStartConfigHandler2) handleConfigAck(_ *packet.ClientConfigurationAck) error {
	h.cancel()
	//h.Remote.SetRemote(h.Player.Conn)
	//h.Player.Conn.SetRemote(h.Remote)
	//
	//h.Player.SetState(packet.Config, NewClientConfigHandler(h.Player))
	//h.Remote.SetState(packet.Config, NewServerConfigHandler(h.Player, h.Remote))
	//
	//h.DoneCh <- true
	return nil
}

type WaitForStartConfigHandler[T any] struct {
	Player *Player[T]
	Remote *proto.Conn
	DoneCh chan bool
}

func (h *WaitForStartConfigHandler[T]) HandlePacket(pp proto.Packet) (err error) {
	switch pp.Id {
	case packet.ClientPlayConfigurationAckID:
		p := new(packet.ClientConfigurationAck)
		if err = pp.Read(p); err != nil {
			return
		}
		return h.handleConfigAck(p)
	default:
		return nil // Eat any other packet for now
	}
}

func (h *WaitForStartConfigHandler[T]) handleConfigAck(_ *packet.ClientConfigurationAck) error {
	h.Remote.SetRemote(h.Player.Conn)
	h.Player.Conn.SetRemote(h.Remote)

	h.Player.SetState(packet.Config, NewClientConfigHandler(h.Player))
	h.Remote.SetState(packet.Config, NewServerConfigHandler(h.Player))

	h.DoneCh <- true
	return nil
}
