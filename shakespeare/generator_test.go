package shakespeare

import "testing"
import "github.com/stretchr/testify/assert"

func TestGenerator(t *testing.T) {
	colA := []string{"a1"}
	colB := []string{"b1", "b1"}
	colC := []string{"c1", "c1", "c1"}

	testcases := []struct {
		name     string
		pre      string
		post     string
		cols     [][]string
		expected string
	}{
		{
			name:     "empty",
			pre:      "",
			post:     "",
			cols:     nil,
			expected: "",
		},
		{
			name:     "pre",
			pre:      "pre",
			post:     "",
			cols:     nil,
			expected: "pre",
		},
		{
			name:     "post",
			pre:      "",
			post:     "post",
			cols:     nil,
			expected: "post",
		},
		{
			name:     "column",
			pre:      "",
			post:     "",
			cols:     [][]string{colA},
			expected: "a1",
		},
		{
			name:     "columns",
			pre:      "",
			post:     "",
			cols:     [][]string{colA, colB, colC},
			expected: "a1 b1 c1",
		},
		{
			name:     "no columns",
			pre:      "pre",
			post:     "fix",
			cols:     nil,
			expected: "pre fix",
		},
		{
			name:     "no prefix",
			pre:      "",
			post:     "End",
			cols:     [][]string{colA},
			expected: "a1 End",
		},
		{
			name:     "no postfix",
			pre:      "Pre",
			post:     "",
			cols:     [][]string{colB},
			expected: "Pre b1",
		},
		{
			name:     "all",
			pre:      "Start",
			post:     "End",
			cols:     [][]string{colA, colB, colC},
			expected: "Start a1 b1 c1 End",
		},
	}

	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			gen := New(test.pre, test.post, test.cols)
			actual := gen.Sentence()
			assert.Equal(t, test.expected, actual)
		})
	}
}
