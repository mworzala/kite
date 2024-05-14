package kite

import "errors"

var (
	ErrAlreadyConnecting = errors.New("already connecting")
)

type ServerInfo struct {
	Address string
	Port    int
	Secret  string // Velocity forwarding secret
}
