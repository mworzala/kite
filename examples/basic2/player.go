package main

import (
	"github.com/google/uuid"
	"github.com/mworzala/kite"
	"github.com/mworzala/kite/pkg/proto"
	"github.com/mworzala/kite/pkg/proto/packet"
)

type Player struct {
	proxy *Proxy
	conn  *kite.Conn

	UUID     uuid.UUID
	Username string
	Profile  *packet.GameProfile

	pendingLoginChan chan error
	remote           *kite.Conn
}

func (p *Player) Disconnect(reason string) {
	//todo take chat component message
	println("disconnect")
	p.conn.Close()
}

func (p *Player) handleClientPacket(pp proto.Packet) error {
	switch pp.State {
	case packet.Handshake:
		return p.handleClientHandshakePacket(pp)
	case packet.Status:
		return p.handleClientStatusPacket(pp)
	case packet.Login:
		return p.handleClientLoginPacket(pp)
	case packet.Config:
		return p.handleClientConfigPacket(pp)
	case packet.Play:
		if p.remote == nil {
			panic("bad state")
		}
		return p.remote.ForwardPacket(pp)
	}

	println("new packet", pp.Id)
	pp.Consume()
	return nil
}

func (p *Player) handleServerPacket(pp proto.Packet) error {
	switch pp.State {
	case packet.Login:
		return p.handleServerLoginPacket(pp)
	case packet.Config:
		return p.handleServerConfigPacket(pp)
	case packet.Play:
		return p.conn.ForwardPacket(pp)
	default:
		return proto.UnknownPacket
	}
}
