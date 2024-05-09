package binary

import (
	"io"

	"golang.org/x/exp/constraints"
)

type (
	ReadFunc[T any]  func(*T, io.Reader) error
	WriteFunc[T any] func(*T, io.Writer) error
)

type Enum interface {
	constraints.Integer
	Validate() bool
}

const (
	varIntSegmentBits = uint32(0x7F) // 0111 1111
	varIntContinueBit = 0x80         // 1000 0000
)

const (
	jsonChatLength = 262144
)

func MeasureVarInt(value int32) int {
	val := uint32(value)
	length := 0
	for {
		length++
		if (val & ^varIntSegmentBits) == 0 {
			break
		}
		val >>= 7
	}

	return length
}
