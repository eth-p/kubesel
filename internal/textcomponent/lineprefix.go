package printer

// LinePrefix adds a prefix to each line.
type LinePrefix struct {
	Prefix Component
	Child  Component
}

// Simplify implements [Component]
func (c *LinePrefix) Simplify(r *Renderer) []BasicComponent {
	if c.Child == nil {
		return []BasicComponent{}
	}

	if c.Prefix == nil {
		return r.Simplify(c.Child)
	}

	// Simplify the prefix.
	prefix := r.Simplify(c.Prefix)
	if len(prefix) == 0 {
		return r.Simplify(c.Child)
	}

	// Simplify the children.
	children := r.Simplify(c.Child)
	if len(children) == 0 {
		return children
	}

	// If we have children, prefix each new line.
	components := make([]BasicComponent, 0, len(children))
	components = append(components, prefix...)
	for _, child := range children {
		components = append(components, child)
		if child == Newline {
			components = append(components, prefix...)
		}
	}

	return components
}
