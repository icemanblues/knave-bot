package karma

import (
	"errors"

	"github.com/icemanblues/knave-bot/slack"

	"github.com/icemanblues/knave-bot/shakespeare"
)

// MockDAO a mock dao for karma whose mock functions can be monkeypatched
type MockDAO struct {
	GetKarmaMock    func(team, user string) (int, error)
	UpdateKarmaMock func(team, user string, delta int) (int, error)
	DeleteKarmaMock func(team, user string) (int, error)
	UsageMock       func(*slack.CommandData, *slack.Response) error
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
func (m MockDAO) Usage(d *slack.CommandData, r *slack.Response) error {
	return m.UsageMock(d, r)
}

// HappyDao factory method for a mock dao that will always succeed
func HappyDao() *MockDAO {
	return &MockDAO{
		GetKarmaMock: func(team, user string) (int, error) {
			return 5, nil
		},
		UpdateKarmaMock: func(team, user string, delta int) (int, error) {
			return delta + 1, nil
		},
		DeleteKarmaMock: func(team, user string) (int, error) {
			return 0, nil
		},
		UsageMock: func(d *slack.CommandData, r *slack.Response) error {
			return nil
		},
	}
}

// SadDao factory method for a mock dao that will always fail with an error
func SadDao() *MockDAO {
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
		UsageMock: func(d *slack.CommandData, r *slack.Response) error {
			return errors.New("UsageMock")
		},
	}
}

func mockProcessor(dao DAO) *SlackProcessor {
	return NewProcessor(dao,
		shakespeare.New("insult", "", nil),
		shakespeare.New("compliment", "", nil))
}
