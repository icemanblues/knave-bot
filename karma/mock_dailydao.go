package karma

import (
	"errors"
	"time"
)

// MockDailyDao a mock for testing DailyDao
type MockDailyDao struct {
	GetDailyMock    func(team, user string, date time.Time) (int, error)
	UpdateDailyMock func(team, user string, date time.Time, karma int) (int, error)
}

// GetDaily .
func (m MockDailyDao) GetDaily(team, user string, date time.Time) (int, error) {
	return m.GetDailyMock(team, user, date)
}

// UpdateDaily .
func (m MockDailyDao) UpdateDaily(team, user string, date time.Time, karma int) (int, error) {
	return m.UpdateDailyMock(team, user, date, karma)
}

// HappyDailyDao .
func HappyDailyDao() *MockDailyDao {
	return &MockDailyDao{
		GetDailyMock: func(team, user string, date time.Time) (int, error) {
			return 0, nil
		},
		UpdateDailyMock: func(team, user string, date time.Time, karma int) (int, error) {
			return karma + 1, nil
		},
	}
}

// SadDailyDao .
func SadDailyDao() *MockDailyDao {
	return &MockDailyDao{
		GetDailyMock: func(team, user string, date time.Time) (int, error) {
			return 0, errors.New("GetDailyMock")
		},
		UpdateDailyMock: func(team, user string, date time.Time, karma int) (int, error) {
			return 0, errors.New("UpdateDailyMock")
		},
	}
}
