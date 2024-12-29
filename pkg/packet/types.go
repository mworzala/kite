package packet

import (
	"io"

	"github.com/mworzala/kite/pkg/proto/binary"
)

type GameProfile struct {
	UUID       string
	Username   string
	Properties []ProfileProperty
}

func (p *GameProfile) Read(r io.Reader) (err error) {
	if p.UUID, err = binary.ReadUUID(r); err != nil {
		return
	}
	if p.Username, err = binary.ReadString(r); err != nil {
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
	if err = binary.WriteString(w, p.Username); err != nil {
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

type ResourcePackStatus int

const (
	ResourcePackSuccessfullyLoaded ResourcePackStatus = iota
	ResourcePackDeclined
	ResourcePackFailedDownload
	ResourcePackAccepted
	ResourcePackDownloaded
	ResourcePackInvalidURL
	ResourcePackFailedReload
	ResourcePackDiscarded
)

func (s ResourcePackStatus) Validate() bool {
	return s >= ResourcePackSuccessfullyLoaded && s <= ResourcePackDiscarded
}

func (s ResourcePackStatus) String() string {
	switch s {
	case ResourcePackSuccessfullyLoaded:
		return "successfully_loaded"
	case ResourcePackDeclined:
		return "declined"
	case ResourcePackFailedDownload:
		return "failed_download"
	case ResourcePackAccepted:
		return "accepted"
	case ResourcePackDownloaded:
		return "downloaded"
	case ResourcePackInvalidURL:
		return "invalid_url"
	case ResourcePackFailedReload:
		return "failed_reload"
	case ResourcePackDiscarded:
		return "discarded"
	}
	return "unknown"
}
