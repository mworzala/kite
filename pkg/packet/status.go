package packet

import (
	"encoding/json"
	"io"

	"github.com/mworzala/kite/pkg/buffer"
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
	p.Payload, err = buffer.Long.Read(r)
	return
}
func (p *ClientStatusPingRequest) Write(w io.Writer) (err error) {
	return buffer.Long.Write(w, p.Payload)
}

// Server

const (
	ServerStatusStatusResponseID = iota
	ServerStatusPingResponseID
)

type (
	ServerStatusResponse struct {
		Payload StatusResponse
	}
	StatusResponse struct {
		Version           ServerVersion    `json:"version"`
		Players           ServerPlayerList `json:"players"`
		Description       json.RawMessage  `json:"description"` //todo this is a text component
		Favicon           string           `json:"favicon"`
		EnforceSecureChat bool             `json:"enforcesSecureChat"`
	}
	ServerVersion struct {
		Name     string `json:"name"`
		Protocol int    `json:"protocol"`
	}
	ServerPlayerList struct {
		Max    int             `json:"max"`
		Online int             `json:"online"`
		Sample []*ServerPlayer `json:"sample"`
	}
	ServerPlayer struct {
		Name string `json:"name"`
		ID   string `json:"id"`
	}
)

func (p *ServerStatusResponse) Direction() Direction { return Clientbound }
func (p *ServerStatusResponse) ID(state State) int {
	return stateId1(state, Status, ServerStatusStatusResponseID)
}
func (p *ServerStatusResponse) Read(r io.Reader) (err error) {
	p.Payload, err = buffer.JSON[StatusResponse]().Read(r)
	return
}
func (p *ServerStatusResponse) Write(w io.Writer) (err error) {
	return buffer.JSON[StatusResponse]().Write(w, p.Payload)
}

type ServerStatusPingResponse struct {
	Payload int64
}

func (p *ServerStatusPingResponse) Direction() Direction { return Clientbound }
func (p *ServerStatusPingResponse) ID(state State) int {
	return stateId1(state, Status, ServerStatusPingResponseID)
}
func (p *ServerStatusPingResponse) Read(r io.Reader) (err error) {
	p.Payload, err = buffer.Long.Read(r)
	return
}
func (p *ServerStatusPingResponse) Write(w io.Writer) (err error) {
	return buffer.Long.Write(w, p.Payload)
}

var (
	_ Packet = (*ClientStatusRequest)(nil)
	_ Packet = (*ClientStatusPingRequest)(nil)

	_ Packet = (*ServerStatusResponse)(nil)
	_ Packet = (*ServerStatusPingResponse)(nil)
)
