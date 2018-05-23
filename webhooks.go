package gotwilio

import (
	"encoding/json"
	"net/url"
	"time"

	"github.com/gorilla/schema"
)

var formDecoder *schema.Decoder

func init() {
	formDecoder = schema.NewDecoder()
	formDecoder.SetAliasTag("form")
}

func DecodeWebhook(data url.Values, out interface{}) error {
	return formDecoder.Decode(out, data)
}

type InteractionData struct {
	Body string `json:"body"`
}

// https://www.twilio.com/docs/proxy/api/proxy-webhooks

// https://www.twilio.com/docs/proxy/api/proxy-webhooks#callbackurl
// These webhooks are fired for each new interaction and are informational only.
type ProxyCallbackWebhook struct {
	OutboundResourceStatus string    `form:"outboundResourceStatus"`
	OutboundResourceType   string    `form:"outboundResourceType"`
	InteractionDateUpdated time.Time `form:"interactionDateUpdated"`
	InteractionData        string    `form:"interactionData"`
	InteractionDateCreated time.Time `form:"interactionDateCreated"`
	InboundResourceURL     string    `form:"inboundResourceUrl"`
	InteractionServiceSid  string    `form:"interactionServiceSid"`
	OutboundParticipantSid string    `form:"outboundParticipantSid"`
	InteractionType        string    `form:"interactionType"`
	InteractionAccountSid  string    `form:"interactionAccountSid"`
	InboundParticipantSid  string    `form:"inboundParticipantSid"`
	InboundResourceStatus  string    `form:"inboundResourceStatus"`
	OutboundResourceSid    string    `form:"outboundResourceSid"`
	OutboundResourceURL    string    `form:"outboundResourceUrl"`
	InboundResourceType    string    `form:"inboundResourceType"`
	InboundResourceSid     string    `form:"inboundResourceSid"`
	InteractionSessionSid  string    `form:"interactionSessionSid"`
	InteractionSid         string    `form:"interactionSid"`
}

func (p ProxyCallbackWebhook) GetInteractionData() (InteractionData, error) {
	var out InteractionData
	err := json.Unmarshal([]byte(p.InteractionData), &out)
	return out, err
}

// https://www.twilio.com/docs/proxy/api/proxy-webhooks#interceptcallbackurl
// Fires on each interaction. If responded to with a 403 to this webhook we
// will abort/block the interaction. Any other status or timeout the interaction continues
type ProxyInterceptCallbackWebhook struct {
	InteractionDateUpdated time.Time `form:"interactionDateUpdated"`
	InteractionData        string    `form:"interactionData"`
	InteractionDateCreated time.Time `form:"interactionDateCreated"`
	InboundResourceURL     string    `form:"inboundResourceUrl"`
	InteractionServiceSid  string    `form:"interactionServiceSid"`
	InteractionType        string    `form:"interactionType"`
	InteractionAccountSid  string    `form:"interactionAccountSid"`
	InboundParticipantSid  string    `form:"inboundParticipantSid"`
	InboundResourceStatus  string    `form:"inboundResourceStatus"`
	InboundResourceType    string    `form:"inboundResourceType"`
	InboundResourceSid     string    `form:"inboundResourceSid"`
	InteractionSessionSid  string    `form:"interactionSessionSid"`
	InteractionSid         string    `form:"interactionSid"`
}

func (p ProxyInterceptCallbackWebhook) GetInteractionData() (InteractionData, error) {
	var out InteractionData
	err := json.Unmarshal([]byte(p.InteractionData), &out)
	return out, err
}

