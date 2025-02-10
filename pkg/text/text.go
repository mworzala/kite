package text

type Text struct {
	Text string `json:"text" nbt:"text"`
	//	S    Style
	//Extra []Component
}

func (t Text) component() {}
