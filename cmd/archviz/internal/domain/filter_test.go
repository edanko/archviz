package domain

import (
	"math"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

var makeSampleGraph = sync.OnceValue(func() ComponentGraph {
	g := NewComponentGraph()

	g.AddComponent(
		"user",
		"User",
	)

	g.AddComponent(
		"bs",
		"Big System",
	)
	g.AddComponent(
		"bs.auth",
		"Auth Subsystem",
	)

	g.AddComponent(
		"bs.auth.frontend",
		"Auth Frontend",
		WithType(TypeUI),
	)

	g.AddLink("user", "bs.auth.frontend", "logins")

	g.AddComponent(
		"bs.auth.gateway",
		"Auth Gateway",
	)
	g.AddLink("bs.auth.frontend", "bs.auth.gateway", "uses")

	g.AddComponent(
		"bs.auth.backend",
		"Auth Backend",
		WithType(TypeService),
		WithTags(map[string]string{
			"role": "backend",
		}),
	)
	g.AddLink("bs.auth.gateway", "bs.auth.backend", "uses")

	g.AddComponent(
		"bs.auth.backend.db",
		"Auth Database",
		WithType(TypeDatabase),
	)
	g.AddLink("bs.auth.backend", "bs.auth.backend.db", "store credentials")

	g.AddComponent(
		"bs.orders",
		"Orders Subsystem",
	)
	g.AddComponent(
		"bs.orders.frontend",
		"Orders Frontend",
		WithType(TypeUI),
	)
	g.AddLink("user", "bs.orders.frontend", "orders")

	g.AddComponent(
		"bs.orders.gateway",
		"Orders Gateway",
	)
	g.AddLink("bs.orders.frontend", "bs.orders.gateway", "uses")
	g.AddLink("bs.orders.gateway", "bs.auth.gateway", "authorizes")

	g.AddComponent(
		"bs.orders.backend",
		"Orders Backend",
		WithType(TypeService),
		WithTags(map[string]string{
			"role": "backend",
		}),
	)
	g.AddLink("bs.orders.gateway", "bs.orders.backend", "uses")

	g.AddComponent(
		"bs.orders.backend.db",
		"Orders Database",
		WithType(TypeDatabase),
	)
	g.AddLink("bs.orders.backend", "bs.orders.backend.db", "store orders")

	return *g
})

func Test_match(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		patterns []string
		input    string
		want     bool
	}{
		{
			name: "exact match",
			patterns: []string{
				"a.b.c.d",
			},
			input: "a.b.c.d",
			want:  true,
		},
		{
			name: "wildcard match in last part",
			patterns: []string{
				"a.b.c.*",
			},
			input: "a.b.c.d",
			want:  true,
		},
		{
			name: "wildcard match in middle part",
			patterns: []string{
				"a.*.c.d",
			},
			input: "a.b.c.d",
			want:  true,
		},
		{
			name: "no match",
			patterns: []string{
				"a.b.c.d",
			},
			input: "a.b.c.e",
			want:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// act
			got := matchIDs(tt.patterns, tt.input)

			// assert
			require.Equal(t, tt.want, got)
		})
	}
}

func TestComponentGraph_walk(t *testing.T) {
	t.Parallel()

	t.Run("successfully walk whole component graph", func(t *testing.T) {
		t.Parallel()

		// arrange
		g := makeSampleGraph()

		// act
		var gotIDs []string
		g.walk(func(path []string, c *Component) {
			gotIDs = append(gotIDs, c.ID)
		})

		// assert
		wantIDs := []string{
			"user",
			"bs",
			"bs.auth",
			"bs.auth.frontend",
			"bs.auth.gateway",
			"bs.auth.backend",
			"bs.auth.backend.db",
			"bs.orders",
			"bs.orders.frontend",
			"bs.orders.gateway",
			"bs.orders.backend",
			"bs.orders.backend.db",
		}
		require.ElementsMatch(t, wantIDs, gotIDs)
	})
}

