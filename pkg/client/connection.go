package client

import (
	"bytes"
	"fmt"
	"net"

	buffer2 "github.com/mworzala/kite/internal/pkg/buffer"
	"github.com/mworzala/kite/pkg/proto/binary"
	"github.com/mworzala/kite/pkg/proto/packet"
)

type Connection struct {
	delegate net.Conn
	state    packet.State

	readBuffer  []byte
	cacheBuffer []byte // Used for caching partially read packets
}

func NewConnection(conn net.Conn) (*Connection, error) {
	c := &Connection{
		delegate: conn,
		state:    packet.Handshake,

		readBuffer: make([]byte, 1024*4),
	}

	go c.readLoop()

	return c, nil
}

func (c *Connection) State() packet.State {
	return c.state
}

func (c *Connection) SetState(state packet.State) {
	c.state = state
}

func (c *Connection) IsOpen() bool {
	return c.delegate != nil
}

func (c *Connection) Close() error {
	if c.delegate == nil {
		return nil
	}
	d := c.delegate
	c.delegate = nil
	return d.Close()
}

func (c *Connection) WritePacket(pkt packet.Packet) error {
	if !c.IsOpen() {
		return fmt.Errorf("connection closed")
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

func (c *Connection) readLoop() {
	defer c.Close()
	for {
		var start int
		if c.cacheBuffer != nil {
			start = copy(c.readBuffer, c.cacheBuffer)
			println("copied", start, "bytes from cache")
			c.cacheBuffer = nil
		}

		n, err := c.delegate.Read(c.readBuffer[start:])
		if n > 0 {
			println("read", n, "bytes")
			buffer := buffer2.NewPacketBuffer(c.readBuffer[:start+n])

			c.processPackets(buffer)
			if !c.IsOpen() {
				return
			}

			c.cacheBuffer = buffer.AllocRemainder()
			println("cached", len(c.cacheBuffer), "bytes")
		}
		if err != nil {
			panic(err)
		}
	}
}

func (c *Connection) processPackets(buffer *buffer2.PacketBuffer) {
	for {
		if buffer.Remaining() == 0 {
			return
		}

		m := buffer.Mark()

		length, err := binary.ReadVarInt(buffer)
		if err != nil {
			panic(err)
		}
		println("length:", length)

		if c.state == packet.Handshake && length == 0xFE {
			println("got legacy ping")
			c.Close()
			return
		}

		if int(length) > buffer.Remaining() {
			println("not enough data for packet, caching")
			buffer.Reset(m)
			return
		}

		packetID, err := binary.ReadVarInt(buffer)
		if err != nil {
			panic(err)
		}
		println("packetID:", packetID)

		pkt := createPacket(c.state, packetID)
		if pkt != nil {
			println("created packet", fmt.Sprintf("%T", pkt))
			if err := pkt.Read(buffer); err != nil {
				panic(err)
			}

			println("packet data:", fmt.Sprintf("%#v", pkt), "remaining =", buffer.Remaining())
			c.processSingle(pkt)
		}

	}
}

func (c *Connection) processSingle(pkt packet.Packet) {
	switch pkt := pkt.(type) {
	case *packet.ClientHandshake:
		c.state = pkt.NextState
	case *packet.ClientStatusRequest:
		if err := c.WritePacket(&packet.ServerStatusResponse{
			Status: `
{
    "version": {
        "name": "1.20.4",
        "protocol": 765
    },
    "players": {
        "max": 100,
        "online": 5,
        "sample": [
            {
                "name": "thinkofdeath",
                "id": "4566e69f-c907-48ee-8d71-d7ba5aa00d20"
            }
        ]
    },
    "description": {
        "text": "Hello world"
    },
    "favicon": "data:image/png;base64,<data>",
    "enforcesSecureChat": false
}
`,
		}); err != nil {
			panic(err)
		}
	case *packet.ClientStatusPingRequest:
		if err := c.WritePacket(&packet.ServerStatusPingResponse{
			Payload: pkt.Payload,
		}); err != nil {
			panic(err)
		}

		c.Close()

		//c.Close()
	}
}

func createPacket(state packet.State, packetId int32) packet.Packet {
	switch state {
	case packet.Handshake:
		switch packetId {
		case packet.ClientHandshakeHandshakeID:
			return new(packet.ClientHandshake)
		default:
			return nil
		}
	case packet.Status:
		switch packetId {
		case packet.ClientStatusStatusRequestID:
			return new(packet.ClientStatusRequest)
		case packet.ClientStatusPingRequestID:
			return new(packet.ClientStatusPingRequest)
		default:
			return nil
		}
	default:
		return nil
	}
}
