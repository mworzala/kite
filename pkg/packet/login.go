package packet

import (
	"github.com/mworzala/kite/pkg/mojang"
	"io"

	"github.com/google/uuid"
	"github.com/mworzala/kite/pkg/buffer"
	"github.com/mworzala/kite/pkg/text"
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
	UUID uuid.UUID
}

func (p *ClientLoginStart) Direction() Direction { return Serverbound }
func (p *ClientLoginStart) ID(state State) int {
	return stateId1(state, Login, ClientLoginLoginStartID)
}
func (p *ClientLoginStart) Read(r io.Reader) (err error) {
	p.Name, p.UUID, err = buffer.Read2(r, buffer.String, buffer.UUID)
	return
}
func (p *ClientLoginStart) Write(w io.Writer) (err error) {
	return buffer.Write2(w, buffer.String, p.Name, buffer.UUID, p.UUID)
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
	p.SharedSecret, p.VerifyToken, err = buffer.Read2(r, buffer.ByteArray, buffer.ByteArray)
	return
}
func (p *ClientEncryptionResponse) Write(w io.Writer) (err error) {
	return buffer.Write2(w, buffer.ByteArray, p.SharedSecret, buffer.ByteArray, p.VerifyToken)
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
	if p.MessageID, err = buffer.VarInt.Read(r); err != nil {
		return
	}
	var successful bool
	if successful, err = buffer.Bool.Read(r); err != nil || !successful {
		return
	}
	if p.Data, err = buffer.RawBytes.Read(r); err != nil {
		return
	}
	return nil
}
func (p *ClientLoginPluginResponse) Write(w io.Writer) (err error) {
	if err = buffer.VarInt.Write(w, p.MessageID); err != nil {
		return
	}
	if err = buffer.Bool.Write(w, p.Data != nil); err != nil || p.Data == nil {
		return
	}
	if err = buffer.RawBytes.Write(w, p.Data); err != nil {
		return
	}
	return nil
}

type ClientLoginAcknowledged struct{}

func (p *ClientLoginAcknowledged) Direction() Direction { return Serverbound }
func (p *ClientLoginAcknowledged) ID(state State) int {
	return stateId1(state, Login, ClientLoginLoginAcknowledgedID)
}

func (p *ClientLoginAcknowledged) Read(_ io.Reader) (err error) {
	return nil
}
func (p *ClientLoginAcknowledged) Write(_ io.Writer) (err error) {
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
	Reason text.Component
}

func (p *ServerLoginDisconnect) Direction() Direction { return Clientbound }
func (p *ServerLoginDisconnect) ID(state State) int {
	return stateId1(state, Login, ServerLoginDisconnectID)
}
func (p *ServerLoginDisconnect) Read(r io.Reader) (err error) {
	p.Reason, err = buffer.TextComponentJSON.Read(r)
	return
}
func (p *ServerLoginDisconnect) Write(w io.Writer) (err error) {
	return buffer.TextComponentJSON.Write(w, p.Reason)
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
	p.ServerID, p.PublicKey, p.VerifyToken, p.ShouldAuthenticate, err = buffer.Read4(r,
		buffer.String, buffer.ByteArray, buffer.ByteArray, buffer.Bool)
	return
}
func (p *ServerEncryptionRequest) Write(w io.Writer) (err error) {
	return buffer.Write4(w, buffer.String, p.ServerID, buffer.ByteArray, p.PublicKey,
		buffer.ByteArray, p.VerifyToken, buffer.Bool, p.ShouldAuthenticate)
}

type ServerLoginSuccess struct {
	mojang.GameProfile
}

func (p *ServerLoginSuccess) Direction() Direction { return Clientbound }
func (p *ServerLoginSuccess) ID(state State) int {
	return stateId1(state, Login, ServerLoginLoginSuccessID)
}
func (p *ServerLoginSuccess) Read(r io.Reader) (err error) {
	return p.GameProfile.Read(r)
}
func (p *ServerLoginSuccess) Write(w io.Writer) (err error) {
	return p.GameProfile.Write(w)
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
	p.MessageID, p.Channel, p.Data, err = buffer.Read3(r, buffer.VarInt, buffer.String, buffer.RawBytes)
	return
}
func (p *ServerLoginPluginRequest) Write(w io.Writer) (err error) {
	return buffer.Write3(w, buffer.VarInt, p.MessageID, buffer.String, p.Channel, buffer.RawBytes, p.Data)
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
