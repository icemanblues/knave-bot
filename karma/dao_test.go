package karma

import (
	"testing"
	"time"

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

func TestIsoDate(t *testing.T) {
	testcases := []struct {
		name     string
		date     time.Time
		expected string
	}{
		{
			name:     "November 9th",
			date:     time.Date(2019, time.November, 9, 0, 0, 0, 0, time.Local),
			expected: "2019-11-09",
		},
		{
			name:     "New Year",
			date:     time.Date(2020, time.January, 1, 0, 0, 0, 0, time.Local),
			expected: "2020-01-01",
		},
	}

	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			actual := IsoDate(test.date)
			assert.Equal(t, test.expected, actual)
		})
	}
}
