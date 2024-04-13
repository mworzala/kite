package handler

import (
	"fmt"

	"github.com/mworzala/kite/pkg/proto"
	"github.com/mworzala/kite/pkg/proto/packet"
	"github.com/mworzala/kite/pkg/proxy"
)

var _ proto.Handler = (*ServerboundHandshakeHandler)(nil)

type ServerboundHandshakeHandler struct {
	*proxy.Player
}

func NewServerboundHandshakeHandler(p *proxy.Player) proto.Handler {
	return &ServerboundHandshakeHandler{p}
}

func (h *ServerboundHandshakeHandler) HandlePacket(pp proto.Packet) error {
	var err error
	if pp.Id == packet.ClientHandshakeHandshakeID {
		p := new(packet.ClientHandshake)
		if err = pp.Read(p); err != nil {
			return err
		}
		return h.handleHandshake(p)
	}

	return proto.UnknownPacket
}

func (h *ServerboundHandshakeHandler) handleHandshake(p *packet.ClientHandshake) error {
	switch p.NextState {
	case packet.Status:
		h.SetState(packet.Status, NewServerboundStatusHandler(h.Player))
	case packet.Login:
		h.SetState(packet.Login, NewServerboundLoginHandler(h.Player))
	default:
		return fmt.Errorf("unknown state %v", p.NextState)
	}
	return nil
}
