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

// Field tabular like data in an attachment
type Field struct {
	Title string `json:"title,omitempty"`
	Value string `json:"value,omitempty"`
	Short bool   `json:"short,omitempty"`
}

// Attachments any additional "attachments" for a slash-command response
type Attachments struct {
	Fallback string `json:"fallback,omitempty"`
	// Color      string `json:"color,omitempty"`
	Pretext string `json:"pretext,omitempty"`
	// AuthorLink string `json:"author_link,omitempty"`
	// AuthorName string `json:"author_name,omitempty"`
	// AuthorIcon string `json:"author_icon,omitempty"`
	Title     string  `json:"title,omitempty"`
	TitleLink string  `json:"title_link,omitempty"`
	Text      string  `json:"text,omitempty"`
	Fields    []Field `json:"fields,omitempty"`
	// ImageURL   string `json:"image_url,omitempty"`
	// ThumbURL   string `json:"thumb_url,omitempty"`
	// Footer     string `json:"footer,omitempty"`
	// FooterIcon string `json:"footer_icon,omitempty"`
	// Timestamp  string `json:"ts,omitempty"`
}

// Response a slack slash-command response
type Response struct {
	ResponseType string        `json:"response_type,omitempty"`
	Text         string        `json:"text,omitempty"`
	Attachments  []Attachments `json:"attachments,omitempty"`
}

// ResponseType simple string enum for slash-command responses
var ResponseType = struct {
	InChannel string
	Ephemeral string
}{
	InChannel: "in_channel",
	Ephemeral: "ephemeral",
}

// NewAttachments factory method to create a simple []Attachments
func NewAttachments(s string) []Attachments {
	if s == "" {
		return nil
	}

	return []Attachments{Attachments{Text: s}}
}

// ChannelResponse factory method for a response that should be displayed to the entire channel
func ChannelResponse(msg string) *Response {
	return &Response{
		ResponseType: ResponseType.InChannel,
		Text:         msg,
	}
}

// ChannelAttachmentsResponse factory method for a response (with attachment) to the channel
func ChannelAttachmentsResponse(msg, att string) *Response {
	return &Response{
		ResponseType: ResponseType.InChannel,
		Text:         msg,
		Attachments:  NewAttachments(att),
	}
}

// DirectResponse factory method for a response (with attachment) to the callee
func DirectResponse(msg, att string) *Response {
	return &Response{
		ResponseType: ResponseType.Ephemeral,
		Text:         msg,
		Attachments:  NewAttachments(att),
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
var escaped = regexp.MustCompile("<@(U[A-Z0-9]+).*>")

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
