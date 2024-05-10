package kite

import (
	"github.com/mworzala/kite/pkg/proto"
	"github.com/mworzala/kite/pkg/proto/packet"
)

type ClientConfigHandler struct {
	Player *Player
}

func NewClientConfigHandler(p *Player) proto.Handler {
	return &ClientConfigHandler{p}
}

func (h *ClientConfigHandler) HandlePacket(pp proto.Packet) (err error) {
	switch pp.Id {
	case packet.ClientConfigPluginMessageID:
		p := new(packet.ClientConfigPluginMessage)
		if err = pp.Read(p); err != nil {
			return
		}
		return h.handlePluginMessage(p)
	case packet.ClientConfigFinishConfigurationID:
		p := new(packet.ClientConfigFinishConfiguration)
		if err = pp.Read(p); err != nil {
			return
		}
		return h.handleFinishConfiguration(p)
	default:
		return proto.Forward
	}
}

func (h *ClientConfigHandler) handlePluginMessage(p *packet.ClientConfigPluginMessage) error {
	return proto.Forward
}

func (h *ClientConfigHandler) handleFinishConfiguration(_ *packet.ClientConfigFinishConfiguration) error {
	h.Player.Conn.SetState(packet.Play, NewClientPlayHandler(h.Player))
	return proto.Forward
}

var _ proto.Handler = (*ClientConfigHandler)(nil)
