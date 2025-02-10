package domain

import "strings"

const locationSeparator = "/"

// View represents a view in the architecture.
type View struct {
	// ID is a unique identifier of the view.
	ID string
	// Title is a name of the view.
	Title string
	// Location is a path to the view in the menu.
	Location string
	// MaxDepth is the maximum depth of the components will be displayed in the view (0 for unlimited).
	MaxDepth int
	// Components is a list of components that are part of the view.
	Components []string
}

// locationParts returns a list of parts of the location split by locationSeparator.
func (v *View) locationParts() []string {
	return strings.Split(v.Location, locationSeparator)
}

// Views is a map that holds views by their unique identifiers.
type Views map[string]View

// BuildHierarchy constructs a hierarchical map from the given views.
// Currently it only supports a single level of hierarchy.
func (vs Views) BuildHierarchy() map[string][]View {
	result := make(map[string][]View)

	for _, view := range vs {
		locationParts := view.locationParts()
		if len(locationParts) == 1 && locationParts[0] == "" {
			continue
		}

		key := locationParts[0]

		result[key] = append(result[key], view)
	}

	return result
}
