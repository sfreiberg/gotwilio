// Package gotwilio is a library for interacting with http://www.twilio.com/ API.
package gotwilio

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// https://www.twilio.com/docs/proxy/api/pb-proxy-service
const ProxyBaseUrl = "https://proxy.twilio.com/v1"

type ProxyServiceRequest struct {
	UniqueName              string // optional
	CallbackURL             string // optional
	OutOfSessionCallbackURL string // optional
	InterceptCallbackURL    string // optional
	GeoMatchLevel           string // Default: country Options: country, area-code, extended-area-code	optional
	NumberSelectionBehavior string // Default: prefer-sticky Options: prefer-sticky, avoid-sticky	optional
	DefaultTtl              int    // (seconds) Default: 0	optional
}

type ProxyService struct {
	Sid                     string      `json:"sid"`
	UniqueName              string      `json:"unique_name"`
	AccountSid              string      `json:"account_sid"`
	CallbackURL             string      `json:"callback_url"`
	DefaultTTL              int         `json:"default_ttl"`
	NumberSelectionBehavior string      `json:"number_selection_behavior"`
	GeoMatchLevel           string      `json:"geo_match_level"`
	InterceptCallbackURL    interface{} `json:"intercept_callback_url"`
	OutOfSessionCallbackURL interface{} `json:"out_of_session_callback_url"`
	DateCreated             time.Time   `json:"date_created"`
	DateUpdated             time.Time   `json:"date_updated"`
	URL                     string      `json:"url"`
	Links                   struct {
		PhoneNumbers string `json:"phone_numbers"`
		ShortCodes   string `json:"short_codes"`
		Sessions     string `json:"sessions"`
	} `json:"links"`
}

// https://www.twilio.com/docs/proxy/api/proxy-webhooks

// https://www.twilio.com/docs/proxy/api/proxy-webhooks#callbackurl
// These webhooks are fired for each new interaction and are informational only.
type ProxyCallbackWebhook struct {
	OutboundResourceStatus string    `json:"outboundResourceStatus"`
	OutboundResourceType   string    `json:"outboundResourceType"`
	InteractionDateUpdated time.Time `json:"interactionDateUpdated"`
	InteractionData        string    `json:"interactionData"`
	InteractionDateCreated time.Time `json:"interactionDateCreated"`
	InboundResourceURL     string    `json:"inboundResourceUrl"`
	InteractionServiceSid  string    `json:"interactionServiceSid"`
	OutboundParticipantSid string    `json:"outboundParticipantSid"`
	InteractionType        string    `json:"interactionType"`
	InteractionAccountSid  string    `json:"interactionAccountSid"`
	InboundParticipantSid  string    `json:"inboundParticipantSid"`
	InboundResourceStatus  string    `json:"inboundResourceStatus"`
	OutboundResourceSid    string    `json:"outboundResourceSid"`
	OutboundResourceURL    string    `json:"outboundResourceUrl"`
	InboundResourceType    string    `json:"inboundResourceType"`
	InboundResourceSid     string    `json:"inboundResourceSid"`
	InteractionSessionSid  string    `json:"interactionSessionSid"`
	InteractionSid         string    `json:"interactionSid"`
}

// https://www.twilio.com/docs/proxy/api/proxy-webhooks#interceptcallbackurl
// Fires on each interaction. If responded to with a 403 to this webhook we
// will abort/block the interaction. Any other status or timeout the interaction continues
type ProxyInterceptCallbackWebhook struct {
	InteractionDateUpdated time.Time `json:"interactionDateUpdated"`
	InteractionData        string    `json:"interactionData"`
	InteractionDateCreated time.Time `json:"interactionDateCreated"`
	InboundResourceURL     string    `json:"inboundResourceUrl"`
	InteractionServiceSid  string    `json:"interactionServiceSid"`
	InteractionType        string    `json:"interactionType"`
	InteractionAccountSid  string    `json:"interactionAccountSid"`
	InboundParticipantSid  string    `json:"inboundParticipantSid"`
	InboundResourceStatus  string    `json:"inboundResourceStatus"`
	InboundResourceType    string    `json:"inboundResourceType"`
	InboundResourceSid     string    `json:"inboundResourceSid"`
	InteractionSessionSid  string    `json:"interactionSessionSid"`
	InteractionSid         string    `json:"interactionSid"`
}

