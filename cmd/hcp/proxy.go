package main

import (
	"context"

	"github.com/mworzala/kite/pkg/client"
	"github.com/mworzala/kite/pkg/proto/packet"
	"github.com/mworzala/kite/pkg/proxy"
)

type ProxyImpl struct {
	proxy *proxy.Proxy
}

func NewProxyImpl() (*ProxyImpl, error) {
	p := &ProxyImpl{proxy: proxy.New()}

	proxy.SetClientPacketHandler(p.proxy, packet.Handshake, p.handleHandshake)
	proxy.SetClientPacketHandler(p.proxy, packet.Status, p.handleStatusRequest)
	proxy.SetClientPacketHandler(p.proxy, packet.Status, p.handlePingRequest)

	return p, nil
}

func (p *ProxyImpl) Start() error {
	return p.proxy.Start(context.Background())
}

func (p *ProxyImpl) handleHandshake(c *client.Connection, pkt *packet.ClientHandshake) error {
	c.SetState(pkt.NextState)
	return nil
}

func (p *ProxyImpl) handleStatusRequest(c *client.Connection, pkt *packet.ClientStatusRequest) error {
	err := c.WritePacket(&packet.ServerStatusResponse{
		Status: `
{
    "version": {
        "name": "1.20.4",
        "protocol": 765
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
	})
	if err != nil {
		panic(err)
	}
	return nil
}

func (p *ProxyImpl) handlePingRequest(c *client.Connection, pkt *packet.ClientStatusPingRequest) error {
	if err := c.WritePacket(&packet.ServerStatusPingResponse{
		Payload: pkt.Payload,
	}); err != nil {
		panic(err)
	}

	return c.Close()
}
