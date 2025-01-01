package buffer

import (
	"io"

	"github.com/valyala/bytebufferpool"
)

var ErrBufferOverflow = io.ErrShortBuffer

type Buffer struct {
	delegate []byte
	position int
	limit    int
}

func Wrap(delegate []byte) *Buffer {
	return &Buffer{
		delegate: delegate,
		limit:    len(delegate),
	}
}

// Read implements the io.Reader interface
func (b *Buffer) Read(p []byte) (n int, err error) {
	if b.position+len(p) > b.limit {
		return 0, ErrBufferOverflow
	}
	n = copy(p, b.delegate[b.position:])
	b.position += n
	return
}

func (b *Buffer) Remaining() int {
	return b.limit - b.position
}

func (b *Buffer) AllocRemainder() []byte {
	rem := b.Remaining()
	if rem == 0 {
		return nil
	}

	buf := make([]byte, rem)
	copy(buf, b.delegate[b.position:])
	b.position += rem
	return buf
}

func (b *Buffer) AllocRemainderTo(other *bytebufferpool.ByteBuffer) {
	rem := b.Remaining()
	if rem == 0 {
		return
	}

	_, _ = other.Write(b.delegate[b.position:])
	b.position += rem
}

func (b *Buffer) RemainingSlice() []byte {
	sl := b.delegate[b.position:b.limit]
	b.position = b.limit
	return sl
}

func (b *Buffer) ShrinkRemaining() {
	b.delegate = b.delegate[b.position:]
	b.position = 0
}

func (b *Buffer) Mark() int {
	return b.position
}

func (b *Buffer) Reset(mark int) {
	b.position = mark
}

func (b *Buffer) Limit(limit int) {
	if limit == -1 {
		b.limit = len(b.delegate)
		return
	}
	b.limit = b.position + limit
}