// https://www.twilio.com/docs/proxy/api/proxy-webhooks#outofsessioncallbackurl
// A URL to send webhooks to when an action (inbound call or SMS) occurs where
// there is no session or a closed session. If your server (or a Twilio function)
// responds with valid TwiML, this will be processed.
// This means it is possible to e.g. play a message for a call, send an automated
// text message response, or redirect a call to another number.
type OutOfSessionCallbackWebhook struct {
	AccountSid          string    `json:"AccountSid"`
	SessionUniqueName   string    `json:"sessionUniqueName"`
	SessionAccountSid   string    `json:"sessionAccountSid"`
	SessionServiceSid   string    `json:"sessionServiceSid"`
	SessionSid          string    `json:"sessionSid"`
	SessionStatus       string    `json:"sessionStatus"`
	SessionMode         string    `json:"sessionMode"`
	SessionDateCreated  time.Time `json:"sessionDateCreated"`
	SessionDateUpdated  time.Time `json:"sessionDateUpdated"`
	SessionDateEnded    time.Time `json:"sessionDateEnded"`
	SessionClosedReason string    `json:"sessionClosedReason"`

	To          string `json:"To"`
	ToCity      string `json:"ToCity"`
	ToState     string `json:"ToState"`
	ToZip       string `json:"ToZip"`
	ToCountry   string `json:"ToCountry"`
	From        string `json:"From"`
	FromCity    string `json:"FromCity"`
	FromState   string `json:"FromState"`
	FromZip     string `json:"FromZip"`
	FromCountry string `json:"FromCountry"`

	InboundParticipantSid                string    `json:"inboundParticipantSid"`
	InboundParticipantIdentifier         string    `json:"inboundParticipantIdentifier"`
	InboundParticipantFriendlyName       string    `json:"inboundParticipantFriendlyName"`
	InboundParticipantProxyIdentifier    string    `json:"inboundParticipantProxyIdentifier"`
	InboundParticipantProxyIdentifierSid string    `json:"inboundParticipantProxyIdentifierSid"`
	InboundParticipantAccountSid         string    `json:"inboundParticipantAccountSid"`
	InboundParticipantServiceSid         string    `json:"inboundParticipantServiceSid"`
	InboundParticipantSessionSid         string    `json:"inboundParticipantSessionSid"`
	InboundParticipantDateCreated        time.Time `json:"inboundParticipantDateCreated"`
	InboundParticipantDateUpdated        time.Time `json:"inboundParticipantDateUpdated"`

	OutboundParticipantSid                string    `json:"outboundParticipantSid"`
	OutboundParticipantIdentifier         string    `json:"outboundParticipantIdentifier"`
	OutboundParticipantFriendlyName       string    `json:"outboundParticipantFriendlyName"`
	OutboundParticipantProxyIdentifier    string    `json:"outboundParticipantProxyIdentifier"`
	OutboundParticipantProxyIdentifierSid string    `json:"outboundParticipantProxyIdentifierSid"`
	OutboundParticipantAccountSid         string    `json:"outboundParticipantAccountSid"`
	OutboundParticipantServiceSid         string    `json:"outboundParticipantServiceSid"`
	OutboundParticipantSessionSid         string    `json:"outboundParticipantSessionSid"`
	OutboundParticipantDateCreated        time.Time `json:"outboundParticipantDateCreated"`
	OutboundParticipantDateUpdated        time.Time `json:"outboundParticipantDateUpdated"`

	CallSid    string `json:"CallSid"`
	CallStatus string `json:"CallStatus"`

	Caller        string `json:"Caller"`
	CallerCity    string `json:"CallerCity"`
	CallerState   string `json:"CallerState"`
	CallerZip     string `json:"CallerZip"`
	CallerCountry string `json:"CallerCountry"`

	Called        string `json:"Called"`
	CalledCity    string `json:"CalledCity"`
	CalledState   string `json:"CalledState"`
	CalledZip     string `json:"CalledZip"`
	CalledCountry string `json:"CalledCountry"`

	Direction  string `json:"Direction"`
	AddOns     string `json:"AddOns"`
	APIVersion string `json:"ApiVersion"`
}

