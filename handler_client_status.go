package kite

import (
	"encoding/json"

	"github.com/mworzala/kite/pkg/proto"
	"github.com/mworzala/kite/pkg/proto/packet"
)

type ClientStatusHandler struct {
	Conn *proto.Conn
}

func NewClientStatusHandler(conn *proto.Conn) proto.Handler {
	return &ClientStatusHandler{conn}
}

func (h *ClientStatusHandler) HandlePacket(pp proto.Packet) (err error) {
	switch pp.Id {
	case packet.ClientStatusPingRequestID:
		p := new(packet.ClientStatusPingRequest)
		if err = pp.Read(p); err != nil {
			return err
		}
		return h.handlePingRequest(p)
	case packet.ClientStatusStatusRequestID:
		p := new(packet.ClientStatusRequest)
		if err = pp.Read(p); err != nil {
			return err
		}
		return h.handleStatusRequest(p)
	}
	return proto.UnknownPacket
}

func (h *ClientStatusHandler) handlePingRequest(p *packet.ClientStatusPingRequest) error {
	return h.Conn.SendPacket(&packet.ServerStatusPingResponse{Payload: p.Payload})
}

func (h *ClientStatusHandler) handleStatusRequest(p *packet.ClientStatusRequest) error {
	return h.Conn.SendPacket(&packet.ServerStatusResponse{Payload: packet.StatusResponse{
		Version: packet.ServerVersion{
			Name:     "1.20.6",
			Protocol: 766,
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

var _ proto.Handler = (*ClientStatusHandler)(nil)
