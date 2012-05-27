// Package gotwilio is a library for interacting with http://www.twilio.com/ API.
package gotwilio

import (
	"encoding/xml"
	"net/http"
	"net/url"
	"strings"
)

type Twilio struct {
	AccountSid string
	AuthToken  string
	BaseUrl    string
}

type Exception struct {
	XMLName  xml.Name `xml:"TwilioResponse"`
	Status   int      `xml:"RestException>Status"`
	Message  string   `xml:"RestException>Message"`
	Code     int      `xml:"RestException>Code"`
	MoreInfo string   `xml:"RestException>MoreInfo"`
}

func NewTwilioClient(accountSid, authToken string) *Twilio {
	twilioUrl := "https://api.twilio.com/2010-04-01" // Should this be moved into a constant?
	return &Twilio{accountSid, authToken, twilioUrl}
}

func (twilio *Twilio) post(formValues url.Values, twilioUrl string) (*http.Response, error) {
	req, err := http.NewRequest("POST", twilioUrl, strings.NewReader(formValues.Encode()))
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(twilio.AccountSid, twilio.AuthToken)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	return client.Do(req)
}
