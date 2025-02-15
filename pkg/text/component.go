package text

type Component interface {
	component() // Marker method

	Style() Style
	Children() []Component
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
	Name      string
	Objective string
	S         Style
	Extra     []Component
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
	Separator Component
	Block     string
	Entity    string
	Storage   string
	S         Style
	Extra     []Component
}

func (c *Text) component()      {}
func (c *Translate) component() {}
func (c *Score) component()     {}
func (c *Selector) component()  {}
func (c *Keybind) component()   {}
func (c *NBT) component()       {}

func (c *Text) Style() Style      { return c.S }
func (c *Translate) Style() Style { return c.S }
func (c *Score) Style() Style     { return c.S }
func (c *Selector) Style() Style  { return c.S }
func (c *Keybind) Style() Style   { return c.S }
func (c *NBT) Style() Style       { return c.S }

func (c *Text) Children() []Component      { return c.Extra }
func (c *Translate) Children() []Component { return c.Extra }
func (c *Score) Children() []Component     { return c.Extra }
func (c *Selector) Children() []Component  { return c.Extra }
func (c *Keybind) Children() []Component   { return c.Extra }
func (c *NBT) Children() []Component       { return c.Extra }
