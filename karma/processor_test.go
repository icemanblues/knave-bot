package karma

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAbs(t *testing.T) {
	testcases := []struct {
		name     string
		x        int
		expected int
	}{
		{
			name:     "positive",
			x:        5,
			expected: 5,
		},
		{
			name:     "negative",
			x:        -1,
			expected: 1,
		},
		{
			name:     "zero",
			x:        0,
			expected: 0,
		},
	}

	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			actual := Abs(test.x)
			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestParseArg(t *testing.T) {
	w := []string{"hello", "my", "name", "is", "4"}

	testcases := []struct {
		name     string
		idx      int
		words    []string
		expected string
		ok       bool
	}{
		{"first", 0, w, "hello", true},
		{"middle", 2, w, "name", true},
		{"last", 2, w, "name", true},
		{"negative", -1, w, "", false},
		{"empty", 0, nil, "", false},
		{"over capacity", len(w), w, "", false},
	}

	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			actual, ok := parseArg(test.words, test.idx)
			assert.Equal(t, test.expected, actual)
			assert.Equal(t, test.ok, ok)
		})
	}
}

func TestParseArgInt(t *testing.T) {
	w := []string{"hello", "-1", "4", "", "0"}

	testcases := []struct {
		name       string
		words      []string
		idx        int
		defaultInt int
		expected   int
		ok         bool
	}{
		{
			name:       "negative",
			words:      w,
			idx:        -1,
			defaultInt: 5,
			expected:   5,
			ok:         false,
		},
		{
			name:       "empty",
			words:      nil,
			idx:        0,
			defaultInt: 0,
			expected:   0,
			ok:         false,
		},
		{
			name:     "out of bounds",
			words:    w,
			idx:      len(w) + 5,
			expected: 0,
			ok:       false,
		},
		{
			name:       "happy path",
			words:      w,
			idx:        1,
			defaultInt: 0,
			expected:   -1,
			ok:         true,
		},
		{
			name:       "sad path",
			words:      w,
			idx:        0,
			defaultInt: 0,
			expected:   0,
			ok:         false,
		},
	}

	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			actual, ok := parseArgInt(test.words, test.idx, test.defaultInt)
			assert.Equal(t, test.expected, actual)
			assert.Equal(t, test.ok, ok)
		})
	}
}

func TestParseArgUser(t *testing.T) {
	testcases := []struct {
		name     string
		words    []string
		idx      int
		expected string
		ok       bool
	}{
		{
			name:     "U12345",
			words:    []string{"simon", "U12345"},
			idx:      1,
			expected: "U12345",
			ok:       true,
		},
		{
			name:     "<@U12345>",
			words:    []string{"simon", "<@U12345>"},
			idx:      1,
			expected: "U12345",
			ok:       true,
		},
		{
			name:     "happy path",
			words:    []string{"simon", "U12345"},
			idx:      0,
			expected: "",
			ok:       false,
		},
	}

	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			actual, ok := parseArgUser(test.words, test.idx)
			assert.Equal(t, test.expected, actual)
			assert.Equal(t, test.ok, ok)
		})
	}
}

func TestMe(t *testing.T) {

}

func TestStatus(t *testing.T) {

}

func TestAdd(t *testing.T) {

}

func TestSubtract(t *testing.T) {

}
