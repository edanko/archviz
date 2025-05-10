package domain

import "strings"

// ComponentType is the type of an architectural component.
type ComponentType string

const (
	TypeService  ComponentType = "service"
	TypeDatabase ComponentType = "database"
	TypeQueue    ComponentType = "queue"
	TypeUI       ComponentType = "ui"
)

// componentOption is a function that configures a component.
type componentOption func(*Component)

// Component represents an architectural component.
type Component struct {
	ID          string
	Title       string
	Description string
	Type        ComponentType
	Tags        map[string]string
	Parent      *Component
	Properties  map[string]any
	Depth       int
}

// getNamespace returns the namespace of a component at a given depth.
func (c *Component) getNamespace(depth int) string {
	parts := strings.Split(c.ID, idSeparator)
	if depth < 0 {
		return c.ID
	}
	if depth >= len(parts) {
		return c.ID
	}
	return strings.Join(parts[:depth], idSeparator)
}

// newComponent creates a new component.
func newComponent(id, title string, opts ...componentOption) *Component {
	component := &Component{
		ID:    id,
		Title: title,
		Depth: len(strings.Split(id, idSeparator)),
	}
	for _, opt := range opts {
		opt(component)
	}
	return component
}

// WithType sets the component type.
func WithType(t ComponentType) componentOption {
	return func(c *Component) {
		c.Type = t
	}
}

// WithTags sets the component tags.
func WithTags(tags map[string]string) componentOption {
	return func(c *Component) {
		c.Tags = tags
	}
}

// WithDescription sets the component description.
func WithDescription(desc string) componentOption {
	return func(c *Component) {
		c.Description = desc
	}
}

// WithProperties sets the component properties.
func WithProperties(props map[string]any) componentOption {
	return func(c *Component) {
		c.Properties = props
	}
}
