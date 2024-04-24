package handler

import (
	"github.com/mworzala/kite/pkg/proto"
	"github.com/mworzala/kite/pkg/proto/packet"
	"github.com/mworzala/kite/pkg/proxy"
)

var _ proto.Handler = (*ServerboundStatusHandler)(nil)

type ServerboundStatusHandler struct {
	Player *proxy.Player
}

func NewServerboundStatusHandler(p *proxy.Player) proto.Handler {
	return &ServerboundStatusHandler{p}
}

func (h *ServerboundStatusHandler) HandlePacket(pp proto.Packet) (err error) {
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

func (h *ServerboundStatusHandler) handlePingRequest(p *packet.ClientStatusPingRequest) error {
	resp := &packet.ServerStatusPingResponse{Payload: p.Payload}
	return h.Player.SendPacket(resp)
}

func (h *ServerboundStatusHandler) handleStatusRequest(p *packet.ClientStatusRequest) error {
	resp := &packet.ServerStatusResponse{
		Status: `
{
    "version": {
        "name": "1.20.5",
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
	return h.Player.SendPacket(resp)
}
