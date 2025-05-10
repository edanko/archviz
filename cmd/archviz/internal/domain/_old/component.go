package domain

import "cmp"

// enrichComponent adds missing values to the current node with the values from the other node.
func (n *Component) enrichComponent(other *Component) *Component {
	n.Title = cmp.Or(n.Title, other.Title)
	n.Shape = cmp.Or(n.Shape, other.Shape)
	n.Description = cmp.Or(n.Description, other.Description)
	n.Technologies = mergeUnique(n.Technologies, other.Technologies)
	return n
}

// mergeUnique returns unique elements from the given slices.
func mergeUnique[T comparable](ss ...[]T) []T {
	uniqueMap := make(map[T]struct{})

	for _, s := range ss {
		for _, element := range s {
			uniqueMap[element] = struct{}{}
		}
	}

	result := make([]T, 0, len(uniqueMap))

	for key := range uniqueMap {
		result = append(result, key)
	}

	return result
}
