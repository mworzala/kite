package kite

import (
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
	resp := &packet.ServerStatusResponse{
		Status: `
{
    "version": {
        "name": "1.20.6",
        "protocol": 766
    },
    "players": {
        "max": 100,
        "online": 5,
        "sample": [
            {
                "name": "thinkofdeath",
                "id": "4566e69f-c907-48ee-8d71-d7ba5aa00d20"
            }
        ]
    },
    "description": {
        "text": "Hello world"
    },
    "favicon": "data:image/png;base64,<data>",
    "enforcesSecureChat": false
}
`,
	}
	return h.Conn.SendPacket(resp)
}

var _ proto.Handler = (*ClientStatusHandler)(nil)
