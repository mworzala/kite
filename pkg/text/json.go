package text

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

const (
	jsonStartObject = json.Delim('{')
	jsonEndObject   = json.Delim('}')
	jsonStartArray  = json.Delim('[')
	jsonEndArray    = json.Delim(']')
)

func MarshalJSON(c Component) ([]byte, error) {
	return json.Marshal(marshalComponent(c))
}

func UnmarshalJSON(data []byte) (Component, error) {
	d := json.NewDecoder(bytes.NewReader(data))

	tok, err := d.Token()
	if errors.Is(err, io.EOF) {
		return Text{}, nil
	}

	return unmarshalComponent(d, tok)
}

func marshalComponent(c Component) map[string]any {
	res := make(map[string]any)
	switch c := c.(type) {
	case Text:
		res["type"] = "text"
		res["text"] = c.Text
		marshalCommon(res, c.S, c.Extra)
	case Translate:
		res["type"] = "translate"
		res["translate"] = c.Translate
		if c.Fallback != "" {
			res["fallback"] = c.Fallback
		}
		if len(c.With) > 0 {
			value := make([]map[string]any, len(c.With))
			for i, c := range c.With {
				value[i] = marshalComponent(c)
			}
			res["with"] = value
		}
		marshalCommon(res, c.S, c.Extra)
	case Score:
		res["type"] = "score"
		//todo
		marshalCommon(res, c.S, c.Extra)
	case Selector:
		res["type"] = "selector"
		res["selector"] = c.Selector
		if c.Separator != nil {
			res["separator"] = marshalComponent(c.Separator)
		}
		marshalCommon(res, c.S, c.Extra)
	case Keybind:
		res["type"] = "keybind"
		res["keybind"] = c.Keybind
		marshalCommon(res, c.S, c.Extra)
	case NBT:
		res["type"] = "nbt"
		res["nbt"] = c.NBT
		if c.Source != "" {
			res["source"] = c.Source
		}
		if c.Interpret {
			res["interpret"] = true
		}
		if c.Block != "" {
			res["block"] = c.Block
		}
		if c.Entity != "" {
			res["entity"] = c.Entity
		}
		if c.Storage != "" {
			res["storage"] = c.Storage
		}
		marshalCommon(res, c.S, c.Extra)
	default:
		panic("unreachable")
	}
	return res
}

func marshalCommon(res map[string]any, s Style, extra []Component) {
	if s.Color != 0 {
		res["color"] = s.Color.ToRGB().ToString()
	}
	if s.Font != "" {
		res["font"] = s.Font
	}
	if s.Bold != Unset {
		res["bold"] = s.Bold == True
	}
	if s.Italic != Unset {
		res["italic"] = s.Italic == True
	}
	if s.Underlined != Unset {
		res["underlined"] = s.Underlined == True
	}
	if s.Strikethrough != Unset {
		res["strikethrough"] = s.Strikethrough == True
	}
	if s.Obfuscated != Unset {
		res["obfuscated"] = s.Obfuscated == True
	}
	if len(extra) > 0 {
		value := make([]map[string]any, len(extra))
		for i, c := range extra {
			value[i] = marshalComponent(c)
		}
		res["extra"] = extra
	}
}

func unmarshalComponent(d *json.Decoder, tok json.Token) (Component, error) {
	switch tok := tok.(type) {
	case json.Delim:
		if tok == jsonStartObject {
			return unmarshalComponentInner(d)
		} else if tok == jsonStartArray {
			extra, err := unmarshalComponentList(d)
			return Text{Extra: extra}, err
		} else {
			return nil, fmt.Errorf("unexpected delimiter %v", tok)
		}
	case string:
		// Single string component
		return Text{Text: tok}, nil
	default:
		return nil, fmt.Errorf("unexpected literal %T as component", tok)
	}
}

func unmarshalComponentList(d *json.Decoder) ([]Component, error) {
	var extra []Component
	for d.More() {
		next, err := d.Token()
		if err != nil {
			return nil, err
		} else if next == jsonEndArray {
			break
		}

		comp, err := unmarshalComponent(d, next)
		if err != nil {
			return nil, err
		}
		extra = append(extra, comp)
	}
	return extra, nil
}

