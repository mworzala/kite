package binary

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/google/uuid"
)

func WriteRaw(w io.Writer, value []byte) error {
	n, err := w.Write(value)
	if err != nil {
		return err
	}
	if n != len(value) {
		return fmt.Errorf("expected to write %d bytes but wrote %d", len(value), n)
	}
	return nil
}

func WriteByte(w io.Writer, value byte) error {
	return binary.Write(w, binary.BigEndian, value)
}

func WriteBool(w io.Writer, value bool) error {
	if value {
		return WriteByte(w, 1)
	}
	return WriteByte(w, 0)
}

func WriteUShort(w io.Writer, value uint16) error {
	return binary.Write(w, binary.BigEndian, value)
}

func WriteLong(w io.Writer, value int64) error {
	return binary.Write(w, binary.BigEndian, value)
}

func WriteVarInt(w io.Writer, value int32) error {
	for {
		if (value & ^varIntSegmentBits) == 0 {
			return WriteByte(w, byte(value))
		}

		if err := WriteByte(w, byte(value&varIntSegmentBits)|varIntContinueBit); err != nil {
			return err
		}

		value >>= 7 // Go automatically handles sign bit shift differently from Java's >>>
	}
}

func WriteEnum[T Enum](w io.Writer, value T) error {
	if !value.Validate() {
		return fmt.Errorf("invalid enum value for %T: %d", value, value)
	}

	return WriteVarInt(w, int32(value))
}

func WriteSizedString(w io.Writer, value string, maxLength int) error {
	if len(value) > maxLength {
		return fmt.Errorf("string length %d exceeds max length %d", len(value), maxLength)
	}
	if err := WriteVarInt(w, int32(len(value))); err != nil {
		return err
	}
	if _, err := w.Write([]byte(value)); err != nil {
		return err
	}
	return nil
}

func WriteUUID(w io.Writer, value string) error {
	u, err := uuid.Parse(value)
	if err != nil {
		return err
	}
	//todo is this correct?
	return binary.Write(w, binary.BigEndian, u)
}

func WriteChatString(w io.Writer, value string) error {
	return WriteSizedString(w, value, jsonChatLength)
}
