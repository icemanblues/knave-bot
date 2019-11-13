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
	tmplOverDailyLimit       = "Ah ah ah! The daily limit is %v and you've given/taken %v karma already. Only %v remaining"
)

// MsgOverDailyLimit generates daily limit error message (string)
func MsgOverDailyLimit(limit, usage, remainder int) string {
	return fmt.Sprintf(tmplOverDailyLimit, limit, usage, remainder)
}

// UserStatus appends the User's Karma status
func UserStatus(userID string, k int, sb *strings.Builder) {
	sb.WriteString(fmt.Sprintf("<@%s> has %v karma.", userID, k))
}

// UserDailyLimit appends the remaining daily limits
func UserDailyLimit(usage, remaining int, sb *strings.Builder) {
	sb.WriteString(fmt.Sprintf("You have given/taken %v karma with %v remaining today.", usage, remaining))
}

// Salutation appends a Salutation (insult or compliment)
func (p SlackProcessor) Salutation(k int, sb *strings.Builder) {
	if k > 0 {
		sb.WriteString(p.compliment.Sentence())
		return
	}
	if k == 0 {
		return
	}

	sb.WriteString(p.insult.Sentence())
}
