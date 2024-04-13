package server

import "net"

type Connection struct {
	delegate net.Conn
}

func NewServerConnection(address string) (*Connection, error) {
	delegate, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	return &Connection{delegate: delegate}, nil
}
