package gotwilio

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// SmsResponse is returned after a text/sms message is posted to Twilio
type SmsResponse struct {
	Sid         string   `json:"sid"`
	DateCreated string   `json:"date_created"`
	DateUpdate  string   `json:"date_updated"`
	DateSent    string   `json:"date_sent"`
	AccountSid  string   `json:"account_sid"`
	To          string   `json:"to"`
	From        string   `json:"from"`
	Body        string   `json:"body"`
	Status      string   `json:"status"`
	Direction   string   `json:"direction"`
	ApiVersion  string   `json:"api_version"`
	Price       *float32 `json:"price,omitempty"`
	Url         string   `json:"uri"`
}

// Returns SmsResponse.DateCreated as a time.Time object
// instead of a string.
func (sms *SmsResponse) DateCreatedAsTime() (time.Time, error) {
	return time.Parse(time.RFC1123Z, sms.DateCreated)
}

// Returns SmsResponse.DateUpdate as a time.Time object
// instead of a string.
func (sms *SmsResponse) DateUpdateAsTime() (time.Time, error) {
	return time.Parse(time.RFC1123Z, sms.DateUpdate)
}

// Returns SmsResponse.DateSent as a time.Time object
// instead of a string.
func (sms *SmsResponse) DateSentAsTime() (time.Time, error) {
	return time.Parse(time.RFC1123Z, sms.DateSent)
}

// SendSMS uses Twilio to send a text message.
// See http://www.twilio.com/docs/api/rest/sending-sms for more information.
func (twilio *Twilio) SendSMS(from, to, body, statusCallback, applicationSid string) (*SmsResponse, *Exception, error) {
	var smsResponse *SmsResponse
	var exception *Exception
	twilioUrl := twilio.BaseUrl + "/Accounts/" + twilio.AccountSid + "/Messages.json"

	formValues := url.Values{}
	formValues.Set("From", from)
	formValues.Set("To", to)
	formValues.Set("Body", body)
	if statusCallback != "" {
		formValues.Set("StatusCallback", statusCallback)
	}
	if applicationSid != "" {
		formValues.Set("ApplicationSid", applicationSid)
	}

	res, err := twilio.post(formValues, twilioUrl)
	if err != nil {
		return smsResponse, exception, err
	}
	defer res.Body.Close()

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return smsResponse, exception, err
	}

	if res.StatusCode != http.StatusCreated {
		exception = new(Exception)
		err = json.Unmarshal(responseBody, exception)

		// We aren't checking the error because we don't actually care.
		// It's going to be passed to the client either way.
		return smsResponse, exception, err
	}

	smsResponse = new(SmsResponse)
	err = json.Unmarshal(responseBody, smsResponse)
	return smsResponse, exception, err
}
