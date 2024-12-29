package main

import (
	"github.com/mworzala/kite/pkg/proto"
	"github.com/mworzala/kite/pkg/proto/packet"
)

func (p *Player) handleClientHandshakePacket(pp proto.Packet) (err error) {
	if pp.Id == packet.ClientHandshakeHandshakeID {
		pkt := new(packet.ClientHandshake)
		if err = pp.Read(pkt); err != nil {
			return err
		}
		return p.handleHandshake(pkt)
	}

	return proto.UnknownPacket
}

func (p *Player) handleHandshake(pkt *packet.ClientHandshake) error {
	switch pkt.Intent {
	case packet.IntentStatus:
		p.conn.SetState(packet.Status)
	case packet.IntentLogin:
		p.conn.SetState(packet.Login)
	case packet.IntentTransfer:
		p.Disconnect("Transfer not supported")
	default:
		// An invalid intent would be caught when reading the packet, so this is impossible.
		panic("unreachable")
	}
	return nil
}
