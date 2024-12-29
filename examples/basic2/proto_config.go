package main

import (
	"bytes"
	"fmt"

	buffer2 "github.com/mworzala/kite/internal/pkg/buffer"
	packet2 "github.com/mworzala/kite/pkg/packet"
	"github.com/mworzala/kite/pkg/proto"
	"github.com/mworzala/kite/pkg/proto/binary"
)

func (p *Player) handleClientConfigPacket(pp proto.Packet) (err error) {
	if p.remote == nil {
		panic("bad state")
	}
	switch pp.Id {
	case packet2.ClientConfigFinishConfigurationID:
		pkt := new(packet2.ClientConfigFinishConfiguration)
		if err = pp.Read(pkt); err != nil {
			return
		}
		return p.handleClientConfigFinishConfiguration(pkt)
	default:
		return p.remote.ForwardPacket(pp)
	}
}

func (p *Player) handleClientConfigFinishConfiguration(pkt *packet2.ClientConfigFinishConfiguration) (err error) {
	err = p.remote.SendPacket(pkt)
	p.conn.SetState(packet2.Play)
	p.remote.SetState(packet2.Play)
	return
}

func (p *Player) handleServerConfigPacket(pp proto.Packet) (err error) {
	switch pp.Id {
	case packet2.ServerConfigPluginMessageID:
		pkt := new(packet2.ServerPluginMessage)
		if err = pp.Read(pkt); err != nil {
			return err
		}
		return p.handleServerPluginMessage(pkt)
	case packet2.ServerConfigFinishConfigurationID:
		// Don't need to do anything, just being explicit
		return p.conn.ForwardPacket(pp)
	default:
		return p.conn.ForwardPacket(pp)
	}
}

func (p *Player) handleServerPluginMessage(pkt *packet2.ServerPluginMessage) (err error) {
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

		resp := &packet2.ServerPluginMessage{
			Channel: "minecraft:brand",
			Data:    newPayload.Bytes(),
		}
		return p.conn.SendPacket(resp)
	}

	// Just forward it
	return p.conn.SendPacket(pkt)
}
