package packet

import (
	"io"

	"github.com/mworzala/kite/pkg/proto/binary"
)

const (
	ClientLoginLoginStartID = iota
	ClientLoginEncryptionResponseID
	ClientLoginPluginResponseID
	ClientLoginLoginAcknowledgedID
	ClientLoginCookieResponseID
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

type ClientEncryptionResponse struct {
	SharedSecret []byte
	VerifyToken  []byte
}

func (p *ClientEncryptionResponse) Direction() Direction { return Serverbound }
func (p *ClientEncryptionResponse) ID(state State) int {
	return stateId1(state, Login, ClientLoginEncryptionResponseID)
}

func (p *ClientEncryptionResponse) Read(r io.Reader) (err error) {
	if p.SharedSecret, err = binary.ReadByteArray(r); err != nil {
		return
	}
	if p.VerifyToken, err = binary.ReadByteArray(r); err != nil {
		return
	}
	return nil
}

func (p *ClientEncryptionResponse) Write(w io.Writer) (err error) {
	if err = binary.WriteByteArray(w, p.SharedSecret); err != nil {
		return
	}
	if err = binary.WriteByteArray(w, p.VerifyToken); err != nil {
		return
	}
	return nil
}

type ClientLoginPluginResponse struct {
	MessageID int32
	Data      []byte // nil when unsuccessful
}

func (p *ClientLoginPluginResponse) Direction() Direction { return Serverbound }
func (p *ClientLoginPluginResponse) ID(state State) int {
	return stateId1(state, Login, ClientLoginPluginResponseID)
}

func (p *ClientLoginPluginResponse) Read(r io.Reader) (err error) {
	if p.MessageID, err = binary.ReadVarInt(r); err != nil {
		return
	}
	var successful bool
	if successful, err = binary.ReadBool(r); err != nil || !successful {
		return
	}
	if p.Data, err = binary.ReadRaw(r, -1); err != nil {
		return
	}
	return nil
}

func (p *ClientLoginPluginResponse) Write(w io.Writer) (err error) {
	if err = binary.WriteVarInt(w, p.MessageID); err != nil {
		return
	}
	if err = binary.WriteBool(w, p.Data != nil); err != nil || p.Data == nil {
		return
	}
	if err = binary.WriteRaw(w, p.Data); err != nil {
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
	ServerLoginCookieRequestID
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

type ServerEncryptionRequest struct {
	ServerID           string
	PublicKey          []byte
	VerifyToken        []byte
	ShouldAuthenticate bool
}

func (p *ServerEncryptionRequest) Direction() Direction { return Clientbound }
func (p *ServerEncryptionRequest) ID(state State) int {
	return stateId1(state, Login, ServerLoginEncryptionRequestID)
}

func (p *ServerEncryptionRequest) Read(r io.Reader) (err error) {
	if p.ServerID, err = binary.ReadSizedString(r, 20); err != nil {
		return
	}
	if p.PublicKey, err = binary.ReadByteArray(r); err != nil {
		return
	}
	if p.VerifyToken, err = binary.ReadByteArray(r); err != nil {
		return
	}
	if p.ShouldAuthenticate, err = binary.ReadBool(r); err != nil {
		return
	}
	return nil
}

func (p *ServerEncryptionRequest) Write(w io.Writer) (err error) {
	if err = binary.WriteSizedString(w, p.ServerID, 20); err != nil {
		return
	}
	if err = binary.WriteByteArray(w, p.PublicKey); err != nil {
		return
	}
	if err = binary.WriteByteArray(w, p.VerifyToken); err != nil {
		return
	}
	if err = binary.WriteBool(w, p.ShouldAuthenticate); err != nil {
		return
	}
	return nil
}

type ServerLoginSuccess struct {
	GameProfile
	StrictErrorHandling bool
}

func (p *ServerLoginSuccess) Direction() Direction { return Clientbound }
func (p *ServerLoginSuccess) ID(state State) int {
	return stateId1(state, Login, ServerLoginLoginSuccessID)
}

func (p *ServerLoginSuccess) Read(r io.Reader) (err error) {
	if err = p.GameProfile.Read(r); err != nil {
		return
	}
	if p.StrictErrorHandling, err = binary.ReadBool(r); err != nil {
		return
	}
	return nil
}

func (p *ServerLoginSuccess) Write(w io.Writer) (err error) {
	if err = p.GameProfile.Write(w); err != nil {
		return
	}
	if err = binary.WriteBool(w, p.StrictErrorHandling); err != nil {
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
	_ Packet = (*ClientEncryptionResponse)(nil)
	_ Packet = (*ClientLoginAcknowledged)(nil)

	_ Packet = (*ServerLoginDisconnect)(nil)
	_ Packet = (*ServerEncryptionRequest)(nil)
	_ Packet = (*ServerLoginSuccess)(nil)
	_ Packet = (*ServerLoginPluginRequest)(nil)
)
