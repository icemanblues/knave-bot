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
	dao DAO
}

// NewProcessor factory method
func NewProcessor(dao DAO) *SQLiteProcessor {
	return &SQLiteProcessor{dao}
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
	if idx >= len(words) || idx < 0 {
		return "", false
	}

	return words[idx], true
}

func parseArgInt(words []string, idx int, d int) (int, bool) {
	s, ok := parseArg(words, idx)
	if !ok {
		return d, false
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		return d, false
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
	return responseHelp, nil
}

func (p SQLiteProcessor) me(team, userID string) (*slack.Response, error) {
	k, err := p.dao.GetKarma(team, userID)
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
		return slack.DirectResponse(msgMissingName, cmdStatus), nil
	}

	target, ok := slack.IsSlackUser(name)
	if !ok {
		return slack.DirectResponse(msgInvalidUser, cmdStatus), nil
	}

	k, err := p.dao.GetKarma(team, target)
	if err != nil {
		return nil, err
	}

	msg, att := &strings.Builder{}, &strings.Builder{}
	msg.WriteString(fmt.Sprintf("<@%s> has requested karma total for <@%s>.", callee, target))
	UserStatus(target, k, msg)
	Salutation(k, att)
	return slack.ChannelAttachmentsResponse(msg.String(), att.String()), nil
}

func (p SQLiteProcessor) add(team, callee string, words []string) (*slack.Response, error) {
	name, ok := parseArg(words, 1)
	if !ok {
		return slack.DirectResponse(msgAddMissingTarget, cmdAdd), nil
	}

	target, ok := slack.IsSlackUser(name)
	if !ok {
		return slack.DirectResponse(msgInvalidUser, cmdAdd), nil
	}

	if target == callee {
		return slack.ErrorResponse(msgAddSelfTarget), nil
	}

	delta, _ := parseArgInt(words, 2, 1)
	if delta == 0 {
		return slack.ErrorResponse(msgNoOp), nil
	}
	if delta < 0 {
		return slack.ErrorResponse(msgAddCantRemove), nil
	}
	if delta > 5 {
		return slack.ErrorResponse(msgDeltaLimit), nil
	}

	k, err := p.dao.UpdateKarma(team, target, delta)
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
		return slack.DirectResponse(msgSubtractMissingTarget, cmdSub), nil
	}

	target, ok := slack.IsSlackUser(name)
	if !ok {
		return slack.DirectResponse(msgInvalidUser, cmdSub), nil
	}

	if target == callee {
		return slack.ErrorResponse(msgSubtractSelfTarget), nil
	}

	// optional: see if next parameter is an amount, if so, use it
	delta, _ := parseArgInt(words, 2, 1)
	if delta == 0 {
		return slack.DirectResponse(msgNoOp, cmdSub), nil
	}
	if delta < 0 {
		return slack.DirectResponse(msgSubtractCantAdd, cmdSub), nil
	}
	if delta > 5 {
		return slack.ErrorResponse(msgDeltaLimit), nil
	}

	k, err := p.dao.UpdateKarma(team, target, -delta)
	if err != nil {
		return nil, err
	}

	msg, att := &strings.Builder{}, &strings.Builder{}
	msg.WriteString(fmt.Sprintf("<@%s> is taking away %v karma from <@%s>. ", callee, delta, target))
	UserStatus(target, k, msg)
	Salutation(-delta, att)
	return slack.ChannelAttachmentsResponse(msg.String(), att.String()), nil
}
