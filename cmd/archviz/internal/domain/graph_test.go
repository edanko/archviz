package domain

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestComponentGraph_AddComponent(t *testing.T) {
	t.Parallel()

	t.Run("successfully add component", func(t *testing.T) {
		t.Parallel()

		// arrange
		graph := NewComponentGraph()

		componentID := "a.b.c.d"
		componentTitle := "some title"

		// act
		graph.AddComponent(componentID, componentTitle)

		// assert
		got, err := graph.GetComponent(componentID)
		require.NoError(t, err)

		require.Equal(t, componentID, got.ID)
		require.Equal(t, componentTitle, got.Title)
	})

	t.Run("successfully add component and assign parent", func(t *testing.T) {
		t.Parallel()

		// arrange
		graph := NewComponentGraph()

		parentID := "a"
		parentTitle := "parent"
		graph.AddComponent(parentID, parentTitle)

		// act
		childID := "a.b"
		childTitle := "child"
		graph.AddComponent(childID, childTitle)

		// assert
		child, err := graph.GetComponent(childID)
		require.NoError(t, err)

		parent, err := graph.GetComponent(parentID)
		require.NoError(t, err)

		require.Equal(t, parent, child.Parent)
	})

	t.Run("successfully assign parent to existing component", func(t *testing.T) {
		t.Parallel()

		// arrange
		graph := NewComponentGraph()

		childID := "a.b"
		childTitle := "child"
		graph.AddComponent(childID, childTitle)

		// act
		parentID := "a"
		parentTitle := "parent"
		graph.AddComponent(parentID, parentTitle)

		// assert
		child, err := graph.GetComponent(childID)
		require.NoError(t, err)

		parent, err := graph.GetComponent(parentID)
		require.NoError(t, err)

		require.Equal(t, parent, child.Parent)
	})

	t.Run("resolve component when actual component is added #1", func(t *testing.T) {
		t.Parallel()

		// arrange
		graph := NewComponentGraph()

		component1ID := "a.b.c.d"
		component1Title := "some title"
		graph.AddComponent(component1ID, component1Title)

		component2ID := "a.b.c.e"
		graph.AddLink(component1ID, component2ID, "uses")

		_, err := graph.GetComponent(component2ID)
		require.ErrorIs(t, err, ErrComponentUnresolved)

		component2Title := "some other title"

		// act
		graph.AddComponent(component2ID, component2Title)

		// assert
		got, err := graph.GetComponent(component2ID)
		require.NoError(t, err)

		require.Equal(t, component2ID, got.ID)
		require.Equal(t, component2Title, got.Title)
	})

	t.Run("resolve component when actual component is added #2", func(t *testing.T) {
		t.Parallel()

		// arrange
		graph := NewComponentGraph()

		component1ID := "a.b.c.d"
		component1Title := "some title"
		graph.AddComponent(component1ID, component1Title)

		component2ID := "a.b.c.e"
		graph.AddLink(component2ID, component1ID, "uses")

		_, err := graph.GetComponent(component2ID)
		require.ErrorIs(t, err, ErrComponentUnresolved)

		component2Title := "some other title"

		// act
		graph.AddComponent(component2ID, component2Title)

		// assert
		got, err := graph.GetComponent(component2ID)
		require.NoError(t, err)

		require.Equal(t, component2ID, got.ID)
		require.Equal(t, component2Title, got.Title)
	})

	t.Run("do nothing if component already exists", func(t *testing.T) {
		t.Parallel()

		// arrange
		graph := NewComponentGraph()

		componentID := "a.b.c.d"
		componentTitle := "some title"
		graph.AddComponent(componentID, componentTitle)

		// act
		newTitle := "new title"
		graph.AddComponent(componentID, newTitle)

		// assert
		got, err := graph.GetComponent(componentID)
		require.NoError(t, err)

		require.Equal(t, componentID, got.ID)
		require.Equal(t, componentTitle, got.Title)
	})
}

func TestComponentGraph_GetComponent(t *testing.T) {
	t.Parallel()

	t.Run("successfully get component", func(t *testing.T) {
		t.Parallel()

		// arrange
		graph := NewComponentGraph()

		componentID := "a.b.c.d"
		componentTitle := "some title"
		graph.AddComponent(componentID, componentTitle)

		// act
		got, err := graph.GetComponent(componentID)

		// assert
		require.NoError(t, err)
		require.Equal(t, componentID, got.ID)
		require.Equal(t, componentTitle, got.Title)
	})

	t.Run("return error if component does not exist", func(t *testing.T) {
		t.Parallel()

		// arrange
		graph := NewComponentGraph()

		// act
		_, err := graph.GetComponent("a.b.c.d")

		// assert
		require.ErrorIs(t, err, ErrComponentNotFound)
	})

	t.Run("return error if component has not been resolved yet", func(t *testing.T) {
		t.Parallel()

		// arrange
		graph := NewComponentGraph()

		component1 := "a.b.c.d"
		graph.AddComponent(component1, "")

		component2 := "a.b.c.e"
		graph.AddLink(component1, component2, "")

		// act
		_, err := graph.GetComponent(component2)

		// assert
		require.ErrorIs(t, err, ErrComponentUnresolved)
	})
}

func TestComponentGraph_GetComponentsInNamespace(t *testing.T) {
	t.Parallel()

	t.Run("successfully get components in namespace", func(t *testing.T) {
		t.Parallel()

		// arrange
		graph := NewComponentGraph()

		component1ID := "a.b.c.d"
		component1Title := "some title"
		graph.AddComponent(component1ID, component1Title)

		component2ID := "a.b.c.e"
		component2Title := "some other title"
		graph.AddComponent(component2ID, component2Title)

		// act
		got := graph.GetComponentsInNamespace("a.b.c")

		// assert
		require.Equal(t, 2, len(got))
		require.Contains(t, got[0].ID, component1ID)
		require.Contains(t, got[1].ID, component2ID)
	})

	t.Run("return empty slice if no components in namespace", func(t *testing.T) {
		t.Parallel()

		// arrange
		graph := NewComponentGraph()

		// act
		got := graph.GetComponentsInNamespace("a.b.c")

		// assert
		require.Equal(t, 0, len(got))
	})
}
