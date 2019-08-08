package karma

import (
	"fmt"
	"strings"

	"github.com/icemanblues/knave-bot/shakespeare"
	"github.com/icemanblues/knave-bot/slack"
)

func (kp SQLiteProcessor) karmaStatus(userID string, k int) (*slack.Response, error) {
	if k > 0 {
		msg := fmt.Sprintf("<@%s> has %v karma. Have a compliment.\n_%s_", userID, k, shakespeare.Compliment())
		return slack.ChannelResponse(msg), nil
	}
	if k == 0 {
		msg := fmt.Sprintf("<@%s> has %v karma. Today is a good day to do some good", userID, k)
		return slack.ChannelResponse(msg), nil
	}

	msg := fmt.Sprintf("<@%s> have %v karma. Be a better person or I will insult you again.\n_%s_", userID, k, shakespeare.Insult())
	return slack.ChannelResponse(msg), nil
}

func karmaStatus(userID string, k int, sb *strings.Builder) {
	sb.WriteString(fmt.Sprintf("<@%s> has %v karma.", userID, k))
}

func karmaSalutation(k int, sb *strings.Builder) {
	if k > 0 {
		sb.WriteString(shakespeare.Compliment())
		return
	}
	if k == 0 {
		return
	}

	sb.WriteString(shakespeare.Insult())
}
