package packet

import (
	"io"

	"github.com/mworzala/kite/pkg/proto/binary"
)

// An Intent is the target state when joining a server.
type Intent int

const (
	IntentStatus = iota + 1
	IntentLogin
	IntentTransfer
)

func (s Intent) Validate() bool {
	return s >= IntentStatus && s <= IntentTransfer
}

func (s Intent) String() string {
	switch s {
	case IntentStatus:
		return "status"
	case IntentLogin:
		return "login"
	case IntentTransfer:
		return "transfer"
	}
	return "unknown"
}

type GameProfile struct {
	UUID       string
	Username   string
	Properties []ProfileProperty
}

func (p *GameProfile) Read(r io.Reader) (err error) {
	if p.UUID, err = binary.ReadUUID(r); err != nil {
		return
	}
	if p.Username, err = binary.ReadSizedString(r, 16); err != nil {
		return
	}
	if p.Properties, err = binary.ReadCollection(r, (*ProfileProperty).Read); err != nil {
		return
	}
	return nil
}

func (p *GameProfile) Write(w io.Writer) (err error) {
	if err = binary.WriteUUID(w, p.UUID); err != nil {
		return
	}
	if err = binary.WriteSizedString(w, p.Username, 16); err != nil {
		return
	}
	if err = binary.WriteCollection(w, p.Properties, (*ProfileProperty).Write); err != nil {
		return
	}
	return nil
}

type ProfileProperty struct {
	Name      string
	Value     string
	Signature *string
}

func (p *ProfileProperty) Read(r io.Reader) (err error) {
	if p.Name, err = binary.ReadString(r); err != nil {
		return
	}
	if p.Value, err = binary.ReadString(r); err != nil {
		return
	}
	if p.Signature, err = binary.ReadOptionalFunc(r, binary.ReadString); err != nil {
		return
	}
	return
}

func (p *ProfileProperty) Write(w io.Writer) (err error) {
	if err = binary.WriteString(w, p.Name); err != nil {
		return
	}
	if err = binary.WriteString(w, p.Value); err != nil {
		return
	}
	if err = binary.WriteOptionalFunc(w, p.Signature, binary.WriteString); err != nil {
		return
	}
	return
}
