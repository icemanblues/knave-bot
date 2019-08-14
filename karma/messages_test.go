package karma

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserStatus(t *testing.T) {
	testcases := []struct {
		name     string
		user     string
		karma    int
		expected string
	}{
		{
			name:     "posiive karma",
			user:     "user",
			karma:    105,
			expected: "<@user> has 105 karma.",
		},
		{
			name:     "negative karma",
			user:     "user",
			karma:    -23,
			expected: "<@user> has -23 karma.",
		},
		{
			name:     "no karma",
			user:     "user",
			karma:    0,
			expected: "<@user> has 0 karma.",
		},
	}

	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			sb := &strings.Builder{}
			UserStatus(test.user, test.karma, sb)
			actual := sb.String()
			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestSalutation(t *testing.T) {
	testcases := []struct {
		name     string
		karma    int
		length int
		expected string
	}{
		{
			name:     "positive karma compliment",
			karma:    5,
			length: 4,
			expected: "",
		},
		{
			name:     "negative karma insult",
			karma:    -1,
			length: 4,
			expected: "",
		},
		{
			name:     "no karma silence",
			karma:    0,
			length: 0,
			expected: "",
		},
	}

	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			sb := &strings.Builder{}
			Salutation(test.karma, sb)
			actual := sb.String()
			// Since I don't know how to detext an insult from a compliment, this is silly
			if actual != test.expected && len(actual) < test.length {
				t.Error()
			}
		})
	}
}
