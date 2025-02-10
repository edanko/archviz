package domain

import (
	"fmt"
	"strings"
)

const (
	// idSeparator is the character used to separate component IDs.
	idSeparator = "."
)

// trieNode represents a node in a trie.
type trieNode struct {
	children  map[string]*trieNode
	component *Component
}

// ComponentGraph represents a graph of architectural components.
type ComponentGraph struct {
	root                 *trieNode
	idMap                map[string]*Component
	unresolvedComponents map[string]*Component
	typeIndex            map[ComponentType]map[string]*Component
	tagIndex             map[string]map[string]map[string]struct{}

	// key is component ID from which this link originates.
	links map[string][]Link
}

// NewComponentGraph creates a new ComponentGraph.
func NewComponentGraph() *ComponentGraph {
	return &ComponentGraph{
		root:                 &trieNode{},
		idMap:                make(map[string]*Component),
		unresolvedComponents: make(map[string]*Component),
		typeIndex:            make(map[ComponentType]map[string]*Component),
		tagIndex:             make(map[string]map[string]map[string]struct{}),

		links: make(map[string][]Link),
	}
}

// AddComponent adds a component to the graph.
func (g *ComponentGraph) AddComponent(id, title string, opts ...componentOption) {
	if _, exists := g.idMap[id]; exists {
		return
	}

	component := newComponent(id, title, opts...)

	parts := strings.Split(id, idSeparator)
	current := g.root
	var parent *Component

	for i, part := range parts {
		if current.children == nil {
			current.children = make(map[string]*trieNode)
		}
		if _, ok := current.children[part]; !ok {
			current.children[part] = &trieNode{}
		}
		current = current.children[part]

		if i < len(parts)-1 && current.component != nil {
			parent = current.component
		}
	}

	component.Parent = parent
	current.component = component
	g.idMap[id] = component

	g.resolvePendingLinks(id, component)
	g.updateTypeIndex(component)
	g.updateTagIndex(component)
	g.updateChildren(parts, component)
}

// resolvePendingLinks resolves any pending links pointing to this component.
func (g *ComponentGraph) resolvePendingLinks(id string, _ *Component) {
	_, exists := g.unresolvedComponents[id]
	if !exists {
		return
	}

	// ???

	delete(g.unresolvedComponents, id)
}

// updateTypeIndex updates the type index.
func (g *ComponentGraph) updateTypeIndex(component *Component) {
	if g.typeIndex[component.Type] == nil {
		g.typeIndex[component.Type] = make(map[string]*Component)
	}
	g.typeIndex[component.Type][component.ID] = component
}

// updateTagIndex updates the tag index.
func (g *ComponentGraph) updateTagIndex(component *Component) {
	for k, v := range component.Tags {
		if g.tagIndex[k] == nil {
			g.tagIndex[k] = make(map[string]map[string]struct{})
		}
		if g.tagIndex[k][v] == nil {
			g.tagIndex[k][v] = make(map[string]struct{})
		}
		g.tagIndex[k][v][component.ID] = struct{}{}
	}
}

// updateChildren updates the children of a component.
func (g *ComponentGraph) updateChildren(parts []string, component *Component) {
	newDepth := len(parts)
	for existingID, existingComp := range g.idMap {
		if existingID == component.ID {
			continue
		}

		existingParts := strings.Split(existingID, idSeparator)
		if len(existingParts) != newDepth+1 {
			continue
		}

		match := true
		for i := range newDepth {
			if existingParts[i] != parts[i] {
				match = false
				break
			}
		}

		if match {
			if existingComp.Parent == nil ||
				len(strings.Split(existingComp.Parent.ID, idSeparator)) >= newDepth {
				existingComp.Parent = component
			}
		}
	}
}

// GetComponent returns a component by ID.
func (g *ComponentGraph) GetComponent(id string) (*Component, error) {
	component, exists := g.idMap[id]
	if exists {
		return component, nil
	}

	_, unresolved := g.unresolvedComponents[id]
	if unresolved {
		return nil, ErrComponentUnresolved
	}

	return nil, ErrComponentNotFound
}

// GetComponentsInNamespace returns all components in a namespace.
func (g *ComponentGraph) GetComponentsInNamespace(namespace string) []*Component {
	parts := strings.Split(namespace, idSeparator)
	current := g.root

	for _, part := range parts {
		if current.children == nil {
			return nil
		}
		next, ok := current.children[part]
		if !ok {
			return nil
		}
		current = next
	}

	return g.collectComponents(current)
}

// collectAnscestors recursively collects all anscestors of a component.
func (g *ComponentGraph) collectAnscestors(component *Component) []*Component {
	parents := []*Component{}
	current := component
	for current != nil && current.Parent != nil {
		parent, exists := g.idMap[current.Parent.ID]
		if !exists {
			break
		}
		parents = append(parents, parent)
		current = parent
	}
	return parents
}

// collectComponents recursively collects all components in a trie node.
func (g *ComponentGraph) collectComponents(node *trieNode) []*Component {
	components := make([]*Component, 0)

	if node.component != nil {
		components = append(components, node.component)
	}
	for _, child := range node.children {
		children := g.collectComponents(child)
		components = append(components, children...)
	}

	return components
}

// FIXME: do something with this.
func (g *ComponentGraph) Check() {
	for _, c := range g.unresolvedComponents {
		fmt.Println("Unresolved component:", c.ID)
	}
}
