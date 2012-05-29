package gotwilio

import (
	"encoding/xml"
	"net/http"
	"net/url"
	"strconv"
	"time"
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
	IfMachine            string // False, Continue or Hangup; http://www.twilio.com/docs/errors/21207
	Timeout              int    // Optional
	Record               bool   // Optional
}

type VoiceResponse struct {
	XMLName        xml.Name `xml:"TwilioResponse"`
	Sid            string   `xml:"Call>Sid"`
	DateCreated    string   `xml:"Call>DateCreated"`
	DateUpdated    string   `xml:"Call>DateUpdated"`
	ParentCallSid  string   `xml:"Call>ParentCallSid"`
	AccountSid     string   `xml:"Call>AccountSid"`
	To             string   `xml:"Call>To"`
	ToFormatted    string   `xml:"Call>ToFormatted"`
	From           string   `xml:"Call>From"`
	FromFormatted  string   `xml:"Call>FromFormatted"`
	PhoneNumberSid string   `xml:"Call>PhoneNumberSid"`
	Status         string   `xml:"Call>Status"`
	StartTime      string   `xml:"Call>StartTime"`
	EndTime        string   `xml:"Call>EndTime"`
	Duration       int      `xml:"Call>Duration"`
	Price          float32  `xml:"Call>Price"`
	Direction      string   `xml:"Call>Direction"`
	AnsweredBy     string   `xml:"Call>AnsweredBy"`
	ApiVersion     string   `xml:"Call>ApiVersion"`
	Annotation     string   `xml:"Call>Annotation"`
	ForwardedFrom  string   `xml:"Call>ForwardedFrom"`
	GroupSid       string   `xml:"Call>GroupSid"`
	CallerName     string   `xml:"Call>CallerName"`
	Uri            string   `xml:"Call>Uri"`
	// TODO: handle SubresourceUris
}

func (vr *VoiceResponse) DateCreatedAsTime() (time.Time, error) {
	return time.Parse(time.RFC1123Z, vr.DateCreated)
}

func (vr *VoiceResponse) DateUpdatedAsTime() (time.Time, error) {
	return time.Parse(time.RFC1123Z, vr.DateUpdated)
}

func (vr *VoiceResponse) StartTimeAsTime() (time.Time, error) {
	return time.Parse(time.RFC1123Z, vr.StartTime)
}

func (vr *VoiceResponse) EndTimeAsTime() (time.Time, error) {
	return time.Parse(time.RFC1123Z, vr.EndTime)
}

func NewCallbackParameters(url string) *CallbackParameters {
	return &CallbackParameters{Url: url, Timeout: 60}
}

func (twilio *Twilio) CallWithUrlCallbacks(from, to string, callbackParameters *CallbackParameters) (*VoiceResponse, *Exception, error) {
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

	return twilio.voicePost(formValues)
}

func (twilio *Twilio) CallWithApplicationCallbacks(from, to, applicationSid string) (*VoiceResponse, *Exception, error) {
	formValues := url.Values{}
	formValues.Set("From", from)
	formValues.Set("To", to)
	formValues.Set("ApplicationSid", applicationSid)

	return twilio.voicePost(formValues)
}

func (twilio *Twilio) voicePost(formValues url.Values) (*VoiceResponse, *Exception, error) {
	var voiceResponse *VoiceResponse
	var exception *Exception
	twilioUrl := twilio.BaseUrl + "/Accounts/" + twilio.AccountSid + "/Calls"

	res, err := twilio.post(formValues, twilioUrl)
	if err != nil {
		return voiceResponse, exception, err
	}
	defer res.Body.Close()

	decoder := xml.NewDecoder(res.Body)

	if res.StatusCode != http.StatusCreated {
		exception = new(Exception)
		err = decoder.Decode(exception)

		// We aren't checking the error because we don't actually care.
		// It's going to be passed to the client either way.
		return voiceResponse, exception, err
	}

	voiceResponse = new(VoiceResponse)
	err = decoder.Decode(voiceResponse)
	return voiceResponse, exception, err
}
