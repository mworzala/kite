package main

import (
	"bytes"
	"encoding/binary"
	"errors"
)

type PacketDirection bool

const (
	Clientbound PacketDirection = false
	Serverbound PacketDirection = true
)

type GameState int

const (
	Handshake GameState = iota
	Status
	Login
	Configuration
	Play
)

const (
	segmentBits = 0x7F // 0111 1111
	continueBit = 0x80 // 1000 0000
)

//
// Packet types
//

type Packet struct {
	Length   uint32
	PacketID uint32
	Data     []byte
}

type HandshakePacket struct {
	ProtocolVersion int
	ServerAddress   string
	ServerPort      uint16
	NextState       GameState
}

type LoginStartPacket struct {
	Name       string // 16 bytes!
	PlayerUUID string
}

type Property struct {
	Name      string
	Value     string
	Signature *string
	isSigned  bool
}

type LoginSuccessPacket struct {
	UUID          string
	Username      string
	NumProperties int
	Property      []Property
}

type LoginPluginRequestPacket struct {
	MessageID int
	Channel   string
	Data      []byte
}

//
// Decoding functions
//

// Decode takes a byte array and returns a Packet struct, the number of bytes processed, and an error
func Decode(packetb []byte) (Packet, error) {

	reader := bytes.NewReader(packetb)
	var bytesRead int64 = 0

	// decode the packet Id
	packetId, err := readVarInt(reader)
	if err != nil {
		return Packet{}, err
	}

	packetIdLength := int64(reader.Size() - int64(reader.Len()))
	bytesRead += packetIdLength

	return Packet{
		Length:   uint32(0),
		PacketID: uint32(packetId),
		Data:     packetb[bytesRead:],
	}, nil
}

func Encode(packet Packet) ([]byte, error) {
	// Create a buffer to write our data to
	buf := new(bytes.Buffer)

	// Encode the packet ID
	err := writeVarInt(buf, int(packet.PacketID))
	if err != nil {
		return nil, err
	}

	// Append the data
	_, err = buf.Write(packet.Data)
	if err != nil {
		return nil, err
	}

	// Convert the buffer to a byte slice and return
	return buf.Bytes(), nil
}

// DecodeHandshake takes a Packet struct and returns a HandshakePacket struct and an error
func DecodeHandshake(pkt *Packet) (HandshakePacket, error) {
	reader := bytes.NewReader(pkt.Data)

	// decode the protocol version
	protocolVersion, err := readVarInt(reader)
	if err != nil {
		return HandshakePacket{}, err
	}

	// decode the server address
	bytesServerAddress := make([]byte, 255)
	_, err = reader.Read(bytesServerAddress)
	if err != nil {
		return HandshakePacket{}, err
	}

	var bytesServerPort [2]byte
	// decode the server port
	_, err = reader.Read(bytesServerPort[:])
	if err != nil {
		return HandshakePacket{}, err
	}

	// decode the next state
	nextState, err := readVarInt(reader)
	if err != nil {
		return HandshakePacket{}, err
	}

	return HandshakePacket{
		ProtocolVersion: protocolVersion,
		ServerAddress:   string(bytesServerAddress),
		ServerPort:      binary.BigEndian.Uint16(bytesServerPort[:]),
		NextState:       GameState(nextState),
	}, nil
}

func EncodeHandshake(hsp HandshakePacket) ([]byte, error) {
	var buf bytes.Buffer

	// encode the protocol version
	if err := writeVarInt(&buf, hsp.ProtocolVersion); err != nil {
		return nil, err
	}

	// encode the server address
	if len(hsp.ServerAddress) > 255 {
		return nil, errors.New("server address too long")
	}
	if _, err := buf.WriteString(hsp.ServerAddress); err != nil {
		return nil, err
	}

	// encode the server port
	var bytesServerPort [2]byte
	binary.BigEndian.PutUint16(bytesServerPort[:], hsp.ServerPort)
	if _, err := buf.Write(bytesServerPort[:]); err != nil {
		return nil, err
	}

	// encode the next state
	if err := writeVarInt(&buf, int(hsp.NextState)); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// DecodeLoginStart takes a Packet struct and returns a LoginStartPacket struct and an error
func DecodeLoginStart(pkt *Packet) (LoginStartPacket, error) {
	reader := bytes.NewReader(pkt.Data)

	playerName, err := readString(reader)
	if err != nil {
		return LoginStartPacket{}, err
	}

	uuid, err := readUUID(reader)
	if err != nil {
		return LoginStartPacket{}, err
	}

	return LoginStartPacket{
		Name:       playerName,
		PlayerUUID: uuid,
	}, nil
}

// DecodeLoginSuccess takes a Packet struct and returns a LoginSuccessPacket struct and an error
func DecodeLoginSuccess(pkt *Packet) (LoginSuccessPacket, error) {
	reader := bytes.NewReader(pkt.Data)

	// decode the uuid
	uuid, err := readUUID(reader)
	if err != nil {
		return LoginSuccessPacket{}, err
	}

	username, err := readString(reader) // 16 => 2 bytes
	if err != nil {
		return LoginSuccessPacket{}, err
	}

	// decode the number of properties
	numProperties, err := readVarInt(reader)
	if err != nil {
		return LoginSuccessPacket{}, err
	}

	// decode the properties
	properties := make([]Property, numProperties)
	for i := 0; i < numProperties; i++ {
		// decode the name
		name, err := readString(reader)
		if err != nil {
			return LoginSuccessPacket{}, err
		}

		// decode the value
		value, err := readString(reader)
		if err != nil {
			return LoginSuccessPacket{}, err
		}

		// decode the signature
		var signature *string
		isSigned, err := reader.ReadByte()
		if err != nil {
			return LoginSuccessPacket{}, err
		}

		bIsSigned := isSigned == 1

		if bIsSigned {
			bytesSignature := make([]byte, 16)
			_, err = reader.Read(bytesSignature)
			if err != nil {
				return LoginSuccessPacket{}, err
			}

			signatureStr := string(bytesSignature)
			signature = &signatureStr
		}

		properties[i] = Property{
			Name:      name,
			Value:     value,
			Signature: signature,
			isSigned:  bIsSigned,
		}
	}

	return LoginSuccessPacket{
		UUID:          uuid,
		Username:      username,
		NumProperties: numProperties,
		Property:      properties,
	}, nil
}
