package karma

import (
	"errors"
	"fmt"
	"github.com/icemanblues/knave-bot/slack"
	"time"
)

// MockDAO a mock dao for karma whose mock functions can be monkeypatched
type MockDAO struct {
	GetKarmaMock         func(team, user string) (int, error)
	UpdateKarmaMock      func(team, user string, delta int) (int, error)
	DeleteKarmaMock      func(team, user string) (int, error)
	UsageMock            func(slack.CommandData, slack.Response) error
	TopMock              func(team string, n int) ([]UserKarma, error)
	GetDailyMock         func(team, user string, date time.Time) (int, error)
	UpdateDailyMock      func(team, user string, date time.Time, karma int) (int, error)
	UpdateKarmaDailyMock func(team, callee, target string, delta int, date time.Time) (int, error)
}

// GetKarma .
func (m MockDAO) GetKarma(team, user string) (int, error) {
	return m.GetKarmaMock(team, user)
}

// UpdateKarma .
func (m MockDAO) UpdateKarma(team, user string, delta int) (int, error) {
	return m.UpdateKarmaMock(team, user, delta)
}

// DeleteKarma .
func (m MockDAO) DeleteKarma(team, user string) (int, error) {
	return m.DeleteKarmaMock(team, user)
}

// Usage .
func (m MockDAO) Usage(d slack.CommandData, r slack.Response) error {
	return m.UsageMock(d, r)
}

// Top .
func (m MockDAO) Top(team string, n int) ([]UserKarma, error) {
	return m.TopMock(team, n)
}

// GetDaily .
func (m MockDAO) GetDaily(team, user string, date time.Time) (int, error) {
	return m.GetDailyMock(team, user, date)
}

// UpdateDaily .
func (m MockDAO) UpdateDaily(team, user string, date time.Time, karma int) (int, error) {
	return m.UpdateDailyMock(team, user, date, karma)
}

// UpdateKarmaDaily .
func (m MockDAO) UpdateKarmaDaily(team, callee, target string, delta int, date time.Time) (int, error) {
	return m.UpdateKarmaDailyMock(team, callee, target, delta, date)
}

// NewMockDao constructor func for making mock dao
func NewMockDao(usage int) MockDAO {
	return MockDAO{
		GetKarmaMock: func(team, user string) (int, error) {
			return 5, nil
		},
		UpdateKarmaMock: func(team, user string, delta int) (int, error) {
			return delta + 1, nil
		},
		DeleteKarmaMock: func(team, user string) (int, error) {
			return 0, nil
		},
		UsageMock: func(d slack.CommandData, r slack.Response) error {
			return nil
		},
		TopMock: func(team string, n int) ([]UserKarma, error) {
			r := make([]UserKarma, 0, n)
			for i := 0; i < n; i++ {
				name := fmt.Sprintf("USER%v", i)
				karma := 100 + i
				r = append(r, UserKarma{name, karma})
			}
			return r, nil
		},
		GetDailyMock: func(team, user string, date time.Time) (int, error) {
			return usage, nil
		},
		UpdateDailyMock: func(team, user string, date time.Time, karma int) (int, error) {
			return karma + usage, nil
		},
		UpdateKarmaDailyMock: func(team, callee, target string, delta int, date time.Time) (int, error) {
			return delta + 1, nil
		},
	}
}

// HappyDao factory method for a mock dao that will always succeed
func HappyDao() MockDAO {
	return NewMockDao(0)
}

// FullDailyDao factory method for a mock dao were all users are full on their daily limit
func FullDailyDao() MockDAO {
	return NewMockDao(DefaultConfig.DailyLimit)
}

// SadDao factory method for a mock dao that will always fail with an error
func SadDao() MockDAO {
	return MockDAO{
		GetKarmaMock: func(team, user string) (int, error) {
			return 0, errors.New("GetKarmaMock")
		},
		UpdateKarmaMock: func(team, user string, delta int) (int, error) {
			return 0, errors.New("UpdateKarmaMock")
		},
		DeleteKarmaMock: func(team, user string) (int, error) {
			return 0, errors.New("DeleteKarmaMock")
		},
		UsageMock: func(d slack.CommandData, r slack.Response) error {
			return errors.New("UsageMock")
		},
		TopMock: func(team string, n int) ([]UserKarma, error) {
			return nil, errors.New("TopMock")
		},
		GetDailyMock: func(team, user string, date time.Time) (int, error) {
			return 0, errors.New("GetDailyMock")
		},
		UpdateDailyMock: func(team, user string, date time.Time, karma int) (int, error) {
			return 0, errors.New("UpdateDailyMock")
		},
		UpdateKarmaDailyMock: func(team, callee, target string, delta int, date time.Time) (int, error) {
			return 0, errors.New("UpdateKarmaDailyMock")
		},
	}
}