// https://www.twilio.com/docs/proxy/api/proxy-webhooks#outofsessioncallbackurl
// A URL to send webhooks to when an action (inbound call or SMS) occurs where
// there is no session or a closed session. If your server (or a Twilio function)
// responds with valid TwiML, this will be processed.
// This means it is possible to e.g. play a message for a call, send an automated
// text message response, or redirect a call to another number.
type ProxyOutOfSessionCallbackWebhook struct {
	AccountSid                 string    `form:"AccountSid"`
	SessionUniqueName          string    `form:"sessionUniqueName"`
	SessionAccountSid          string    `form:"sessionAccountSid"`
	SessionServiceSid          string    `form:"sessionServiceSid"`
	SessionSid                 string    `form:"sessionSid"`
	SessionStatus              string    `form:"sessionStatus"`
	SessionMode                string    `form:"sessionMode"`
	SessionDateCreated         time.Time `form:"sessionDateCreated"`
	SessionDateStarted         time.Time `form:"sessionDateStarted"`
	SessionDateUpdated         time.Time `form:"sessionDateUpdated"`
	SessionDateEnded           time.Time `form:"sessionDateEnded"`
	SessionDateLastInteraction time.Time `form:"sessionDateLastInteraction"`
	SessionClosedReason        string    `form:"sessionClosedReason"`
	TTL                        string    `form:"ttl"`

	// SMS Specific
	Body          string `form:"Body"`
	SmsSid        string `form:"SmsSid"`
	MessageSid    string `form:"MessageSid"`
	NumMedia      string `form:"NumMedia"`
	NumSegments   string `form:"NumSegments"`
	SmsStatus     string `form:"SmsStatus"`
	SmsMessageSid string `form:"SmsMessageSid"`

	To          string `form:"To"`
	ToCity      string `form:"ToCity"`
	ToState     string `form:"ToState"`
	ToZip       string `form:"ToZip"`
	ToCountry   string `form:"ToCountry"`
	From        string `form:"From"`
	FromCity    string `form:"FromCity"`
	FromState   string `form:"FromState"`
	FromZip     string `form:"FromZip"`
	FromCountry string `form:"FromCountry"`

	InboundParticipantSid                string    `form:"inboundParticipantSid"`
	InboundParticipantIdentifier         string    `form:"inboundParticipantIdentifier"`
	InboundParticipantFriendlyName       string    `form:"inboundParticipantFriendlyName"`
	InboundParticipantProxyIdentifier    string    `form:"inboundParticipantProxyIdentifier"`
	InboundParticipantProxyIdentifierSid string    `form:"inboundParticipantProxyIdentifierSid"`
	InboundParticipantAccountSid         string    `form:"inboundParticipantAccountSid"`
	InboundParticipantServiceSid         string    `form:"inboundParticipantServiceSid"`
	InboundParticipantSessionSid         string    `form:"inboundParticipantSessionSid"`
	InboundParticipantDateCreated        time.Time `form:"inboundParticipantDateCreated"`
	InboundParticipantDateUpdated        time.Time `form:"inboundParticipantDateUpdated"`

	OutboundParticipantSid                string    `form:"outboundParticipantSid"`
	OutboundParticipantIdentifier         string    `form:"outboundParticipantIdentifier"`
	OutboundParticipantFriendlyName       string    `form:"outboundParticipantFriendlyName"`
	OutboundParticipantProxyIdentifier    string    `form:"outboundParticipantProxyIdentifier"`
	OutboundParticipantProxyIdentifierSid string    `form:"outboundParticipantProxyIdentifierSid"`
	OutboundParticipantAccountSid         string    `form:"outboundParticipantAccountSid"`
	OutboundParticipantServiceSid         string    `form:"outboundParticipantServiceSid"`
	OutboundParticipantSessionSid         string    `form:"outboundParticipantSessionSid"`
	OutboundParticipantDateCreated        time.Time `form:"outboundParticipantDateCreated"`
	OutboundParticipantDateUpdated        time.Time `form:"outboundParticipantDateUpdated"`

	CallSid    string `form:"CallSid"`
	CallStatus string `form:"CallStatus"`

	Caller        string `form:"Caller"`
	CallerCity    string `form:"CallerCity"`
	CallerState   string `form:"CallerState"`
	CallerZip     string `form:"CallerZip"`
	CallerCountry string `form:"CallerCountry"`

	Called        string `form:"Called"`
	CalledCity    string `form:"CalledCity"`
	CalledState   string `form:"CalledState"`
	CalledZip     string `form:"CalledZip"`
	CalledCountry string `form:"CalledCountry"`

	Direction  string `form:"Direction"`
	AddOns     string `form:"AddOns"`
	APIVersion string `form:"ApiVersion"`
}
