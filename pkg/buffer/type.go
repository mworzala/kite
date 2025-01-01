package buffer

import (
	"io"

	"github.com/google/uuid"
	"github.com/mworzala/kite/pkg/text"
	"golang.org/x/exp/constraints"
)

type Type[T any] interface {
	Read(r io.Reader) (T, error)
	Write(w io.Writer, v T) error
}

var (
	Byte      Type[byte]      = byteType{}
	Bool      Type[bool]      = boolType{}
	Uint16    Type[uint16]    = uShortType{}
	VarInt    Type[int32]     = varIntType{}
	Long      Type[int64]     = longType{}
	UUID      Type[uuid.UUID] = uuidType{}
	String    Type[string]    = stringType{}
	ByteArray Type[[]byte]    = byteArrayType{}
	RawBytes  Type[[]byte]    = rawBytesType{}

	TextComponent     Type[text.Component] = textComponentType{}
	TextComponentJSON Type[text.Component] = textComponentJSONType{}
)

// Complex types

func Opt[T comparable](t Type[T]) Type[T] {
	return optType[T]{t}
}

type Enum[T interface {
	constraints.Integer
	Validate() bool
}] struct{}

func JSON[T any]() Type[T] {
	return jsonType[T]{}
}

func Struct[T interface {
	Read(r io.Reader) error
	Write(w io.Writer) error
}]() Type[T] {
	return structType[T]{}
}

func List[T any](t Type[T]) Type[[]T] {
	return listType[T]{t}
}
