package packet

import (
	"io"

	"github.com/mworzala/kite/pkg/buffer"
)

const ClientHandshakeHandshakeID = 0x00

type ClientHandshake struct {
	ProtocolVersion int32
	ServerAddress   string
	ServerPort      uint16
	Intent          Intent
}

func (p *ClientHandshake) Direction() Direction { return Serverbound }
func (p *ClientHandshake) ID(state State) int {
	return stateId1(state, Handshake, ClientHandshakeHandshakeID)
}
func (p *ClientHandshake) Read(r io.Reader) (err error) {
	p.ProtocolVersion, p.ServerAddress, p.ServerPort, p.Intent, err = buffer.Read4(r,
		buffer.VarInt, buffer.String, buffer.Uint16, buffer.Enum[Intent]{})
	return
}
func (p *ClientHandshake) Write(w io.Writer) (err error) {
	return buffer.Write4(w,
		buffer.VarInt, p.ProtocolVersion, buffer.String, p.ServerAddress,
		buffer.Uint16, p.ServerPort, buffer.Enum[Intent]{}, p.Intent)
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
