package karma

import (
	"fmt"
	"strings"

	"github.com/icemanblues/knave-bot/shakespeare"
)

// UserStatus appends the User's Karma status
func UserStatus(userID string, k int, sb *strings.Builder) {
	sb.WriteString(fmt.Sprintf("<@%s> has %v karma.", userID, k))
}

// Salutation appends a Salutation (insult or compliment)
func Salutation(k int, sb *strings.Builder) {
	if k > 0 {
		sb.WriteString(shakespeare.Compliment())
		return
	}
	if k == 0 {
		return
	}

	sb.WriteString(shakespeare.Insult())
}
