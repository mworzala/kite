package proto

import (
	"errors"

	buffer2 "github.com/mworzala/kite/internal/pkg/buffer"
	"github.com/mworzala/kite/pkg/packet"
	"github.com/mworzala/kite/pkg/proto/binary"
)

var (
	UnknownPacket = errors.New("unknown packet")
)

type Packet struct {
	State        packet.State
	Id           int32
	Buf          *buffer2.PacketBuffer
	Start        int
	ReadInternal bool
}

// Consume marks the entire content of the buffer as read.
//
// This is pretty much to require explicit consumption of an unused packet.
func (p Packet) Consume() {
	_, _ = binary.ReadRaw(p.Buf, binary.Remaining)
	p.ReadInternal = true
}

func (p Packet) Read(t packet.Packet) error {
	if p.ReadInternal {
		return errors.New("packet already read")
	}
	p.ReadInternal = true
	return t.Read(p.Buf)
}
