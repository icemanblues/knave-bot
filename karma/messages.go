package karma

import (
	"fmt"
	"strings"

	"github.com/icemanblues/knave-bot/slack"
)

// Command Examples
const (
	cmdMe     = "/karma me"
	cmdStatus = "/karma status @user"
	cmdAdd    = "/karma ++ @user"
	cmdSub    = "/karma -- @user"
	cmdHelp   = "/karma help"
	cmdTop    = "/karma top"
)

// Slack Reponses

// ResponseHelp the slack response for the HELP command
var ResponseHelp = slack.Response{
	ResponseType: slack.ResponseType.Ephemeral,
	Text:         "",
	Attachments: []slack.Attachments{
		{
			Fallback: "*Help* Helpful information on how to manage karma.",
			Title:    "Helpful information on how to manage karma.",
			Text:     "Below are the sub-commands:",
			Fields: []slack.Field{
				{
					Title: cmdMe,
					Value: "Return your karma and daily usage limits.",
					Short: true,
				},
				{
					Title: cmdStatus,
					Value: "Provide a @user and return their karma.",
					Short: true,
				},
				{
					Title: cmdAdd,
					Value: "Provide a @user and increase their karma. Optionally, pass a quantity of karma to give.",
					Short: true,
				},
				{
					Title: cmdSub,
					Value: "Provide a @user and decrease their karma. Optionally, pass a quantity of karma to take.",
					Short: true,
				},
				{
					Title: cmdTop,
					Value: "Return the top 3 users by karma. Optionally, pass a quantity for the top n users",
					Short: true,
				},
				{
					Title: cmdHelp,
					Value: "This helpful dialogue. You're welcome!",
					Short: true,
				},
			},
		},
	},
}

// Re-usable string constants for crafting messages
const (
	msgMissingName           = "I need to know whose karma to retrieve."
	msgNoOp                  = "Don't waste my time. For shame!"
	msgInvalidUser           = "I'm not sure that name is a valid slack user."
	msgDeltaLimit            = "Whoa there! Let's keep the karma swings to 5 and under."
	msgAddMissingTarget      = "To whom do you want to give karma?"
	msgAddSelfTarget         = "Don't be a weasel. For Shame!"
	msgAddCantRemove         = "`++` is used to give karma. Try `--` to take away karma."
	msgSubtractMissingTarget = "From whom do you want to take karma away?"
	msgSubtractSelfTarget    = "Do you have something to confess? Why remove your own karma?"
	msgSubtractCantAdd       = "Negative karma doesn't make sense. Please use positive numbers!"
	msgNoKarmaForTop         = "Um.. is it possible that there are no users with positive karma :("
)

// MsgOverDailyLimit generates daily limit error message (string)
func MsgOverDailyLimit(limit, usage, remainder int) string {
	return fmt.Sprintf("Ah ah ah! The daily limit is %v and you've given/taken %v karma already. Only %v remaining", limit, usage, remainder)
}

// MsgUserStatus the User's Karma status
func MsgUserStatus(userID string, k int) string {
	return fmt.Sprintf("<@%s> has %v karma.", userID, k)
}

// MsgUserDailyLimit the remaining daily limits
func MsgUserDailyLimit(usage, remaining int) string {
	return fmt.Sprintf("You have given/taken %v karma with %v remaining today.", usage, remaining)
}

// MsgUserStatusTarget lets all users know who requested karma totals
func MsgUserStatusTarget(callee, target string) string {
	return fmt.Sprintf("<@%s> has requested karma total for <@%s>. ", callee, target)
}

// MsgGiveKarma announces who gave how much karma to whom
func MsgGiveKarma(callee, target string, delta int) string {
	return fmt.Sprintf("<@%s> is giving %v karma to <@%s>. ", callee, delta, target)
}

// MsgTakeKarma announces who took how much karma from whom
func MsgTakeKarma(callee, target string, delta int) string {
	return fmt.Sprintf("<@%s> is taking away %v karma from <@%s>. ", callee, delta, target)
}

// MsgTopKarma table for viewing top users by karma
func MsgTopKarma(topUsers []UserKarma) string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("The top %v users by karma:\n", len(topUsers)))
	sb.WriteString("Rank\tName\tKarma\n")
	for i, user := range topUsers {
		sb.WriteString(fmt.Sprintf("%v\t<@%v>\t%v\n", i+1, user.User, user.Karma))
	}
	return sb.String()
}

// Salutation appends a Salutation (insult or compliment)
func (p SlackProcessor) Salutation(k int) string {
	if k > 0 {
		return p.compliment.Sentence()
	}
	if k == 0 {
		return ""
	}

	return p.insult.Sentence()
}
