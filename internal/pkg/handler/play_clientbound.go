package handler

import (
	"github.com/mworzala/kite"
	"github.com/mworzala/kite/pkg/proto"
)

var _ proto.Handler = (*ServerboundLoginHandler)(nil)

type ClientboundPlayHandler struct {
	Player *kite.Player
}

func NewClientboundPlayHandler(p *kite.Player) proto.Handler {
	return &ClientboundPlayHandler{p}
}

func (h *ClientboundPlayHandler) HandlePacket(pp proto.Packet) (err error) {
	//println("clientbound play packet", pp.Id)
	return proto.Forward
}
