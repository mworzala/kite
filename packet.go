package kite

import (
	"errors"

	"github.com/mworzala/kite/pkg/buffer"
	"github.com/mworzala/kite/pkg/packet"
)

// PacketBuffer represents a buffer containing a single packet.
// Does not hold the connection state, the state can be inferred from the connection which created the buffer.
//
// Packet buffers must be fully consumed. If not reading the content, the Consume method should be called.
// Forwarding a packet buffer to another connection does count as consumption.
type PacketBuffer struct {
	Id       int            // The protocol ID of the packet
	internal *buffer.Buffer // Delegate buffer containing the packet data with configured mark and limit
	mark     int            // Start location of packet in buffer
	read     bool           // Whether the packet has been read
}

// Consume marks the packet as read without using the data
func (p PacketBuffer) Consume() {
	// We reset to our known start (in case someone read partially), then reset to the end of the packet
	p.internal.Reset(p.mark)
	p.internal.Reset(p.mark + p.internal.Remaining())
	p.read = true
}

func (p PacketBuffer) Read(t packet.Packet) error {
	if p.read {
		return errors.New("packet already read")
	}

	p.read = true
	return t.Read(p.internal)
}
