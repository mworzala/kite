package text

import "fmt"

// Color is a color in ARGB format.
//
// It may be used as both an RGB and RGBA color by ignoring the alpha bits.
type Color uint32

const (
	Black       Color = 0x000000
	DarkBlue    Color = 0x0000AA
	DarkGreen   Color = 0x00AA00
	DarkAqua    Color = 0x00AAAA
	DarkRed     Color = 0xAA0000
	DarkPurple  Color = 0xAA00AA
	Gold        Color = 0xFFAA00
	Gray        Color = 0xAAAAAA
	DarkGray    Color = 0x555555
	Blue        Color = 0x5555FF
	Green       Color = 0x55FF55
	Aqua        Color = 0x55FFFF
	Red         Color = 0xFF5555
	LightPurple Color = 0xFF55FF
	Yellow      Color = 0xFFFF55
	White       Color = 0xFFFFFF
)

func (c Color) RGB() (r, g, b byte) {
	r = byte(c >> 16 & 0xFF)
	g = byte(c >> 8 & 0xFF)
	b = byte(c & 0xFF)
	return
}

// RGBA implements the go color.Color interface
func (c Color) RGBA() (r, g, b, a uint32) {
	r = uint32(c >> 16 & 0xFF)
	g = uint32(c >> 8 & 0xFF)
	b = uint32(c & 0xFF)
	a = uint32(c >> 24 & 0xFF)
	return
}

// ToRGB returns the RGB representation of the color.
//
// This will always have an alpha of 1
func (c Color) ToRGB() Color {
	return c | 0xFF000000
}

func (c Color) ToString() string {
	if n := c.Name(); n != "" {
		return n
	}

	// If we have an alpha value of 0xFF, we can use the RGB representation
	if c&0xFF000000 == 0xFF000000 {
		return fmt.Sprintf("#%06x", c&0xFFFFFF)
	} else {
		return fmt.Sprintf("#%08x", c)
	}
}

// Name returns the named color, or an empty string if the color is not named.
func (c Color) Name() string {
	switch c {
	case Black:
		return "black"
	case DarkBlue:
		return "dark_blue"
	case DarkGreen:
		return "dark_green"
	case DarkAqua:
		return "dark_aqua"
	case DarkRed:
		return "dark_red"
	case DarkPurple:
		return "dark_purple"
	case Gold:
		return "gold"
	case Gray:
		return "gray"
	case DarkGray:
		return "dark_gray"
	case Blue:
		return "blue"
	case Green:
		return "green"
	case Aqua:
		return "aqua"
	case Red:
		return "red"
	case LightPurple:
		return "light_purple"
	case Yellow:
		return "yellow"
	case White:
		return "white"
	default:
		return ""
	}
}

func colorFromHex(s string) (Color, error) {
	if len(s) > 0 && s[0] == '#' {
		s = s[1:]
	}
	if len(s) != 6 && len(s) != 8 {
		return 0, fmt.Errorf("invalid color string: %s", s)
	}

	var c Color
	var err error
	if len(s) == 6 {
		c, err = colorFromHex6(s)
	} else {
		c, err = colorFromHex8(s)
	}

	if err != nil {
		return 0, err
	}

	return c, nil
}

func colorFromHex6(s string) (Color, error) {
	var c Color
	_, err := fmt.Sscanf(s, "%06x", &c)
	if err != nil {
		return 0, err
	}

	return c, nil
}

func colorFromHex8(s string) (Color, error) {
	var c Color
	_, err := fmt.Sscanf(s, "%08x", &c)
	if err != nil {
		return 0, err
	}

	return c, nil
}
