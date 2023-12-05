package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

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

func readString(r *bytes.Reader) (string, error) {

	// decode the name
	lengthName, err := readVarInt(r)
	if err != nil {
		return "", err
	}

	bytesName := make([]byte, lengthName)
	_, err = r.Read(bytesName)
	if err != nil {
		return "", err
	}

	return string(bytesName), nil
}

// readUUID reads a UUID from a bytes.Reader
func readUUID(r *bytes.Reader) (string, error) {
	var mostSigBits, leastSigBits uint64

	// decode the uuid
	err := binary.Read(r, binary.BigEndian, &mostSigBits)
	if err != nil {
		return "", err
	}
	err = binary.Read(r, binary.BigEndian, &leastSigBits)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", mostSigBits>>32, (mostSigBits>>16)&0xFFFF, mostSigBits&0xFFFF, (leastSigBits>>48)&0xFFFF, leastSigBits&0xFFFFFFFFFFFF), err
}
