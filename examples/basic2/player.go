package main

import (
	"github.com/google/uuid"
	"github.com/mworzala/kite"
	packet2 "github.com/mworzala/kite/pkg/packet"
	"github.com/mworzala/kite/pkg/proto"
)

type Player struct {
	proxy *Proxy
	conn  *kite.Conn

	UUID     uuid.UUID
	Username string
	Profile  *packet2.GameProfile

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
	case packet2.Handshake:
		return p.handleClientHandshakePacket(pp)
	case packet2.Status:
		return p.handleClientStatusPacket(pp)
	case packet2.Login:
		return p.handleClientLoginPacket(pp)
	case packet2.Config:
		return p.handleClientConfigPacket(pp)
	case packet2.Play:
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
	case packet2.Login:
		return p.handleServerLoginPacket(pp)
	case packet2.Config:
		return p.handleServerConfigPacket(pp)
	case packet2.Play:
		return p.conn.ForwardPacket(pp)
	default:
		return proto.UnknownPacket
	}
}
