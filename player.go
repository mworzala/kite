package kite

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/mworzala/kite/pkg/proto"
	"github.com/mworzala/kite/pkg/proto/packet"
)

type Player[T any] struct {
	*proto.Conn

	UUID     uuid.UUID
	Username string
	Profile  *packet.GameProfile

	// Generic user data for proxy users to assign
	Data *T

	// Protects against multiple connections using ConnectTo
	connecting atomic.Bool
}

// SendPacket sends a packet to the client connection.
//
// To send a server packet, first get the Server().
func (p *Player[T]) SendPacket(pkt packet.Packet) error {
	return p.Conn.SendPacket(pkt)
}

// Server returns the current server for the player, or nil if they are not connected to a server.
func (p *Player[T]) Server() *proto.Conn {
	return p.Conn.GetRemote()
}

// ConnectTo attempts to connect to the given server, blocking until the connection is successful.
// Returns nil if the connection was successful, the error otherwise.
// If the target already has a pending connection, this method will immediately return ErrAlreadyConnecting.
//
// Note that it is invalid to call this method from the clients reading goroutine. It will block
// waiting for the switch to happen, which will be a deadlock if called from the clients read thread.
// To be safe, this method should always be called from a new goroutine.
//
// This method implements one form of server switching including relevant checks to prevent a forced
// disconnection on failure, however it is possible to reimplement externally if more control is needed.
func (p *Player[T]) ConnectTo(server *ServerInfo, clientConfigHandlerFunc func(*Player[T]) proto.Handler, serverConfigHandlerFunc func(*Player[T]) proto.Handler) (err error) {
	// Prevent multiple connections at the same time by using the atomic bool kinda like a (try)lock.
	if !p.connecting.CompareAndSwap(false, true) {
		return ErrAlreadyConnecting
	}
	defer p.connecting.Store(false)

	//todo need to validate that the player is in at least the config phase also

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var remote *proto.Conn
	ch := make(chan bool)
	ctx, remote, err = proto.CreateServerConn(ctx, server.Address, uint16(server.Port), p.Profile, server.Secret, func(conn *proto.Conn, causeFunc context.CancelCauseFunc, profile *packet.GameProfile, secret string) proto.Handler {
		h := NewServerVelocityLoginHandler(conn, causeFunc, profile, secret)
		return &HoldingTest{
			ServerVelocityLoginHandler: h.(*ServerVelocityLoginHandler),
			doneCh:                     ch,
		}
	})
	if err != nil {
		return err
	}

	<-ctx.Done()
	cause := context.Cause(ctx)
	if errors.Is(cause, context.DeadlineExceeded) {
		// This would trigger from the timeout above. This is a complete fail.
		return fmt.Errorf("timed out connecting to remote server: %w", ctx.Err())
	} else if !errors.Is(cause, context.Canceled) {
		// Otherwise, we failed with the given cause. This should cause a disconnect failure.
		return fmt.Errorf("failed to connect to remote server: %w", cause)
	}

	// A clean cancel is the success case, so we don't need to do anything.
	// We are now connected to the remote.

	// Disconnect from the old server
	oldRemote := p.Conn.GetRemote()
	p.Conn.SetRemote(nil)
	oldRemote.SetRemote(nil)
	oldRemote.Close()

	// Swap to the config phase and connect to the new server
	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	p.SetState(packet.Play, &WaitForStartConfigHandler2{
		ctx:    ctx,
		cancel: cancel,
	})

	//h.Player.SetState(packet.Play, &WaitForStartConfigHandler{
	//	Player: h.Player,
	//	Remote: h.Remote,
	//	DoneCh: doneCh,
	//})
	if err = p.SendPacket(&packet.ServerStartConfiguration{}); err != nil {
		return err
	}
	<-ctx.Done()

	remote.SetRemote(p.Conn)
	p.Conn.SetRemote(remote)

	p.SetState(packet.Config, clientConfigHandlerFunc(p))
	remote.SetState(packet.Config, serverConfigHandlerFunc(p))

	ch <- true

	return nil
}
