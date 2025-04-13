package printer

// Sequence contains zero or more child [Component]s in an ordered sequence.
type Sequence struct {
	Children []Component
}

// Append adds a new child [Component] to the end of the sequence.
func (c *Sequence) Append(children ...Component) {
	c.Children = append(c.Children, children...)
}

// Simplify implements [Component]
func (c *Sequence) Simplify(r *Renderer) []BasicComponent {
	result := make([]BasicComponent, 0, len(c.Children))
	for _, child := range c.Children {
		result = append(result, r.Simplify(child)...)
	}
	return result
}
