package kite

import (
	"github.com/mworzala/kite/pkg/proto"
	"github.com/mworzala/kite/pkg/proto/packet"
)

var _ proto.Handler = (*ClientPlayHandler[any])(nil)

type ClientPlayHandler[T any] struct {
	Player *Player[T]
}

func NewClientPlayHandler[T any](p *Player[T]) proto.Handler {
	return &ClientPlayHandler[T]{p}
}

func (h *ClientPlayHandler[T]) HandlePacket(pp proto.Packet) (err error) {
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
