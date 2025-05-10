package domain

import (
	"fmt"
	"maps"
	"path/filepath"
	"slices"
	"sync"
)

type Filter struct {
	IDs        []string
	Types      []ComponentType
	Tags       map[string]string
	DepthRange [2]int // [min, max]

	ConnectedTo   string
	ConnectedFrom string

	WithChildren    bool
	WithDependents  bool
	IsolateSubgraph bool

	Aggregate bool

	ExcludeIDs   []string
	ExcludeTypes []ComponentType
}

// isDepthRangeSet returns true if the depth range is set.
func (f *Filter) isDepthRangeSet() bool {
	return f.DepthRange[0] != 0 || f.DepthRange[1] != 0
}

// walk traverses the component graph depth-first, calling the provided visitor
// function for each component with its path in the graph.
//
// The visitor function is called with the component path as a slice of strings,
// followed by the component itself.
//
// The traversal is performed using a stack, where each item represents a node in
// the graph and its path. The stack is initialized with the root node, and then
// iteratively processed until it is empty.
//
// When a node is processed, the visitor function is called with the node's path
// and component. Then, the node's children are added to the stack, with their
// paths being the concatenation of the current node's path and the segment that
// leads to the child node.
func (g *ComponentGraph) walk(visitor func(path []string, component *Component)) {
	if g == nil || visitor == nil {
		return
	}

	stack := []struct {
		node *trieNode
		path []string
	}{
		{node: g.root, path: []string{}},
	}

	for len(stack) > 0 {
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if current.node.component != nil {
			visitor(current.path, current.node.component)
		}

		for segment, child := range current.node.children {
			newPath := append(slices.Clone(current.path), segment)
			stack = append(stack, struct {
				node *trieNode
				path []string
			}{child, newPath})
		}
	}
}

// Filter applies the given filter to the graph and returns a new set of components
// and links that match the filter criteria. The filter criteria can be a set of
// IDs, a set of types, a set of tags, or a regular expression. The filter can also
// be configured to only consider components within a certain depth range.
//
// Additionally, the filter can be configured to either aggregate components
// (i.e. collapse them into a single component) or isolate the subgraph (i.e.
// remove all components that are not reachable from the filtered components).
//
// The returned components are the filtered components, and the returned links are
// the links between the filtered components.
func (g *ComponentGraph) Filter(f *Filter) ([]*Component, []Link) {
	if g == nil {
		return nil, nil
	}
	if f == nil {
		return nil, nil
	}

	var components []*Component
	seenComponents := make(map[string]struct{})
	var links []Link

	g.walk(func(path []string, c *Component) {
		if c == nil {
			return
		}

		match := g.match(c, path, len(path), f)
		if !match {
			return
		}

		if f.isDepthRangeSet() {
			currentDepth := len(path)
			if currentDepth < f.DepthRange[0] || currentDepth > f.DepthRange[1] {
				return
			}
		}

		if _, exists := seenComponents[c.ID]; !exists {
			components = append(components, c)
			seenComponents[c.ID] = struct{}{}
		}

		for _, anscestor := range g.collectAnscestors(c) {
			if anscestor == nil || len(path) < 1 {
				continue
			}

			if f.isDepthRangeSet() {
				parentDepth := len(path) - 1
				if parentDepth < f.DepthRange[0] || parentDepth > f.DepthRange[1] {
					continue
				}
			}

			if _, exists := seenComponents[anscestor.ID]; !exists {
				components = append(components, anscestor)
				seenComponents[anscestor.ID] = struct{}{}
			}
		}

		if _, exists := seenComponents[c.ID]; exists {
			return
		}
		components = append(components, c)
		seenComponents[c.ID] = struct{}{}
	})

	if f.Aggregate && f.DepthRange[1] > 0 {
		components, links = g.aggregateComponents(components, f.DepthRange[1])
	}

	if f.IsolateSubgraph {
		components, links = g.isolateSubgraph(components, links)
	}

	return components, links
}

// match returns true if the given component matches the given filter criteria.
//
// The filter criteria are:
//
//   - depth: if the filter specifies a depth range, the component must have a depth
//     within that range.
//   - IDs: if the filter specifies a set of IDs, the component must have one of
//     those IDs.
//   - types: if the filter specifies a set of types, the component must have one of
//     those types.
//   - tags: if the filter specifies a set of tags, the component must have all of
//     those tags.
//   - connected to: if the filter specifies a connected to ID, the component must
//     have a link to that ID.
//   - connected from: if the filter specifies a connected from ID, the component
//     must have a link from that ID.
func (g *ComponentGraph) match(
	component *Component,
	_ []string,
	depth int,
	filter *Filter,
) bool {
	if component == nil || filter == nil {
		return false
	}

	if slices.Contains(filter.ExcludeIDs, component.ID) {
		return false
	}

	if slices.Contains(filter.ExcludeTypes, component.Type) {
		return false
	}

	if (filter.DepthRange[0] != 0 || filter.DepthRange[1] != 0) &&
		(depth < filter.DepthRange[0] || depth > filter.DepthRange[1]) {
		return false
	}

	if len(filter.IDs) > 0 &&
		!matchIDs(filter.IDs, component.ID) {
		return false
	}

	if len(filter.Types) > 0 &&
		!slices.Contains(filter.Types, component.Type) {
		return false
	}

	for tag, tagValue := range filter.Tags {
		if component.Tags == nil || component.Tags[tag] != tagValue {
			return false
		}
	}

	if filter.ConnectedTo != "" &&
		!g.isConnectedTo(component.ID, filter.ConnectedTo) {
		return false
	}

	if filter.ConnectedFrom != "" &&
		!g.isConnectedFrom(component.ID, filter.ConnectedFrom) {
		return false
	}

	return true
}

