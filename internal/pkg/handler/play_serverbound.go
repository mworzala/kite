package handler

import (
	"net"

	"github.com/mworzala/kite/pkg/proto"
	"github.com/mworzala/kite/pkg/proto/packet"
	"github.com/mworzala/kite/pkg/proxy"
)

var _ proto.Handler = (*ServerboundLoginHandler)(nil)

type ServerboundPlayHandler struct {
	Player *proxy.Player
}

func NewServerboundPlayHandler(p *proxy.Player) proto.Handler {
	return &ServerboundPlayHandler{p}
}

func (h *ServerboundPlayHandler) HandlePacket(pp proto.Packet) (err error) {
	//println("serverbound play packet", pp.Id)

	if pp.Id == packet.ServerPlayPlayerChatID {
		println("OPENING NEW SERVER CONN")
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
			NextState:       packet.Login,
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
