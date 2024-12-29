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
	Intent          Intent
}

const serverAddressLength = 255

func (p *ClientHandshake) Direction() Direction { return Serverbound }
func (p *ClientHandshake) ID(state State) int {
	return stateId1(state, Handshake, ClientHandshakeHandshakeID)
}
func (p *ClientHandshake) Read(r io.Reader) (err error) {
	if p.ProtocolVersion, err = binary.ReadVarInt(r); err != nil {
		return
	}
	if p.ServerAddress, err = binary.ReadSizedString(r, serverAddressLength); err != nil {
		return
	}
	if p.ServerPort, err = binary.ReadUShort(r); err != nil {
		return
	}
	if p.Intent, err = binary.ReadEnum[Intent](r); err != nil {
		return
	}
	return nil
}
func (p *ClientHandshake) Write(w io.Writer) (err error) {
	if err = binary.WriteVarInt(w, p.ProtocolVersion); err != nil {
		return
	}
	if err = binary.WriteSizedString(w, p.ServerAddress, serverAddressLength); err != nil {
		return
	}
	if err = binary.WriteUShort(w, p.ServerPort); err != nil {
		return
	}
	if err = binary.WriteEnum(w, p.Intent); err != nil {
		return
	}
	return nil
}

var _ Packet = (*ClientHandshake)(nil)

// An Intent is the target state when joining a server.
type Intent int

const (
	IntentStatus = iota + 1
	IntentLogin
	IntentTransfer
)

func (s Intent) Validate() bool {
	return s >= IntentStatus && s <= IntentTransfer
}

func (s Intent) String() string {
	switch s {
	case IntentStatus:
		return "status"
	case IntentLogin:
		return "login"
	case IntentTransfer:
		return "transfer"
	}
	return "unknown"
}