// matchIDs returns true if the id matches any of the patterns.
func matchIDs(patterns []string, id string) bool {
	for _, pattern := range patterns {
		if pattern == id {
			return true
		}

		matched, err := filepath.Match(pattern, id)
		if err == nil && matched {
			return true
		}
	}
	return false
}

func (g *ComponentGraph) isConnectedTo(sourceID, targetID string) bool {
	for _, link := range g.links[sourceID] {
		if link.TargetID == targetID {
			return true
		}
	}
	return false
}

func (g *ComponentGraph) isConnectedFrom(targetID, sourceID string) bool {
	links, exists := g.links[sourceID]
	if !exists {
		return false
	}
	for _, link := range links {
		if link.TargetID == targetID {
			return true
		}
	}
	return false
}

func (g *ComponentGraph) aggregateComponents(
	components []*Component,
	maxDepth int,
) ([]*Component, []Link) {
	collapsedComponents := make([]*Component, 0)
	seenComponents := make(map[string]struct{})

	linkMap := make(map[string]map[string]int)
	ch := make(chan []AggregatedLink, len(g.idMap))

	processLinks := func(c *Component) {
		localLinks := make(map[string]map[string]int)
		sourceNs := c.getNamespace(maxDepth)

		for _, link := range g.links[c.ID] {
			if targetComp, exists := g.idMap[link.TargetID]; exists {
				targetNs := targetComp.getNamespace(maxDepth)
				if localLinks[sourceNs] == nil {
					localLinks[sourceNs] = make(map[string]int)
				}
				localLinks[sourceNs][targetNs]++
			}
		}

		ch <- convertToLinks(localLinks)
	}

	for _, comp := range components {
		ns := comp.getNamespace(maxDepth)
		_, exists := seenComponents[ns]
		if exists {
			continue
		}

		seenComponents[ns] = struct{}{}

		collapsedComponents = append(collapsedComponents, &Component{
			ID:          ns,
			Title:       comp.Title,
			Description: comp.Description,
			Type:        comp.Type,
			Tags:        maps.Clone(comp.Tags),
			Parent:      comp.Parent,
			Depth:       comp.Depth,
			Properties:  maps.Clone(comp.Properties),
		})
	}

	var wg sync.WaitGroup
	for _, comp := range components {
		wg.Add(1)
		go func(c *Component) {
			defer wg.Done()
			processLinks(c)
		}(comp)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for aggregatedLinks := range ch {
		for _, link := range aggregatedLinks {
			if link.Source == link.Target {
				continue // Skip self-links
			}
			if _, exists := linkMap[link.Source]; !exists {
				linkMap[link.Source] = make(map[string]int)
			}
			linkMap[link.Source][link.Target] += link.TotalLinks
		}
	}

	// Collect all unique links
	processedLinks := make([]Link, 0)
	for src, targets := range linkMap {
		for tgt, count := range targets {
			if src != tgt { // Ensure we don't add self-links
				processedLinks = append(processedLinks, Link{
					SourceID:    src,
					TargetID:    tgt,
					Description: fmt.Sprintf("Links count: %d", count),
				})
			}
		}
	}

	return collapsedComponents, processedLinks
}

type AggregatedLink struct {
	Source      string
	Target      string
	TotalLinks  int
	LinkSamples []Link // Optional: keep sample links
}

func convertToLinks(local map[string]map[string]int) []AggregatedLink {
	var links []AggregatedLink
	for src, targets := range local {
		for tgt, count := range targets {
			links = append(links, AggregatedLink{
				Source:     src,
				Target:     tgt,
				TotalLinks: count,
			})
		}
	}
	return links
}

func (g *ComponentGraph) isolateSubgraph(
	components []*Component,
	links []Link,
) ([]*Component, []Link) {
	componentIDs := make(map[string]struct{})
	for _, c := range components {
		if c != nil {
			componentIDs[c.ID] = struct{}{}
		}
	}

	var filteredComponents []*Component
	for _, c := range components {
		if c != nil && containsKey(componentIDs, c.ID) {
			filteredComponents = append(filteredComponents, c)
		}
	}

	var filteredLinks []Link
	for _, link := range links {
		if link.SourceID == "" || link.TargetID == "" {
			continue
		}
		if containsKey(componentIDs, link.SourceID) &&
			containsKey(componentIDs, link.TargetID) {
			filteredLinks = append(filteredLinks, link)
		}
	}

	return filteredComponents, filteredLinks
}

// Helper function to check if a key exists in a map
func containsKey(m map[string]struct{}, key string) bool {
	_, exists := m[key]
	return exists
}
