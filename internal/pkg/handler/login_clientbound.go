package handler

import (
	"github.com/mworzala/kite/pkg/proto"
	"github.com/mworzala/kite/pkg/proto/packet"
	"github.com/mworzala/kite/pkg/proxy"
)

var _ proto.Handler = (*ServerboundLoginHandler)(nil)

type ClientboundLoginHandler struct {
	Player *proxy.Player
	Remote *proto.Conn
	DoneCh chan bool
}

func NewClientboundLoginHandler(p *proxy.Player, remote *proto.Conn, doneCh chan bool) proto.Handler {
	return &ClientboundLoginHandler{p, remote, doneCh}
}

func (h *ClientboundLoginHandler) HandlePacket(pp proto.Packet) (err error) {
	switch pp.Id {
	case packet.ServerLoginLoginSuccessID:
		p := new(packet.ServerLoginSuccess)
		if err = pp.Read(p); err != nil {
			return err
		}
		return h.handleLoginSuccess(p)
	}
	return proto.UnknownPacket
}

func (h *ClientboundLoginHandler) handleLoginSuccess(p *packet.ServerLoginSuccess) error {
	err := h.Remote.SendPacket(&packet.ClientLoginAcknowledged{})
	if err != nil {
		return err
	}

	h.Remote.SetRemote(h.Player.Conn)
	h.Player.Conn.SetRemote(h.Remote)
	h.Remote.SetState(packet.Config, NewClientboundConfigurationHandler(h.Player, h.Remote))
	h.DoneCh <- true

	return nil
}
