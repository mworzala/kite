package proto

import (
	"context"
	"fmt"
	"net"

	"github.com/mworzala/kite/pkg/proto/packet"
)

//todo this has turned into a mess...

// CreateServerConn connects to a remote server and returns a context which completes either when
// an error occurs or the connection is ready to enter the config phase. The returned context has
// a cause (see context.Cause). If the cause is not context.Canceled, the connection has been closed.
func CreateServerConn(ctx context.Context, address string, port uint16, profile *packet.GameProfile, handleFunc func(*Conn, context.CancelCauseFunc, *packet.GameProfile) Handler) (context.Context, *Conn, error) {
	var cancel context.CancelCauseFunc
	ctx, cancel = context.WithCancelCause(ctx)

	serverConn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", address, port))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to dial remote: %w", err)
	}
	remote, readLoop := NewConn(packet.Clientbound, serverConn)

	// Handshake immediately, then we are in login.
	handshake := &packet.ClientHandshake{
		ProtocolVersion: 766,
		ServerAddress:   address,
		ServerPort:      port,
		Intent:          packet.IntentLogin,
	}
	if err = remote.SendPacket(handshake); err != nil {
		return nil, nil, err
	}

	// Setup velocity forwarding handler & begin login.
	remote.SetState(packet.Login, handleFunc(remote, cancel, profile))
	err = remote.SendPacket(&packet.ClientLoginStart{
		Name: profile.Username,
		UUID: profile.UUID,
	})
	if err != nil {
		panic(err)
	}

	// Start reading from the remote connection
	go readLoop()

	return ctx, remote, nil
}
