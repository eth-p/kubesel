package printer

import (
	"strings"
)

// Component is an abstract representation of some unit of printable text.
type Component interface {

	// Simplify returns a flattened and and simplified the component in an inside-out order.
	// This is similar to how a AST optimizer would do it.
	Simplify(*Renderer) []BasicComponent
}

// BasicComponent is either a [newlineType] or a [Text].
type BasicComponent interface {
	Component
	isBasic()

	// Render renders the component to string buffer.
	Render(*Renderer)
}

// Renderer is responsible for rendering [Component] into strings.
// To prevent unnecessary computation, this is caches simplification steps.
type Renderer struct {
	simplifyCache map[Component][]BasicComponent
	output        strings.Builder
}

func NewRenderer() *Renderer {
	return &Renderer{
		simplifyCache: make(map[Component][]BasicComponent),
	}
}

func (r *Renderer) Render(c Component) {
	if bc, ok := c.(BasicComponent); ok {
		bc.Render(r)
	}

	for _, component := range c.Simplify(r) {
		component.Render(r)
	}
}

func (r *Renderer) String() string {
	return r.output.String()
}

// WriteString should be used by [Component]s to render their contents.
func (r *Renderer) WriteString(s string) {
	r.output.WriteString(s)
}

// Simplify should be used by [Component]s to render their contents.
func (r *Renderer) Simplify(c Component) []BasicComponent {
	cached, ok := r.simplifyCache[c]
	if !ok {
		cached = c.Simplify(r)
		r.simplifyCache[c] = cached
	}

	return cached
}
