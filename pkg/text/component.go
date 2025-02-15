package text

type Component interface {
	component() // Marker method
}

type Text struct {
	Text  string
	S     Style
	Extra []Component
}

type Translate struct {
	Translate string
	Fallback  string
	With      []Component
	S         Style
	Extra     []Component
}

type Score struct {
	//todo
	S     Style
	Extra []Component
}

type Selector struct {
	Selector  string
	Separator Component
	S         Style
	Extra     []Component
}

type Keybind struct {
	Keybind string
	S       Style
	Extra   []Component
}

type NBT struct {
	Source    string
	NBT       string
	Interpret bool
	Separator bool
	Block     string
	Entity    string
	Storage   string
	S         Style
	Extra     []Component
}

func (c Text) component()      {}
func (c Translate) component() {}
func (c Score) component()     {}
func (c Selector) component()  {}
func (c Keybind) component()   {}
func (c NBT) component()       {}
