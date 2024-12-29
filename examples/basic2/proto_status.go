package main

import (
	"encoding/json"

	"github.com/mworzala/kite/pkg/packet"
	"github.com/mworzala/kite/pkg/proto"
)

func (p *Player) handleClientStatusPacket(pp proto.Packet) (err error) {
	switch pp.Id {
	case packet.ClientStatusPingRequestID:
		pkt := new(packet.ClientStatusPingRequest)
		if err = pp.Read(pkt); err != nil {
			return err
		}
		return p.handlePingRequest(pkt)
	case packet.ClientStatusStatusRequestID:
		pkt := new(packet.ClientStatusRequest)
		if err = pp.Read(pkt); err != nil {
			return err
		}
		return p.handleStatusRequest(pkt)
	}
	return proto.UnknownPacket
}

func (p *Player) handlePingRequest(pkt *packet.ClientStatusPingRequest) error {
	return p.conn.SendPacket(&packet.ServerStatusPingResponse{Payload: pkt.Payload})
}

func (p *Player) handleStatusRequest(pkt *packet.ClientStatusRequest) error {
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
