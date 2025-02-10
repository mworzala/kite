package buffer

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/google/uuid"
	"github.com/mworzala/kite/pkg/text"
)

const (
	varIntSegmentBits = uint32(0x7F) // 0111 1111
	varIntContinueBit = 0x80         // 1000 0000
)

type byteType struct{}

func (byteType) Read(r io.Reader) (byte, error) {
	var value byte
	err := binary.Read(r, binary.BigEndian, &value)
	return value, err
}
func (byteType) Write(w io.Writer, v byte) error {
	return binary.Write(w, binary.BigEndian, v)
}

type boolType struct{}

func (boolType) Read(r io.Reader) (bool, error) {
	var value byte
	err := binary.Read(r, binary.BigEndian, &value)
	return value != 0, err
}
func (boolType) Write(w io.Writer, v bool) error {
	if v {
		return binary.Write(w, binary.BigEndian, byte(1))
	}
	return binary.Write(w, binary.BigEndian, byte(0))
}

type uShortType struct{}

func (uShortType) Read(r io.Reader) (uint16, error) {
	var value uint16
	err := binary.Read(r, binary.BigEndian, &value)
	return value, err
}
func (uShortType) Write(w io.Writer, v uint16) error {
	return binary.Write(w, binary.BigEndian, v)
}

type varIntType struct{}

