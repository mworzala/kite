package component

type Text struct {
	Content  string
	S        Style
	Children []Component
}

func (t *Text) component() {}
