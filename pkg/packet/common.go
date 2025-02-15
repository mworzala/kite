package packet

import (
	"io"

	"github.com/google/uuid"
	"github.com/mworzala/kite/pkg/buffer"
	"github.com/mworzala/kite/pkg/text"
)

type ClientResourcePackStatus struct {
	UUID   uuid.UUID
	Status ResourcePackStatus
}

func (p *ClientResourcePackStatus) Direction() Direction { return Serverbound }
func (p *ClientResourcePackStatus) ID(state State) int {
	return stateId2(state, Config, Play, ClientConfigResourcePackResponseID, ClientPlayResourcePackStatusID)
}
func (p *ClientResourcePackStatus) Read(r io.Reader) (err error) {
	p.UUID, p.Status, err = buffer.Read2(r, buffer.UUID, buffer.Enum[ResourcePackStatus]{})
	return
}
func (p *ClientResourcePackStatus) Write(w io.Writer) (err error) {
	return buffer.Write2(w, buffer.UUID, p.UUID, buffer.Enum[ResourcePackStatus]{}, p.Status)
}

type ClientPluginMessage struct {
	Channel string
	Data    []byte
}

func (p *ClientPluginMessage) Direction() Direction { return Serverbound }
func (p *ClientPluginMessage) ID(state State) int {
	return stateId2(state, Config, Play, ClientConfigPluginMessageID, ClientPlayPluginMessageID)
}
func (p *ClientPluginMessage) Read(r io.Reader) (err error) {
	p.Channel, p.Data, err = buffer.Read2(r, buffer.String, buffer.RawBytes)
	return
}
func (p *ClientPluginMessage) Write(w io.Writer) (err error) {
	return buffer.Write2(w, buffer.String, p.Channel, buffer.RawBytes, p.Data)
}

type ServerResourcePackPush struct {
	Id     string
	Url    string
	Hash   string
	Forced bool
	Prompt text.Component // Optional
}

func (p *ServerResourcePackPush) Direction() Direction { return Clientbound }
func (p *ServerResourcePackPush) ID(state State) int {
	return stateId2(state, Config, Play, ServerConfigAddResourcePackID, ServerPlayResourcePackPushID)
}
func (p *ServerResourcePackPush) Read(r io.Reader) (err error) {
	p.Id, p.Url, p.Hash, p.Forced, p.Prompt, err = buffer.Read5(r,
		buffer.String, buffer.String, buffer.String,
		buffer.Bool, buffer.Opt(buffer.TextComponent))
	return
}
func (p *ServerResourcePackPush) Write(w io.Writer) (err error) {
	return buffer.Write5(w,
		buffer.String, p.Id, buffer.String, p.Url, buffer.String, p.Hash,
		buffer.Bool, p.Forced, buffer.Opt(buffer.TextComponent), p.Prompt)
}

type ServerResourcePackPop struct {
	Id uuid.UUID // Optional
}

func (p *ServerResourcePackPop) Direction() Direction { return Clientbound }
func (p *ServerResourcePackPop) ID(state State) int {
	return stateId2(state, Config, Play, ServerConfigRemoveResourcePackID, ServerPlayResourcePackPopID)
}
func (p *ServerResourcePackPop) Read(r io.Reader) (err error) {
	p.Id, err = buffer.Opt(buffer.UUID).Read(r)
	return
}
func (p *ServerResourcePackPop) Write(w io.Writer) (err error) {
	return buffer.Opt(buffer.UUID).Write(w, p.Id)
}

type ServerPluginMessage struct {
	Channel string
	Data    []byte
}

func (p *ServerPluginMessage) Direction() Direction { return Clientbound }
func (p *ServerPluginMessage) ID(state State) int {
	return stateId2(state, Config, Play, ServerConfigPluginMessageID, ServerPlayPluginMessageID)
}
func (p *ServerPluginMessage) Read(r io.Reader) (err error) {
	p.Channel, p.Data, err = buffer.Read2(r, buffer.String, buffer.RawBytes)
	return
}
func (p *ServerPluginMessage) Write(w io.Writer) (err error) {
	return buffer.Write2(w, buffer.String, p.Channel, buffer.RawBytes, p.Data)
}

type ServerDisconnect struct {
	Reason text.Component
}

func (p *ServerDisconnect) Direction() Direction { return Clientbound }
func (p *ServerDisconnect) ID(state State) int {
	return stateId2(state, Config, Play, ServerConfigDisconnectID, ServerPlayDisconnectID)
}
func (p *ServerDisconnect) Read(r io.Reader) (err error) {
	p.Reason, err = buffer.TextComponent.Read(r)
	return
}
func (p *ServerDisconnect) Write(w io.Writer) (err error) {
	return buffer.TextComponent.Write(w, p.Reason)
}

var (
	_ Packet = (*ClientResourcePackStatus)(nil)
	_ Packet = (*ClientPluginMessage)(nil)

	_ Packet = (*ServerResourcePackPush)(nil)
	_ Packet = (*ServerResourcePackPop)(nil)
	_ Packet = (*ServerPluginMessage)(nil)
	_ Packet = (*ServerDisconnect)(nil)
)
