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
	varIntSegmentBits = 0x7F // 0111 1111
	varIntContinueBit = 0x80 // 1000 0000
)

const (
	jsonChatLength = 262144
)
