package main

import (
	"errors"

	"github.com/mworzala/kite"
	packet2 "github.com/mworzala/kite/pkg/packet"
)

func (p *Player) handleClientHandshakePacket(pb kite.PacketBuffer) (err error) {
	if pb.Id == packet2.ClientHandshakeHandshakeID {
		pkt := new(packet2.ClientHandshake)
		if err = pb.Read(pkt); err != nil {
			return err
		}
		return p.handleHandshake(pkt)
	}

	return errors.New("unexpected handshake packet")
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
