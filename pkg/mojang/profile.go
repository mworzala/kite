package mojang

import (
	"github.com/google/uuid"
	"github.com/mworzala/kite/pkg/buffer"
	"io"
)

type GameProfile struct {
	ID             uuid.UUID         `json:"id"`
	Name           string            `json:"name"`
	Properties     []ProfileProperty `json:"properties"`
	ProfileActions []any             `json:"profileActions"` //todo what actually goes in here and when?
}

type ProfileProperty struct {
	Name      string `json:"name"`
	Value     string `json:"value"`
	Signature string `json:"signature,omitempty"`
}

func (gp *GameProfile) Read(r io.Reader) (err error) {
	gp.ID, gp.Name, err = buffer.Read2(r, buffer.UUID, buffer.String)
	if err != nil {
		return err
	}
	gp.Properties, err = buffer.ReadList(r, func() (pp ProfileProperty, err error) {
		pp.Name, pp.Value, pp.Signature, err = buffer.Read3(r,
			buffer.String, buffer.String, buffer.Opt(buffer.String))
		return pp, err
	})
	if err != nil {
		return err
	}
	return nil
}

func (gp *GameProfile) Write(w io.Writer) (err error) {
	err = buffer.Write2(w, buffer.UUID, gp.ID, buffer.String, gp.Name)
	if err != nil {
		return err
	}
	err = buffer.WriteList(w, gp.Properties, func(pp ProfileProperty) error {
		return buffer.Write3(w, buffer.String, pp.Name, buffer.String,
			pp.Value, buffer.Opt(buffer.String), pp.Signature)
	})
	return nil
}
