// Package gotwilio is a library for interacting with http://www.twilio.com/ API.
package gotwilio

import (
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	baseURL       = "https://api.twilio.com/2010-04-01"
	clientTimeout = time.Second * 30
)

// Twilio stores basic information important for connecting to the
// twilio.com REST api such as AccountSid and AuthToken.
type Twilio struct {
	AccountSid string
	AuthToken  string
	BaseUrl    string
	HTTPClient *http.Client
}

// Exception is a representation of a twilio exception.
type Exception struct {
	Status   int    `json:"status"`    // HTTP specific error code
	Message  string `json:"message"`   // HTTP error message
	Code     int    `json:"code"`      // Twilio specific error code
	MoreInfo string `json:"more_info"` // Additional info from Twilio
}

// Create a new Twilio struct.
func NewTwilioClient(accountSid, authToken string) *Twilio {
	return NewTwilioClientCustomHTTP(accountSid, authToken, nil)
}

// Create a new Twilio client, optionally using a custom http.Client
func NewTwilioClientCustomHTTP(accountSid, authToken string, HTTPClient *http.Client) *Twilio {
	return &Twilio{accountSid, authToken, baseURL, HTTPClient}
}

func (twilio *Twilio) post(formValues url.Values, twilioUrl string) (*http.Response, error) {
	req, err := http.NewRequest("POST", twilioUrl, strings.NewReader(formValues.Encode()))
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(twilio.AccountSid, twilio.AuthToken)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	return twilio.do(req)
}

func (twilio *Twilio) get(twilioUrl string) (*http.Response, error) {
	req, err := http.NewRequest("GET", twilioUrl, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(twilio.AccountSid, twilio.AuthToken)

	return twilio.do(req)
}

func (twilio *Twilio) delete(twilioUrl string) (*http.Response, error) {
	req, err := http.NewRequest("DELETE", twilioUrl, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(twilio.AccountSid, twilio.AuthToken)

	return twilio.do(req)
}

func (twilio *Twilio) do(req *http.Request) (*http.Response, error) {
	client := twilio.HTTPClient
	if client == nil {
		client = defaultClient()
	}

	return client.Do(req)
}

func defaultClient() *http.Client {
	client := http.Client{
		Timeout: clientTimeout,
	}

	return &client
}
