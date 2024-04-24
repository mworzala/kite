package packet

import (
	"io"

	"github.com/mworzala/kite/pkg/proto/binary"
)

const (
	ClientLoginLoginStartID = iota
	ClientLoginLoginEncryptionResponseID
	ClientLoginPluginResponseID
	ClientLoginLoginAcknowledgedID
)

type ClientLoginStart struct {
	Name string
	UUID string
}

const nameLength = 16

func (p *ClientLoginStart) Direction() Direction { return Serverbound }
func (p *ClientLoginStart) ID(state State) int {
	return stateId1(state, Login, ClientLoginLoginStartID)
}

func (p *ClientLoginStart) Read(r io.Reader) (err error) {
	if p.Name, err = binary.ReadSizedString(r, nameLength); err != nil {
		return
	}
	if p.UUID, err = binary.ReadUUID(r); err != nil {
		return
	}
	return nil
}

func (p *ClientLoginStart) Write(w io.Writer) (err error) {
	if err = binary.WriteSizedString(w, p.Name, nameLength); err != nil {
		return
	}
	if err = binary.WriteUUID(w, p.UUID); err != nil {
		return
	}
	return nil
}

type ClientLoginAcknowledged struct{}

func (p *ClientLoginAcknowledged) Direction() Direction { return Serverbound }
func (p *ClientLoginAcknowledged) ID(state State) int {
	return stateId1(state, Login, ClientLoginLoginAcknowledgedID)
}

func (p *ClientLoginAcknowledged) Read(r io.Reader) (err error) {
	return nil
}
func (p *ClientLoginAcknowledged) Write(w io.Writer) (err error) {
	return nil
}

const (
	ServerLoginDisconnectID = iota
	ServerLoginEncryptionRequestID
	ServerLoginLoginSuccessID
	ServerLoginSetCompressionID
	ServerLoginPluginRequestID
)

type ServerLoginDisconnect struct {
	Reason string
}

func (p *ServerLoginDisconnect) Direction() Direction { return Clientbound }
func (p *ServerLoginDisconnect) ID(state State) int {
	return stateId1(state, Login, ServerLoginDisconnectID)
}

func (p *ServerLoginDisconnect) Read(r io.Reader) (err error) {
	if p.Reason, err = binary.ReadChatString(r); err != nil {
		return
	}
	return nil
}

func (p *ServerLoginDisconnect) Write(w io.Writer) (err error) {
	if err = binary.WriteChatString(w, p.Reason); err != nil {
		return
	}
	return nil
}

type ServerLoginSuccess struct {
	UUID     string
	Username string
	//todo properties
	Abc bool
}

func (p *ServerLoginSuccess) Direction() Direction { return Clientbound }
func (p *ServerLoginSuccess) ID(state State) int {
	return stateId1(state, Login, ServerLoginLoginSuccessID)
}

func (p *ServerLoginSuccess) Read(r io.Reader) (err error) {
	if p.UUID, err = binary.ReadUUID(r); err != nil {
		return
	}
	if p.Username, err = binary.ReadSizedString(r, 16); err != nil {
		return
	}
	_, _ = binary.ReadVarInt(r) //todo properties
	if p.Abc, err = binary.ReadBool(r); err != nil {
		return
	}
	return nil
}

func (p *ServerLoginSuccess) Write(w io.Writer) (err error) {
	if err = binary.WriteUUID(w, p.UUID); err != nil {
		return
	}
	if err = binary.WriteSizedString(w, p.Username, 16); err != nil {
		return
	}
	_ = binary.WriteVarInt(w, 0) //todo properties
	if err = binary.WriteBool(w, p.Abc); err != nil {
		return
	}
	return nil
}

type ServerLoginPluginRequest struct {
	MessageID int32
	Channel   string
	Data      []byte
}

func (p *ServerLoginPluginRequest) Direction() Direction { return Clientbound }
func (p *ServerLoginPluginRequest) ID(state State) int {
	return stateId1(state, Login, ServerLoginPluginRequestID)
}

func (p *ServerLoginPluginRequest) Read(r io.Reader) (err error) {
	if p.MessageID, err = binary.ReadVarInt(r); err != nil {
		return
	}
	if p.Channel, err = binary.ReadSizedString(r, 20); err != nil {
		return
	}
	if p.Data, err = binary.ReadRaw(r, -1); err != nil {
		return
	}
	return nil
}

func (p *ServerLoginPluginRequest) Write(w io.Writer) (err error) {
	if err = binary.WriteVarInt(w, p.MessageID); err != nil {
		return
	}
	if err = binary.WriteSizedString(w, p.Channel, 20); err != nil {
		return
	}
	if err = binary.WriteRaw(w, p.Data); err != nil {
		return
	}
	return nil
}

var (
	_ Packet = (*ClientLoginStart)(nil)
	_ Packet = (*ClientLoginAcknowledged)(nil)

	_ Packet = (*ServerLoginDisconnect)(nil)
)
