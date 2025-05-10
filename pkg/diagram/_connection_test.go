package diagram

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDiagram_AddSimpleConnection(t *testing.T) {
	t.Parallel()

	t.Run("should add connection", func(t *testing.T) {
		t.Parallel()
		// arrange
		d, err := New()
		require.NoError(t, err)

		// act
		err = d.AddSimpleConnection("a", "b")

		// assert
		require.NoError(t, err)

		got, err := d.D2()
		require.NoError(t, err)

		want := "a -- b\n"

		require.Equal(t, want, got)
	})

}

func TestDiagram_AddDirectedConnection(t *testing.T) {
	t.Run("should add connection", func(t *testing.T) {
		t.Parallel()
		// arrange
		d, err := New()
		require.NoError(t, err)

		// act
		err = d.AddDirectedConnection("a", "b")

		// assert
		require.NoError(t, err)

		got, err := d.D2()
		require.NoError(t, err)

		want := "a -> b\n"

		require.Equal(t, want, got)
	})
}

func TestDiagram_AddBidirectionalConnection(t *testing.T) {
	t.Run("should add connection", func(t *testing.T) {
		t.Parallel()
		// arrange
		d, err := New()
		require.NoError(t, err)

		// act
		err = d.AddBidirectionalConnection("a", "b")

		// assert
		require.NoError(t, err)

		got, err := d.D2()
		require.NoError(t, err)

		want := "a <-> b\n"

		require.Equal(t, want, got)
	})
}

func TestDiagram_addConnection_Options(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		firstKey  string
		secondKey string
		options   []GenericOption
		want      string
	}{
		{
			name:      "with label",
			firstKey:  "a",
			secondKey: "b",
			options: []GenericOption{
				WithLabel("uses"),
			},
			want: "a -- b: uses\n",
		},
		{
			name:      "with empty label",
			firstKey:  "a",
			secondKey: "b",
			options: []GenericOption{
				WithEmptyLabel(),
			},
			want: "a -- b: \"\"\n",
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
			err = d.AddSimpleConnection(tt.firstKey, tt.secondKey, tt.options...)

			// assert
			require.NoError(t, err)

			got, err := d.D2()
			require.NoError(t, err)

			require.Equal(t, tt.want, got)
		})
	}
}

func TestDiagram_AddSimpleConnection_InvalidOptions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		opt  GenericOption
	}{
		{
			name: "with icon",
			opt:  WithIcon("https://example.com/test.svg"),
		},
		{
			name: "with shape",
			opt:  WithShape(NodeShapeSquare),
		},
		{
			name: "with code",
			opt:  WithCode("go", "fmt.Println(\"hello world\")"),
		},
		{
			name: "with markdown",
			opt:  WithMarkdown("# hello world"),
		},
		{
			name: "with link",
			opt:  WithLink("https://example.com"),
		},
		{
			name: "with width",
			opt:  WithWidth(100),
		},
		{
			name: "with label position",
			opt:  WithLabelPosition(TopCenter),
		},
		{
			name: "with icon position",
			opt:  WithIconPosition(TopCenter),
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
			err = d.AddSimpleConnection("a", "b", tt.opt)

			// assert
			require.Error(t, err)
		})
	}
}
