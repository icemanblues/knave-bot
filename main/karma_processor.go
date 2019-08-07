package main

import (
	"fmt"
	"strconv"
	"strings"
)

type KarmaProcessor interface {
	Process(cd *CommandData) (*Response, error)
}

type KdbProcessor struct {
	kdb KarmaDB
}

func NewKdbProcessor(kdb KarmaDB) *KdbProcessor {
	return &KdbProcessor{kdb}
}

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
	case "++":
		if len(words) == 1 {
			return ErrorResponse("I need to know whose karma to update.\n\"/karma ++ @name\""), nil
		}

		target := words[1]
		t, ok := IsSlackUser(target)
		if !ok {
			return ErrorResponse("I'm not sure that id a valid slack user.\n\"/karma ++ @name\""), nil
		}

		if t == c.UserID {
			return ErrorResponse("Don't be a weasel. For Shame!"), nil
		}

		// optional: see if next parameter is an amount, if so, use it
		delta := 1
		if len(words) > 2 {
			if d, err := strconv.Atoi(words[2]); err == nil {
				// Might need to do a ABS here
				delta = d
			}
		}

		return kp.delta(c.TeamID, target, delta)
	default:
		return kp.help()
	}
}

func (kp KdbProcessor) help() (*Response, error) {
	return ErrorResponse(`
	*Help* This will provide you with additional information on how to work with Karma.
	me - This will return your current karma.
	++ - Give it a user and it will increase their karma. Optionally, pass a number as well. /karma ++ @name 5
	help - this helpful dialogue. You're welcome!
	`), nil
}

func (kp KdbProcessor) me(team, userID string) (*Response, error) {
	k := kp.kdb.GetKarma(team, userID)
	return kp.status(userID, k)
}

func (kp KdbProcessor) delta(team, userID string, delta int) (*Response, error) {
	if delta == 0 {
		return ErrorResponse("Don't waste my time. For shame!"), nil
	}

	k := kp.kdb.UpdateKarma(team, userID, delta)
	return kp.status(userID, k)
}

func (kp KdbProcessor) status(userID string, k int) (*Response, error) {
	if k > 0 {
		msg := fmt.Sprintf("%s has %v karma. Have a compliment.\n%s", userID, k, Compliment())
		return ChannelResponse(msg), nil
	}
	if k == 0 {
		msg := fmt.Sprintf("%s has %v karma. Today is a good day, be well", userID, k)
		return ChannelResponse(msg), nil
	}

	msg := fmt.Sprintf("%s have %v karma. Be a better person or I will insult you again.\n%s", userID, k, Insult())
	return ChannelResponse(msg), nil
}
