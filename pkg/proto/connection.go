package proto

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"fmt"
	"io"
	"net"
	"runtime/debug"
	"sync"
	"syscall"

	buffer2 "github.com/mworzala/kite/internal/pkg/buffer"
	"github.com/mworzala/kite/internal/pkg/crypto"
	"github.com/mworzala/kite/pkg/proto/binary"
	"github.com/mworzala/kite/pkg/proto/packet"
	"github.com/valyala/bytebufferpool"
)

// We do not allow any packets bigger than 5kb (the max cookie size)
// This prevents an attack vector of creating a bunch of connections and sending huge packets with
// a length 1 more than the actual data size. This forces the server to cache the packet, using up
// to the max packet size in memory (for each connection).
const preConfigMaxPacketSize = 5 * 1024

var writePool bytebufferpool.Pool

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

	reader io.Reader
	// writer protected by wlock
	writer io.Writer
	wlock  sync.Mutex

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

		reader: conn,
		writer: conn,
		wlock:  sync.Mutex{},

		id: idCounter,

		state:   packet.Handshake,
		handler: nil,

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
	//println(fmt.Sprintf("setting %s state to %s", c.direction.String(), state.String()))
	c.state = state
	c.handler = handler
}

func (c *Conn) SendPacket(pkt packet.Packet) (err error) {
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

	// Write the packet (without length prefix, it will be written depending on compression during write)
	buf := writePool.Get()
	defer writePool.Put(buf) // Always return the buffer to the pool

	if err = binary.WriteVarInt(buf, int32(pktId)); err != nil {
		return
	}
	if err = pkt.Write(buf); err != nil {
		return
	}

	if err = c.writePacketSync(buf.B); err != nil {
		return
	}

	return nil
}

func (c *Conn) writePacketSync(buffer []byte) (err error) {
	c.wlock.Lock()
	defer c.wlock.Unlock()

	if err = binary.WriteVarInt(c.writer, int32(len(buffer))); err != nil {
		return
	}

	n := len(buffer)
	for n > 0 {
		written, err := c.writer.Write(buffer)
		n -= written
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Conn) GetRemote() *Conn {
	return c.remote
}

func (c *Conn) SetRemote(remote *Conn) {
	c.remote = remote
}

func (c *Conn) EnableEncryption(sharedSecret []byte) error {
	block, err := aes.NewCipher(sharedSecret)
	if err != nil {
		return err
	}

	cfb := crypto.NewCFB8Decrypt(block, sharedSecret)
	c.reader = &cipher.StreamReader{S: cfb, R: c.reader}

	cfb = crypto.NewCFB8Encrypt(block, sharedSecret)
	c.writer = &cipher.StreamWriter{S: cfb, W: c.writer}

	return nil
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

		n, err := c.reader.Read(c.readBuffer[start:])
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

		// See comment on preConfigMaxPacketSize
		if c.state <= packet.Login && length > preConfigMaxPacketSize {
			c.Close()
			return
		}

		// If the packet contains more data than is available in the buffer, cache the remainder.
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
				_, err = c.remote.writer.Write(buffer.RemainingSlice())
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
