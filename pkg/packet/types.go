package packet

import (
	"io"

	"github.com/google/uuid"
	"github.com/mworzala/kite/pkg/buffer"
)

type GameProfile struct {
	UUID       uuid.UUID
	Username   string
	Properties []ProfileProperty
}

func (p *GameProfile) Read(r io.Reader) (err error) {
	p.UUID, p.Username, p.Properties, err = buffer.Read3(r,
		buffer.UUID, buffer.String, buffer.List(buffer.Struct[ProfileProperty]()))
	return nil
}

func (p *GameProfile) Write(w io.Writer) (err error) {
	return buffer.Write3(w, buffer.UUID, p.UUID, buffer.String, p.Username,
		buffer.List(buffer.Struct[ProfileProperty]()), p.Properties)
}

type ProfileProperty struct {
	Name      string
	Value     string
	Signature string // Optional
}

func (p ProfileProperty) Read(r io.Reader) (err error) {
	p.Name, p.Value, p.Signature, err = buffer.Read3(r,
		buffer.String, buffer.String, buffer.Opt(buffer.String))
	return
}

func (p ProfileProperty) Write(w io.Writer) (err error) {
	return buffer.Write3(w, buffer.String, p.Name,
		buffer.String, p.Value, buffer.Opt(buffer.String), p.Signature)
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
