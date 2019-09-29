package karma

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/icemanblues/knave-bot/shakespeare"
	"github.com/icemanblues/knave-bot/slack"
)

// Commands a set of the support commands by this processor
var commands = map[string]struct{}{
	"help":   struct{}{},
	"me":     struct{}{},
	"status": struct{}{},
	"++":     struct{}{},
	"--":     struct{}{},
}

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

// SlackProcessor an implementation of KarmaProcessor that uses SQLite
type SlackProcessor struct {
	dao        DAO
	insult     shakespeare.Generator
	compliment shakespeare.Generator
}

// NewProcessor factory method
func NewProcessor(dao DAO, insult, compliment shakespeare.Generator) *SlackProcessor {
	return &SlackProcessor{dao, insult, compliment}
}

// Process handles Karma processing from slack API
func (p SlackProcessor) Process(c *slack.CommandData) (*slack.Response, error) {
	if len(c.Text) == 0 {
		return p.help()
	}

	words := strings.Fields(c.Text)
	if len(words) == 0 {
		return p.help()
	}

	if _, ok := commands[words[0]]; ok {
		return p.processCommand(words, c)
	}

	words = userCmdAlias(words)
	if _, ok := commands[words[0]]; ok {
		return p.processCommand(words, c)
	}

	words = addSubCmdAlias(words)
	if _, ok := commands[words[0]]; ok {
		return p.processCommand(words, c)
	}

	return p.help()
}

func (p SlackProcessor) processCommand(words []string, c *slack.CommandData) (*slack.Response, error) {
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
	}

	return p.help()
}

// alias for /karma +5 @user => /karma ++ @user 5
func addSubCmdAlias(words []string) []string {
	if len(words) < 2 {
		return words
	}

	delta, err := strconv.Atoi(words[0])
	if err != nil {
		return words
	}

	if delta == 0 {
		return words
	}

	d := strconv.Itoa(Abs(delta))

	if delta > 0 {
		w := []string{"++", words[1], d}
		w = append(w, words[2:]...)
		return w
	}

	// must be negative
	w := []string{"--", words[1], d}
	w = append(w, words[2:]...)
	return w
}

// alias for /karma @user cmd => /karma cmd @user
func userCmdAlias(words []string) []string {
	if len(words) < 2 {
		return words
	}

	target, tok := parseArgUser(words, 0)
	cmd, cok := parseArg(words, 1)
	if tok && cok {
		w := []string{cmd, target}
		if len(words) > 2 {
			w = append(w, words[2:]...)
		}
		return w
	}

	return words
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

func (p SlackProcessor) help() (*slack.Response, error) {
	return responseHelp, nil
}

func (p SlackProcessor) me(team, userID string) (*slack.Response, error) {
	k, err := p.dao.GetKarma(team, userID)
	if err != nil {
		return nil, err
	}

	msg, att := &strings.Builder{}, &strings.Builder{}
	UserStatus(userID, k, msg)
	p.Salutation(k, att)
	return slack.DirectResponse(msg.String(), att.String()), nil
}

func (p SlackProcessor) status(team, callee string, words []string) (*slack.Response, error) {
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
	msg.WriteString(fmt.Sprintf("<@%s> has requested karma total for <@%s>. ", callee, target))
	UserStatus(target, k, msg)
	p.Salutation(k, att)
	return slack.ChannelAttachmentsResponse(msg.String(), att.String()), nil
}

func (p SlackProcessor) add(team, callee string, words []string) (*slack.Response, error) {
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
	p.Salutation(delta, att)
	return slack.ChannelAttachmentsResponse(msg.String(), att.String()), nil
}

func (p SlackProcessor) subtract(team, callee string, words []string) (*slack.Response, error) {
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
	p.Salutation(-delta, att)
	return slack.ChannelAttachmentsResponse(msg.String(), att.String()), nil
}
