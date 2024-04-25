package binary

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	buffer2 "github.com/mworzala/kite/internal/pkg/buffer"
)

const Remaining = -1

// ReadRaw reads a raw byte array from the reader.
// The length may be Remaining to read whatever is present
func ReadRaw(r io.Reader, length int) ([]byte, error) {
	if length == Remaining {
		return ReadRemaining(r)
	}

	bytes := make([]byte, length)
	n, err := r.Read(bytes)
	if err != nil && (length < 0 || !errors.Is(err, io.EOF)) {
		return nil, err
	}
	if length >= 0 && n != length {
		return nil, fmt.Errorf("expected %d bytes but read %d", length, n)
	}
	return bytes, err
}

func ReadByte(r io.Reader) (byte, error) {
	var value byte
	err := binary.Read(r, binary.BigEndian, &value)
	return value, err
}

func ReadBool(r io.Reader) (bool, error) {
	value, err := ReadByte(r)
	return value != 0, err
}

func ReadUShort(r io.Reader) (uint16, error) {
	var value uint16
	err := binary.Read(r, binary.BigEndian, &value)
	return value, err
}

func ReadLong(r io.Reader) (int64, error) {
	var value int64
	err := binary.Read(r, binary.BigEndian, &value)
	return value, err
}

func ReadVarInt(r io.Reader) (int32, error) {
	var value int
	var position uint

	for {
		currentByte, err := ReadByte(r)
		if err != nil {
			return 0, err
		}

		value |= int(currentByte&varIntSegmentBits) << position

		if currentByte&varIntContinueBit == 0 {
			break
		}

		position += 7

		if position >= 32 {
			return 0, fmt.Errorf("VarInt is too big")
		}
	}

	return int32(value), nil
}

func ReadEnum[T Enum](r io.Reader) (value T, err error) {
	var raw int32
	if raw, err = ReadVarInt(r); err != nil {
		return
	}

	value = T(raw)
	if !value.Validate() {
		err = fmt.Errorf("invalid enum value for %T: %d", T(0), raw)
	}
	return
}

func ReadSizedString(r io.Reader, maxLength int) (string, error) {
	lengthName, err := ReadVarInt(r)
	if err != nil {
		return "", err
	}
	if lengthName > int32(maxLength) {
		return "", fmt.Errorf("string length %d exceeds maximum %d", lengthName, maxLength)
	}

	bytesName := make([]byte, lengthName)
	_, err = r.Read(bytesName)
	if err != nil {
		return "", err
	}

	return string(bytesName), nil
}

func ReadUUID(r io.Reader) (string, error) {
	var mostSigBits, leastSigBits uint64
	err := binary.Read(r, binary.BigEndian, &mostSigBits)
	if err != nil {
		return "", err
	}
	err = binary.Read(r, binary.BigEndian, &leastSigBits)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		mostSigBits>>32, (mostSigBits>>16)&0xFFFF, mostSigBits&0xFFFF,
		(leastSigBits>>48)&0xFFFF, leastSigBits&0xFFFFFFFFFFFF,
	), err
}

func ReadByteArray(r io.Reader) ([]byte, error) {
	length, err := ReadVarInt(r)
	if err != nil {
		return nil, err
	}
	return ReadRaw(r, int(length))
}

func ReadChatString(r io.Reader) (string, error) {
	return ReadSizedString(r, jsonChatLength)
}

func ReadCollection[T any](r io.Reader, read func(io.Reader) (T, error)) (values []T, err error) {
	var length int32
	if length, err = ReadVarInt(r); err != nil {
		return nil, err
	}
	values = make([]T, length)
	for i := range values {
		if values[i], err = read(r); err != nil {
			return nil, err
		}
	}
	return values, nil
}

func ReadRemaining(r io.Reader) ([]byte, error) {
	if remaining, ok := r.(*buffer2.PacketBuffer); ok {
		return remaining.AllocRemainder(), nil
	}
	return io.ReadAll(r)
}
