package gotwilio

import (
	"encoding/xml"
	"errors"
	"strings"
)

const (
	// opt-out keywords
	SMSKeywordEnd         = "end"
	SMSKeywordQuit        = "quit"
	SMSKeywordStop        = "stop"
	SMSKeywordCancel      = "cancel"
	SMSKeywordStopAll     = "stopall"
	SMSKeywordUnsubscribe = "unsubscribe"

	// opt-in keywords
	SMSKeywordYes    = "yes"
	SMSKeywordStart  = "start"
	SMSKeywordUnstop = "unstop"
)

var (
	// SMSOptOutKeywords Twilio support for opt-out keywords
	SMSOptOutKeywords = []string{
		SMSKeywordEnd,
		SMSKeywordQuit,
		SMSKeywordStop,
		SMSKeywordCancel,
		SMSKeywordStopAll,
		SMSKeywordUnsubscribe,
	}

	// SMSOptInKeywords Twilio support for opt-in keywords
	SMSOptInKeywords = []string{
		SMSKeywordYes,
		SMSKeywordStart,
		SMSKeywordUnstop,
	}
)

// IsSMSOptOutKeyword check if given keyword is an opt-out keyword supported by Twilio
func IsSMSOptOutKeyword(s string) bool {
	s = strings.ToLower(s)
	for _, k := range SMSOptOutKeywords {
		if k == s {
			return true
		}
	}
	return false
}

// IsSMSOptInKeyword check if given keyword is an opt-in keyword supported by Twilio
func IsSMSOptInKeyword(s string) bool {
	s = strings.ToLower(s)
	for _, k := range SMSOptInKeywords {
		if k == s {
			return true
		}
	}
	return false
}

// MessagingResponse Twilio's TWiML sms response
type MessagingResponse struct {
	XMLName  xml.Name           `xml:"Response"`
	Messages []*TWiMLSmsMessage `xml:"Message"`
}

// TWiMLSmsMessage response content
type TWiMLSmsMessage struct {
	Message string `xml:",chardata"`

	// nouns
	Body     *string `xml:"Body,omitempty"`
	Media    *string `xml:"Media,omitempty"`
	Redirect *string `xml:"Redirect,omitempty"`

	// verbs - xml attributes
	Action *string `xml:"Action,attr,omitempty"`
	Method *string `xml:"Method,attr,omitempty"`
}

// Message add message to TMiML response
func (r *MessagingResponse) Message(msg *TWiMLSmsMessage) (*MessagingResponse, error) {
	if (msg.Body != nil || msg.Media != nil) && (msg.Action != nil || msg.Method != nil) {
		return r, errors.New("can't nest verbs within Message and can't net Message in any other verb")
	}

	// twilio doesn't allow message when body is set
	if msg.Body != nil && msg.Message != "" {
		msg.Message = ""
	}

	r.Messages = append(r.Messages, msg)
	return r, nil
}

// TWiMLSmsRender render XML response to send to Twilio
func (r *MessagingResponse) TWiMLSmsRender() (string, error) {
	output, err := xml.MarshalIndent(r, "  ", "   ")
	if err != nil {
		return "", err
	}
	return xml.Header + string(output), nil
}
