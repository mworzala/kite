package handler

import (
	"github.com/mworzala/kite"
	"github.com/mworzala/kite/pkg/proto"
	"github.com/mworzala/kite/pkg/proto/binary"
	"github.com/mworzala/kite/pkg/proto/packet"
)

var _ proto.Handler = (*ClientboundPlayHandler)(nil)

type ClientboundPlayHandler struct {
	Player *kite.Player
}

func NewClientboundPlayHandler(p *kite.Player) proto.Handler {
	return &ClientboundPlayHandler{p}
}

func (h *ClientboundPlayHandler) HandlePacket(pp proto.Packet) (err error) {

	if pp.Id == packet.ClientPlayChatSessionUpdateID {
		_, _ = binary.ReadRaw(pp.Buf(), binary.Remaining)

		println("DROPPING CHAT SESSION UPDATE")
		return nil
	}
	//println("clientbound play packet", pp.Id)
	return proto.Forward
}
