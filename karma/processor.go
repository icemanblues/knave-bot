package karma

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/icemanblues/knave-bot/slack"
)

// Abs absolute value of an int
func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// Processor processes slash-commands into slack responses
type Processor interface {
	Process(cd *slack.CommandData) (*slack.Response, error)
}

// SQLiteProcessor an implementation of KarmaProcessor that uses SQLite
type SQLiteProcessor struct {
	kdb DAO
}

// NewProcessor factory method
func NewProcessor(kdb DAO) *SQLiteProcessor {
	return &SQLiteProcessor{kdb}
}

// Process handles Karma processing from slack API
func (p SQLiteProcessor) Process(c *slack.CommandData) (*slack.Response, error) {
	if len(c.Text) == 0 {
		return p.help()
	}

	words := strings.Fields(c.Text)
	if len(words) == 0 {
		return p.help()
	}

	switch words[0] {
	case "help":
		return p.help()

	case "me":
		return p.me(c.TeamID, c.UserID)

	case "status":
		return p.status(c.TeamID, c.UserID, words)

	case "++":
		return p.add(c.TeamID, c.UserID, words)

	case "--":
		return p.subtract(c.TeamID, c.UserID, words)

	default:
		return p.help()
	}
}

func parseArg(words []string, idx int) (string, bool) {
	if idx >= len(words) {
		return "", false
	}

	return words[idx], true
}

func parseArgInt(words []string, idx int) (int, bool) {
	s, ok := parseArg(words, idx)
	if !ok {
		return 0, false
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, false
	}

	return i, true
}

func parseArgUser(words []string, idx int) (string, bool) {
	s, ok := parseArg(words, idx)
	if !ok {
		return "", false
	}

	return slack.IsSlackUser(s)
}

func (p SQLiteProcessor) help() (*slack.Response, error) {
	return slack.ErrorResponse(`
	*Help* This will provide you with additional information on how to work with Karma.
	* _me_ This will return your current karma.
	* _status_ Provide a user, and it will return their current karma
	* _++_ Provide a user and it will increase their karma. Optionally, pass a quantity of karma to give.
	* _--_ Provide a user and it will decrease their karma. Optionally, pass a quantity of karma to take.
	* _help_ this helpful dialogue. You're welcome!
	`), nil
}

func (p SQLiteProcessor) me(team, userID string) (*slack.Response, error) {
	k, err := p.kdb.GetKarma(team, userID)
	if err != nil {
		return nil, err
	}

	msg, att := &strings.Builder{}, &strings.Builder{}
	UserStatus(userID, k, msg)
	Salutation(k, att)
	return slack.DirectResponse(msg.String(), att.String()), nil
}

func (p SQLiteProcessor) status(team, callee string, words []string) (*slack.Response, error) {
	name, ok := parseArg(words, 1)
	if !ok {
		return slack.ErrorResponse("I need to know whose karma to retrieve.\n`/karma status @name`"), nil
	}

	target, ok := slack.IsSlackUser(name)
	if !ok {
		return slack.ErrorResponse("I'm not sure that name is a valid slack user.\n`/karma status @name`"), nil
	}

	k, err := p.kdb.GetKarma(team, target)
	if err != nil {
		return nil, err
	}

	msg, att := &strings.Builder{}, &strings.Builder{}
	msg.WriteString(fmt.Sprintf("<@%s> has requested <@%s> karma. ", callee, target))
	UserStatus(target, k, msg)
	Salutation(k, att)
	return slack.ChannelAttachmentsResponse(msg.String(), att.String()), nil
}

func (p SQLiteProcessor) add(team, callee string, words []string) (*slack.Response, error) {
	name, ok := parseArg(words, 1)
	if !ok {
		return slack.ErrorResponse("To whom do you want to give karma?\n`/karma ++ @name`"), nil
	}

	target, ok := slack.IsSlackUser(name)
	if !ok {
		return slack.ErrorResponse("I'm not sure that is a valid slack user.\n`/karma ++ @name`"), nil
	}

	if target == callee {
		return slack.ErrorResponse("Don't be a weasel. For Shame!"), nil
	}

	// optional: see if next parameter is an amount, if so, use it
	delta, ok := parseArgInt(words, 2)
	if !ok {
		delta = 1
	}

	if delta == 0 {
		return slack.ErrorResponse("Don't waste my time. For shame!"), nil
	}
	if delta < 0 {
		return slack.ErrorResponse("`++` is used to give karma. Try `--` to take away karma."), nil
	}

	k, err := p.kdb.UpdateKarma(team, target, delta)
	if err != nil {
		return nil, err
	}

	msg, att := &strings.Builder{}, &strings.Builder{}
	msg.WriteString(fmt.Sprintf("<@%s> is giving %v karma to <@%s>. ", callee, delta, target))
	UserStatus(target, k, msg)
	Salutation(delta, att)
	return slack.ChannelAttachmentsResponse(msg.String(), att.String()), nil
}

func (p SQLiteProcessor) subtract(team, callee string, words []string) (*slack.Response, error) {
	name, ok := parseArg(words, 1)
	if !ok {
		return slack.ErrorResponse("To whom do you want to take karma?\n`/karma -- @name`"), nil
	}

	target, ok := slack.IsSlackUser(name)
	if !ok {
		return slack.ErrorResponse("I'm not sure that is a valid slack user.\n`/karma -- @name`"), nil
	}

	if target == callee {
		return slack.ErrorResponse("Do you have something to confess? Why remove your own karma?"), nil
	}

	// optional: see if next parameter is an amount, if so, use it
	delta, ok := parseArgInt(words, 2)
	if !ok {
		delta = 1
	}

	if delta == 0 {
		return slack.ErrorResponse("Don't waste my time. For shame!"), nil
	}
	if delta < 0 {
		return slack.ErrorResponse("Negative karma doesn't make sense. Please use positive numbers!"), nil
	}

	k, err := p.kdb.UpdateKarma(team, target, -delta)
	if err != nil {
		return nil, err
	}

	msg, att := &strings.Builder{}, &strings.Builder{}
	msg.WriteString(fmt.Sprintf("<@%s> is taking away %v karma from <@%s>. ", callee, delta, target))
	UserStatus(target, k, msg)
	Salutation(-delta, att)
	return slack.ChannelAttachmentsResponse(msg.String(), att.String()), nil
}
