package diagram

import (
	"fmt"
	"maps"
	"strings"
)

// Node represents a D2 diagram node
type Node struct {
	ID    string
	Label string
	Tag   string
	Code  string
	// SQLFields  []sqlField
	Attributes map[string]string
}

// NodeOption configures a Node
type NodeOption func(*Node)

// AddNode adds a new node to the diagram
func (d *Diagram) AddNode(id string, opts ...NodeOption) *Diagram {
	if _, seen := d.seenNodes[id]; seen {
		return d
	}

	node := &Node{
		ID:         id,
		Attributes: make(map[string]string),
	}

	for _, opt := range opts {
		opt(node)
	}

	d.nodes = append(d.nodes, node)
	return d
}

func (n *Node) render() string {
	b := builderPool.Get().(*strings.Builder)
	defer func() {
		b.Reset()
		builderPool.Put(b)
	}()

	attrs := make(map[string]string)
	maps.Copy(attrs, n.Attributes)

	labelInAttrs := attrs["label"]
	switch {
	case labelInAttrs != "":
		b.WriteString(n.ID)

	case n.Label != "" && n.Label != n.ID:
		if n.Label == "" {
			n.Label = `""`
		}
		fmt.Fprintf(b, "%s: %s", n.ID, n.Label)

	case n.Tag != "" && n.Code != "":
		fmt.Fprintf(b, "%s: |%s\n%s|\n\n", n.ID, n.Tag, n.Code)

	default:
		b.WriteString(n.ID)
	}
	n.writeAttributes(b, attrs)
	return b.String()
}

func (n *Node) writeAttributes(b *strings.Builder, attrs map[string]string) {
	if len(attrs) == 0 {
		return
	}

	b.WriteString(" {\n")
	for name, value := range attrs {
		if name == "label" && value == n.Label {
			continue
		}

		if strings.HasPrefix(value, "#") {
			value = `"` + value + `"`
		}

		fmt.Fprintf(b, "    %s: %s\n", name, value)
	}
	b.WriteString("}")
}

func WithLabel(label string) NodeOption {
	return func(n *Node) {
		n.Label = label
	}
}

func WithShape(shape string) NodeOption {
	return func(n *Node) {
		n.Attributes["shape"] = shape
	}
}

func WithNodeAttribute(key, value string) NodeOption {
	return func(n *Node) {
		n.Attributes[key] = value
	}
}

func WithCode(tag, content string) NodeOption {
	if tag == "latex" || tag == "tex" {
		content = strings.ReplaceAll(content, `\`, `\\`)
	}

	return func(n *Node) {
		n.Tag = tag
		n.Code = content
	}
}

func WithClass(classes ...string) NodeOption {
	return func(n *Node) {
		prev, exists := n.Attributes["class"]
		if exists {
			prev = strings.Trim(prev, "[]")
			classes = append(strings.Split(prev, "; "), classes...)
		}
		n.Attributes["class"] = "[" + strings.Join(classes, "; ") + "]"
	}
}

// type sqlField struct {
// 	key        string
// 	value      string
// 	constraint string
// }
