package karma

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/icemanblues/knave-bot/slack"
	"github.com/stretchr/testify/assert"
)

func setup(h *SQLiteHandler) *gin.Engine {
	r := gin.Default()

	r.GET("/karma/:team/:user", h.GetKarma)
	r.PUT("/karma/:team/:user", h.AddKarma)
	r.DELETE("/karma/:team/:user", h.DelKarma)
	r.POST("/slash/karma", h.SlashKarma)

	return r
}

// MockDAO a mock dao for karma whose mock functions can be monkeypatched
type MockDAO struct {
	GetKarmaMock    func(team, user string) (int, error)
	UpdateKarmaMock func(team, user string, delta int) (int, error)
	DeleteKarmaMock func(team, user string) (int, error)
}

func (m MockDAO) GetKarma(team, user string) (int, error) {
	return m.GetKarmaMock(team, user)
}

func (m MockDAO) UpdateKarma(team, user string, delta int) (int, error) {
	return m.UpdateKarmaMock(team, user, delta)
}

func (m MockDAO) DeleteKarma(team, user string) (int, error) {
	return m.DeleteKarmaMock(team, user)
}

func happyDao() *MockDAO {
	return &MockDAO{
		GetKarmaMock: func(team, user string) (int, error) {
			return 10, nil
		},
		UpdateKarmaMock: func(team, user string, delta int) (int, error) {
			return delta + 1, nil
		},
		DeleteKarmaMock: func(team, user string) (int, error) {
			return 0, nil
		},
	}
}

func sadDao() *MockDAO {
	return &MockDAO{
		GetKarmaMock: func(team, user string) (int, error) {
			return 0, errors.New("GetKarmaMock")
		},
		UpdateKarmaMock: func(team, user string, delta int) (int, error) {
			return 0, errors.New("UpdateKarmaMock")
		},
		DeleteKarmaMock: func(team, user string) (int, error) {
			return 0, errors.New("DeleteKarmaMock")
		},
	}
}

func TestSlashKarma(t *testing.T) {
	dao := happyDao()
	proc := NewProcessor(dao)
	h := NewHandler(proc, dao)
	r := setup(h)

	testcases := []struct {
		name     string
		form     string
		code     int
		expected *slack.Response
	}{
		{
			name:     "help",
			form:     "",
			code:     200,
			expected: responseHelp,
		},
	}

	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			form := strings.NewReader(test.form)
			req, _ := http.NewRequest("POST", "/slash/karma", form)
			r.ServeHTTP(w, req)

			var actual slack.Response
			err := json.Unmarshal(w.Body.Bytes(), &actual)
			assert.Nil(t, err)

			assert.Equal(t, 200, w.Code)
			assert.Equal(t, test.expected, &actual)
		})
	}
}

func TestSlashKarmaError(t *testing.T) {
	dao := sadDao()
	proc := NewProcessor(dao)
	h := NewHandler(proc, dao)
	r := setup(h)

	testcases := []struct {
		name     string
		form     string
		code     int
		expected *slack.Response
	}{
		{
			name:     "error status",
			form:     "text=me&user_id=UCALLER&team_id=nycfc",
			code:     200,
			expected: responseUnknownError,
		},
	}

	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			form := strings.NewReader(test.form)
			req, _ := http.NewRequest("POST", "/slash/karma", form)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			r.ServeHTTP(w, req)

			var actual slack.Response
			err := json.Unmarshal(w.Body.Bytes(), &actual)
			assert.Nil(t, err)

			assert.Equal(t, 200, w.Code)
			assert.Equal(t, test.expected, &actual)
		})
	}
}

func TestGetKarma(t *testing.T) {
	dao := happyDao()
	proc := NewProcessor(dao)
	h := NewHandler(proc, dao)
	r := setup(h)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/karma/nycfc/davidvilla", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "10", w.Body.String())
}
