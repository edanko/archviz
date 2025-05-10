package diagram

import (
	"fmt"
	"strings"
)

// Connection represents a D2 diagram connection between nodes
type Connection struct {
	From       string
	Direction  string
	To         string
	Label      string
	Attributes map[string]string
}

// ConnectionOption configures a Connection
type ConnectionOption func(*Connection)

// AddConnection adds a new connection between nodes
func (d *Diagram) addConnection(from, direction, to string, opts ...ConnectionOption) *Diagram {
	conn := &Connection{
		From:       from,
		Direction:  direction,
		To:         to,
		Attributes: make(map[string]string),
	}
	for _, opt := range opts {
		opt(conn)
	}
	d.connections = append(d.connections, conn)
	return d
}

func (c *Connection) render() string {
	b := builderPool.Get().(*strings.Builder)
	defer func() {
		b.Reset()
		builderPool.Put(b)
	}()

	fmt.Fprintf(b, "%s -> %s", c.From, c.To)

	if c.Label != "" {
		fmt.Fprintf(b, ": \"%s\"", c.Label)
	}

	if len(c.Attributes) > 0 {
		b.WriteString(" {\n")
		for name, value := range c.Attributes {
			fmt.Fprintf(b, "    %s: %s\n", name, value)
		}
		b.WriteString("}")
	}
	return b.String()
}

func WithConnectionClass(classes ...string) ConnectionOption {
	return func(c *Connection) {
		prev, exists := c.Attributes["class"]
		if exists {
			classes = append(strings.Split(prev, "; "), classes...)
		}
		c.Attributes["class"] = strings.Join(classes, "; ")
	}
}

// Connection options
func WithConnectionLabel(label string) ConnectionOption {
	return func(c *Connection) {
		c.Label = label
	}
}

func WithConnectionAttribute(key, value string) ConnectionOption {
	return func(c *Connection) {
		c.Attributes[key] = value
	}
}

// AddSimpleConnection adds a simple (non-directional) connection between two nodes.
func (d *Diagram) AddSimpleConnection(
	from, to string,
	opts ...ConnectionOption,
) *Diagram {
	return d.addConnection(
		from, "--", to,
		opts...,
	)
}

// AddDirectedConnection adds a directed connection between two nodes.
func (d *Diagram) AddDirectedConnection(
	from, to string,
	opts ...ConnectionOption,
) *Diagram {
	return d.addConnection(
		from, "->", to,
		opts...,
	)
}

// AddBidirectionalConnection adds a bidirectional connection between two nodes.
func (d *Diagram) AddBidirectionalConnection(
	from, to string,
	opts ...ConnectionOption,
) *Diagram {
	return d.addConnection(
		from, "<->", to,
		opts...,
	)
}
