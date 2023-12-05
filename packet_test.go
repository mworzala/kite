package main

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"testing"
)

func TestDecodeLoginStartPacket(t *testing.T) {
	testCases := []struct {
		name     string
		packet   Packet
		expected LoginStartPacket
		wantErr  bool
	}{
		{
			name: "Valid packet",
			packet: Packet{
				Length:   uint32(32), // 16 for player name + 16 bytes for uuid
				PacketID: 1,          // Example PacketID, adjust as needed
				Data:     append([]byte{0x0a}, append([]byte("PlayerName"), mockPlayerUUIDData()...)...),
			},
			expected: LoginStartPacket{
				Name:       "PlayerName",
				PlayerUUID: "00000000-075b-cd15-0000-00003ade68b1",
			},
			wantErr: false,
		},
		{
			name: "Invalid packet - short data",
			packet: Packet{
				Length:   5,
				PacketID: 1,
				Data:     []byte("Short"),
			},
			wantErr: true,
		},
		// Additional test cases go here
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := DecodeLoginStartPacket(&tc.packet)
			if (err != nil) != tc.wantErr {
				t.Errorf("DecodeLoginStartPacket() error = %v, wantErr %v for case %s", err, tc.wantErr, tc.name)
				return
			}
			if !tc.wantErr && !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("DecodeLoginStartPacket() = %v, want %v for case %s", result, tc.expected, tc.name)
			}
		})
	}
}

// Mock function to create player UUID data
func mockPlayerUUIDData() []byte {
	buf := new(bytes.Buffer)
	// UUID most significant bits (mock value)
	binary.Write(buf, binary.BigEndian, uint64(123456789))
	// UUID least significant bits (mock value)
	binary.Write(buf, binary.BigEndian, uint64(987654321))
	return buf.Bytes()
}
