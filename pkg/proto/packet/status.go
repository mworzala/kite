package packet

import (
	"io"

	"github.com/mworzala/kite/pkg/proto/binary"
)

const (
	ClientStatusStatusRequestID = iota
	ClientStatusPingRequestID
)

type ClientStatusRequest struct{}

func (p *ClientStatusRequest) Direction() Direction { return Serverbound }
func (p *ClientStatusRequest) ID(state State) int {
	return stateId1(state, Status, ClientStatusStatusRequestID)
}

func (p *ClientStatusRequest) Read(r io.Reader) error  { return nil }
func (p *ClientStatusRequest) Write(w io.Writer) error { return nil }

type ClientStatusPingRequest struct {
	Payload int64
}

func (p *ClientStatusPingRequest) Direction() Direction { return Serverbound }
func (p *ClientStatusPingRequest) ID(state State) int {
	return stateId1(state, Status, ClientStatusPingRequestID)
}

func (p *ClientStatusPingRequest) Read(r io.Reader) (err error) {
	if p.Payload, err = binary.ReadLong(r); err != nil {
		return
	}
	return nil
}

func (p *ClientStatusPingRequest) Write(w io.Writer) (err error) {
	if err = binary.WriteLong(w, p.Payload); err != nil {
		return
	}
	return nil
}

// Server

const (
	ServerStatusStatusResponseID = iota
	ServerStatusPingResponseID
)

type ServerStatusResponse struct {
	Payload StatusResponse
}

func (p *ServerStatusResponse) Direction() Direction { return Clientbound }
func (p *ServerStatusResponse) ID(state State) int {
	return stateId1(state, Status, ServerStatusStatusResponseID)
}

func (p *ServerStatusResponse) Read(r io.Reader) (err error) {
	if p.Payload, err = binary.ReadTypedJSON[StatusResponse](r); err != nil {
		return
	}
	return
}

func (p *ServerStatusResponse) Write(w io.Writer) (err error) {
	if err = binary.WriteTypedJSON(w, p.Payload); err != nil {
		return
	}
	return nil
}

type ServerStatusPingResponse struct {
	Payload int64
}

func (p *ServerStatusPingResponse) Direction() Direction { return Clientbound }
func (p *ServerStatusPingResponse) ID(state State) int {
	return stateId1(state, Status, ServerStatusPingResponseID)
}

func (p *ServerStatusPingResponse) Read(r io.Reader) (err error) {
	if p.Payload, err = binary.ReadLong(r); err != nil {
		return
	}
	return nil
}

func (p *ServerStatusPingResponse) Write(w io.Writer) (err error) {
	if err = binary.WriteLong(w, p.Payload); err != nil {
		return
	}
	return nil
}

var (
	_ Packet = (*ClientStatusRequest)(nil)
	_ Packet = (*ClientStatusPingRequest)(nil)

	_ Packet = (*ServerStatusResponse)(nil)
	_ Packet = (*ServerStatusPingResponse)(nil)
)
