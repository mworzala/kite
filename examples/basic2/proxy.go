package main

import (
	"github.com/mworzala/kite/pkg/mojang"
)

type Proxy struct {
	MojKeyPair mojang.KeyPair

	VelocitySecret string
}
