package kite

import (
	"bytes"
	"fmt"

	buffer2 "github.com/mworzala/kite/internal/pkg/buffer"
	"github.com/mworzala/kite/pkg/proto"
	"github.com/mworzala/kite/pkg/proto/binary"
	"github.com/mworzala/kite/pkg/proto/packet"
)

var _ proto.Handler = (*ServerConfigHandler)(nil)

type ServerConfigHandler struct {
	Player *Player
	Remote *proto.Conn
}

func NewServerConfigHandler(p *Player, remote *proto.Conn) proto.Handler {
	return &ServerConfigHandler{p, remote}
}

func (h *ServerConfigHandler) HandlePacket(pp proto.Packet) (err error) {
	switch pp.Id {
	case packet.ServerConfigPluginMessageID:
		p := new(packet.ClientConfigPluginMessage)
		if err = pp.Read(p); err != nil {
			return err
		}
		return h.handlePluginMessage(p)
	case packet.ServerConfigFinishConfigurationID:
		println("server finished configuration")
		h.Remote.SetState(packet.Play, NewServerPlayHandler(h.Player))
		return proto.Forward
	}
	return proto.Forward
}

func (h *ServerConfigHandler) handlePluginMessage(p *packet.ClientConfigPluginMessage) error {
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

		resp := &packet.ServerConfigPluginMessage{
			Channel: "minecraft:brand",
			Data:    newPayload.Bytes(),
		}
		return h.Player.SendPacket(resp)
	}

	return proto.Forward
}
