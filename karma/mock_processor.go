package karma

import "github.com/icemanblues/knave-bot/shakespeare"

func mockProcessor(dao DAO) SlackProcessor {
	return NewProcessor(DefaultConfig, dao,
		shakespeare.New("insult", "", nil),
		shakespeare.New("compliment", "", nil))
}

func happyMockProcessor() SlackProcessor {
	return mockProcessor(HappyDao())
}

func fullUsageMockProcessor() SlackProcessor {
	return mockProcessor(NewMockDao(DefaultConfig.DailyLimit))
}

func sadMockProcessor() SlackProcessor {
	return mockProcessor(SadDao())
}
