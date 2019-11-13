package karma

import (
	"strconv"
	"strings"
	"time"

	"github.com/icemanblues/knave-bot/shakespeare"
	"github.com/icemanblues/knave-bot/slack"
	log "github.com/sirupsen/logrus"
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
	Process(cd slack.CommandData) (slack.Response, error)
}

// SlackProcessor an implementation of KarmaProcessor that uses SQLite
type SlackProcessor struct {
	config     ProcConfig
	dao        DAO
	insult     shakespeare.Generator
	compliment shakespeare.Generator
}

// NewProcessor factory method
func NewProcessor(config ProcConfig, dao DAO, insult, compliment shakespeare.Generator) SlackProcessor {
	return SlackProcessor{config, dao, insult, compliment}
}

// Process handles Karma processing from slack API
func (p SlackProcessor) Process(c slack.CommandData) (slack.Response, error) {
	if len(c.Text) == 0 {
		return p.help()
	}

	words := strings.Fields(c.Text)
	if len(words) == 0 {
		return p.help()
	}

	if _, ok := Commands[words[0]]; ok {
		return p.processCommand(words, c)
	}

	words = userCmdAlias(words)
	if _, ok := Commands[words[0]]; ok {
		return p.processCommand(words, c)
	}

	words = addSubCmdAlias(words)
	if _, ok := Commands[words[0]]; ok {
		return p.processCommand(words, c)
	}

	return p.help()
}

func (p SlackProcessor) processCommand(words []string, c slack.CommandData) (slack.Response, error) {
	switch words[0] {
	case help:
		return p.help()

	case me:
		return p.me(c.TeamID, c.UserID)

	case status:
		return p.status(c.TeamID, c.UserID, words)

	case add:
		return p.add(c.TeamID, c.UserID, words)

	case sub:
		return p.subtract(c.TeamID, c.UserID, words)

	case top:
		return p.top(c.TeamID, words)
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
		w := []string{add, words[1], d}
		w = append(w, words[2:]...)
		return w
	}

	// must be negative
	w := []string{sub, words[1], d}
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

func (p SlackProcessor) help() (slack.Response, error) {
	return ResponseHelp, nil
}

func (p SlackProcessor) me(team, userID string) (slack.Response, error) {
	k, err := p.dao.GetKarma(team, userID)
	if err != nil {
		return slack.Response{}, err
	}

	// daily usage check
	usage, err := p.dao.GetDaily(team, userID, time.Now())
	if err != nil {
		return slack.Response{}, err
	}
	available := p.config.DailyLimit - usage

	msg, att := &strings.Builder{}, &strings.Builder{}
	msg.WriteString(MsgUserStatus(userID, k))
	msg.WriteString("\n")
	msg.WriteString(MsgUserDailyLimit(usage, available))
	att.WriteString(p.Salutation(k))
	return slack.DirectResponse(msg.String(), att.String()), nil
}

func (p SlackProcessor) status(team, callee string, words []string) (slack.Response, error) {
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
		return slack.Response{}, err
	}

	msg, att := &strings.Builder{}, &strings.Builder{}
	msg.WriteString(MsgUserStatusTarget(callee, target))
	msg.WriteString(MsgUserStatus(target, k))
	att.WriteString(p.Salutation(k))
	return slack.ChannelAttachmentsResponse(msg.String(), att.String()), nil
}

func (p SlackProcessor) top(team string, words []string) (slack.Response, error) {
	n, _ := parseArgInt(words, 1, p.config.TopUserDefault)

	// no negatives are allowed
	if n <= 0 {
		n = p.config.TopUserDefault
	}
	// anything larger than 10 would look funny
	if n > p.config.TopUserMax {
		n = p.config.TopUserMax
	}

	topUsers, err := p.dao.Top(team, n)
	if err != nil {
		return slack.Response{}, err
	}

	if len(topUsers) == 0 {
		return slack.DirectResponse(msgNoKarmaForTop, ""), nil
	}

	msg, att := &strings.Builder{}, &strings.Builder{}
	msg.WriteString(MsgTopKarma(topUsers))
	att.WriteString(p.compliment.Sentence())
	return slack.ChannelAttachmentsResponse(msg.String(), att.String()), nil
}

func (p SlackProcessor) add(team, callee string, words []string) (slack.Response, error) {
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
	if delta > p.config.SingleLimit {
		return slack.ErrorResponse(msgDeltaLimit), nil
	}

	// daily usage check
	usage, err := p.dao.GetDaily(team, callee, time.Now())
	if err != nil {
		return slack.Response{}, err
	}
	available := p.config.DailyLimit - usage
	if available < delta {
		return slack.ErrorResponse(MsgOverDailyLimit(p.config.DailyLimit, usage, available)), nil
	}

	// TODO: Might want to combine these two dao statements so that they are atomic (in one transaction)
	k, err := p.dao.UpdateKarma(team, target, delta)
	if err != nil {
		return slack.Response{}, err
	}
	_, err = p.dao.UpdateDaily(team, callee, time.Now(), delta)
	if err != nil {
		log.Errorf("Was able to update the karma but not the daily usage. utoh! %v %v %v", team, callee, err)
		return slack.Response{}, err
	}

	msg, att := &strings.Builder{}, &strings.Builder{}
	msg.WriteString(MsgGiveKarma(callee, target, delta))
	msg.WriteString(MsgUserStatus(target, k))
	att.WriteString(p.Salutation(delta))
	return slack.ChannelAttachmentsResponse(msg.String(), att.String()), nil
}

func (p SlackProcessor) subtract(team, callee string, words []string) (slack.Response, error) {
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
	if delta > p.config.SingleLimit {
		return slack.ErrorResponse(msgDeltaLimit), nil
	}

	// daily usage check
	usage, err := p.dao.GetDaily(team, callee, time.Now())
	if err != nil {
		return slack.Response{}, err
	}
	available := p.config.DailyLimit - usage
	if available < delta {
		return slack.ErrorResponse(MsgOverDailyLimit(p.config.DailyLimit, usage, available)), nil
	}

	// TODO: Might want to combine these two dao statements so that they are atomic (in one transaction)
	k, err := p.dao.UpdateKarma(team, target, -delta)
	if err != nil {
		return slack.Response{}, err
	}
	_, err = p.dao.UpdateDaily(team, callee, time.Now(), delta)
	if err != nil {
		log.Errorf("Was able to update the karma but not the daily usage. utoh! %v %v %v", team, callee, err)
		return slack.Response{}, err
	}

	msg, att := &strings.Builder{}, &strings.Builder{}
	msg.WriteString(MsgTakeKarma(callee, target, delta))
	msg.WriteString(MsgUserStatus(target, k))
	att.WriteString(p.Salutation(-delta))
	return slack.ChannelAttachmentsResponse(msg.String(), att.String()), nil
}