func (varIntType) Read(r io.Reader) (int32, error) {
	var value int
	var position uint

	for {
		currentByte, err := Byte.Read(r)
		if err != nil {
			return 0, err
		}

		value |= int(currentByte&byte(varIntSegmentBits)) << position

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
func (varIntType) Write(w io.Writer, v int32) error {
	value := uint32(v)
	for {
		if (value & ^varIntSegmentBits) == 0 {
			return Byte.Write(w, byte(value))
		}

		if err := Byte.Write(w, byte(value&varIntSegmentBits)|varIntContinueBit); err != nil {
			return err
		}

		value >>= 7 // Go automatically handles sign bit shift differently from Java's >>>
	}
}

type longType struct{}

func (longType) Read(r io.Reader) (int64, error) {
	var value int64
	err := binary.Read(r, binary.BigEndian, &value)
	return value, err
}
func (longType) Write(w io.Writer, v int64) error {
	return binary.Write(w, binary.BigEndian, v)
}

type uuidType struct{}

func (uuidType) Read(r io.Reader) (_ uuid.UUID, err error) {
	var mostSigBits, leastSigBits uint64
	if err = binary.Read(r, binary.BigEndian, &mostSigBits); err != nil {
		return
	}
	if err = binary.Read(r, binary.BigEndian, &leastSigBits); err != nil {
		return
	}
	return uuid.UUID{
		byte(mostSigBits >> 56), byte(mostSigBits >> 48), byte(mostSigBits >> 40), byte(mostSigBits >> 32),
		byte(mostSigBits >> 24), byte(mostSigBits >> 16), byte(mostSigBits >> 8), byte(mostSigBits),
		byte(leastSigBits >> 56), byte(leastSigBits >> 48), byte(leastSigBits >> 40), byte(leastSigBits >> 32),
		byte(leastSigBits >> 24), byte(leastSigBits >> 16), byte(leastSigBits >> 8), byte(leastSigBits),
	}, nil
}
func (uuidType) Write(w io.Writer, v uuid.UUID) error {
	mostSigBits := binary.BigEndian.Uint64(v[:8])
	leastSigBits := binary.BigEndian.Uint64(v[8:])
	if err := binary.Write(w, binary.BigEndian, mostSigBits); err != nil {
		return err
	}
	return binary.Write(w, binary.BigEndian, leastSigBits)
}

type stringType struct{}

func (stringType) Read(r io.Reader) (string, error) {
	length, err := VarInt.Read(r)
	if err != nil {
		return "", err
	}
	remainder := r.(*Buffer).Remaining()
	if length > int32(remainder) {
		return "", fmt.Errorf("string length %d exceeds maximum %d", length, remainder)
	}
	value := make([]byte, length)
	_, err = r.Read(value)
	if err != nil {
		return "", err
	}
	return string(value), nil
}
func (stringType) Write(w io.Writer, v string) error {
	if err := VarInt.Write(w, int32(len(v))); err != nil {
		return err
	}
	if _, err := w.Write([]byte(v)); err != nil {
		return err
	}
	return nil
}

type byteArrayType struct{}

func (byteArrayType) Read(r io.Reader) ([]byte, error) {
	length, err := VarInt.Read(r)
	if err != nil {
		return nil, err
	}
	bytes := make([]byte, length)
	n, err := r.Read(bytes)
	if err != nil && (length < 0 || !errors.Is(err, io.EOF)) {
		return nil, err
	}
	if length >= 0 && n != int(length) {
		return nil, fmt.Errorf("expected %d bytes but read %d", length, n)
	}
	return bytes, err
}
func (byteArrayType) Write(w io.Writer, v []byte) error {
	if err := VarInt.Write(w, int32(len(v))); err != nil {
		return err
	}
	n, err := w.Write(v)
	if err != nil {
		return err
	}
	if n != len(v) {
		return fmt.Errorf("expected to write %d bytes but wrote %d", len(v), n)
	}
	return nil
}

type rawBytesType struct{}

func (rawBytesType) Read(r io.Reader) ([]byte, error) {
	if remaining, ok := r.(*Buffer); ok {
		return remaining.AllocRemainder(), nil
	}
	return io.ReadAll(r)
}
func (rawBytesType) Write(w io.Writer, v []byte) error {
	n, err := w.Write(v)
	if err != nil {
		return err
	}
	if n != len(v) {
		return fmt.Errorf("expected to write %d bytes but wrote %d", len(v), n)
	}
	return nil
}

type textComponentType struct{}

func (textComponentType) Read(r io.Reader) (text.Component, error) {
	// TODO
	return nil, nil
}
func (textComponentType) Write(w io.Writer, v text.Component) error {
	// TODO
	return nil
}

type textComponentJSONType struct{}

func (textComponentJSONType) Read(r io.Reader) (text.Component, error) {
	_, _ = String.Read(r)
	return nil, nil
}
func (textComponentJSONType) Write(w io.Writer, v text.Component) error {
	raw, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return String.Write(w, string(raw))
}

// Enum

func (e Enum[T]) Read(r io.Reader) (value T, err error) {
	var raw int32
	if raw, err = VarInt.Read(r); err != nil {
		return
	}

	value = T(raw)
	if !value.Validate() {
		err = fmt.Errorf("invalid enum value for %T: %d", T(0), raw)
	}
	return
}
func (e Enum[T]) Write(w io.Writer, v T) error {
	if !v.Validate() {
		return fmt.Errorf("invalid enum value for %T: %d", v, v)
	}
	return VarInt.Write(w, int32(v))
}

type optType[T comparable] struct {
	t Type[T]
}

func (o optType[T]) Read(r io.Reader) (t T, err error) {
	present, err := Bool.Read(r)
	if err != nil || !present {
		return t, err
	}
	return o.t.Read(r)
}
func (o optType[T]) Write(w io.Writer, v T) error {
	var zero T
	isZero := v == zero
	if err := Bool.Write(w, !isZero); err != nil || isZero {
		return err
	}
	return o.t.Write(w, v)
}

type jsonType[T any] struct{}

func (jsonType[T]) Read(r io.Reader) (t T, err error) {
	var bytes []byte
	if bytes, err = ByteArray.Read(r); err != nil {
		return
	}
	err = json.Unmarshal(bytes, &t)
	return
}
func (jsonType[T]) Write(w io.Writer, v T) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return ByteArray.Write(w, data)
}

type listType[T any] struct {
	t Type[T]
}

func (l listType[T]) Read(r io.Reader) (t []T, err error) {
	length, err := VarInt.Read(r)
	if err != nil {
		return
	}
	t = make([]T, length)
	for i := range t {
		if t[i], err = l.t.Read(r); err != nil {
			return
		}
	}
	return
}
func (l listType[T]) Write(w io.Writer, v []T) error {
	if err := VarInt.Write(w, int32(len(v))); err != nil {
		return err
	}
	for _, value := range v {
		if err := l.t.Write(w, value); err != nil {
			return err
		}
	}
	return nil
}

type structType[T interface {
	Read(r io.Reader) error
	Write(w io.Writer) error
}] struct {
}

func (s structType[T]) Read(r io.Reader) (t T, err error) {
	v := *new(T)
	if err = v.Read(r); err != nil {
		return
	}
	return
}

func (s structType[T]) Write(w io.Writer, v T) error {
	return v.Write(w)
}
