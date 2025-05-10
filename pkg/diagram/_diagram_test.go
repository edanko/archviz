package diagram

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"oss.terrastruct.com/d2/lib/log"
)

func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("should create sample diagram", func(t *testing.T) {
		t.Parallel()

		// arrange
		d, err := New()
		require.NoError(t, err)

		// act
		_, err = d.AddNode(
			"a",
			WithShape("square"),
			WithStyle(
				With3D(true),
				WithFill("red"),
			),
		)
		require.NoError(t, err)

		_, err = d.AddNode(
			"b",
			WithShape("stored_data"),
		)
		require.NoError(t, err)

		err = d.AddDirectedConnection(
			"a", "b",
			WithLabel("uses"),
			WithStyle(
				WithStrokeDash(4),
			),
		)
		require.NoError(t, err)

		_, err = d.AddNode(
			"c",
			WithEmptyLabel(),
			WithIcon("https://example.com/test.svg"),
			WithShape("image"),
		)
		require.NoError(t, err)

		_, err = d.AddNode(
			"d",
			WithCode("go", "fmt.Println(\"hello\")"),
		)
		require.NoError(t, err)

		_, err = d.AddNode(
			"e",
			WithEmptyLabel(),
		)
		require.NoError(t, err)

		_, err = d.AddNode(
			"e.explaination",
			WithMarkdown("# hello\n\n**world**"),
		)
		require.NoError(t, err)

		_, err = d.AddNode("f")
		require.NoError(t, err)

		// assert
		want := `a: {
  shape: square
  style.3d: true
  style.fill: red
}
b: {shape: stored_data}
a -> b: uses {style.stroke-dash: 4}
c: "" {
  icon: https://example.com/test.svg
  shape: image
}
d: |go fmt.Println("hello") |
e: "" {
  explaination: |md
    # hello

    **world**
  |
}
f
`

		got, err := d.D2()
		require.NoError(t, err)

		require.Equal(t, want, got)
	})
}

func TestDiagram_Render(t *testing.T) {
	t.Parallel()

	t.Run("should create non-empty svg", func(t *testing.T) {
		t.Parallel()

		// arrange
		d, err := New()
		require.NoError(t, err)

		// act
		_, err = d.AddNode("a")
		require.NoError(t, err)

		_, err = d.AddNode("b")
		require.NoError(t, err)

		err = d.AddSimpleConnection("a", "b")
		require.NoError(t, err)

		ctx := context.Background()
		ctx = log.WithTB(ctx, t)

		// assert
		got, err := d.Render(ctx)
		require.NoError(t, err)

		require.NotEmpty(t, got)
	})
}

func TestDiagram_PNG(t *testing.T) {
	t.Skip("test may download playwright in background, there is no need to run it every time")
	t.Parallel()

	t.Run("should create non-empty png", func(t *testing.T) {
		t.Parallel()

		// arrange
		d, err := New()
		require.NoError(t, err)

		// act
		_, err = d.AddNode("a")
		require.NoError(t, err)

		_, err = d.AddNode("b")
		require.NoError(t, err)

		err = d.AddSimpleConnection("a", "b")
		require.NoError(t, err)

		ctx := context.Background()
		ctx = log.WithTB(ctx, t)

		// assert
		got, err := d.PNG(ctx)
		require.NoError(t, err)

		require.NotEmpty(t, got)
	})
}
