package gotwilio

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
	"syscall"
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
	Status   int    `json:"status"`    // HTTP specific error code
	Message  string `json:"message"`   // HTTP error message
	Code     int    `json:"code"`      // Twilio specific error code
	MoreInfo string `json:"more_info"` // Additional info from Twilio
}

const twilioUrl = "https://api.twilio.com/2010-04-01"

// Create a new Twilio struct from provided credentials.
func NewTwilioClient(accountSid, authToken string) *Twilio {
	return &Twilio{accountSid, authToken, twilioUrl}
}

func NewTwilioClientFromEnvironment() (*Twilio, error) {
	accountSid, sidFound := syscall.Getenv("TWILIO_ACCOUNT_SID")
	authToken, authFound := syscall.Getenv("TWILIO_AUTH_TOKEN")
	if sidFound && authFound {
		return &Twilio{accountSid, authToken, twilioUrl}, nil
	} else {
		return nil, errors.New("Could not find required environment variables")
	}
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
