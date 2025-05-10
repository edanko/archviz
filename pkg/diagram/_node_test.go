package diagram

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDiagram_AddNode(t *testing.T) {
	t.Parallel()

	t.Run("successfully add node", func(t *testing.T) {
		t.Parallel()
		// arrange
		d, err := New()
		require.NoError(t, err)

		// act
		_, err = d.AddNode("a")

		// assert
		require.NoError(t, err)

		got, err := d.D2()
		require.NoError(t, err)

		want := "a\n"

		require.Equal(t, want, got)
	})

	t.Run("return error on duplicate node key", func(t *testing.T) {
		t.Parallel()
		// arrange
		d, err := New()
		require.NoError(t, err)

		_, err = d.AddNode("a")
		require.NoError(t, err)

		// act
		_, err = d.AddNode("a")

		// assert
		require.Error(t, err)
	})
}

func TestDiagram_AddNode_Options(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		nodeKey string
		options []GenericOption
		want    string
	}{
		{
			name:    "with label",
			nodeKey: "a",
			options: []GenericOption{
				WithLabel("some title"),
			},
			want: "a: some title\n",
		},
		{
			name:    "with empty label",
			nodeKey: "a",
			options: []GenericOption{
				WithEmptyLabel(),
			},
			want: "a: \"\"\n",
		},
		{
			name:    "with icon",
			nodeKey: "a",
			options: []GenericOption{
				WithIcon("https://example.com/test.svg"),
			},
			want: "a: {icon: https://example.com/test.svg}\n",
		},
		{
			name:    "with shape",
			nodeKey: "a",
			options: []GenericOption{
				WithShape(NodeShapeSquare),
			},
			want: "a: {shape: square}\n",
		},
		{
			name:    "with code",
			nodeKey: "a",
			options: []GenericOption{
				WithCode("go", "fmt.Println(\"hello world\")"),
			},
			want: "a: |go fmt.Println(\"hello world\") |\n",
		},
		{
			name:    "with markdown",
			nodeKey: "a",
			options: []GenericOption{
				WithMarkdown("# hello world"),
			},
			want: "a: |md # hello world |\n",
		},
		{
			name:    "with link",
			nodeKey: "a",
			options: []GenericOption{
				WithLink("https://example.com"),
			},
			want: "a: {link: https://example.com}\n",
		},
		{
			name:    "with width",
			nodeKey: "a",
			options: []GenericOption{
				WithWidth(100),
			},
			want: "a: {width: 100}\n",
		},
		{
			name:    "with label position",
			nodeKey: "a",
			options: []GenericOption{
				WithLabelPosition(TopCenter),
			},
			want: "a: {label.near: top-center}\n",
		},
		{
			name:    "with icon position",
			nodeKey: "a",
			options: []GenericOption{
				WithIconPosition(TopCenter),
			},
			want: "a: {icon.near: top-center}\n",
		},
		{
			name:    "multiple options",
			nodeKey: "a",
			options: []GenericOption{
				WithLabel("some title"),
				WithIcon("https://example.com/test.svg"),
				WithIconPosition(TopCenter),
				WithLink("https://example.com"),
			},
			want: `a: some title {
  icon: https://example.com/test.svg
  icon.near: top-center
  link: https://example.com
}
`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// arrange
			d, err := New()
			require.NoError(t, err)

			// act
			_, err = d.AddNode(tt.nodeKey, tt.options...)

			// assert
			require.NoError(t, err)

			got, err := d.D2()
			require.NoError(t, err)

			require.Equal(t, tt.want, got)
		})
	}
}
