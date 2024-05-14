package kite

import (
	"context"
	"fmt"

	"github.com/mworzala/kite/internal/pkg/velocity"
	"github.com/mworzala/kite/pkg/proto"
	"github.com/mworzala/kite/pkg/proto/packet"
)

type HoldingTest struct {
	*ServerVelocityLoginHandler
	doneCh chan bool
}

func (h *HoldingTest) HandlePacket(pp proto.Packet) (err error) {
	switch pp.Id {
	case packet.ServerLoginLoginSuccessID:
		p := new(packet.ServerLoginSuccess)
		if err = pp.Read(p); err != nil {
			return err
		}
		if err = h.ServerVelocityLoginHandler.handleLoginSuccess(p); err != nil {
			return err
		}
		<-h.doneCh
		return nil
	default:
		return h.ServerVelocityLoginHandler.HandlePacket(pp)
	}
}

type ServerVelocityLoginHandler struct {
	remote   *proto.Conn
	complete context.CancelCauseFunc // nil is success, otherwise fail reason

	profile *packet.GameProfile // The profile of the connecting player
	secret  string
}

func NewServerVelocityLoginHandler(remote *proto.Conn, complete context.CancelCauseFunc, profile *packet.GameProfile, secret string) proto.Handler {
	return &ServerVelocityLoginHandler{remote, complete, profile, secret}
}

func (h *ServerVelocityLoginHandler) HandlePacket(pp proto.Packet) (err error) {
	//todo we should handle encryption request here to create an error that the backend server is in online-mode
	switch pp.Id {
	case packet.ServerLoginDisconnectID:
		p := new(packet.ServerLoginDisconnect)
		if err = pp.Read(p); err != nil {
			return err
		}
		return h.handleDisconnect(p)
	case packet.ServerLoginPluginRequestID:
		p := new(packet.ServerLoginPluginRequest)
		if err = pp.Read(p); err != nil {
			return err
		}
		return h.handlePluginRequest(p)
	case packet.ServerLoginLoginSuccessID:
		p := new(packet.ServerLoginSuccess)
		if err = pp.Read(p); err != nil {
			return err
		}
		return h.handleLoginSuccess(p)
	default:
		return proto.UnknownPacket
	}
}

func (h *ServerVelocityLoginHandler) handleDisconnect(p *packet.ServerLoginDisconnect) error {
	h.complete(fmt.Errorf("disconnect: %s", p.Reason))
	return nil
}

func (h *ServerVelocityLoginHandler) handlePluginRequest(p *packet.ServerLoginPluginRequest) error {
	if p.Channel != "velocity:player_info" {
		return h.remote.SendPacket(&packet.ClientLoginPluginResponse{
			MessageID: p.MessageID,
			Data:      nil, // Unhandled message
		})
	}

	requestVersion := velocity.DefaultForwardingVersion
	if len(p.Data) > 0 {
		requestVersion = int(p.Data[0])
	}
	forward, err := velocity.CreateSignedForwardingData([]byte(h.secret), h.profile, requestVersion)
	if err != nil {
		return err
	}

	return h.remote.SendPacket(&packet.ClientLoginPluginResponse{
		MessageID: p.MessageID,
		Data:      forward,
	})
}

func (h *ServerVelocityLoginHandler) handleLoginSuccess(p *packet.ServerLoginSuccess) error {
	err := h.remote.SendPacket(&packet.ClientLoginAcknowledged{})
	if err != nil {
		return err
	}

	// Yay! We are connected to the remote server
	h.complete(nil)
	return nil
}

var _ proto.Handler = (*ServerVelocityLoginHandler)(nil)
