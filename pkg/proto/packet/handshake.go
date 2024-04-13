package packet

import (
	"io"

	"github.com/mworzala/kite/pkg/proto/binary"
)

const ClientHandshakeHandshakeID = 0x00

type ClientHandshake struct {
	ProtocolVersion int32
	ServerAddress   string
	ServerPort      uint16
	NextState       State
}

const serverAddressLength = 255

func (p *ClientHandshake) Direction() Direction { return Serverbound }
func (c *ClientHandshake) ID(state State) int {
	return stateId1(state, Handshake, ClientHandshakeHandshakeID)
}

func (c *ClientHandshake) Read(r io.Reader) (err error) {
	if c.ProtocolVersion, err = binary.ReadVarInt(r); err != nil {
		return
	}
	if c.ServerAddress, err = binary.ReadSizedString(r, serverAddressLength); err != nil {
		return
	}
	if c.ServerPort, err = binary.ReadUShort(r); err != nil {
		return
	}
	if c.NextState, err = binary.ReadEnum[State](r); err != nil {
		return
	}
	return nil
}

func (c *ClientHandshake) Write(w io.Writer) (err error) {
	if err = binary.WriteVarInt(w, c.ProtocolVersion); err != nil {
		return
	}
	if err = binary.WriteSizedString(w, c.ServerAddress, serverAddressLength); err != nil {
		return
	}
	if err = binary.WriteUShort(w, c.ServerPort); err != nil {
		return
	}
	if err = binary.WriteEnum(w, c.NextState); err != nil {
		return
	}
	return nil
}

var _ Packet = (*ClientHandshake)(nil)
