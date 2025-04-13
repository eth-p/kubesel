package printer

// newlineType is a text [Component] that inserts a new line.
type newlineType struct{}

// Newline is a text [Component] that inserts a new line.
var Newline = &newlineType{}

func (c *newlineType) Render(r *Renderer) {
	r.WriteString("\n")
}

func (c *newlineType) Simplify(r *Renderer) []BasicComponent {
	return []BasicComponent{c}
}

func (c *newlineType) isBasic() {}

// Trim is a text [Component] that can trim leading or trailing newlines.
type Trim struct {
	Leading  bool
	Trailing bool
	Child    Component
}

func (c *Trim) Simplify(r *Renderer) []BasicComponent {
	children := r.Simplify(c.Child)

	// Trim leading newlines.
	if c.Leading {
		firstNonNewline := 0
		anyNonNewline := false
		for i, child := range children {
			if child != Newline {
				firstNonNewline = i
				anyNonNewline = true
				break
			}
		}

		if !anyNonNewline {
			return []BasicComponent{}
		}

		children = children[firstNonNewline:]
	}

	// Trim trailing newlines.
	if c.Trailing {
		lastNonNewline := 0
		for i := len(children) - 1; i >= 0; i-- {
			child := children[i]
			if child != Newline {
				lastNonNewline = i + 1
				break
			}
		}

		children = children[0:lastNonNewline]
	}

	return children
}