func TestComponentGraph_Filter(t *testing.T) {
	t.Parallel()

	t.Run("successfully filter by component ID", func(t *testing.T) {
		t.Parallel()

		// arrange
		g := makeSampleGraph()
		componentID := "bs.orders.backend"

		// act
		filter := &Filter{
			IDs: []string{
				componentID,
			},
		}
		gotComponents, gotLinks := g.Filter(filter)

		// assert
		wantComponents := []*Component{
			must(g.GetComponent("bs")),
			must(g.GetComponent("bs.orders")),

			must(g.GetComponent(componentID)),
		}
		require.ElementsMatch(t, wantComponents, gotComponents)

		require.Empty(t, gotLinks)
	})

	t.Run("successfully filter by component ID (with pattern matching)", func(t *testing.T) {
		t.Parallel()

		// arrange
		g := makeSampleGraph()

		// act
		filter := &Filter{
			IDs: []string{
				"bs.orders.*",
			},
		}
		gotComponents, gotLinks := g.Filter(filter)

		// assert
		wantComponents := []*Component{
			must(g.GetComponent("bs")),
			must(g.GetComponent("bs.orders")),

			must(g.GetComponent("bs.orders.gateway")),
			must(g.GetComponent("bs.orders.frontend")),
			must(g.GetComponent("bs.orders.backend")),
			must(g.GetComponent("bs.orders.backend.db")),
		}
		require.ElementsMatch(t, wantComponents, gotComponents)

		require.Empty(t, gotLinks)
	})

	t.Run("successfully filter by component IDs when there is an exact match and a pattern match", func(t *testing.T) {
		t.Parallel()

		// arrange
		g := makeSampleGraph()

		// act
		filter := &Filter{
			IDs: []string{
				"user",
				"bs.orders.*",
			},
		}
		gotComponents, gotLinks := g.Filter(filter)

		// assert
		wantComponents := []*Component{
			must(g.GetComponent("user")),
			must(g.GetComponent("bs")),
			must(g.GetComponent("bs.orders")),

			must(g.GetComponent("bs.orders.gateway")),
			must(g.GetComponent("bs.orders.frontend")),
			must(g.GetComponent("bs.orders.backend")),
			must(g.GetComponent("bs.orders.backend.db")),
		}
		require.ElementsMatch(t, wantComponents, gotComponents)

		require.Empty(t, gotLinks)
	})

	t.Run("successfully filter by depth range", func(t *testing.T) {
		t.Parallel()

		// arrange
		g := makeSampleGraph()

		// act
		filter := &Filter{
			DepthRange: [2]int{0, 2},
		}
		gotComponents, gotLinks := g.Filter(filter)

		// assert
		wantComponents := []*Component{
			must(g.GetComponent("user")),
			must(g.GetComponent("bs")),
			must(g.GetComponent("bs.auth")),
			must(g.GetComponent("bs.orders")),
		}
		require.ElementsMatch(t, wantComponents, gotComponents)

		require.Empty(t, gotLinks)
	})

	t.Run("successfully filter by component type", func(t *testing.T) {
		t.Parallel()

		// arrange
		g := makeSampleGraph()

		// act
		filter := &Filter{
			Types: []ComponentType{
				TypeService,
			},
		}
		gotComponents, gotLinks := g.Filter(filter)

		// assert
		wantComponents := []*Component{
			must(g.GetComponent("bs")),
			must(g.GetComponent("bs.auth")),
			must(g.GetComponent("bs.orders")),

			must(g.GetComponent("bs.orders.backend")),
			must(g.GetComponent("bs.auth.backend")),
		}
		require.ElementsMatch(t, wantComponents, gotComponents)

		require.Empty(t, gotLinks)
	})

	t.Run("successfully filter by component tags", func(t *testing.T) {
		t.Parallel()

		// arrange
		g := makeSampleGraph()

		// act
		filter := &Filter{
			Tags: map[string]string{
				"role": "backend",
			},
		}
		gotComponents, gotLinks := g.Filter(filter)

		// assert
		wantComponents := []*Component{
			must(g.GetComponent("bs")),
			must(g.GetComponent("bs.auth")),
			must(g.GetComponent("bs.orders")),

			must(g.GetComponent("bs.orders.backend")),
			must(g.GetComponent("bs.auth.backend")),
		}
		require.ElementsMatch(t, wantComponents, gotComponents)

		require.Empty(t, gotLinks)
	})

	t.Run("successfully filter by connected to", func(t *testing.T) {
		t.Parallel()

		// arrange
		g := makeSampleGraph()

		// act
		filter := &Filter{
			ConnectedTo: "bs.auth.backend",
		}
		gotComponents, gotLinks := g.Filter(filter)

		// assert
		wantComponents := []*Component{
			must(g.GetComponent("bs")),
			must(g.GetComponent("bs.auth")),

			must(g.GetComponent("bs.auth.gateway")),
		}
		require.ElementsMatch(t, wantComponents, gotComponents)

		require.Empty(t, gotLinks)
	})

	t.Run("successfully filter by connected from", func(t *testing.T) {
		t.Parallel()

		// arrange
		g := makeSampleGraph()

		// act
		filter := &Filter{
			ConnectedFrom: "user",
		}
		gotComponents, gotLinks := g.Filter(filter)

		// assert
		wantComponents := []*Component{
			must(g.GetComponent("bs")),
			must(g.GetComponent("bs.auth")),
			must(g.GetComponent("bs.orders")),

			must(g.GetComponent("bs.orders.frontend")),
			must(g.GetComponent("bs.auth.frontend")),
		}
		require.ElementsMatch(t, wantComponents, gotComponents)

		require.Empty(t, gotLinks)
	})

	t.Run("successfully aggregate components and links by max depth", func(t *testing.T) {
		t.Parallel()

		// arrange
		g := makeSampleGraph()

		// act
		filter := &Filter{
			DepthRange: [2]int{0, 2},
			Aggregate:  true,
		}
		gotComponents, gotLinks := g.Filter(filter)

		// assert
		wantComponents := []*Component{
			must(g.GetComponent("bs")),
			must(g.GetComponent("bs.orders")),
			must(g.GetComponent("bs.auth")),
			must(g.GetComponent("user")),
		}
		require.ElementsMatch(t, wantComponents, gotComponents)

		wantLinks := []Link{
			{
				SourceID:    "user",
				TargetID:    "bs.auth",
				Description: "Links count: 1",
			},
			{
				SourceID:    "user",
				TargetID:    "bs.orders",
				Description: "Links count: 1",
			},
		}
		require.ElementsMatch(t, wantLinks, gotLinks)
	})

	t.Run("successfully aggregate components and links by max depth with IDs filter", func(t *testing.T) {
		t.Parallel()

		// arrange
		g := makeSampleGraph()

		// act
		filter := &Filter{
			IDs: []string{
				"user",
				"bs",
				"bs.auth",
				"bs.auth.*",
			},
			DepthRange:      [2]int{0, math.MaxInt},
			Aggregate:       true,
			IsolateSubgraph: true,
		}
		gotComponents, gotLinks := g.Filter(filter)

		// assert
		wantComponents := []*Component{
			must(g.GetComponent("user")),
			must(g.GetComponent("bs")),
			must(g.GetComponent("bs.auth")),
			must(g.GetComponent("bs.auth.gateway")),
			must(g.GetComponent("bs.auth.frontend")),
			must(g.GetComponent("bs.auth.backend")),
			must(g.GetComponent("bs.auth.backend.db")),
		}
		require.ElementsMatch(t, wantComponents, gotComponents)

		wantLinks := []Link{
			{
				SourceID:    "user",
				TargetID:    "bs.auth.frontend",
				Description: "Links count: 1",
			},
			{
				SourceID:    "bs.auth.frontend",
				TargetID:    "bs.auth.gateway",
				Description: "Links count: 1",
			},
			{
				SourceID:    "bs.auth.gateway",
				TargetID:    "bs.auth.backend",
				Description: "Links count: 1",
			},
			{
				SourceID:    "bs.auth.backend",
				TargetID:    "bs.auth.backend.db",
				Description: "Links count: 1",
			},
		}
		require.ElementsMatch(t, wantLinks, gotLinks)
	})
}
