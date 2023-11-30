package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
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

func readVarInt(r *bytes.Reader) (int, error) {
	var value int
	var position uint

	for {
		currentByte, err := r.ReadByte()
		if err != nil {
			return 0, err
		}

		value |= int(currentByte&segmentBits) << position

		if currentByte&continueBit == 0 {
			break
		}

		position += 7

		if position >= 32 {
			return 0, fmt.Errorf("VarInt is too big")
		}
	}

	return value, nil
}

func writeVarInt(w *bytes.Buffer, value int) error {
	for {
		if (value & ^segmentBits) == 0 {
			return w.WriteByte(byte(value))
		}

		if err := w.WriteByte(byte(value&segmentBits) | continueBit); err != nil {
			return err
		}

		value >>= 7 // Go automatically handles sign bit shift differently from Java's >>>
	}
}

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

// BuildUUID creates a UUID from a two 64bit integers
func BuildUUID(mostSigBits, leastSigBits uint64) string {
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", mostSigBits>>32, (mostSigBits>>16)&0xFFFF, mostSigBits&0xFFFF, (leastSigBits>>48)&0xFFFF, leastSigBits&0xFFFFFFFFFFFF)
}

// DecodeLoginStartPacket takes a Packet struct and returns a LoginStartPacket struct and an error
func DecodeLoginStartPacket(pkt *Packet) (LoginStartPacket, error) {
	reader := bytes.NewReader(pkt.Data)

	// decode the name
	bytesName := make([]byte, 16)
	_, err := reader.Read(bytesName)
	if err != nil {
		return LoginStartPacket{}, err
	}

	// decode the player uuid
	var uuidMSD, uuidLSD uint64
	err = binary.Read(reader, binary.BigEndian, &uuidMSD)
	if err != nil {
		return LoginStartPacket{}, err
	}
	err = binary.Read(reader, binary.BigEndian, &uuidLSD)
	if err != nil {
		return LoginStartPacket{}, err
	}

	return LoginStartPacket{
		Name:       string(bytesName),
		PlayerUUID: BuildUUID(uuidMSD, uuidLSD),
	}, nil
}
