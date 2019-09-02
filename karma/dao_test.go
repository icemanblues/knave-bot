package karma

import (
	"testing"

	"github.com/icemanblues/knave-bot/slack"
	"github.com/stretchr/testify/assert"
)

func TestStringAttachment(t *testing.T) {
	var defaultAtt slack.Attachments
	sd := "[{}]"

	att := slack.Attachments{
		Text: "text",
	}
	sAtt := "[{\"text\":\"text\"}]"

	testcases := []struct {
		name     string
		arg      []slack.Attachments
		expected *string
	}{
		{
			name:     "nil",
			arg:      nil,
			expected: nil,
		},
		{
			name:     "empty",
			arg:      []slack.Attachments{},
			expected: nil,
		},
		{
			name:     "default",
			arg:      []slack.Attachments{defaultAtt},
			expected: &sd,
		},
		{
			name:     "single",
			arg:      []slack.Attachments{att},
			expected: &sAtt,
		},
	}

	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			actual := stringAttachment(test.arg)

			// testify.assert gets confused with string pointers
			if test.expected == nil {
				assert.Nil(t, actual)
			} else {
				assert.NotNil(t, actual)
				assert.Equal(t, *test.expected, *actual)
			}
		})
	}
}
