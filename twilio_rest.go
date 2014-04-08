package gotwilio

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// Twilio stores basic information important for connecting to the
// twilio.com REST api such as AccountSid and AuthToken.
type Twilio struct {
	AccountSid, AuthToken, BaseUrl string
}

// Exception is a representation of a twilio exception.
type Exception struct {
	Status   int    `json:"status"`    // HTTP specific error code
	Message  string `json:"message"`   // HTTP error message
	Code     int    `json:"code"`      // Twilio specific error code
	MoreInfo string `json:"more_info"` // Additional info from Twilio
}

func (exception *Exception) Error() string {
	return fmt.Sprintf("%s: status code: %d, error code: %d, info: %s", exception.Message, exception.Status, exception.Code, exception.MoreInfo)
}

const twilioUrl = "https://api.twilio.com/2010-04-01"

// NewTwilioClient creates a new Twilio struct from provided credentials.
// Not recommended for use in public code, see TwilioClientFromEnvironment
func NewTwilioClient(accountSid, authToken string) *Twilio {
	return &Twilio{accountSid, authToken, twilioUrl}
}

// NewTwilioClientFromEnvironment creates a new Twilio struct from environment variables.
// Recommended for use in public code
func NewTwilioClientFromEnv() (*Twilio, error) {
	accountSid := os.Getenv("TWILIO_ACCOUNT_SID")
	authToken := os.Getenv("TWILIO_AUTH_TOKEN")
	if accountSid != "" && authToken != "" {
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

func (twilio *Twilio) get(twilioUrl string) (*http.Response, error) {
	req, err := http.NewRequest("GET", twilioUrl, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(twilio.AccountSid, twilio.AuthToken)

	client := &http.Client{}
	return client.Do(req)
}
