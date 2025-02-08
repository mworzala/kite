package main

import (
	"errors"
	"github.com/mworzala/kite/pkg/mojang"

	"github.com/google/uuid"
	"github.com/mworzala/kite"
	"github.com/mworzala/kite/pkg/packet"
)

type Player struct {
	proxy *Proxy
	conn  *kite.Conn

	UUID     uuid.UUID
	Username string
	Profile  *mojang.GameProfile

	pendingLoginChan chan error
	remote           *kite.Conn
}

func (p *Player) Disconnect(reason string) {
	//todo take chat component message
	println("disconnect")
	p.conn.Close()
}

func (p *Player) handleClientPacket(pp kite.PacketBuffer) error {
	switch p.conn.GetState() {
	case packet.Handshake:
		return p.handleClientHandshakePacket(pp)
	case packet.Status:
		return p.handleClientStatusPacket(pp)
	case packet.Login:
		return p.handleClientLoginPacket(pp)
	case packet.Config:
		return p.handleClientConfigPacket(pp)
	case packet.Play:
		return p.handleClientPlayPacket(pp)
	default:
		return errors.New("unexpected client state")
	}
}

func (p *Player) handleServerPacket(pp kite.PacketBuffer) error {
	switch p.remote.GetState() {
	case packet.Login:
		return p.handleServerLoginPacket(pp)
	case packet.Config:
		return p.handleServerConfigPacket(pp)
	case packet.Play:
		return p.handleServerPlayPacket(pp)
	default:
		return errors.New("unexpected server state")
	}
}
