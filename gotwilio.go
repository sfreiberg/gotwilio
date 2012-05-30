// Package gotwilio is a library for interacting with http://www.twilio.com/ API.
package gotwilio

import (
	"encoding/xml"
	"net/http"
	"net/url"
	"strings"
)

// Twilio stores basic information important for connecting to the
// twilio.com REST api such as AccountSid and AuthToken.
type Twilio struct {
	AccountSid string
	AuthToken  string
	BaseUrl    string
}

// Exception is a representation of a twilio exception.
type Exception struct {
	XMLName  xml.Name `xml:"TwilioResponse"`
	Status   int      `xml:"RestException>Status"` // HTTP specific error code
	Message  string   `xml:"RestException>Message"` // HTTP error message
	Code     int      `xml:"RestException>Code"` // Twilio specific error code
	MoreInfo string   `xml:"RestException>MoreInfo"` // Additional info from Twilio
}

// Create a new Twilio struct.
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
