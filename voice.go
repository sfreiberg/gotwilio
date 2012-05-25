package gotwilio

import (
	"net/url"
	"strconv"
)

// These are the paramters to use when you want Twilio to use callback urls.
// See http://www.twilio.com/docs/api/rest/making-calls for more info.
type CallbackParameters struct {
	Url                  string // Required
	Method               string // Optional
	FallbackUrl          string // Optional
	FallbackMethod       string // Optional
	StatusCallback       string // Optional
	StatusCallbackMethod string // Optional
	SendDigits           string // Optional
	// http://www.twilio.com/docs/errors/21207
	IfMachine string // False, Continue or Hangup
	Timeout   int    // Optional
	Record    bool   // Optional
}

func NewCallbackParameters(url string) *CallbackParameters {
	return &CallbackParameters{Url: url, Timeout: 60}
}

func (twilio *Twilio) CallWithUrlCallbacks(from, to string, callbackParameters *CallbackParameters) (string, error) {
	twilioUrl := twilio.BaseUrl + "/Accounts/" + twilio.AccountSid + "/Calls.json"

	formValues := url.Values{}
	formValues.Set("From", from)
	formValues.Set("To", to)
	formValues.Set("Url", callbackParameters.Url)
	formValues.Set("Method", callbackParameters.Method)
	formValues.Set("FallbackUrl", callbackParameters.FallbackUrl)
	formValues.Set("FallbackMethod", callbackParameters.FallbackMethod)
	formValues.Set("StatusCallback", callbackParameters.StatusCallback)
	formValues.Set("StatusCallbackMethod", callbackParameters.StatusCallbackMethod)
	formValues.Set("SendDigits", callbackParameters.SendDigits)
	formValues.Set("IfMachine", callbackParameters.IfMachine)
	formValues.Set("Timeout", strconv.Itoa(callbackParameters.Timeout))
	if callbackParameters.Record {
		formValues.Set("Record", "true")
	} else {
		formValues.Set("Record", "false")
	}

	return twilio.post(formValues, twilioUrl)
}

// TODO: Needs a better function name
func (twilio *Twilio) CallWithApplicationCallbacks(from, to, applicationSid string) (string, error) {
	twilioUrl := twilio.BaseUrl + "/Accounts/" + twilio.AccountSid + "/Calls.json"
	formValues := url.Values{}
	formValues.Set("From", from)
	formValues.Set("To", to)
	formValues.Set("ApplicationSid", applicationSid)
	return twilio.post(formValues, twilioUrl)
}
