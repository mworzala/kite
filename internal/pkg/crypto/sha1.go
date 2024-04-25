package crypto

import (
	"crypto/sha1"
	"encoding/hex"
	"strings"
)

// Sha1 computes a Minecraft protocol compatible SHA-1 hash, which is notable different
// from a standard implementation.
//
// Based on https://gist.github.com/toqueteos/5372776
func Sha1(entries ...[]byte) string {
	h := sha1.New()
	for _, entry := range entries {
		if _, err := h.Write(entry); err != nil {
			panic(err)
		}
	}
	hash := h.Sum(nil)

	// Check for negative hashes
	negative := (hash[0] & 0x80) == 0x80
	if negative {
		hash = twosComplement(hash)
	}

	// Trim away zeroes
	res := strings.TrimLeft(hex.EncodeToString(hash), "0")
	if negative {
		res = "-" + res
	}

	return res
}

// little endian
func twosComplement(p []byte) []byte {
	carry := true
	for i := len(p) - 1; i >= 0; i-- {
		p[i] = byte(^p[i])
		if carry {
			carry = p[i] == 0xff
			p[i]++
		}
	}
	return p
}
