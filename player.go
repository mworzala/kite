package kite

import (
	"github.com/mworzala/kite/pkg/proto"
	"github.com/mworzala/kite/pkg/proto/packet"
)

type Player struct {
	*proto.Conn
}

// SendPacket sends a packet to the client connection.
//
// To send a server packet, first get the Server().
func (p *Player) SendPacket(pkt packet.Packet) error {
	return p.Conn.SendPacket(pkt)
}

// Server returns the current server for the player, or nil if they are not connected to a server.
// todo
func (p *Player) Server() any {
	return nil
}
