package kite

import (
	"github.com/mworzala/kite/pkg/proto"
	"github.com/mworzala/kite/pkg/proto/packet"
)

type ClientConfigHandler struct {
	Player *Player

	PlayHandlerFunc func(*Player) proto.Handler
}

func NewClientConfigHandler(p *Player) proto.Handler {
	return &ClientConfigHandler{Player: p}
}

func (h *ClientConfigHandler) HandlePacket(pp proto.Packet) (err error) {
	switch pp.Id {
	case packet.ClientConfigPluginMessageID:
		p := new(packet.ClientConfigPluginMessage)
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

func (h *ClientConfigHandler) HandlePluginMessage(p *packet.ClientConfigPluginMessage) error {
	return proto.Forward
}

func (h *ClientConfigHandler) HandleFinishConfiguration(_ *packet.ClientConfigFinishConfiguration) error {
	if h.PlayHandlerFunc != nil {
		h.Player.SetState(packet.Play, h.PlayHandlerFunc(h.Player))
	} else {
		h.Player.SetState(packet.Play, NewClientPlayHandler(h.Player))
	}
	return proto.Forward
}

var _ proto.Handler = (*ClientConfigHandler)(nil)
