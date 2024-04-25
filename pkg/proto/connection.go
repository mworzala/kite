package proto

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"runtime/debug"
	"syscall"

	buffer2 "github.com/mworzala/kite/internal/pkg/buffer"
	"github.com/mworzala/kite/pkg/proto/binary"
	"github.com/mworzala/kite/pkg/proto/packet"
)

var idCounter int

// A Conn is a connection that can send and process packets.
// Typically, a Conn exists for both the client<>proxy and proxy<>server connections.
// When connecting to a new server, a new Conn will be created before the old one is closed.
//
// A Handler is used to process packets received by the Conn.
//
// A Conn should be informed of the remote connection to allow packet forwarding.
type Conn struct {
	direction packet.Direction
	delegate  net.Conn
	closed    bool

	id int

	//todo both of these buffers should be pooled
	readBuffer  []byte
	cacheBuffer []byte // Used for caching partially read packets

	state   packet.State
	handler Handler

	remote *Conn
}

func NewConn(direction packet.Direction, conn net.Conn) (*Conn, func()) {
	c := &Conn{
		direction: direction,
		delegate:  conn,

		id: idCounter,

		readBuffer: make([]byte, 1024*1024*4),
	}
	idCounter++
	//todo weird to return readLoop, not a big fan.
	return c, c.readLoop
}

func (c *Conn) Close() {
	if c == nil || c.closed {
		return
	}

	c.closed = true
	c.delegate.Close()
	if c.remote != nil {
		c.remote.Close()
	}
}

func (c *Conn) SetState(state packet.State, handler Handler) {
	println(fmt.Sprintf("setting %s state to %s", c.direction.String(), state.String()))
	c.state = state
	c.handler = handler
}

func (c *Conn) SendPacket(pkt packet.Packet) error {
	if c.closed {
		return io.EOF
	}

	if pkt.Direction() == c.direction {
		return fmt.Errorf("packet %T has wrong direction", pkt)
	}
	pktId := pkt.ID(c.state)
	if pktId < 0 {
		return fmt.Errorf("packet %T is not applicable to state %s", pkt, c.state.String())
	}

	// Frame the packet
	temp := new(bytes.Buffer)
	if err := binary.WriteVarInt(temp, int32(pktId)); err != nil {
		return err
	}
	if err := pkt.Write(temp); err != nil {
		return err
	}

	buffer := bytes.NewBuffer(nil)
	if err := binary.WriteVarInt(buffer, int32(temp.Len())); err != nil {
		return err
	}
	if _, err := buffer.Write(temp.Bytes()); err != nil {
		return err
	}

	// Write the packet
	if _, err := c.delegate.Write(buffer.Bytes()); err != nil {
		return err
	}

	return nil
}

func (c *Conn) GetRemote() *Conn {
	return c.remote
}

func (c *Conn) SetRemote(remote *Conn) {
	c.remote = remote
}

func (c *Conn) readLoop() {
	defer func() {
		if r := recover(); r != nil {
			println(fmt.Sprintf("panic in readLoop: %v\n%s", r, string(debug.Stack())))
			c.Close()
		}
	}()
	for {
		var start int
		if c.cacheBuffer != nil {
			start = copy(c.readBuffer, c.cacheBuffer)
			c.cacheBuffer = nil
		}

		n, err := c.delegate.Read(c.readBuffer[start:])
		if n > 0 {
			buffer := buffer2.NewPacketBuffer(c.readBuffer[:start+n])

			c.processPackets(buffer)
			if c.closed {
				return
			}

			c.cacheBuffer = buffer.AllocRemainder()
		}
		if errors.Is(err, io.EOF) || errors.Is(err, net.ErrClosed) {
			return
		} else if err != nil {
			panic(err)
		}
	}
}

func (c *Conn) processPackets(buffer *buffer2.PacketBuffer) {
	for {
		if buffer.Remaining() == 0 {
			return
		}

		packetStart := buffer.Mark()
		length, err := binary.ReadVarInt(buffer)
		if err != nil {
			panic(err)
		}

		// We just do not support legacy clients for now, just close the connection
		if c.state == packet.Handshake && length == 0xFE {
			c.Close()
			return
		}

		if int(length) > buffer.Remaining() {
			buffer.Reset(packetStart)
			buffer.Limit(-1)
			return
		}
		buffer.Limit(int(length)) // Cap the read buffer to the packet length

		packetID, err := binary.ReadVarInt(buffer)
		if errors.Is(err, net.ErrClosed) {
			c.Close()
			return
		} else if err != nil {
			panic(err)
		}

		if c.handler != nil {
			err = c.handler.HandlePacket(Packet{Id: packetID, buf: buffer})
			if errors.Is(err, Forward) {
				if c.remote == nil {
					panic("no remote")
				}

				// Write the raw packet to the remote
				buffer.Reset(packetStart)
				_, err = c.remote.delegate.Write(buffer.RemainingSlice())
				if errors.Is(err, syscall.EPIPE) || errors.Is(err, net.ErrClosed) {
					c.Close()
					return
				} else if err != nil {
					panic(err)
				}
			} else if err != nil {
				println(fmt.Errorf("packet processing failed on %d: %w (%s/%s/%d)", c.id, err, c.direction.String(), c.state.String(), packetID).Error())
				c.Close()
				return
			}
		}

		if buffer.Remaining() > 0 {
			panic(fmt.Errorf("%s: %s/%s/%d", "packet not fully read", c.direction.String(), c.state.String(), packetID))
		}

		buffer.Limit(-1)
	}
}
