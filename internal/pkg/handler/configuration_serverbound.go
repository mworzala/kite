package handler

import (
	"github.com/mworzala/kite"
	"github.com/mworzala/kite/pkg/proto"
	"github.com/mworzala/kite/pkg/proto/packet"
)

var _ proto.Handler = (*ServerboundConfigurationHandler)(nil)

type ServerboundConfigurationHandler struct {
	Player *kite.Player
}

func NewServerboundConfigurationHandler(p *kite.Player) proto.Handler {
	return &ServerboundConfigurationHandler{p}
}

func (h *ServerboundConfigurationHandler) HandlePacket(pp proto.Packet) (err error) {
	switch pp.Id {
	case packet.ClientConfigPluginMessageID:
		p := new(packet.ClientConfigPluginMessage)
		if err = pp.Read(p); err != nil {
			return err
		}
		return h.handlePluginMessage(p)
	case packet.ClientConfigFinishConfigurationID:
		println("client finished configuration")
		h.Player.Conn.SetState(packet.Play, NewClientboundPlayHandler(h.Player))
		return proto.Forward
	}
	return proto.Forward
}

func (h *ServerboundConfigurationHandler) handlePluginMessage(p *packet.ClientConfigPluginMessage) error {
	//println("PLUGIN MESSAGE ", p.Channel, string(p.Data))

	return proto.Forward
}
