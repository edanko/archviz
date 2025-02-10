package domain

import "strings"

// Link represents a link between two components.
type Link struct {
	// SourceID and TargetID are component IDs.
	SourceID, TargetID string
	// Description is a human-readable description of the link.
	Description string
}

func (l Link) String() string {
	elements := []string{l.SourceID, l.TargetID}
	return strings.Join(elements, " -> ")
}

// AddLink adds a link between two components.
// Method allows to add link if source or target component is not added to the graph yet.
// In this case, the component will be added to the unresolved components.
func (g *ComponentGraph) AddLink(sourceID, targetID, description string) {
	g.maybeAddToUnresolved(sourceID)
	g.maybeAddToUnresolved(targetID)

	g.links[sourceID] = append(g.links[sourceID], Link{
		SourceID:    sourceID,
		TargetID:    targetID,
		Description: description,
	})
}

// maybeAddToUnresolved adds a component to the unresolved components
// if it doesn't already exist.
func (g *ComponentGraph) maybeAddToUnresolved(id string) {
	_, exists := g.idMap[id]
	if exists {
		return
	}

	g.unresolvedComponents[id] = &Component{
		ID: id,
	}
}
