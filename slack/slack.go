package slack

import "regexp"

// CommandData data payload for a slash command
type CommandData struct {
	Command      string `json:"command,omitempty"`
	Text         string `json:"text,omitempty"`
	ResponseURL  string `json:"response_url,omitempty"`
	EnterpriseID string `json:"enterprise_id,omitempty"`
	TeamID       string `json:"team_id,omitempty"`
	ChannelID    string `json:"channel_id,omitempty"`
	UserID       string `json:"user_id,omitempty"`
}

// Attachments any additional "attachments" for a slash-command response
type Attachments struct {
	Text string `json:"text,omitempty"`
}

// Response a slack slash-command response
type Response struct {
	ResponseType string      `json:"response_type,omitempty"`
	Text         string      `json:"text,omitempty"`
	Attachments  Attachments `json:"attachments,omitempty"`
}

// ResponseType simple string enum for slash-command responses
var ResponseType = struct {
	InChannel string
	Ephemeral string
}{
	InChannel: "in_channel",
	Ephemeral: "ephemeral",
}

// ChannelResponse factory method for a response that should be displayed to the entire channel
func ChannelResponse(msg string) *Response {
	return &Response{
		ResponseType: ResponseType.InChannel,
		Text:         msg,
	}
}

// ErrorResponse factory method for a response that should be displayed only to the callee
func ErrorResponse(msg string) *Response {
	return &Response{
		ResponseType: ResponseType.Ephemeral,
		Text:         msg,
	}
}

var canonical = regexp.MustCompile("U[A-Z0-9]+")
var escaped = regexp.MustCompile("<@(U[A-Z0-9]+)|.*>")

// IsSlackUser returns the canonical slack user id. <@UAWQFTRT7|roland.kluge> => UAWQFTRT7
func IsSlackUser(userID string) (string, bool) {
	if escaped.MatchString(userID) {
		return escaped.ReplaceAllString(userID, "$1"), true
	}
	if canonical.MatchString(userID) {
		return userID, true
	}

	return "", false
}
