package karma

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/icemanblues/knave-bot/slack"
	"github.com/stretchr/testify/assert"
)

func setup(dao DAO) *gin.Engine {
	proc := mockProcessor(dao)
	h := NewHandler(proc, dao)

	r := gin.Default()

	knaveRouter := r.Group("/knavebot")
	karmaRouter := r.Group("/karmabot")
	BindRoutes(karmaRouter, knaveRouter, h)

	return r
}

func TestGetKarma(t *testing.T) {
	testcases := []struct {
		name     string
		dao      DAO
		code     int
		expected string
	}{
		{
			name:     "GetKarma",
			dao:      HappyDao(),
			code:     200,
			expected: "5",
		},
		{
			name:     "GetKarma error",
			dao:      SadDao(),
			code:     500,
			expected: "GetKarmaMock",
		},
	}

	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			// setup
			r := setup(test.dao)

			// undertest
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/karmabot/v1/team/nycfc/davidvilla", nil)
			r.ServeHTTP(w, req)

			// assert
			assert.Equal(t, test.code, w.Code)
			assert.Equal(t, test.expected, w.Body.String())
		})
	}
}

func TestAddKarma(t *testing.T) {
	testcases := []struct {
		name     string
		dao      DAO
		delta    string
		code     int
		expected string
	}{
		{
			name:     "AddKarma",
			dao:      HappyDao(),
			delta:    "5",
			code:     200,
			expected: "6",
		},
		{
			name:     "AddKarma malformed delta",
			dao:      HappyDao(),
			delta:    "Not-A-Number",
			code:     400,
			expected: "Please pass a valid integer. Not-A-Number",
		},
		{
			name:     "AddKarma negative delta",
			dao:      HappyDao(),
			delta:    "-2",
			code:     200,
			expected: "-1",
		},
		{
			name:     "AddKarma error",
			dao:      SadDao(),
			delta:    "5",
			code:     500,
			expected: "UpdateKarmaMock",
		},
	}

	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			// setup
			r := setup(test.dao)

			// undertest
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("PUT", "/karmabot/v1/team/nycfc/davidvilla?delta="+test.delta, nil)
			r.ServeHTTP(w, req)

			// assert
			assert.Equal(t, test.code, w.Code)
			assert.Equal(t, test.expected, w.Body.String())
		})
	}
}

func TestDelKarma(t *testing.T) {
	testcases := []struct {
		name     string
		dao      DAO
		code     int
		expected string
	}{
		{
			name:     "DeleteKarma",
			dao:      HappyDao(),
			code:     200,
			expected: "0",
		},
		{
			name:     "DeleteKarma error",
			dao:      SadDao(),
			code:     500,
			expected: "DeleteKarmaMock",
		},
	}

	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			// setup
			r := setup(test.dao)

			// undertest
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("DELETE", "/karmabot/v1/team/nycfc/davidvilla", nil)
			r.ServeHTTP(w, req)

			// assert
			assert.Equal(t, test.code, w.Code)
			assert.Equal(t, test.expected, w.Body.String())
		})
	}
}

type KarmaTestCase struct {
	name     string
	dao      DAO
	form     url.Values
	code     int
	expected *slack.Response
}

func makeForm(text string) url.Values {
	return url.Values{
		"text":    []string{text},
		"user_id": []string{"UCALLER"},
		"team_id": []string{"nycfc"},
	}
}

func karmaTestRunner(t *testing.T, testcases []KarmaTestCase) {
	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			// setup
			r := setup(test.dao)

			// undertest
			w := httptest.NewRecorder()
			body := strings.NewReader(test.form.Encode())
			req, _ := http.NewRequest("POST", "/knavebot/v1/cmd/karma", body)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			r.ServeHTTP(w, req)

			var actual slack.Response
			err := json.Unmarshal(w.Body.Bytes(), &actual)

			// assert
			assert.Nil(t, err)
			assert.Equal(t, test.code, w.Code)
			assert.Equal(t, test.expected, &actual)
		})
	}
}

func TestHelp(t *testing.T) {
	testcases := []KarmaTestCase{
		{
			name:     "no text",
			dao:      HappyDao(),
			form:     makeForm(""),
			code:     200,
			expected: responseHelp,
		},
		{
			name:     "help",
			dao:      HappyDao(),
			form:     makeForm("help"),
			code:     200,
			expected: responseHelp,
		},
		{
			name:     "error help",
			dao:      SadDao(),
			form:     makeForm("help"),
			code:     200,
			expected: responseHelp,
		},
	}

	karmaTestRunner(t, testcases)
}

func TestStatus(t *testing.T) {
	testcases := []KarmaTestCase{
		{
			name: "status",
			dao:  HappyDao(),
			form: makeForm("status <@USER>"),
			code: 200,
			expected: slack.ChannelAttachmentsResponse(
				"<@UCALLER> has requested karma total for <@USER>. <@USER> has 5 karma.",
				"compliment"),
		},
		{
			name:     "error status",
			dao:      SadDao(),
			form:     makeForm("status <@USER>"),
			code:     200,
			expected: responseUnknownError,
		},
	}

	karmaTestRunner(t, testcases)
}

func TestMe(t *testing.T) {
	testcases := []KarmaTestCase{
		{
			name: "me",
			dao:  HappyDao(),
			form: makeForm("me"),
			code: 200,
			expected: slack.DirectResponse(
				"<@UCALLER> has 5 karma.",
				"compliment"),
		},
		{
			name:     "error me",
			dao:      SadDao(),
			form:     makeForm("me"),
			code:     200,
			expected: responseUnknownError,
		},
	}

	karmaTestRunner(t, testcases)
}

func TestAdd(t *testing.T) {
	testcases := []KarmaTestCase{
		{
			name: "add",
			dao:  HappyDao(),
			form: makeForm("++ <@USER>"),
			code: 200,
			expected: slack.ChannelAttachmentsResponse(
				"<@UCALLER> is giving 1 karma to <@USER>. <@USER> has 2 karma.",
				"compliment"),
		},
		{
			name: "add 2",
			dao:  HappyDao(),
			form: makeForm("++ <@USER> 2"),
			code: 200,
			expected: slack.ChannelAttachmentsResponse(
				"<@UCALLER> is giving 2 karma to <@USER>. <@USER> has 3 karma.",
				"compliment"),
		},
		{
			name:     "error add",
			dao:      SadDao(),
			form:     makeForm("++ <@USER>"),
			code:     200,
			expected: responseUnknownError,
		},
	}

	karmaTestRunner(t, testcases)
}

func TestSubtract(t *testing.T) {
	testcases := []KarmaTestCase{
		{
			name: "subtract",
			dao:  HappyDao(),
			form: makeForm("-- <@USER>"),
			code: 200,
			expected: slack.ChannelAttachmentsResponse(
				"<@UCALLER> is taking away 1 karma from <@USER>. <@USER> has 0 karma.",
				"insult"),
		},
		{
			name: "subtract 3",
			dao:  HappyDao(),
			form: makeForm("-- <@USER> 3"),
			code: 200,
			expected: slack.ChannelAttachmentsResponse(
				"<@UCALLER> is taking away 3 karma from <@USER>. <@USER> has -2 karma.",
				"insult"),
		},
		{
			name:     "error subtract",
			dao:      SadDao(),
			form:     makeForm("-- <@USER>"),
			code:     200,
			expected: responseUnknownError,
		},
	}

	karmaTestRunner(t, testcases)
}
