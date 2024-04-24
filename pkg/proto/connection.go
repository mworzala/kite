package proto

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"

	buffer2 "github.com/mworzala/kite/internal/pkg/buffer"
	"github.com/mworzala/kite/pkg/proto/binary"
	"github.com/mworzala/kite/pkg/proto/packet"
)

type Conn struct {
	direction packet.Direction
	delegate  net.Conn
	closed    bool

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

		readBuffer: make([]byte, 1024*1024*4),
	}
	//todo weird to return readLoop, not a big fan.
	return c, c.readLoop
}

func (c *Conn) Close() {
	c.closed = true
	c.delegate.Close()
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

	pktId := pkt.ID(c.state)
	if pktId < 0 {
		return fmt.Errorf("packet id < 0")
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
		if c.delegate != nil {
			c.delegate.Close()
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

		m := buffer.Mark()

		length, err := binary.ReadVarInt(buffer)
		if err != nil {
			panic(err)
		}

		m2 := buffer.Mark()

		if c.state == packet.Handshake && length == 0xFE {
			println("got legacy ping")
			c.Close()
			return
		}

		if int(length) > buffer.Remaining() {
			println(fmt.Sprintf("direction=%s, state=%s not enough data for packet, caching: len=%d, rem=%d", c.direction, c.state, length, buffer.Remaining()))
			buffer.Reset(m)
			buffer.Limit(-1)
			return
		}

		buffer.Limit(int(length))

		packetID, err := binary.ReadVarInt(buffer)
		if err != nil {
			panic(err)
		}
		//println(fmt.Sprintf("INCOMING direction=%s, state=%s, id=%d", c.direction.String(), c.state.String(), packetID))

		if c.handler != nil {
			err = c.handler.HandlePacket(Packet{Id: packetID, buf: buffer})
			if errors.Is(err, UnknownPacket) {
				err = fmt.Errorf("%w: %s/%s/%d", err, c.direction.String(), c.state.String(), packetID)
			}
			if errors.Is(err, Forward) {
				if c.remote == nil {
					panic("no remote")
				}

				//panic("a")
				buffer.Reset(m2)

				// Write the raw packet to the remote
				//outBuffer := bytes.NewBuffer(nil)
				////todo +1 isnt valid because it actually needs to be the length of the varint
				if err = binary.WriteVarInt(c.remote.delegate, length); err != nil {
					panic(err)
				}
				//if err = binary.WriteVarInt(outBuffer, int32(packetID)); err != nil {
				//	panic(err)
				//}
				//if _, err = c.remote.delegate.Write(outBuffer.Bytes()); err != nil {
				//	panic(err)
				//}

				if _, err = c.remote.delegate.Write(buffer.AllocRemainder()); err != nil {
					//if _, err = io.CopyN(c.remote.delegate, buffer, int64(length)); err != nil {
					panic(err)
				}
			} else if err != nil {
				c.Close()
				println(fmt.Errorf("error handling packet: %w", err).Error())
				return
			}
		}

		if buffer.Remaining() > 0 {
			panic(fmt.Errorf("%s: %s/%s/%d", "remaining data in buffer", c.direction.String(), c.state.String(), packetID))
		}

		buffer.Limit(-1)
	}
}
