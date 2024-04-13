package handler

import (
	"github.com/mworzala/kite/pkg/proto"
	"github.com/mworzala/kite/pkg/proxy"
)

var _ proto.Handler = (*ServerboundLoginHandler)(nil)

type ClientboundPlayHandler struct {
	Player *proxy.Player
}

func NewClientboundPlayHandler(p *proxy.Player) proto.Handler {
	return &ClientboundPlayHandler{p}
}

func (h *ClientboundPlayHandler) HandlePacket(pp proto.Packet) (err error) {
	println("clientbound play packet", pp.Id)
	return proto.Forward
}
