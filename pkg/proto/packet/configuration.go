package packet

import (
	"io"

	"github.com/mworzala/kite/pkg/proto/binary"
)

const (
	ClientConfigClientInformationID = iota
	ClientConfigCookieResponseID
	ClientConfigPluginMessageID
	ClientConfigFinishConfigurationID
	ClientConfigKeepAliveID
	ClientConfigPongID
	ClientConfigResourcePackResponseID
	ClientConfigKnownPacksID
)

type ClientConfigPluginMessage struct {
	Channel string
	Data    []byte
}

func (p *ClientConfigPluginMessage) Direction() Direction { return Serverbound }
func (p *ClientConfigPluginMessage) ID(state State) int {
	return stateId1(state, Config, ClientConfigPluginMessageID)
}

func (p *ClientConfigPluginMessage) Read(r io.Reader) (err error) {
	if p.Channel, err = binary.ReadSizedString(r, nameLength); err != nil {
		return
	}
	if p.Data, err = binary.ReadRaw(r, binary.Remaining); err != nil {
		return
	}
	return nil
}

func (p *ClientConfigPluginMessage) Write(w io.Writer) (err error) {
	if err = binary.WriteSizedString(w, p.Channel, 32767); err != nil {
		return //todo what actually is max length
	}
	if err = binary.WriteRaw(w, p.Data); err != nil {
		return
	}
	return nil
}

type ClientConfigFinishConfiguration struct{}

func (p *ClientConfigFinishConfiguration) Direction() Direction { return Serverbound }
func (p *ClientConfigFinishConfiguration) ID(state State) int {
	return stateId1(state, Config, ClientConfigFinishConfigurationID)
}

func (p *ClientConfigFinishConfiguration) Read(r io.Reader) (err error) {
	return nil
}

func (p *ClientConfigFinishConfiguration) Write(w io.Writer) (err error) {
	return nil
}

const (
	ServerConfigCookieRequestID = iota
	ServerConfigPluginMessageID
	ServerConfigDisconnectID
	ServerConfigFinishConfigurationID
	ServerConfigKeepAliveID
	ServerConfigPingID
	ServerConfigResetChatID
	ServerConfigRegistryDataID
	ServerConfigRemoveResourcePackID
	ServerConfigAddResourcePackID
	ServerConfigStoreCookieID
	ServerConfigTransferID
	ServerConfigFeatureFlagsID
	ServerConfigUpdateTagsID
	ServerConfigKnownPacksID
)

type ServerConfigPluginMessage struct {
	Channel string
	Data    []byte
}

func (p *ServerConfigPluginMessage) Direction() Direction { return Clientbound }
func (p *ServerConfigPluginMessage) ID(state State) int {
	return stateId1(state, Config, ServerConfigPluginMessageID)
}

func (p *ServerConfigPluginMessage) Read(r io.Reader) (err error) {
	if p.Channel, err = binary.ReadSizedString(r, nameLength); err != nil {
		return
	}
	if p.Data, err = binary.ReadRaw(r, binary.Remaining); err != nil {
		return
	}
	return nil
}

func (p *ServerConfigPluginMessage) Write(w io.Writer) (err error) {
	if err = binary.WriteSizedString(w, p.Channel, 32767); err != nil {
		return //todo what actually is max length
	}
	if err = binary.WriteRaw(w, p.Data); err != nil {
		return
	}
	return nil
}

type ServerConfigDisconnect struct {
	//todo nbt
}

var (
	_ Packet = (*ClientConfigPluginMessage)(nil)
	_ Packet = (*ClientConfigFinishConfiguration)(nil)

	_ Packet = (*ServerConfigPluginMessage)(nil)
)
