package text

type Style struct {
	Color         Color
	ShadowColor   Color
	Font          string
	Bold          Tristate
	Italic        Tristate
	Underlined    Tristate
	Strikethrough Tristate
	Obfuscated    Tristate
}

type Tristate int8

const (
	Unset Tristate = iota
	True
	False
)
