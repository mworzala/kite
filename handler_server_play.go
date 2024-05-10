package kite

import (
	"github.com/mworzala/kite/pkg/proto"
	"github.com/mworzala/kite/pkg/proto/packet"
)

type ServerPlayHandler struct {
	Player *Player
}

func NewServerPlayHandler(p *Player) proto.Handler {
	return &ServerPlayHandler{p}
}

func (h *ServerPlayHandler) HandlePacket(pp proto.Packet) (err error) {
	//println("serverbound play packet", pp.Id)

	if pp.Id == packet.ServerPlayPlayerChatID {
		go func() {
			err := h.Player.ConnectTo(&ServerInfo{
				Address: "localhost",
				Port:    25566,
			})
			if err != nil {
				println("failed to connect to server", err.Error())
				return
			}

			println("connected to new server!!")
		}()
	}

	return proto.Forward
}

var _ proto.Handler = (*ServerPlayHandler)(nil)
