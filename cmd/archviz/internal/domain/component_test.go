package domain

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_newComponent(t *testing.T) {
	t.Parallel()

	t.Run("successfully create component", func(t *testing.T) {
		t.Parallel()

		// arrange
		id := "a.b.c.d"
		title := "some title"

		// act
		got := newComponent(id, title)

		// assert
		want := &Component{
			ID:    id,
			Title: title,
			Depth: 4,
		}
		require.Equal(t, want, got)
	})

	t.Run("successfully create component with type", func(t *testing.T) {
		t.Parallel()

		// arrange
		id := "a.b.c.d"
		title := "some title"
		componentType := TypeService

		// act
		got := newComponent(id, title, WithType(componentType))

		// assert
		want := &Component{
			ID:    id,
			Title: title,
			Type:  componentType,
			Depth: 4,
		}
		require.Equal(t, want, got)
	})

	t.Run("successfully create component with tags", func(t *testing.T) {
		t.Parallel()

		// arrange
		id := "a.b.c.d"
		title := "some title"
		tags := map[string]string{
			"tag1": "value1",
			"tag2": "value2",
		}

		// act
		got := newComponent(id, title, WithTags(tags))

		// assert
		want := &Component{
			ID:    id,
			Title: title,
			Tags:  tags,
			Depth: 4,
		}
		require.Equal(t, want, got)
	})

	t.Run("successfully create component with description", func(t *testing.T) {
		t.Parallel()

		// arrange
		id := "a.b.c.d"
		title := "some title"
		description := "some description"

		// act
		got := newComponent(id, title, WithDescription(description))

		// assert
		want := &Component{
			ID:          id,
			Title:       title,
			Description: description,
			Depth:       4,
		}
		require.Equal(t, want, got)
	})

	t.Run("successfully create component with properties", func(t *testing.T) {
		t.Parallel()

		// arrange
		id := "a.b.c.d"
		title := "some title"
		properties := map[string]any{
			"property1": "value1",
			"property2": "value2",
		}

		// act
		got := newComponent(id, title, WithProperties(properties))

		// assert
		want := &Component{
			ID:         id,
			Title:      title,
			Properties: properties,
			Depth:      4,
		}
		require.Equal(t, want, got)
	})
}

func TestComponent_getNamespace(t *testing.T) {
	t.Parallel()

	t.Run("should get namespace at depth", func(t *testing.T) {
		t.Parallel()

		// arrange
		component := &Component{
			ID: "a.b.c.d",
		}

		// act
		got := component.getNamespace(2)

		// assert
		want := "a.b"
		require.Equal(t, want, got)
	})

	t.Run("should return actual ID if depth is negative", func(t *testing.T) {
		t.Parallel()

		// arrange
		component := &Component{
			ID: "a.b.c.d",
		}

		// act
		got := component.getNamespace(-1)

		// assert
		want := component.ID
		require.Equal(t, want, got)
	})

	t.Run("should return empty string if depth is zero", func(t *testing.T) {
		t.Parallel()

		// arrange
		component := &Component{
			ID: "a.b.c.d",
		}

		// act
		got := component.getNamespace(0)

		// assert
		want := ""
		require.Equal(t, want, got)
	})

	t.Run("should return empty string if depth is greater than depth", func(t *testing.T) {
		t.Parallel()

		// arrange
		component := &Component{
			ID: "a.b.c.d",
		}

		// act
		got := component.getNamespace(5)

		// assert
		want := component.ID
		require.Equal(t, want, got)
	})
}
