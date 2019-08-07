package main

import (
	"strconv"
	"strings"
)

// Abs absolute value of an int
func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// KarmaProcessor processes slash-commands into slack responses
type KarmaProcessor interface {
	Process(cd *CommandData) (*Response, error)
}

// KdbProcessor an implementation of KarmaProcessor that uses SQLite
type KdbProcessor struct {
	kdb KarmaDB
}

// NewKdbProcessor factory method
func NewKdbProcessor(kdb KarmaDB) *KdbProcessor {
	return &KdbProcessor{kdb}
}

// Process handles Karma processing from slack API
func (kp KdbProcessor) Process(c *CommandData) (*Response, error) {
	if len(c.Text) == 0 {
		return kp.help()
	}

	words := strings.Fields(c.Text)
	if len(words) == 0 {
		return kp.help()
	}

	switch words[0] {
	case "help":
		return kp.help()

	case "me":
		return kp.me(c.TeamID, c.UserID)

	case "status":
		if len(words) == 1 {
			return ErrorResponse("I need to know whose karma to update.\n`/karma ++ @name`"), nil
		}

		target := words[1]
		t, ok := IsSlackUser(target)
		if !ok {
			return ErrorResponse("I'm not sure that id a valid slack user.\n`/karma ++ @name`"), nil
		}
		return kp.me(c.TeamID, t)

	case "++":
		if len(words) == 1 {
			return ErrorResponse("I need to know whose karma to update.\n`/karma ++ @name`"), nil
		}

		target := words[1]
		t, ok := IsSlackUser(target)
		if !ok {
			return ErrorResponse("I'm not sure that id a valid slack user.\n`/karma ++ @name`"), nil
		}

		if t == c.UserID {
			return ErrorResponse("Don't be a weasel. For Shame!"), nil
		}

		// optional: see if next parameter is an amount, if so, use it
		delta := 1
		if len(words) > 2 {
			if d, err := strconv.Atoi(words[2]); err == nil {
				// Might need to do a ABS here
				delta = Abs(d)
			}
		}

		if delta == 0 {
			return ErrorResponse("Don't waste my time. For shame!"), nil
		}

		return kp.delta(c.TeamID, t, delta)

	default:
		return kp.help()
	}
}

func (kp KdbProcessor) help() (*Response, error) {
	return ErrorResponse(`
	*Help* This will provide you with additional information on how to work with Karma.
	* _me_ This will return your current karma.
	* _status_ Provide a user, and it will return their current karma
	* _++_ Provide a user and it will increase their karma. Optionally, pass a quantity of karma to give.
	* _help_ this helpful dialogue. You're welcome!
	`), nil
}

func (kp KdbProcessor) me(team, userID string) (*Response, error) {
	k, err := kp.kdb.GetKarma(team, userID)
	if err != nil {
		return nil, err
	}

	return kp.karmaStatus(userID, k)
}

func (kp KdbProcessor) delta(team, userID string, delta int) (*Response, error) {
	k, err := kp.kdb.UpdateKarma(team, userID, delta)
	if err != nil {
		return nil, err
	}

	return kp.karmaStatus(userID, k)
}
