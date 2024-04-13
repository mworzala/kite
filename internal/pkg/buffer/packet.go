package buffer

import "io"

var (
	ErrBufferOverflow = io.ErrShortBuffer
)

type PacketBuffer struct {
	delegate []byte
	position int
	limit    int
}

func NewPacketBuffer(delegate []byte) *PacketBuffer {
	return &PacketBuffer{
		delegate: delegate,
		limit:    len(delegate),
	}
}

func (p2 *PacketBuffer) Read(p []byte) (n int, err error) {
	if p2.position+len(p) > p2.limit {
		return 0, ErrBufferOverflow
	}
	n = copy(p, p2.delegate[p2.position:])
	p2.position += n
	return
}

func (p2 *PacketBuffer) Remaining() int {
	return p2.limit - p2.position
}

func (p2 *PacketBuffer) AllocRemainder() []byte {
	rem := p2.Remaining()
	if rem == 0 {
		return nil
	}

	buf := make([]byte, rem)
	copy(buf, p2.delegate[p2.position:])
	p2.position += rem
	return buf
}

func (p2 *PacketBuffer) ShrinkRemaining() {
	p2.delegate = p2.delegate[p2.position:]
	p2.position = 0
}

func (p2 *PacketBuffer) Mark() int {
	return p2.position
}

func (p2 *PacketBuffer) Reset(mark int) {
	p2.position = mark
}

func (p2 *PacketBuffer) Limit(limit int) {
	if limit == -1 {
		p2.limit = len(p2.delegate)
		return
	}
	p2.limit = p2.position + limit
}
