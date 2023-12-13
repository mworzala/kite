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

	// Play
)
