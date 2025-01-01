package main

import (
	"bytes"
	"fmt"

	"github.com/mworzala/kite"
	"github.com/mworzala/kite/pkg/buffer"
	"github.com/mworzala/kite/pkg/packet"
)

func (p *Player) handleClientConfigPacket(pb kite.PacketBuffer) (err error) {
	if p.remote == nil {
		panic("bad state")
	}
	switch pb.Id {
	case packet.ClientConfigFinishConfigurationID:
		pkt := new(packet.ClientConfigFinishConfiguration)
		if err = pb.Read(pkt); err != nil {
			return
		}
		return p.handleClientConfigFinishConfiguration(pkt)
	default:
		return p.remote.ForwardPacket(pb)
	}
}

func (p *Player) handleClientConfigFinishConfiguration(pkt *packet.ClientConfigFinishConfiguration) (err error) {
	err = p.remote.SendPacket(pkt)
	p.conn.SetState(packet.Play)
	p.remote.SetState(packet.Play)
	return
}

func (p *Player) handleServerConfigPacket(pb kite.PacketBuffer) (err error) {
	switch pb.Id {
	case packet.ServerConfigPluginMessageID:
		pkt := new(packet.ServerPluginMessage)
		if err = pb.Read(pkt); err != nil {
			return err
		}
		return p.handleServerPluginMessage(pkt)
	case packet.ServerConfigFinishConfigurationID:
		// Don't need to do anything, just being explicit
		return p.conn.ForwardPacket(pb)
	default:
		return p.conn.ForwardPacket(pb)
	}
}

func (p *Player) handleServerPluginMessage(pkt *packet.ServerPluginMessage) (err error) {
	if pkt.Channel == "minecraft:brand" {
		oldPayload := buffer.Wrap(pkt.Data)
		oldBrand, err := buffer.String.Read(oldPayload)
		if err != nil {
			return err
		}

		var newPayload bytes.Buffer
		err = buffer.String.Write(&newPayload, fmt.Sprintf("%s // %s", oldBrand, "kite"))
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
