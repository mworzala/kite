package main

import (
	"bytes"
	"fmt"

	buffer2 "github.com/mworzala/kite/internal/pkg/buffer"
	"github.com/mworzala/kite/pkg/proto"
	"github.com/mworzala/kite/pkg/proto/binary"
	"github.com/mworzala/kite/pkg/proto/packet"
)

func (p *Player) handleClientConfigPacket(pp proto.Packet) (err error) {
	if p.remote == nil {
		panic("bad state")
	}
	println("client config packet", pp.Id)
	switch pp.Id {
	case packet.ClientConfigFinishConfigurationID:
		pkt := new(packet.ClientConfigFinishConfiguration)
		if err = pp.Read(pkt); err != nil {
			return
		}
		return p.handleClientConfigFinishConfiguration(pkt)
	default:
		return p.remote.ForwardPacket(pp)
	}
}

func (p *Player) handleClientConfigFinishConfiguration(pkt *packet.ClientConfigFinishConfiguration) (err error) {
	println("client moving to play")
	p.remote.SetState(packet.Play)
	return p.remote.SendPacket(pkt)
}

func (p *Player) handleServerConfigPacket(pp proto.Packet) (err error) {
	println("server config packet", pp.Id)
	switch pp.Id {
	case packet.ServerConfigPluginMessageID:
		pkt := new(packet.ServerPluginMessage)
		if err = pp.Read(pkt); err != nil {
			return err
		}
		return p.handleServerPluginMessage(pkt)
	case packet.ServerConfigFinishConfigurationID:
		//todo make actual handler
		println("server moving to play")
		p.conn.SetState(packet.Play)
		return p.conn.ForwardPacket(pp)
	default:
		return p.conn.ForwardPacket(pp)
	}
}

func (p *Player) handleServerPluginMessage(pkt *packet.ServerPluginMessage) (err error) {
	if pkt.Channel == "minecraft:brand" {
		oldPayload := buffer2.NewPacketBuffer(pkt.Data)
		oldBrand, err := binary.ReadSizedString(oldPayload, 32767)
		if err != nil {
			return err
		}

		var newPayload bytes.Buffer
		err = binary.WriteSizedString(&newPayload, fmt.Sprintf("%s // %s", oldBrand, "kite"), 32767)
		if err != nil {
			return err
		}

		resp := &packet.ServerPluginMessage{
			Channel: "minecraft:brand",
			Data:    newPayload.Bytes(),
		}
		return p.conn.SendPacket(resp)
	}

	// Just forward it
	return p.conn.SendPacket(pkt)
}
