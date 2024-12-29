package main

import (
	packet2 "github.com/mworzala/kite/pkg/packet"
	"github.com/mworzala/kite/pkg/proto"
)

func (p *Player) handleClientHandshakePacket(pp proto.Packet) (err error) {
	if pp.Id == packet2.ClientHandshakeHandshakeID {
		pkt := new(packet2.ClientHandshake)
		if err = pp.Read(pkt); err != nil {
			return err
		}
		return p.handleHandshake(pkt)
	}

	return proto.UnknownPacket
}

func (p *Player) handleHandshake(pkt *packet2.ClientHandshake) error {
	switch pkt.Intent {
	case packet2.IntentStatus:
		p.conn.SetState(packet2.Status)
	case packet2.IntentLogin:
		p.conn.SetState(packet2.Login)
	case packet2.IntentTransfer:
		p.Disconnect("Transfer not supported")
	default:
		// An invalid intent would be caught when reading the packet, so this is impossible.
		panic("unreachable")
	}
	return nil
}
