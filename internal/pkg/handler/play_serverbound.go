package handler

import (
	"github.com/mworzala/kite/pkg/proto"
	"github.com/mworzala/kite/pkg/proxy"
)

var _ proto.Handler = (*ServerboundLoginHandler)(nil)

type ServerboundPlayHandler struct {
	Player *proxy.Player
}

func NewServerboundPlayHandler(p *proxy.Player) proto.Handler {
	return &ServerboundPlayHandler{p}
}

func (h *ServerboundPlayHandler) HandlePacket(pp proto.Packet) (err error) {
	println("serverbound play packet", pp.Id)
	return proto.Forward
}
