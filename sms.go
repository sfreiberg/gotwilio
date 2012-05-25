package gotwilio

import (
	"net/url"
)

type SmsResponse struct {
	Sid         string
	DateCreated string // TODO: Change this to date type if possible
	DateUpdate  string // TODO: Change this to date type if possible
	DateSent    string // TODO: Change this to date type if possible
	AccountSid  string
	To          string
	From        string
	Body        string
	Status      string
	Direction   string
	ApiVersion  string
	Price       string // TODO: need to find out what this returns. My example is null
	Url         string
}

// SendTextMessage uses Twilio to send a text message.
// See http://www.twilio.com/docs/api/rest/sending-sms for more information.
func (twilio *Twilio) SendSMS(from, to, body, statusCallback, applicationSid string) (string, error) {
	twilioUrl := twilio.BaseUrl + "/Accounts/" + twilio.AccountSid + "/SMS/Messages.json" // needs a better variable name

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

	return twilio.post(formValues, twilioUrl)
}
