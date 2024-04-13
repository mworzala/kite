package handler

import (
	"net"

	"github.com/mworzala/kite/pkg/proto"
	"github.com/mworzala/kite/pkg/proto/packet"
	"github.com/mworzala/kite/pkg/proxy"
)

var _ proto.Handler = (*ServerboundLoginHandler)(nil)

type ServerboundLoginHandler struct {
	Player *proxy.Player
}

func NewServerboundLoginHandler(p *proxy.Player) proto.Handler {
	return &ServerboundLoginHandler{p}
}

func (h *ServerboundLoginHandler) HandlePacket(pp proto.Packet) (err error) {
	switch pp.Id {
	case packet.ClientLoginLoginStartID:
		p := new(packet.ClientLoginStart)
		if err = pp.Read(p); err != nil {
			return err
		}
		return h.handleLoginStart(p)
	case packet.ClientLoginLoginAcknowledgedID:
		p := new(packet.ClientLoginAcknowledged)
		if err = pp.Read(p); err != nil {
			return err
		}
		return h.handleLoginAcknowledged(p)
	}
	return proto.UnknownPacket
}

func (h *ServerboundLoginHandler) handleLoginStart(p *packet.ClientLoginStart) error {
	resp := &packet.ServerLoginSuccess{
		UUID:     "3bc51b9d-52be-4c4a-a3d6-7cc0bd6e6ea8",
		Username: "notmattw",
	}
	return h.Player.SendPacket(resp)
}

func (h *ServerboundLoginHandler) handleLoginAcknowledged(p *packet.ClientLoginAcknowledged) error {

	serverConn, err := net.Dial("tcp", "localhost:25565")
	if err != nil {
		panic(err)
	}

	doneCh := make(chan bool)
	remote, readLoop := proto.NewConn(packet.Clientbound, serverConn)
	remote.SetState(packet.Handshake, nil)
	err = remote.SendPacket(&packet.ClientHandshake{
		ProtocolVersion: 1073742009,
		ServerAddress:   "localhost:25577",
		ServerPort:      25577,
		NextState:       packet.Login,
	})
	if err != nil {
		panic(err)
	}
	remote.SetState(packet.Login, NewClientboundLoginHandler(h.Player, remote, doneCh))
	err = remote.SendPacket(&packet.ClientLoginStart{
		Name: "notmattw",
		UUID: "3bc51b9d-52be-4c4a-a3d6-7cc0bd6e6ea8",
	})
	if err != nil {
		panic(err)
	}
	go readLoop()

	<-doneCh

	h.Player.SetState(packet.Config, NewServerboundConfigurationHandler(h.Player))
	return nil
}
