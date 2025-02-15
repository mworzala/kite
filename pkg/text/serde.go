package text

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/Tnze/go-mc/nbt"
)

func MarshalJSON(c Component) ([]byte, error) {
	return json.Marshal(marshalTree(c))
}

func UnmarshalJSON(r io.Reader) (Component, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	var tree any
	if err = json.Unmarshal(data, &tree); err != nil {
		return nil, err
	}
	return unmarshalTree(tree.(map[string]any))
}

func MarshalNBT(w io.Writer, c Component, networkFormat bool) error {
	enc := nbt.NewEncoder(w)
	enc.NetworkFormat(networkFormat)
	return enc.Encode(marshalTree(c), "")
}

func UnmarshalNBT(r io.Reader, networkFormat bool) (Component, error) {
	dec := nbt.NewDecoder(r)
	dec.NetworkFormat(networkFormat)

	var tree any
	if _, err := dec.Decode(&tree); err != nil {
		return nil, err
	}
	return unmarshalTree(tree.(map[string]any))
}

func MarshalPlain(c Component) string {
	var b strings.Builder
	marshalPlain(c, &b)
	return b.String()
}

func marshalTree(c Component) map[string]any {
	result := make(map[string]any)
	switch c := c.(type) {
	case *Text:
		result["type"] = "text"
		result["text"] = c.Text
	case *Translate:
		result["type"] = "translate"
		result["translate"] = c.Translate
		if c.Fallback != "" {
			result["fallback"] = c.Fallback
		}
		with := marshalChildrenTree(c.With)
		if len(with) > 0 {
			result["with"] = c.With
		}
	case *Score:
		result["type"] = "score"
		score := make(map[string]any)
		score["name"] = c.Name
		score["objective"] = c.Objective
		result["score"] = score
	case *Selector:
		result["type"] = "selector"
		result["selector"] = c.Selector
		if c.Separator != nil {
			result["separator"] = marshalTree(c.Separator)
		}
	case *Keybind:
		result["type"] = "keybind"
		result["keybind"] = c.Keybind
	case *NBT:
		result["type"] = "nbt"
		if c.Source != "" {
			result["source"] = c.Source
		}
		if c.NBT != "" {
			result["nbt"] = c.NBT
		}
		if c.Interpret {
			result["interpret"] = true
		}
		if c.Separator != nil {
			result["separator"] = marshalTree(c.Separator)
		}
		if c.Block != "" {
			result["block"] = c.Block
		}
		if c.Entity != "" {
			result["entity"] = c.Entity
		}
		if c.Storage != "" {
			result["storage"] = c.Storage
		}
	}

	appendStyleTree(result, c.Style())
	children := marshalChildrenTree(c.Children())
	if len(children) > 0 {
		result["extra"] = children
	}

	return result
}

func marshalChildrenTree(cs []Component) []map[string]any {
	if len(cs) == 0 {
		return nil
	}

	var result []map[string]any
	for _, c := range cs {
		result = append(result, marshalTree(c))
	}
	return result
}

func appendStyleTree(result map[string]any, s Style) {
	if s == (Style{}) {
		return
	}

	result["color"] = s.Color.ToRGB().ToString()
	if s.Font != "" {
		result["font"] = s.Font
	}
	if s.Bold != Unset {
		result["bold"] = s.Bold == True
	}
	if s.Italic != Unset {
		result["italic"] = s.Italic == True
	}
	if s.Underlined != Unset {
		result["underlined"] = s.Underlined == True
	}
	if s.Strikethrough != Unset {
		result["strikethrough"] = s.Strikethrough == True
	}
	if s.Obfuscated != Unset {
		result["obfuscated"] = s.Obfuscated == True
	}
	if s.Insertion != "" {
		result["insertion"] = s.Insertion
	}
	if s.ClickEvent != nil {
		result["clickEvent"] = map[string]any{
			"action": s.ClickEvent.Action,
			"value":  s.ClickEvent.Value,
		}
	}
	if s.HoverEvent != nil {
		result["hoverEvent"] = map[string]any{
			"action": s.HoverEvent.Action,
			"value":  s.HoverEvent.Value,
		}
	}
}

func marshalPlain(c Component, b *strings.Builder) {
	switch c := c.(type) {
	case *Text:
		b.WriteString(c.Text)
	case *Translate:
		if c.Fallback != "" {
			b.WriteString(c.Fallback)
		} else {
			b.WriteString(c.Translate)
		}
	case *Selector:
		b.WriteString(c.Selector)
	case *Keybind:
		b.WriteString(c.Keybind)
	}

	for _, e := range c.Children() {
		marshalPlain(e, b)
	}
}

