package karma

import (
	"testing"

	"github.com/icemanblues/knave-bot/slack"
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

func command(text string) *slack.CommandData {
	return &slack.CommandData{
		Command: "karma",
		UserID:  "UCALLER",
		Text:    text,
	}
}

func TestProcess(t *testing.T) {
	testcases := []struct {
		name         string
		command      *slack.CommandData
		responseType string
		text         string
		attach       bool
	}{
		// STATUS
		{
			name:         "status",
			command:      command("status <@USER>"),
			responseType: slack.ResponseType.InChannel,
			text:         "<@UCALLER> has requested karma total for <@USER>.<@USER> has 5 karma.",
		},
		{
			name:         "status no user",
			command:      command("status"),
			responseType: slack.ResponseType.Ephemeral,
			text:         msgMissingName,
		},
		{
			name:         "status malformed user",
			command:      command("status blah"),
			responseType: slack.ResponseType.Ephemeral,
			text:         msgInvalidUser,
		},
		// ME
		{
			name:         "me",
			command:      command("me"),
			responseType: slack.ResponseType.Ephemeral,
			text:         "<@UCALLER> has 5 karma.",
		},
		// ADD
		{
			name:         "++",
			command:      command("++ <@USER>"),
			responseType: slack.ResponseType.InChannel,
			text:         "<@UCALLER> is giving 1 karma to <@USER>. <@USER> has 2 karma.",
		},
		{
			name:         "++ quantity",
			command:      command("++ <@USER> 3"),
			responseType: slack.ResponseType.InChannel,
			text:         "<@UCALLER> is giving 3 karma to <@USER>. <@USER> has 4 karma.",
		},
		{
			name:         "++ quantity out-of-bounds",
			command:      command("++ <@USER> 9000"),
			responseType: slack.ResponseType.Ephemeral,
			text:         msgDeltaLimit,
		},
		{
			name:         "++ quantity message",
			command:      command("++ <@USER> thanks you so much"),
			responseType: slack.ResponseType.InChannel,
			text:         "<@UCALLER> is giving 1 karma to <@USER>. <@USER> has 2 karma.",
		},
		{
			name:         "++ quantity negative",
			command:      command("++ <@USER> -2"),
			responseType: slack.ResponseType.Ephemeral,
			text:         msgAddCantRemove,
		},
		{
			name:         "++ quantity zero",
			command:      command("++ <@USER> 0"),
			responseType: slack.ResponseType.Ephemeral,
			text:         msgNoOp,
		},
		{
			name:         "++ self target",
			command:      command("++ <@UCALLER> 5"),
			responseType: slack.ResponseType.Ephemeral,
			text:         msgAddSelfTarget,
		},
		{
			name:         "++ missing target",
			command:      command("++"),
			responseType: slack.ResponseType.Ephemeral,
			text:         msgAddMissingTarget,
		},
		{
			name:         "++ malformed target",
			command:      command("++ yikes"),
			responseType: slack.ResponseType.Ephemeral,
			text:         msgInvalidUser,
		},
		// SUBTRACT
		{
			name:         "--",
			command:      command("-- <@USER>"),
			responseType: slack.ResponseType.InChannel,
			text:         "<@UCALLER> is taking away 1 karma from <@USER>. <@USER> has 0 karma.",
		},
		{
			name:         "-- quantity",
			command:      command("-- <@USER> 3"),
			responseType: slack.ResponseType.InChannel,
			text:         "<@UCALLER> is taking away 3 karma from <@USER>. <@USER> has -2 karma.",
		},
		{
			name:         "-- quantity out-of-bounds",
			command:      command("-- <@USER> 9000"),
			responseType: slack.ResponseType.Ephemeral,
			text:         msgDeltaLimit,
		},
		{
			name:         "-- quantity message",
			command:      command("-- <@USER> be better next time"),
			responseType: slack.ResponseType.InChannel,
			text:         "<@UCALLER> is taking away 1 karma from <@USER>. <@USER> has 0 karma.",
		},
		{
			name:         "-- quantity negative",
			command:      command("-- <@USER> -1"),
			responseType: slack.ResponseType.Ephemeral,
			text:         msgSubtractCantAdd,
		},
		{
			name:         "-- quantity zero",
			command:      command("-- <@USER> 0"),
			responseType: slack.ResponseType.Ephemeral,
			text:         msgNoOp,
		},
		{
			name:         "-- self target",
			command:      command("-- <@UCALLER>"),
			responseType: slack.ResponseType.Ephemeral,
			text:         msgSubtractSelfTarget,
		},
		{
			name:         "-- missing target",
			command:      command("--"),
			responseType: slack.ResponseType.Ephemeral,
			text:         msgSubtractMissingTarget,
		},
		{
			name:         "-- malformed target",
			command:      command("-- yikes"),
			responseType: slack.ResponseType.Ephemeral,
			text:         msgInvalidUser,
		},
		// HELP
		{
			name:         "help",
			command:      command("help"),
			responseType: slack.ResponseType.Ephemeral,
			text:         responseHelp.Text,
		},
		{
			name:         "help extra text",
			command:      command("help extra text"),
			responseType: slack.ResponseType.Ephemeral,
			text:         responseHelp.Text,
		},
		{
			name:         "help empty",
			command:      command(""),
			responseType: slack.ResponseType.Ephemeral,
			text:         responseHelp.Text,
		},
	}

	for _, test := range testcases {
		p := NewProcessor(HappyDao())

		t.Run(test.name, func(t *testing.T) {
			actual, err := p.Process(test.command)
			assert.Nil(t, err)
			assert.NotNil(t, actual)

			assert.Equal(t, test.responseType, actual.ResponseType)
			assert.Equal(t, test.text, actual.Text)

			if test.attach {
				assert.NotNil(t, actual.Attachments)
				assert.Len(t, actual.Attachments, 1)
				assert.NotEmpty(t, actual.Attachments[0].Text)
			}
		})
	}
}

func TestProcessError(t *testing.T) {
	testcases := []struct {
		name    string
		command *slack.CommandData
	}{
		{
			name: "status",
			command: &slack.CommandData{
				Text:   "status <@USER>",
				UserID: "UCALLER",
			},
		},
	}

	for _, test := range testcases {
		p := NewProcessor(SadDao())
		t.Run(test.name, func(t *testing.T) {
			actual, err := p.Process(test.command)
			assert.Nil(t, actual)
			assert.NotNil(t, err)
		})
	}
}
