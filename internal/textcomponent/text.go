package printer

import (
	"strings"

	"github.com/charmbracelet/x/ansi"
)

var ansiResetCode = ansi.SGR(ansi.ResetAttr)

// Text is for a regular string of text.
// This may or may not have color.
type Text struct {
	Color string
	Text  string
}

// Render implements [Component]
func (c *Text) Render(r *Renderer) {
	if c.Color == "" {
		r.WriteString(c.Text)
	} else {
		r.WriteString(c.Color)
		r.WriteString(c.Text)
		r.WriteString(ansiResetCode)
	}
}

// Simplify implements [Component]
//
// When [Text] is simplified, it is split up into individual lines.
func (c *Text) Simplify(r *Renderer) []BasicComponent {
	text := c.Text
	index := strings.Index(text, "\n")
	if index == -1 {
		return []BasicComponent{c}
	}

	// There is at least one newline character.
	var split []BasicComponent
	lines := strings.Split(c.Text, "\n")
	last := len(lines) - 1

	for i, text := range lines {
		if i > 0 {
			split = append(split, Newline)
		}

		if i == last && text == "" {
			break
		}

		split = append(split, &Text{
			Text:  text,
			Color: c.Color,
		})
	}

	return split
}

func (c *Text) isBasic() {}
