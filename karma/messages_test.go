package karma

import (
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
			name:     "positive karma",
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
			actual := MsgUserStatus(test.user, test.karma)
			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestSalutation(t *testing.T) {
	testcases := []struct {
		name     string
		karma    int
		expected string
	}{
		{
			name:     "positive karma compliment",
			karma:    5,
			expected: "compliment",
		},
		{
			name:     "negative karma insult",
			karma:    -1,
			expected: "insult",
		},
		{
			name:     "no karma silence",
			karma:    0,
			expected: "",
		},
	}

	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			p := happyMockProcessor()
			actual := p.Salutation(test.karma)
			assert.Equal(t, test.expected, actual)
		})
	}
}
