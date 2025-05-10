package app

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWordWrap(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		text       string
		lineLength int
		want       string
	}{
		{
			name:       "empty input",
			text:       "",
			lineLength: 10,
			want:       "\n",
		},
		{
			name:       "single line under limit",
			text:       "hello",
			lineLength: 10,
			want:       "hello\n",
		},
		{
			name:       "single line exactly limit",
			text:       "1234567890",
			lineLength: 10,
			want:       "1234567890\n",
		},
		{
			name:       "single line needs wrapping",
			text:       "one two three",
			lineLength: 7,
			want:       "one two \\\nthree\n",
		},
		{
			name:       "word exceeds line length",
			text:       "longword",
			lineLength: 5,
			want:       "longword\n",
		},
		{
			name:       "multiple lines",
			text:       "hello\nworld",
			lineLength: 10,
			want:       "hello\n\nworld\n",
		},
		{
			name:       "line with leading and trailing spaces",
			text:       "   hello   \n   world   ",
			lineLength: 10,
			want:       "hello\nworld\n",
		},
		{
			name:       "empty line",
			text:       "\n",
			lineLength: 5,
			want:       "\n\n",
		},
		{
			name:       "whitespace only line",
			text:       "   \t   \n",
			lineLength: 5,
			want:       "\n",
		},
		{
			name:       "mixed lines with wrap and no wrap",
			text:       "short\nthis is a longer line that needs wrapping",
			lineLength: 10,
			want:       "short\n\nthis is a \\\nlonger \\\nline that \\\nneeds \\\nwrapping\n",
		},
		{
			name:       "exact fit after space",
			text:       "1234 5",
			lineLength: 5,
			want:       "1234 \\\n5\n",
		},
		{
			name:       "multiple splits required",
			text:       "a bb ccc dddd eeeee",
			lineLength: 5,
			want:       "a bb \\\nccc \\\ndddd \\\neeeee\n",
		},
		{
			name:       "minimum line length",
			text:       "a b c d e",
			lineLength: 1,
			want:       "a \\\nb \\\nc \\\nd \\\ne\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// act
			got := WordWrap(tt.text, tt.lineLength)

			// assert
			require.Equal(t, tt.want, got)
		})
	}
}
