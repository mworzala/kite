package main

import (
	"log"
	"net"
)

func doLoginStuff() {
	addr := "localhost:25565"

	server, err := net.Dial("tcp", addr)
	if err != nil {
		log.Printf("Failed to connect to target: %s\n", err)
		return
	}
	defer server.Close()

	// Write the handshake packet
	handshakePacket := HandshakePacket{
		ProtocolVersion: 764,
		ServerAddress:   "localhost",
		ServerPort:      25565,
		NextState:       Login,
	}
	if err = writePacket(server, C2SHandshakeHandshake, handshakePacket, EncodeHandshake); err != nil {
		log.Printf("Failed to write handshake packet: %s\n", err)
		return
	}

	// Write the login start packet
	loginStartPacket := LoginStartPacket{
		Name:       "notmattw",
		PlayerUUID: "",
	}

}

func writePacket[T any](conn net.Conn, packetId int, packet T, encode func(T) ([]byte, error)) error {
	raw, err := encode(packet)
	if err != nil {
		return err
	}

	p := Packet{
		PacketID: uint32(packetId),
		Data:     raw,
		Length:   uint32(len(raw)),
	}

	return writePacketRaw(conn, p)
}

func writePacketRaw(conn net.Conn, p Packet) error {
	raw, err := Encode(p)
	if err != nil {
		return err
	}

	_, err = conn.Write(raw)
	return err
}
