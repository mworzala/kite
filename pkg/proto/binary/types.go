package binary

import "golang.org/x/exp/constraints"

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
