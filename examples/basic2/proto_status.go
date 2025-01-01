package main

import (
	"encoding/json"
	"errors"

	"github.com/mworzala/kite"
	"github.com/mworzala/kite/pkg/packet"
)

func (p *Player) handleClientStatusPacket(pb kite.PacketBuffer) (err error) {
	switch pb.Id {
	case packet.ClientStatusPingRequestID:
		pkt := new(packet.ClientStatusPingRequest)
		if err = pb.Read(pkt); err != nil {
			return err
		}
		return p.handlePingRequest(pkt)
	case packet.ClientStatusStatusRequestID:
		pkt := new(packet.ClientStatusRequest)
		if err = pb.Read(pkt); err != nil {
			return err
		}
		return p.handleStatusRequest(pkt)
	}
	return errors.New("unexpected status packet")
}

func (p *Player) handlePingRequest(pkt *packet.ClientStatusPingRequest) error {
	return p.conn.SendPacket(&packet.ServerStatusPingResponse{Payload: pkt.Payload})
}

func (p *Player) handleStatusRequest(_ *packet.ClientStatusRequest) error {
	return p.conn.SendPacket(&packet.ServerStatusResponse{Payload: packet.StatusResponse{
		Version: packet.ServerVersion{
			Name:     "1.21.3",
			Protocol: 768,
		},
		Players: packet.ServerPlayerList{
			Max:    1000,
			Online: 0,
		},
		Description:       json.RawMessage(`{"text": "Hello, Kite"}`),
		Favicon:           "",
		EnforceSecureChat: true,
	}})
}
