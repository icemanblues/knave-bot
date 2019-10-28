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

func TestUserCmdAlias(t *testing.T) {
	testcases := []struct {
		name     string
		args     []string
		expected []string
	}{
		{
			name:     "nil",
			args:     nil,
			expected: nil,
		},
		{
			name:     "empty",
			args:     []string{},
			expected: []string{},
		},
		{
			name:     "one",
			args:     []string{"one"},
			expected: []string{"one"},
		},
		{
			name:     "happy",
			args:     []string{"USER", "++"},
			expected: []string{"++", "USER"},
		},
		{
			name:     "happy msg",
			args:     []string{"USER", "++", "you", "go", "girl!"},
			expected: []string{"++", "USER", "you", "go", "girl!"},
		},
		{
			name:     "no change",
			args:     []string{"status", "USER", "karma,", "baby!"},
			expected: []string{"status", "USER", "karma,", "baby!"},
		},
	}

	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			actual := userCmdAlias(test.args)
			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestAddSubCmdAlias(t *testing.T) {
	testcases := []struct {
		name     string
		args     []string
		expected []string
	}{
		{
			name:     "nil",
			args:     nil,
			expected: nil,
		},
		{
			name:     "empty",
			args:     []string{},
			expected: []string{},
		},
		{
			name:     "one",
			args:     []string{"one"},
			expected: []string{"one"},
		},
		{
			name:     "plus 3",
			args:     []string{"+3", "USER"},
			expected: []string{"++", "USER", "3"},
		},
		{
			name:     "minus 3",
			args:     []string{"-3", "USER"},
			expected: []string{"--", "USER", "3"},
		},
		{
			name:     "plus 5 msg",
			args:     []string{"+5", "USER", "you", "go", "girl!"},
			expected: []string{"++", "USER", "5", "you", "go", "girl!"},
		},
		{
			name:     "minus 2 msg",
			args:     []string{"-2", "USER", "you", "go", "girl!"},
			expected: []string{"--", "USER", "2", "you", "go", "girl!"},
		},
		{
			name:     "no change",
			args:     []string{"++", "USER", "2", "go", "girl!"},
			expected: []string{"++", "USER", "2", "go", "girl!"},
		},
		{
			name:     "four (no plus)",
			args:     []string{"4", "USER"},
			expected: []string{"++", "USER", "4"},
		},
		{
			name:     "zero",
			args:     []string{"0", "USER"},
			expected: []string{"0", "USER"},
		},
	}

	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			actual := addSubCmdAlias((test.args))
			assert.Equal(t, test.expected, actual)
		})
	}
}

type ProcessTestCase struct {
	name         string
	command      *slack.CommandData
	responseType string
	text         string
	attach       bool
}

