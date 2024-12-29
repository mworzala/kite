package packet

import (
	"io"
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

type ClientConfigFinishConfiguration struct{}

func (p *ClientConfigFinishConfiguration) Direction() Direction { return Serverbound }
func (p *ClientConfigFinishConfiguration) ID(state State) int {
	return stateId1(state, Config, ClientConfigFinishConfigurationID)
}
func (p *ClientConfigFinishConfiguration) Read(_ io.Reader) (err error) {
	return nil
}
func (p *ClientConfigFinishConfiguration) Write(_ io.Writer) (err error) {
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

var (
	_ Packet = (*ClientConfigFinishConfiguration)(nil)

	//_ Packet = (*ServerConfigPluginMessage)(nil)
)