// Create a new Twilio Service
func (twilio *Twilio) NewProxyService(service ProxyServiceRequest) (response *ProxyService, exception *Exception, err error) {

	twilioUrl := ProxyBaseUrl + "/Services"

	res, err := twilio.post(proxyServiceFormValues(service), twilioUrl)
	if err != nil {
		return response, exception, err
	}
	defer res.Body.Close()

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return response, exception, err
	}

	if res.StatusCode != http.StatusCreated {
		exception = new(Exception)
		err = json.Unmarshal(responseBody, exception)

		// We aren't checking the error because we don't actually care.
		// It's going to be passed to the client either way.
		return response, exception, err
	}

	response = new(ProxyService)
	err = json.Unmarshal(responseBody, response)
	return response, exception, err

}

func (twilio *Twilio) GetProxyService(sid string) (response *ProxyService, exception *Exception, err error) {

	twilioUrl := fmt.Sprintf("%s/%s/%s", ProxyBaseUrl, "Services", sid)

	res, err := twilio.get(twilioUrl)
	if err != nil {
		return response, exception, err
	}
	defer res.Body.Close()

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return response, exception, err
	}

	if res.StatusCode != http.StatusOK {
		exception = new(Exception)
		err = json.Unmarshal(responseBody, exception)

		// We aren't checking the error because we don't actually care.
		// It's going to be passed to the client either way.
		return response, exception, err
	}

	response = new(ProxyService)
	err = json.Unmarshal(responseBody, response)
	return response, exception, err

}

func (twilio *Twilio) UpdateProxyService(sid string, service ProxyServiceRequest) (response *ProxyService, exception *Exception, err error) {

	twilioUrl := fmt.Sprintf("%s/%s/%s", ProxyBaseUrl, "Services", sid)

	res, err := twilio.post(proxyServiceFormValues(service), twilioUrl)
	if err != nil {
		return response, exception, err
	}
	defer res.Body.Close()

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return response, exception, err
	}

	if res.StatusCode != http.StatusOK {
		exception = new(Exception)
		err = json.Unmarshal(responseBody, exception)

		// We aren't checking the error because we don't actually care.
		// It's going to be passed to the client either way.
		return response, exception, err
	}

	response = new(ProxyService)
	err = json.Unmarshal(responseBody, response)
	return response, exception, err

}

func (twilio *Twilio) DeleteProxyService(sid string) (exception *Exception, err error) {

	twilioUrl := fmt.Sprintf("%s/%s/%s", ProxyBaseUrl, "Services", sid)

	res, err := twilio.delete(twilioUrl)
	if err != nil {
		return exception, err
	}

	respBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusNoContent {
		exc := new(Exception)
		err = json.Unmarshal(respBody, exc)
		return exc, err
	}
	return nil, nil
}

// Form values initialization
func proxyServiceFormValues(service ProxyServiceRequest) url.Values {
	formValues := url.Values{}

	formValues.Set("UniqueName", service.UniqueName)
	formValues.Set("CallbackUrl", service.CallbackURL)
	formValues.Set("OutOfSessionCallbackUrl", service.OutOfSessionCallbackURL)
	formValues.Set("InterceptCallbackUrl", service.InterceptCallbackURL)

	if service.GeoMatchLevel != "" {
		formValues.Set("GeoMatchLevel", service.GeoMatchLevel)
	}
	if service.NumberSelectionBehavior != "" {
		formValues.Set("NumberSelectionBehavior", service.NumberSelectionBehavior)
	}
	formValues.Set("DefaultTtl", fmt.Sprintf("%d", service.DefaultTtl))

	return formValues
}