func processHelper(t *testing.T, p Processor, test ProcessTestCase) {
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

func TestProcessStatus(t *testing.T) {
	p := mockProcessor(HappyDao(), HappyDailyDao())
	testcases := []ProcessTestCase{
		{
			name:         "status",
			command:      command("status <@USER>"),
			responseType: slack.ResponseType.InChannel,
			text:         "<@UCALLER> has requested karma total for <@USER>. <@USER> has 5 karma.",
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
	}

	for _, test := range testcases {
		processHelper(t, p, test)
	}
}

func TestProcessTop(t *testing.T) {
	p := mockProcessor(HappyDao())
	testcases := []ProcessTestCase{
		{
			name:         "top",
			command:      command("top"),
			responseType: slack.ResponseType.InChannel,
			text:         "The top 3 users by karma:\nRank\tName\tKarma\n1\t<@USER0>\t100\n2\t<@USER1>\t101\n3\t<@USER2>\t102\n",
		},
		{
			name:         "top 5",
			command:      command("top 5"),
			responseType: slack.ResponseType.InChannel,
			text:         "The top 5 users by karma:\nRank\tName\tKarma\n1\t<@USER0>\t100\n2\t<@USER1>\t101\n3\t<@USER2>\t102\n4\t<@USER3>\t103\n5\t<@USER4>\t104\n",
		},
		{
			name:         "top negative",
			command:      command("top -5"),
			responseType: slack.ResponseType.InChannel,
			text:         "The top 3 users by karma:\nRank\tName\tKarma\n1\t<@USER0>\t100\n2\t<@USER1>\t101\n3\t<@USER2>\t102\n",
		},
		{
			name:         "top over max",
			command:      command("top 100"),
			responseType: slack.ResponseType.InChannel,
			text:         "The top 10 users by karma:\nRank\tName\tKarma\n1\t<@USER0>\t100\n2\t<@USER1>\t101\n3\t<@USER2>\t102\n4\t<@USER3>\t103\n5\t<@USER4>\t104\n6\t<@USER5>\t105\n7\t<@USER6>\t106\n8\t<@USER7>\t107\n9\t<@USER8>\t108\n10\t<@USER9>\t109\n",
		},
	}

	for _, test := range testcases {
		processHelper(t, p, test)
	}
}

func TestProcessMe(t *testing.T) {
	p := mockProcessor(HappyDao())
	testcases := []ProcessTestCase{
		{
			name:         "me",
			command:      command("me"),
			responseType: slack.ResponseType.Ephemeral,
			text:         "<@UCALLER> has 5 karma.",
		},
	}

	for _, test := range testcases {
		processHelper(t, p, test)
	}
}

func TestProcessAdd(t *testing.T) {
	p := mockProcessor(HappyDao())
	testcases := []ProcessTestCase{
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
	}

	for _, test := range testcases {
		processHelper(t, p, test)
	}
}

func TestProcessSubtract(t *testing.T) {
	p := mockProcessor(HappyDao())
	testcases := []ProcessTestCase{
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
	}

	for _, test := range testcases {
		processHelper(t, p, test)
	}
}

func TestProcessHelp(t *testing.T) {
	p := mockProcessor(HappyDao())
	testcases := []ProcessTestCase{
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
		}}

	for _, test := range testcases {
		processHelper(t, p, test)
	}
}

func TestProcessAliasPlusPlus(t *testing.T) {
	p := mockProcessor(HappyDao())
	testcases := []ProcessTestCase{
		{
			name:         "USER ++",
			command:      command("<@USER> ++"),
			responseType: slack.ResponseType.InChannel,
			text:         "<@UCALLER> is giving 1 karma to <@USER>. <@USER> has 2 karma.",
		},
		{
			name:         "USER ++ quantity",
			command:      command("<@USER> ++ 3"),
			responseType: slack.ResponseType.InChannel,
			text:         "<@UCALLER> is giving 3 karma to <@USER>. <@USER> has 4 karma.",
		},
	}

	for _, test := range testcases {
		processHelper(t, p, test)
	}
}

func TestProcessAliasMinusMinus(t *testing.T) {
	p := mockProcessor(HappyDao())
	testcases := []ProcessTestCase{
		{
			name:         "USER -- quantity",
			command:      command("<@USER> -- 3"),
			responseType: slack.ResponseType.InChannel,
			text:         "<@UCALLER> is taking away 3 karma from <@USER>. <@USER> has -2 karma.",
		},
		{
			name:         "USER -- quantity out-of-bounds",
			command:      command("<@USER> -- 9000"),
			responseType: slack.ResponseType.Ephemeral,
			text:         msgDeltaLimit,
		},
		{
			name:         "USER -- quantity message",
			command:      command("<@USER> -- 2 be better next time"),
			responseType: slack.ResponseType.InChannel,
			text:         "<@UCALLER> is taking away 2 karma from <@USER>. <@USER> has -1 karma.",
		},
	}

	for _, test := range testcases {
		processHelper(t, p, test)
	}
}

func TestProcessAliasAddSubNumber(t *testing.T) {
	p := mockProcessor(HappyDao())
	testcases := []ProcessTestCase{
		{
			name:         "USER three",
			command:      command("<@USER> +3"),
			responseType: slack.ResponseType.InChannel,
			text:         "<@UCALLER> is giving 3 karma to <@USER>. <@USER> has 4 karma.",
		},
		{
			name:         "USER three (no plus)",
			command:      command("<@USER> 3"),
			responseType: slack.ResponseType.InChannel,
			text:         "<@UCALLER> is giving 3 karma to <@USER>. <@USER> has 4 karma.",
		},
		{
			name:         "three USER",
			command:      command("3 <@USER>"),
			responseType: slack.ResponseType.InChannel,
			text:         "<@UCALLER> is giving 3 karma to <@USER>. <@USER> has 4 karma.",
		},
		{
			name:         "USER minus three",
			command:      command("<@USER> -3"),
			responseType: slack.ResponseType.InChannel,
			text:         "<@UCALLER> is taking away 3 karma from <@USER>. <@USER> has -2 karma.",
		},
		{
			name:         "minus three USER",
			command:      command("-3 <@USER>"),
			responseType: slack.ResponseType.InChannel,
			text:         "<@UCALLER> is taking away 3 karma from <@USER>. <@USER> has -2 karma.",
		},
		{
			name:         "USER 0",
			command:      command("<@USER> 0"),
			responseType: responseHelp.ResponseType,
			text:         responseHelp.Text,
		},
		{
			name:         "0 USER",
			command:      command("0 <@USER>"),
			responseType: responseHelp.ResponseType,
			text:         responseHelp.Text,
		},
	}

	for _, test := range testcases {
		processHelper(t, p, test)
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
		p := mockProcessor(SadDao())
		t.Run(test.name, func(t *testing.T) {
			actual, err := p.Process(test.command)
			assert.Nil(t, actual)
			assert.NotNil(t, err)
		})
	}
}
