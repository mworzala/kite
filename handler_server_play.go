package kite

import (
	"github.com/mworzala/kite/pkg/proto"
)

type ServerPlayHandler[T any] struct {
	Player *Player[T]
}

func NewServerPlayHandler[T any](p *Player[T]) proto.Handler {
	return &ServerPlayHandler[T]{p}
}

func (h *ServerPlayHandler[T]) HandlePacket(pp proto.Packet) (err error) {
	return proto.Forward
}

var _ proto.Handler = (*ServerPlayHandler[any])(nil)
