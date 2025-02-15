package text

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	insert = "insert me"
	txt    = &Text{
		Text: "Hello",
		Extra: []Component{
			&Text{Text: " there!", S: Style{
				Color:      Red,
				Italic:     True,
				Obfuscated: False,
			}},
		},
		S: Style{
			Obfuscated:    True,
			Bold:          False,
			Strikethrough: Unset,
			Underlined:    True,
			Italic:        False,
			Font:          "minecraft:default",
			Color:         Aqua,
			//ClickEvent:    SuggestCommand("/help"),
			//HoverEvent: ShowText(&Text{
			//	Content: " world",
			//	Extra: []Component{
			//		&Text{Content: "!"},
			//	},
			//}),
			//Insertion: insert,
		}}
	jsonTxt = `{"bold":false,"color":"#55ffff","extra":[{"color":"#ff5555","italic":true,"obfuscated":false,"text":" there!","type":"text"}],"font":"minecraft:default","italic":false,"obfuscated":true,"text":"Hello","type":"text","underlined":true}`
	//jsonTxt = `{"bold":false,"clickEvent":{"action":"suggest_command","value":"/help"},"color":"#55ffff","extra":[{"color":"#ff5555","italic":true,"obfuscated":false,"text":" there!"}],"font":"minecraft:default","hoverEvent":{"action":"show_text","contents":{"extra":[{"text":"!"}],"text":" world"}},"insertion":"insert me","italic":false,"obfuscated":true,"text":"Hello","underlined":true}`
)

func TestJson_Marshal_text(t *testing.T) {
	s, err := MarshalJSON(txt)
	require.NoError(t, err)
	require.Equal(t, jsonTxt, string(s))
}

func TestJson_Unmarshal_text(t *testing.T) {
	c, err := UnmarshalJSON(bytes.NewReader([]byte(jsonTxt)))
	require.NoError(t, err)
	require.Equal(t, txt, c)
}

func TestJson_Marshal_Regression(t *testing.T) {
	c := &Text{Text: "Hello", S: Style{Color: Red}}
	s, err := MarshalJSON(c)
	require.NoError(t, err)
	require.Equal(t, `{"color":"#ff5555","text":"Hello","type":"text"}`, string(s))
}
