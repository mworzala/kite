package handler

import (
	"net"

	"github.com/mworzala/kite"
	"github.com/mworzala/kite/pkg/proto"
	"github.com/mworzala/kite/pkg/proto/packet"
)

type ServerPlayHandler struct {
	Player *kite.Player
}

func NewServerPlayHandler(p *kite.Player) proto.Handler {
	return &ServerPlayHandler{p}
}

func (h *ServerPlayHandler) HandlePacket(pp proto.Packet) (err error) {
	//println("serverbound play packet", pp.Id)

	if pp.Id == packet.ServerPlayPlayerChatID {
		serverConn, err := net.Dial("tcp", "localhost:25566")
		if err != nil {
			panic(err)
		}

		remote, readLoop := proto.NewConn(packet.Clientbound, serverConn)
		remote.SetState(packet.Handshake, nil)
		err = remote.SendPacket(&packet.ClientHandshake{
			ProtocolVersion: 766,
			ServerAddress:   "localhost:25577",
			ServerPort:      25577,
			Intent:          packet.IntentLogin,
		})
		if err != nil {
			panic(err)
		}
		remote.SetState(packet.Login, NewClientboundLoginHandler2(h.Player, remote))
		err = remote.SendPacket(&packet.ClientLoginStart{
			Name: "notmattw",
			UUID: "3bc51b9d-52be-4c4a-a3d6-7cc0bd6e6ea8",
		})
		if err != nil {
			panic(err)
		}
		go readLoop()
	}

	return proto.Forward
}

var _ proto.Handler = (*ServerPlayHandler)(nil)
