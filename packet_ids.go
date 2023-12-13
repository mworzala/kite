package main

const (
	// Handshake
	C2SHandshakeHandshake = 0x00

	// Status
	C2SStatusStatusRequest = 0x00
	C2SStatusPingRequest   = 0x01

	// Login
	C2SLoginLoginStart        = 0x00
	C2SLoginLoginEncryption   = 0x01
	C2SLoginLoginAcknowledged = 0x03

	// Configuration
	C2SConfClientInfo           = 0x00
	C2SConfPluginMsg            = 0x01
	C2SConfFinishConfig         = 0x02
	C2SConfKeepAlive            = 0x03
	C2SConfPong                 = 0x04
	C2SConfResourcePackResponse = 0x05

	// Play
)

const (
	// Status
	S2CStatusStatusResponse = 0x00
	S2CStatusPingResponse   = 0x01

	// Login
	S2CLoginDisconnect   = 0x00
	S2CLoginEncryption   = 0x01
	S2CLoginLoginSuccess = 0x02

	// Configuration
	S2CConfPluginMsg          = 0x00
	S2CConfDisconnect         = 0x01
	S2CConfFinish             = 0x02
	S2CConfKeepAlive          = 0x03
	S2CConfPing               = 0x04
	S2CConfRegistryData       = 0x05
	S2CConfRemoveResourcePack = 0x06
	S2CConfAddResourcePack    = 0x07
	S2CConfAddFeatureFlags    = 0x08
	S2CConfUpdateTags         = 0x09

	// Play
)
