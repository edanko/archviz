package domain

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestViews_BuildHierarchy(t *testing.T) {
	t.Parallel()

	t.Run("should build sample hierarchy", func(t *testing.T) {
		t.Parallel()

		// arrange
		views := Views{
			"view1": {ID: "view1", Title: "Example View 1", Location: "S1", Components: []string{"componentA"}},
			"view2": {ID: "view2", Title: "Example View 2", Location: "S1", Components: []string{"componentB"}},
			"view3": {ID: "view3", Title: "Example View 3", Location: "S2", Components: []string{"componentC"}},
			"view4": {ID: "view4", Title: "Example View 4", Location: "", Components: []string{"componentD"}},
		}

		// act
		got := views.BuildHierarchy()

		// assert
		want := map[string][]View{
			"S1": {
				views["view1"],
				views["view2"],
			},
			"S2": {
				views["view3"],
			},
		}
		require.Equal(t, want, got)
	})
}
