package diagram

import (
	"fmt"
	"strings"
	"sync"

	"github.com/go-faster/errors"
)

var builderPool = sync.Pool{
	New: func() any { return new(strings.Builder) },
}

// Diagram is a diagram.
type Diagram struct {
	nodes       []*Node
	connections []*Connection
	globals     map[string]string
	classes     map[string]map[string]string

	seenNodes map[string]struct{}
}

// DiagramOption is a function that applies an option to a diagram.
type DiagramOption func(*Diagram) error

// NewDiagram creates a new diagram.
func NewDiagram(opts ...DiagramOption) (*Diagram, error) {
	d := &Diagram{
		nodes:       make([]*Node, 0),
		connections: make([]*Connection, 0),
		globals:     make(map[string]string),
		classes:     make(map[string]map[string]string),

		seenNodes: make(map[string]struct{}),
	}

	for _, opt := range opts {
		if err := opt(d); err != nil {
			return nil, errors.Wrap(err, "failed to apply diagram option")
		}
	}

	return d, nil
}

// AddClass defines a global class with specified attributes
func (d *Diagram) AddClass(name string, attributes map[string]string) *Diagram {
	d.classes[name] = attributes
	return d
}

// SetGlobal sets a global diagram property
func (d *Diagram) SetGlobal(key, value string) *Diagram {
	d.globals[key] = value
	return d
}

// Render generates the D2 script string
func (d *Diagram) D2() string {
	b := builderPool.Get().(*strings.Builder)
	defer func() {
		b.Reset()
		builderPool.Put(b)
	}()

	// Add classes
	if len(d.classes) > 0 {
		b.WriteString("classes: {\n")
		for className, attrs := range d.classes {
			fmt.Fprintf(b, "  %s: {\n", className)
			for name, value := range attrs {
				if strings.HasPrefix(value, "#") {
					value = `"` + value + `"`
				}
				fmt.Fprintf(b, "    %s: %s\n", name, value)
			}
			b.WriteString("  }\n")
		}
		b.WriteString("}\n")
	}

	// Add globals
	for name, value := range d.globals {
		fmt.Fprintf(b, "%s: %s\n", name, value)
	}
	if len(d.globals) > 0 {
		b.WriteString("\n")
	}

	// Add nodes
	for _, node := range d.nodes {
		b.WriteString(node.render())
		b.WriteString("\n")
	}

	// Add connections
	for _, conn := range d.connections {
		b.WriteString(conn.render())
		b.WriteString("\n")
	}

	return strings.TrimSpace(b.String())
}
