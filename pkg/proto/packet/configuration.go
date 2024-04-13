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

//type ClientLoginAcknowledged struct{}
//
//func (p *ClientLoginAcknowledged) Direction() Direction { return Serverbound }
//func (p *ClientLoginAcknowledged) ID(state State) int {
//	return stateId1(state, Login, ClientLoginLoginAcknowledgedID)
//}
//
//func (p *ClientLoginAcknowledged) Read(r io.Reader) (err error) {
//	return nil
//}
//func (p *ClientLoginAcknowledged) Write(w io.Writer) (err error) {
//	return nil
//}

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

func (p *ServerConfigPluginMessage) Direction() Direction { return Serverbound }
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

//type ServerLoginDisconnect struct {
//	Reason string
//}
//
//func (p *ServerLoginDisconnect) Direction() Direction { return Clientbound }
//func (p *ServerLoginDisconnect) ID(state State) int {
//	return stateId1(state, Login, ServerLoginDisconnectID)
//}
//
//func (p *ServerLoginDisconnect) Read(r io.Reader) (err error) {
//	if p.Reason, err = binary.ReadChatString(r); err != nil {
//		return
//	}
//	return nil
//}
//
//func (p *ServerLoginDisconnect) Write(w io.Writer) (err error) {
//	if err = binary.WriteChatString(w, p.Reason); err != nil {
//		return
//	}
//	return nil
//}
//
//type ServerLoginSuccess struct {
//	UUID     string
//	Username string
//	//todo properties
//}
//
//func (p *ServerLoginSuccess) Direction() Direction { return Clientbound }
//func (p *ServerLoginSuccess) ID(state State) int {
//	return stateId1(state, Login, ServerLoginLoginSuccessID)
//}
//
//func (p *ServerLoginSuccess) Read(r io.Reader) (err error) {
//	if p.UUID, err = binary.ReadUUID(r); err != nil {
//		return
//	}
//	if p.Username, err = binary.ReadSizedString(r, 16); err != nil {
//		return
//	}
//	_, _ = binary.ReadVarInt(r) //todo properties
//	return nil
//}
//
//func (p *ServerLoginSuccess) Write(w io.Writer) (err error) {
//	if err = binary.WriteUUID(w, p.UUID); err != nil {
//		return
//	}
//	if err = binary.WriteSizedString(w, p.Username, 16); err != nil {
//		return
//	}
//	_ = binary.WriteVarInt(w, 0) //todo properties
//	return nil
//}
//
//type ServerLoginPluginRequest struct {
//	MessageID int32
//	Channel   string
//	Data      []byte
//}
//
//func (p *ServerLoginPluginRequest) Direction() Direction { return Clientbound }
//func (p *ServerLoginPluginRequest) ID(state State) int {
//	return stateId1(state, Login, ServerLoginPluginRequestID)
//}
//
//func (p *ServerLoginPluginRequest) Read(r io.Reader) (err error) {
//	if p.MessageID, err = binary.ReadVarInt(r); err != nil {
//		return
//	}
//	if p.Channel, err = binary.ReadSizedString(r, 20); err != nil {
//		return
//	}
//	if p.Data, err = binary.ReadRaw(r, -1); err != nil {
//		return
//	}
//	return nil
//}
//
//func (p *ServerLoginPluginRequest) Write(w io.Writer) (err error) {
//	if err = binary.WriteVarInt(w, p.MessageID); err != nil {
//		return
//	}
//	if err = binary.WriteSizedString(w, p.Channel, 20); err != nil {
//		return
//	}
//	if err = binary.WriteRaw(w, p.Data); err != nil {
//		return
//	}
//	return nil
//}

var (
	_ Packet = (*ClientConfigPluginMessage)(nil)

	_ Packet = (*ServerConfigPluginMessage)(nil)
)