func unmarshalTree(m any) (c Component, err error) {
	if s, ok := m.(string); ok {
		return &Text{Text: s}, nil
	} else if l, ok := m.([]any); ok {
		children := make([]Component, len(l))
		for i, e := range l {
			var err error
			if children[i], err = unmarshalTree(e); err != nil {
				return nil, err
			}
		}
		return &Text{Extra: children}, nil
	}
	obj, ok := m.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid component: %v", m)
	}

	if typ, ok := obj["type"].(string); ok {
		switch typ {
		case "text":
			return unmarshalText(obj)
		case "translate":
			return unmarshalTranslate(obj)
		case "score":
			return unmarshalScore(obj)
		case "selector":
			return unmarshalSelector(obj)
		case "keybind":
			return unmarshalKeybind(obj)
		case "nbt":
			return unmarshalNBT(obj)
		}
	}

	// Still don't know the type, try to guess
	if _, ok := obj["text"]; ok {
		return unmarshalText(obj)
	} else if _, ok := obj["translate"]; ok {
		return unmarshalTranslate(obj)
	} else if _, ok := obj["score"]; ok {
		return unmarshalScore(obj)
	} else if _, ok := obj["selector"]; ok {
		return unmarshalSelector(obj)
	} else if _, ok := obj["keybind"]; ok {
		return unmarshalKeybind(obj)
	} else if _, ok := obj["nbt"]; ok {
		return unmarshalNBT(obj)
	}

	// If we _still_ dont know the type, fail.
	return nil, fmt.Errorf("invalid component: %v", m)
}

func unmarshalText(obj map[string]any) (c Component, err error) {
	t := &Text{}
	if text, ok := obj["text"].(string); ok {
		t.Text = text
	}
	if err = unmarshalStyle(obj, &t.S); err != nil {
		return nil, err
	}
	t.Extra, err = unmarshalExtra(obj)
	return t, err
}

func unmarshalTranslate(obj map[string]any) (c Component, err error) {
	t := &Translate{}
	if translate, ok := obj["translate"].(string); ok {
		t.Translate = translate
	}
	if fallback, ok := obj["fallback"].(string); ok {
		t.Fallback = fallback
	}
	if with, ok := obj["with"].([]any); ok {
		t.With = make([]Component, len(with))
		for i, e := range with {
			if t.With[i], err = unmarshalTree(e); err != nil {
				return nil, err
			}
		}
	}
	if err = unmarshalStyle(obj, &t.S); err != nil {
		return nil, err
	}
	t.Extra, err = unmarshalExtra(obj)
	return t, err
}

func unmarshalScore(obj map[string]any) (c Component, err error) {
	panic("todo")
}

func unmarshalSelector(obj map[string]any) (c Component, err error) {
	panic("todo")
}

func unmarshalKeybind(obj map[string]any) (c Component, err error) {
	panic("todo")
}

func unmarshalNBT(obj map[string]any) (c Component, err error) {
	panic("todo")
}

func unmarshalStyle(obj map[string]any, s *Style) (err error) {
	if color, ok := obj["color"].(string); ok {
		s.Color, err = colorFromHex(color)
		if err != nil {
			return err
		}
	}
	if font, ok := obj["font"].(string); ok {
		s.Font = font
	}
	if bold, ok := obj["bold"].(bool); ok {
		s.Bold = BoolToTristate(bold)
	}
	if italic, ok := obj["italic"].(bool); ok {
		s.Italic = BoolToTristate(italic)
	}
	if underlined, ok := obj["underlined"].(bool); ok {
		s.Underlined = BoolToTristate(underlined)
	}
	if strikethrough, ok := obj["strikethrough"].(bool); ok {
		s.Strikethrough = BoolToTristate(strikethrough)
	}
	if obfuscated, ok := obj["obfuscated"].(bool); ok {
		s.Obfuscated = BoolToTristate(obfuscated)
	}
	if insertion, ok := obj["insertion"].(string); ok {
		s.Insertion = insertion
	}
	if clickEvent, ok := obj["clickEvent"].(map[string]any); ok {
		s.ClickEvent = &ClickEvent{}
		if action, ok := clickEvent["action"].(string); ok {
			s.ClickEvent.Action = action
		}
		if value, ok := clickEvent["value"].(string); ok {
			s.ClickEvent.Value = value
		}
	}
	if hoverEvent, ok := obj["hoverEvent"].(map[string]any); ok {
		s.HoverEvent = &HoverEvent{}
		if action, ok := hoverEvent["action"].(string); ok {
			s.HoverEvent.Action = action
		}
		if value, ok := hoverEvent["value"].(any); ok {
			s.HoverEvent.Value = value
		}
	}
	return nil
}

func unmarshalExtra(obj map[string]any) ([]Component, error) {
	extra, ok := obj["extra"].([]any)
	if !ok || len(extra) == 0 {
		return nil, nil
	}

	children := make([]Component, len(extra))
	for i, e := range extra {
		var err error
		if children[i], err = unmarshalTree(e); err != nil {
			return nil, err
		}
	}

	return children, nil
}

func BoolToTristate(b bool) Tristate {
	if b {
		return True
	}
	return False
}
