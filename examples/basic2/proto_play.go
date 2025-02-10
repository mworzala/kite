package main

import (
    "github.com/mworzala/kite"
    "github.com/mworzala/kite/pkg/packet"
    "github.com/mworzala/kite/pkg/text"
)

func (p *Player) handleClientPlayPacket(pb kite.PacketBuffer) (err error) {
    if p.remote == nil {
        panic("bad state")
    }
    switch pb.Id {
    case packet.ClientPlayChatID:
        pkt := new(packet.ClientPlayChat)
        if err = pb.Read(pkt); err != nil {
            return
        }
        pb.Consume() // Chat packet is only partially implemented
        if err = p.handleClientPlayChat(pkt); err != nil {
            return err
        }
        return p.remote.ForwardPacket(pb)
    default:
        return p.remote.ForwardPacket(pb)
    }
}

func (p *Player) handleClientPlayChat(pkt *packet.ClientPlayChat) (err error) {
    println("new message", pkt.Message)
    if pkt.Message == "kick" {
        p.Disconnect2(text.Text{Text: "Hello"})
    }
    return
}

func (p *Player) handleServerPlayPacket(pb kite.PacketBuffer) (err error) {
    switch pb.Id {
    default:
        return p.conn.ForwardPacket(pb)
    }
}
