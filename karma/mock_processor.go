package karma

import "github.com/icemanblues/knave-bot/shakespeare"

func mockProcessor(dao DAO, dailyDao DailyDao) *SlackProcessor {
	return NewProcessor(dao, dailyDao,
		shakespeare.New("insult", "", nil),
		shakespeare.New("compliment", "", nil))
}

func happyMockProcessor() *SlackProcessor {
	return mockProcessor(HappyDao(), HappyDailyDao())
}

func sadMockProcessor() *SlackProcessor {
	return mockProcessor(SadDao(), SadDailyDao())
}
