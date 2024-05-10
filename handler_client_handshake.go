package kite

import (
	"fmt"

	"github.com/mworzala/kite/pkg/proto"
	"github.com/mworzala/kite/pkg/proto/packet"
)

// ClientHandshakeHandler implements the client->proxy handshake logic.
// Handles simple routing to the provided status and login handlers.
//
// If a handler function is not provided, the default handler will be used.
type ClientHandshakeHandler struct {
	ClientHandshakeHandlerOpts

	Conn *proto.Conn
}

type ClientHandshakeHandlerOpts struct {
	StatusHandlerFunc func(*proto.Conn) proto.Handler
	LoginHandlerFunc  func(*proto.Conn) proto.Handler
}

func MakeClientHandshakeHandler(opts ClientHandshakeHandlerOpts) func(*proto.Conn) proto.Handler {
	return func(conn *proto.Conn) proto.Handler {
		return &ClientHandshakeHandler{
			ClientHandshakeHandlerOpts: opts,
			Conn:                       conn,
		}
	}
}

func (h *ClientHandshakeHandler) HandlePacket(pp proto.Packet) error {
	var err error
	if pp.Id == packet.ClientHandshakeHandshakeID {
		p := new(packet.ClientHandshake)
		if err = pp.Read(p); err != nil {
			return err
		}
		return h.handleHandshake(p)
	}

	return proto.UnknownPacket
}

func (h *ClientHandshakeHandler) handleHandshake(p *packet.ClientHandshake) error {
	switch p.Intent {
	case packet.IntentStatus:
		var next proto.Handler
		if h.StatusHandlerFunc != nil {
			next = h.StatusHandlerFunc(h.Conn)
		} else {
			next = NewClientStatusHandler(h.Conn)
		}
		h.Conn.SetState(packet.Status, next)
	case packet.IntentLogin:
		var next proto.Handler
		if h.LoginHandlerFunc != nil {
			next = h.LoginHandlerFunc(h.Conn)
		} else {
			panic("todo implement no auth login")
			//next = NewServerboundLoginHandler(h.Conn)
		}
		h.Conn.SetState(packet.Login, next)
	case packet.IntentTransfer:
		return fmt.Errorf("transfer not supported")
	default:
		// An invalid intent would be caught when reading the packet, so this is impossible.
		panic("unreachable")
	}
	return nil
}

var _ proto.Handler = (*ClientHandshakeHandler)(nil)
