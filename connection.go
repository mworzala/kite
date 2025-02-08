package kite

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"net"
	"runtime/debug"
	"sync"
	"syscall"

	"github.com/mworzala/kite/internal/pkg/crypto"
	"github.com/mworzala/kite/pkg/buffer"
	"github.com/mworzala/kite/pkg/packet"
	"github.com/valyala/bytebufferpool"
)

const (
	maxPacketSize = 2_097_151
	// Before the config state we do not allow any packets bigger than 5kb (the max cookie size)
	// This prevents an attack vector of creating a bunch of connections and sending huge packets with
	// a length 1 more than the actual data size. This forces the server to cache the packet, using up
	// to the max packet size in memory (for each connection).
	maxPacketSizePreConfig = 5 * 1024

	nonceLength = 16
)

var (
	writePool     bytebufferpool.Pool
	readCachePool bytebufferpool.Pool
)

// A Conn represents a connection to a Minecraft server or client. It wraps the
// underlying net.Conn and provides utilities for processing and writing packets.
type Conn struct {
	direction packet.Direction
	delegate  net.Conn
	closed    bool

	// reader and writer should be used instead of directly accessing the delegate.
	reader io.Reader
	writer io.Writer
	wlock  sync.Mutex

	readBuffer  []byte
	cacheBuffer *bytebufferpool.ByteBuffer // Used for caching partially read packets. Pooled in readCachePool

	state   packet.State
	handler func(pb PacketBuffer) error

	nonce []byte // Login state
}

func NewConn(direction packet.Direction, conn net.Conn, handler func(pb PacketBuffer) error) *Conn {
	if handler == nil {
		panic("handler must not be nil")
	}
	c := &Conn{
		direction: direction,
		delegate:  conn,

		reader: conn,
		writer: conn,
		wlock:  sync.Mutex{},

		// TODO: don't always allocate 4mb
		readBuffer: make([]byte, 1024*1024*4),

		state:   packet.Handshake,
		handler: handler,
	}
	return c
}

func (c *Conn) RemoteAddr() net.Addr {
	return c.delegate.RemoteAddr()
}

func (c *Conn) Close() {
	if c == nil || c.closed {
		return
	}

	c.closed = true
	c.delegate.Close()
}

func (c *Conn) GetState() packet.State {
	return c.state
}

func (c *Conn) SetState(state packet.State) {
	c.state = state
}

func (c *Conn) GetNonce() []byte {
	if c.nonce == nil {
		c.nonce = make([]byte, nonceLength)
		if _, err := rand.Read(c.nonce); err != nil {
			panic(fmt.Errorf("failed to generate random nonce: %w", err))
		}
	}
	return c.nonce
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

func (c *Conn) ForwardPacket(pb PacketBuffer) (err error) {
	// Write the raw packet to the remote
	pb.internal.Reset(pb.mark)
	_, err = c.writer.Write(pb.internal.RemainingSlice())
	if errors.Is(err, syscall.EPIPE) || errors.Is(err, net.ErrClosed) {
		c.Close()
		return
	}
	return err
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

	if err = buffer.VarInt.Write(buf, int32(pktId)); err != nil {
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

func (c *Conn) writePacketSync(buf []byte) (err error) {
	c.wlock.Lock()
	defer c.wlock.Unlock()

	if err = buffer.VarInt.Write(c.writer, int32(len(buf))); err != nil {
		return
	}

	n := len(buf)
	for n > 0 {
		written, err := c.writer.Write(buf)
		n -= written
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Conn) ReadLoop() {
	defer func() {
		if r := recover(); r != nil {
			println(fmt.Sprintf("panic in readLoop: %v\n", r))
			debug.PrintStack()
			if c != nil {
				c.Close()
			}
		}
	}()
	for {
		var start int
		if c.cacheBuffer != nil {
			start = copy(c.readBuffer, c.cacheBuffer.B)
			readCachePool.Put(c.cacheBuffer)
			c.cacheBuffer = nil
		}

		n, err := c.reader.Read(c.readBuffer[start:])
		if n > 0 {
			buf := buffer.Wrap(c.readBuffer[:start+n])

			c.processPackets(buf)
			if c.closed {
				return
			}

			if buf.Remaining() > 0 {
				c.cacheBuffer = readCachePool.Get()
				buf.AllocRemainderTo(c.cacheBuffer)
			}
		}
		if errors.Is(err, io.EOF) || errors.Is(err, net.ErrClosed) {
			c.Close()
			return
		} else if err != nil {
			//todo handle appropriately and close connection.
			panic(err)
		}
	}
}

func (c *Conn) processPackets(buf *buffer.Buffer) {
	for {
		if buf.Remaining() == 0 {
			return
		}

		packetStart := buf.Mark()
		length, err := buffer.VarInt.Read(buf)
		if err != nil {
			panic(err)
		}

		// We just do not support legacy clients for now, just close the connection
		if c.state == packet.Handshake && length == 0xFE {
			c.Close()
			return
		}

		// See comment on maxPacketSizePreConfig
		if length > maxPacketSize {
			//todo this should disconnect with a message most likely.
			c.Close()
			return
		} else if c.state <= packet.Login && length > maxPacketSizePreConfig {
			c.Close()
			return
		}

		// If the packet contains more data than is available in the buffer, cache the remainder.
		if int(length) > buf.Remaining() {
			buf.Reset(packetStart)
			buf.Limit(-1)
			return
		}
		buf.Limit(int(length)) // Cap the read buffer to the packet length

		packetID, err := buffer.VarInt.Read(buf)
		if errors.Is(err, net.ErrClosed) {
			c.Close()
			return
		} else if err != nil {
			panic(err)
		}

		err = c.handler(PacketBuffer{Id: int(packetID), internal: buf, mark: packetStart})
		if err != nil {
			println(fmt.Errorf("packet processing failed: %w (%s/%s/%d)", err, c.direction.String(), c.state.String(), packetID).Error())
			c.Close()
			return
		}
		if buf.Remaining() > 0 {
			panic(fmt.Errorf("%s: %s/%s/%d", "packet not fully read", c.direction.String(), c.state.String(), packetID))
		}

		buf.Limit(-1)
	}
}
