package text

type Style struct {
	Color Color
	Font  string

	Bold          Tristate
	Italic        Tristate
	Underlined    Tristate
	Strikethrough Tristate
	Obfuscated    Tristate

	Insertion  string
	ClickEvent *ClickEvent
	HoverEvent *HoverEvent
}

type Tristate uint8

const (
	Unset Tristate = iota
	True
	False
)

type ClickEvent struct {
	Action string
	Value  string
}

type HoverEvent struct {
	Action   string
	Contents any
	Value    any
}
