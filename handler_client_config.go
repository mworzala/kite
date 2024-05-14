package kite

import (
	"github.com/mworzala/kite/pkg/proto"
	"github.com/mworzala/kite/pkg/proto/packet"
)

type ClientConfigHandler[T any] struct {
	Player *Player[T]

	PlayHandlerFunc func(*Player[T]) proto.Handler
}

func NewClientConfigHandler[T any](p *Player[T]) proto.Handler {
	return &ClientConfigHandler[T]{Player: p}
}

func (h *ClientConfigHandler[T]) HandlePacket(pp proto.Packet) (err error) {
	switch pp.Id {
	case packet.ClientConfigPluginMessageID:
		p := new(packet.ClientPluginMessage)
		if err = pp.Read(p); err != nil {
			return
		}
		return h.HandlePluginMessage(p)
	case packet.ClientConfigFinishConfigurationID:
		p := new(packet.ClientConfigFinishConfiguration)
		if err = pp.Read(p); err != nil {
			return
		}
		return h.HandleFinishConfiguration(p)
	default:
		return proto.Forward
	}
}

func (h *ClientConfigHandler[T]) HandlePluginMessage(p *packet.ClientPluginMessage) error {
	return proto.Forward
}

func (h *ClientConfigHandler[T]) HandleFinishConfiguration(_ *packet.ClientConfigFinishConfiguration) error {
	if h.PlayHandlerFunc != nil {
		h.Player.SetState(packet.Play, h.PlayHandlerFunc(h.Player))
	} else {
		h.Player.SetState(packet.Play, NewClientPlayHandler(h.Player))
	}
	return proto.Forward
}

var _ proto.Handler = (*ClientConfigHandler[any])(nil)
