package kite

import (
	"bytes"
	"fmt"

	buffer2 "github.com/mworzala/kite/internal/pkg/buffer"
	"github.com/mworzala/kite/pkg/proto"
	"github.com/mworzala/kite/pkg/proto/binary"
	"github.com/mworzala/kite/pkg/proto/packet"
)

var _ proto.Handler = (*ServerConfigHandler[any])(nil)

type ServerConfigHandler[T any] struct {
	Player *Player[T]

	PlayHandlerFunc func(*Player[T]) proto.Handler
}

func NewServerConfigHandler[T any](p *Player[T]) proto.Handler {
	return &ServerConfigHandler[T]{
		Player: p,
	}
}

func (h *ServerConfigHandler[T]) HandlePacket(pp proto.Packet) (err error) {
	switch pp.Id {
	case packet.ServerConfigPluginMessageID:
		p := new(packet.ServerPluginMessage)
		if err = pp.Read(p); err != nil {
			return err
		}
		return h.HandlePluginMessage(p)
	case packet.ServerConfigFinishConfigurationID:
		println("server moving to play")
		if h.PlayHandlerFunc != nil {
			h.Player.Server().SetState(packet.Play, h.PlayHandlerFunc(h.Player))
		} else {
			h.Player.Server().SetState(packet.Play, NewServerPlayHandler(h.Player))
		}
		return proto.Forward
	default:
		println("unhandled config packet from server", pp.Id)
		return proto.Forward
	}
}

func (h *ServerConfigHandler[T]) HandlePluginMessage(p *packet.ServerPluginMessage) error {
	if p.Channel == "minecraft:brand" {
		oldPayload := buffer2.NewPacketBuffer(p.Data)
		oldBrand, err := binary.ReadSizedString(oldPayload, 32767)
		if err != nil {
			return err
		}

		var newPayload bytes.Buffer
		err = binary.WriteSizedString(&newPayload, fmt.Sprintf("%s // %s", oldBrand, "kite"), 32767)
		if err != nil {
			return err
		}

		resp := &packet.ServerPluginMessage{
			Channel: "minecraft:brand",
			Data:    newPayload.Bytes(),
		}
		return h.Player.SendPacket(resp)
	}

	return proto.Forward
}
