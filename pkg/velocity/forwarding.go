package velocity

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"github.com/mworzala/kite/pkg/mojang"

	"github.com/mworzala/kite/pkg/buffer"
)

const (
	DefaultForwardingVersion = 1

	expectedBufferSize = 2048
)

func CreateSignedForwardingData(requestVersion int, secret []byte, address string, profile *mojang.GameProfile) (result []byte, err error) {
	mac := hmac.New(sha256.New, secret)

	// Currently we only support version 1. I guess we should support the others (relevant >= 1.20.6) in the future.
	version := DefaultForwardingVersion

	buf := new(bytes.Buffer)
	buf.Grow(mac.Size() + expectedBufferSize)
	buf.Write(make([]byte, mac.Size())) // Reserve space for the HMAC

	if err = buffer.VarInt.Write(buf, int32(version)); err != nil {
		return
	}
	if err = buffer.String.Write(buf, address); err != nil {
		return
	}
	if err = profile.Write(buf); err != nil {
		return
	}

	result = buf.Bytes()
	if _, err = mac.Write(result[mac.Size():]); err != nil {
		return
	}
	copy(result, mac.Sum(nil))

	return
}
