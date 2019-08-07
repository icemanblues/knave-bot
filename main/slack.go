package main

import "regexp"

type CommandData struct {
	Command      string `json:"command,omitempty"`
	Text         string `json:"text,omitempty"`
	ResponseURL  string `json:"response_url,omitempty"`
	EnterpriseID string `json:"enterprise_id,omitempty"`
	TeamID       string `json:"team_id,omitempty"`
	ChannelID    string `json:"channel_id,omitempty"`
	UserID       string `json:"user_id,omitempty"`
}

type Attachments struct {
	Text string `json:"text,omitempty"`
}

type Response struct {
	ResponseType string      `json:"response_type,omitempty"`
	Text         string      `json:"text,omitempty"`
	Attachments  Attachments `json:"attachments,omitempty"`
}

var ResponseType = struct {
	InChannel string
	Ephemeral string
}{
	InChannel: "in_channel",
	Ephemeral: "ephemeral",
}

func ChannelResponse(msg string) *Response {
	return &Response{
		ResponseType: ResponseType.InChannel,
		Text:         msg,
	}
}

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
