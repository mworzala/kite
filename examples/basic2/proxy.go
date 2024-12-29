package main

import "crypto/rsa"

type Proxy struct {
	PrivateKey  *rsa.PrivateKey
	PublicKey   []byte
	VerifyToken []byte

	VelocitySecret string
}
