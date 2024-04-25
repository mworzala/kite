package nbt

type TagType int

const (
	TagEnd TagType = iota
	TagByte
	TagShort
	TagInt
	TagLong
	TagFloat
	TagDouble
	TagByteArray
	TagString
	TagList
	TagCompound
	TagIntArray
	TagLongArray
)

type Tag interface {
	Type() TagType
}

type (
	End         struct{}
	Byte        int8
	Short       int16
	Int         int32
	Long        int64
	Float       float32
	Double      float64
	ByteArray   []byte
	String      string
	List[T Tag] struct {
		Ty   TagType
		Data []T
	}
	Compound  map[string]Tag
	IntArray  []int32
	LongArray []int64
)

func (End) Type() TagType       { return TagEnd }
func (Byte) Type() TagType      { return TagByte }
func (Short) Type() TagType     { return TagShort }
func (Int) Type() TagType       { return TagInt }
func (Long) Type() TagType      { return TagLong }
func (Float) Type() TagType     { return TagFloat }
func (Double) Type() TagType    { return TagDouble }
func (ByteArray) Type() TagType { return TagByteArray }
func (String) Type() TagType    { return TagString }
func (List[T]) Type() TagType   { return TagList }
func (Compound) Type() TagType  { return TagCompound }
func (IntArray) Type() TagType  { return TagIntArray }
func (LongArray) Type() TagType { return TagLongArray }
