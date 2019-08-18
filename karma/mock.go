package karma

import "errors"

// MockDAO a mock dao for karma whose mock functions can be monkeypatched
type MockDAO struct {
	GetKarmaMock    func(team, user string) (int, error)
	UpdateKarmaMock func(team, user string, delta int) (int, error)
	DeleteKarmaMock func(team, user string) (int, error)
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
	}
}
