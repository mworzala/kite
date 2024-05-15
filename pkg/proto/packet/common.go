package packet

import (
	"io"

	"github.com/mworzala/kite/pkg/proto/binary"
)

type ResourcePackStatus int

const (
	ResourcePackSuccessfullyLoaded ResourcePackStatus = iota
	ResourcePackDeclined
	ResourcePackFailedDownload
	ResourcePackAccepted
	ResourcePackDownloaded
	ResourcePackInvalidURL
	ResourcePackFailedReload
	ResourcePackDiscarded
)

func (s ResourcePackStatus) Validate() bool {
	return s >= ResourcePackSuccessfullyLoaded && s <= ResourcePackDiscarded
}

func (s ResourcePackStatus) String() string {
	switch s {
	case ResourcePackSuccessfullyLoaded:
		return "successfully_loaded"
	case ResourcePackDeclined:
		return "declined"
	case ResourcePackFailedDownload:
		return "failed_download"
	case ResourcePackAccepted:
		return "accepted"
	case ResourcePackDownloaded:
		return "downloaded"
	case ResourcePackInvalidURL:
		return "invalid_url"
	case ResourcePackFailedReload:
		return "failed_reload"
	case ResourcePackDiscarded:
		return "discarded"
	}
	return "unknown"
}

type ClientResourcePackStatus struct {
	UUID   string
	Status ResourcePackStatus
}

func (p *ClientResourcePackStatus) Direction() Direction { return Serverbound }
func (p *ClientResourcePackStatus) ID(state State) int {
	return stateId2(state, Config, Play, ClientConfigResourcePackResponseID, ClientPlayResourcePackStatusID)
}

func (p *ClientResourcePackStatus) Read(r io.Reader) (err error) {
	if p.UUID, err = binary.ReadString(r); err != nil {
		return
	}
	if p.Status, err = binary.ReadEnum[ResourcePackStatus](r); err != nil {
		return
	}
	return nil
}

func (p *ClientResourcePackStatus) Write(w io.Writer) (err error) {
	if err = binary.WriteString(w, p.UUID); err != nil {
		return
	}
	if err = binary.WriteEnum[ResourcePackStatus](w, p.Status); err != nil {
		return
	}
	return nil
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
	if p.Channel, err = binary.ReadString(r); err != nil {
		return
	}
	if p.Data, err = binary.ReadRaw(r, binary.Remaining); err != nil {
		return
	}
	return nil
}

func (p *ClientPluginMessage) Write(w io.Writer) (err error) {
	if err = binary.WriteString(w, p.Channel); err != nil {
		return //todo what actually is max length
	}
	if err = binary.WriteRaw(w, p.Data); err != nil {
		return
	}
	return nil
}

type ServerResourcePackPush struct {
	Id     string
	Url    string
	Hash   string
	Forced bool
	//Prompt *string //todo needs to be a component, never write
}

func (p *ServerResourcePackPush) Direction() Direction { return Clientbound }
func (p *ServerResourcePackPush) ID(state State) int {
	return stateId2(state, Config, Play, ServerConfigAddResourcePackID, ServerPlayResourcePackPushID)
}

func (p *ServerResourcePackPush) Read(r io.Reader) (err error) {
	if p.Id, err = binary.ReadString(r); err != nil {
		return
	}
	if p.Url, err = binary.ReadString(r); err != nil {
		return
	}
	if p.Hash, err = binary.ReadString(r); err != nil {
		return
	}
	if p.Forced, err = binary.ReadBool(r); err != nil {
		return
	}
	hasPrompt, err := binary.ReadBool(r)
	if err != nil {
		return err
	}
	if hasPrompt {
		panic("cannot read resource pack push with prompt")
	}
	return nil
}

func (p *ServerResourcePackPush) Write(w io.Writer) (err error) {
	if err = binary.WriteString(w, p.Id); err != nil {
		return
	}
	if err = binary.WriteString(w, p.Url); err != nil {
		return
	}
	if err = binary.WriteString(w, p.Hash); err != nil {
		return
	}
	if err = binary.WriteBool(w, p.Forced); err != nil {
		return
	}
	if err = binary.WriteBool(w, false); err != nil {
		return
	}
	return nil
}

type ServerResourcePackPop struct {
	Id *string
}

func (p *ServerResourcePackPop) Direction() Direction { return Clientbound }
func (p *ServerResourcePackPop) ID(state State) int {
	return stateId2(state, Config, Play, ServerConfigRemoveResourcePackID, ServerPlayResourcePackPopID)
}

func (p *ServerResourcePackPop) Read(r io.Reader) (err error) {
	if p.Id, err = binary.ReadOptionalFunc(r, binary.ReadUUID); err != nil {
		return
	}
	return nil
}

func (p *ServerResourcePackPop) Write(w io.Writer) (err error) {
	if err = binary.WriteOptionalFunc(w, p.Id, binary.WriteUUID); err != nil {
		return
	}
	return nil
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
	if p.Channel, err = binary.ReadString(r); err != nil {
		return
	}
	if p.Data, err = binary.ReadRaw(r, binary.Remaining); err != nil {
		return
	}
	return nil
}

func (p *ServerPluginMessage) Write(w io.Writer) (err error) {
	if err = binary.WriteSizedString(w, p.Channel, 32767); err != nil {
		return //todo what actually is max length
	}
	if err = binary.WriteRaw(w, p.Data); err != nil {
		return
	}
	return nil
}

var (
	_ Packet = (*ClientResourcePackStatus)(nil)
	_ Packet = (*ClientPluginMessage)(nil)

	_ Packet = (*ServerResourcePackPush)(nil)
	_ Packet = (*ServerResourcePackPop)(nil)
	_ Packet = (*ServerPluginMessage)(nil)
)
