package proto

import (
	"errors"
	"io"

	"github.com/mworzala/kite/pkg/proto/binary"
	"github.com/mworzala/kite/pkg/proto/packet"
)

var (
	Forward       = errors.New("forward")
	UnknownPacket = errors.New("unknown packet")
)

// Handler is a basic interface for something which can handle a protocol state.
// There is always a single handler set for each player at a given time.
//
// A handler is specific to one side of the protocol traffic, for example there will
// be separate a "play" handler for both client packets and server packets.
type Handler interface {
	// HandlePacket is called for each packet received from the protocol side.
	//
	// The handler should return nil if the packet was handled successfully. It may also return
	// Forward to indicate that the packet should be forwarded as is to the other side of the connection.
	HandlePacket(pp Packet) error
}

type ConnHandlerFunc func(*Conn) Handler

type Packet struct {
	Id   int32
	buf  io.Reader
	read bool
}

// Consume marks the entire content of the buffer as read.
//
// This is pretty much to require explicit consumption of an unused packet.
func (p Packet) Consume() {
	_, _ = binary.ReadRaw(p.buf, binary.Remaining)
	p.read = true
}

func (p Packet) Read(t packet.Packet) error {
	if p.read {
		return errors.New("packet already read")
	}
	p.read = true
	return t.Read(p.buf)
}