func unmarshalComponentInner(d *json.Decoder) (Component, error) {
	var knownType string
	var text *string                  // Text
	var translate *string             // Translatable
	var fallback *string              // Translatable
	var with []Component              // Translatable
	var score *string                 // Score todo
	var selector *string              // Selector
	var separator Component           // Selector
	var keybind *string               // Keybind
	var nbt *string                   // NBT
	var source *string                // NBT
	var interpret bool                // NBT
	var block, entity, storage string // NBT
	var extra []Component
	var style Style

	// TODO: interactivity

	// At this point we are known to be in an object.
	for d.More() {
		tok, err := d.Token()
		if err != nil {
			return nil, err
		}
		if tok == jsonEndObject {
			break
		}

		key := tok.(string)

		tok, err = d.Token()
		if err != nil {
			return nil, err
		}
		switch key {
		case "type":
			knownType, err = assertString2(key, tok)
		case "text":
			text, err = assertString(key, tok)
		case "translate":
			translate, err = assertString(key, tok)
		case "fallback":
			fallback, err = assertString(key, tok)
		case "with":
			with, err = unmarshalComponentList(d)
		case "score":
			score = new(string)
			*score = tok.(string) //todo
		case "selector":
			selector, err = assertString(key, tok)
		case "separator":
			separator, err = unmarshalComponent(d, tok)
		case "keybind":
			keybind, err = assertString(key, tok)
		case "nbt":
			nbt, err = assertString(key, tok)
		case "block":
			block, err = assertString2(key, tok)
		case "entity":
			entity, err = assertString2(key, tok)
		case "storage":
			storage, err = assertString2(key, tok)
		case "interpret":
			interpret = tok.(bool)
		case "extra":
			extra, err = unmarshalComponentList(d)
		case "color":
			var hex *string
			hex, err = assertString(key, tok)
			if err != nil {
				return nil, err
			}
			style.Color, err = colorFromHex(*hex)
		case "font":
			style.Font, err = assertString2(key, tok)
		case "bold":
			style.Bold, err = assertTristate(key, tok)
		case "italic":
			style.Italic, err = assertTristate(key, tok)
		case "underlined":
			style.Underlined, err = assertTristate(key, tok)
		case "strikethrough":
			style.Strikethrough, err = assertTristate(key, tok)
		case "obfuscated":
			style.Obfuscated, err = assertTristate(key, tok)
		default:
			err = skipValue(d, tok)
		}
		if err != nil {
			return nil, err
		}

	}

	// Now we need to construct the component
	switch knownType {
	case "text":
		return Text{Text: or(text, ""), S: style, Extra: extra}, nil
	case "translate":
		return Translate{Translate: or(translate, ""), S: style, Extra: extra}, nil
	case "score":
		return Score{S: style, Extra: extra}, nil
	case "selector":
		return Selector{Selector: or(selector, ""), Separator: separator, S: style, Extra: extra}, nil
	case "keybind":
		return Keybind{Keybind: or(keybind, ""), S: style, Extra: extra}, nil
	case "nbt":
		return NBT{Source: or(source, ""), NBT: or(nbt, ""), Interpret: interpret, Block: block, Entity: entity, Storage: storage, S: style, Extra: extra}, nil
	}

	// Guess based on fields
	if text != nil {
		return Text{Text: *text, S: style, Extra: extra}, nil
	} else if translate != nil {
		return Translate{Translate: *translate, Fallback: or(fallback, ""), With: with, S: style, Extra: extra}, nil
	} else if score != nil {
		return Score{S: style, Extra: extra}, nil
	} else if selector != nil {
		return Selector{Selector: *selector, Separator: separator, S: style, Extra: extra}, nil
	} else if keybind != nil {
		return Keybind{Keybind: *keybind, S: style, Extra: extra}, nil
	} else if nbt != nil {
		return NBT{Source: or(source, ""), NBT: *nbt, Interpret: interpret, Block: block, Entity: entity, Storage: storage, S: style, Extra: extra}, nil
	}

	return Text{Text: or(text, ""), S: style, Extra: extra}, nil
}

func assertString(key string, tok json.Token) (*string, error) {
	s, ok := tok.(string)
	if !ok {
		return nil, fmt.Errorf("expected string for %s, got %T", key, tok)
	}
	return &s, nil
}

func assertString2(key string, tok json.Token) (string, error) {
	s, ok := tok.(string)
	if !ok {
		return "", fmt.Errorf("expected string for %s, got %T", key, tok)
	}
	return s, nil
}

func assertTristate(key string, tok json.Token) (Tristate, error) {
	b, ok := tok.(bool)
	if !ok {
		return Unset, fmt.Errorf("expected bool for %s, got %T", key, tok)
	}
	if b {
		return True, nil
	} else {
		return False, nil
	}
}

func skipValue(d *json.Decoder, tok json.Token) error {
	switch tok {
	case jsonStartArray, jsonStartObject:
		for d.More() {
			tok, err := d.Token()
			if err != nil {
				return err
			}
			if tok == jsonEndArray || tok == jsonEndObject {
				break
			}
			if err = skipValue(d, tok); err != nil {
				return err
			}
		}
		return nil
	default: // Anything else is a value type so no need to read more
		return nil
	}
}

func or[T any](a *T, b T) T {
	if a != nil {
		return *a
	}
	return b
}
