package text

type Style struct {
	Color         Color
	Font          string
	Bold          Tristate
	Italic        Tristate
	Underlined    Tristate
	Strikethrough Tristate
	Obfuscated    Tristate
}

type Tristate uint8

const (
	Unset Tristate = iota
	True
	False
)
