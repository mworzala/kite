package kite

import (
	"github.com/mworzala/kite/pkg/proto"
	"github.com/mworzala/kite/pkg/proto/packet"
)

var _ proto.Handler = (*ClientPlayHandler)(nil)

type ClientPlayHandler struct {
	Player *Player
}

func NewClientPlayHandler(p *Player) proto.Handler {
	return &ClientPlayHandler{p}
}

func (h *ClientPlayHandler) HandlePacket(pp proto.Packet) (err error) {
	switch pp.Id {
	case packet.ClientPlayChatSessionUpdateID:
		// We don't handle signed chat and if this is sent to the server it disconnects the client.
		// I'm probably not properly disabling chat signing somewhere, but this also solves the problem.
		pp.Consume()
		return nil
	default:
		return proto.Forward
	}
}
