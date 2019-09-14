// Package gotwilio is a library for interacting with http://www.twilio.com/ API.
package gotwilio

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

const (
	baseURL       = "https://api.twilio.com/2010-04-01"
	videoURL      = "https://video.twilio.com"
	lookupURL     = "https://lookups.twilio.com/v1" // https://www.twilio.com/docs/lookup/api
	clientTimeout = time.Second * 30
)

// The default http.Client that is used if none is specified
var defaultClient = &http.Client{
	Timeout: time.Second * 30,
}

// Twilio stores basic information important for connecting to the
// twilio.com REST api such as AccountSid and AuthToken.
type Twilio struct {
	AccountSid string
	AuthToken  string
	BaseUrl    string
	VideoUrl   string
	LookupURL  string
	HTTPClient *http.Client

	APIKeySid    string
	APIKeySecret string
}

// Exception is a representation of a twilio exception.
type Exception struct {
	Status   int           `json:"status"`    // HTTP specific error code
	Message  string        `json:"message"`   // HTTP error message
	Code     ExceptionCode `json:"code"`      // Twilio specific error code
	MoreInfo string        `json:"more_info"` // Additional info from Twilio
}

// Print the RESTException in a human-readable form.
func (r Exception) Error() string {
	var errorCode ExceptionCode
	var status int
	if r.Code != errorCode {
		return fmt.Sprintf("Code %d: %s", r.Code, r.Message)
	} else if r.Status != status {
		return fmt.Sprintf("Status %d: %s", r.Status, r.Message)
	}
	return r.Message
}

// Create a new Twilio struct.
func NewTwilioClient(accountSid, authToken string) *Twilio {
	return NewTwilioClientCustomHTTP(accountSid, authToken, nil)
}

// Create a new Twilio client, optionally using a custom http.Client
func NewTwilioClientCustomHTTP(accountSid, authToken string, HTTPClient *http.Client) *Twilio {
	if HTTPClient == nil {
		HTTPClient = defaultClient
	}

	return &Twilio{
		AccountSid: accountSid,
		AuthToken:  authToken,
		BaseUrl:    baseURL,
		VideoUrl:   videoURL,
		LookupURL:  lookupURL,
		HTTPClient: HTTPClient,
	}
}

func (twilio *Twilio) WithAPIKey(apiKeySid string, apiKeySecret string) *Twilio {
	twilio.APIKeySid = apiKeySid
	twilio.APIKeySecret = apiKeySecret
	return twilio
}

func (twilio *Twilio) getJSON(url string, result interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.SetBasicAuth(twilio.getBasicAuthCredentials())
	resp, err := twilio.do(req)
	if err != nil {
		return fmt.Errorf("failed to submit HTTP request: %v", err)
	}

	if resp.StatusCode != 200 {
		re := Exception{}
		json.NewDecoder(resp.Body).Decode(&re)
		return re
	}
	return json.NewDecoder(resp.Body).Decode(&result)
}

func (twilio *Twilio) getBasicAuthCredentials() (string, string) {
	if twilio.APIKeySid != "" {
		return twilio.APIKeySid, twilio.APIKeySecret
	}

	return twilio.AccountSid, twilio.AuthToken
}

func (twilio *Twilio) post(formValues url.Values, twilioUrl string) (*http.Response, error) {
	req, err := http.NewRequest("POST", twilioUrl, strings.NewReader(formValues.Encode()))
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(twilio.getBasicAuthCredentials())
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	return twilio.do(req)
}

func (twilio *Twilio) get(twilioUrl string) (*http.Response, error) {
	req, err := http.NewRequest("GET", twilioUrl, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(twilio.getBasicAuthCredentials())

	return twilio.do(req)
}

func (twilio *Twilio) delete(twilioUrl string) (*http.Response, error) {
	req, err := http.NewRequest("DELETE", twilioUrl, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(twilio.getBasicAuthCredentials())

	return twilio.do(req)
}

func (twilio *Twilio) do(req *http.Request) (*http.Response, error) {
	client := twilio.HTTPClient
	if client == nil {
		client = defaultClient
	}

	return client.Do(req)
}

// Build path to a resource within the Twilio account
func (twilio *Twilio) buildUrl(resourcePath string) string {
	return twilio.BaseUrl + "/" + path.Join("Accounts", twilio.AccountSid, resourcePath)
}
