package gotwilio

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// These are the paramters to use when you want Twilio to use callback urls.
// See http://www.twilio.com/docs/api/rest/making-calls for more info.
type CallbackParameters struct {
	Url                                string   // Required
	Method                             string   // Optional
	FallbackUrl                        string   // Optional
	FallbackMethod                     string   // Optional
	StatusCallback                     string   // Optional
	StatusCallbackMethod               string   // Optional
	StatusCallbackEvent                []string // Optional
	SendDigits                         string   // Optional
	Timeout                            int      // Optional
	Record                             bool     // Optional
	RecordingChannels                  string   // Optional
	RecordingStatusCallback            string   // Optional
	RecordingStatusCallbackMethod      string   // Optional
	MachineDetection                   string   // Optional
	MachineDetectionTimeout            int      // Optional
	MachineDetectionSpeechThreshold    int      // Optional
	MachineDetectionSpeechEndThreshold int      // Optional
	MachineDetectionSilenceTimeout     int      // Optional
	AsyncAmd                           bool     // Optional
	AsyncAmdStatusCallback             string   // Optional
	AsyncAmdStatusCallbackMethod       string   // Optional
}

// VoiceResponse contains the details about successful voice calls.
type VoiceResponse struct {
	Sid            string  `json:"sid"`
	DateCreated    string  `json:"date_created"`
	DateUpdated    string  `json:"date_updated"`
	ParentCallSid  string  `json:"parent_call_sid"`
	AccountSid     string  `json:"account_sid"`
	To             string  `json:"to"`
	ToFormatted    string  `json:"to_formatted"`
	From           string  `json:"from"`
	FromFormatted  string  `json:"from_formatted"`
	PhoneNumberSid string  `json:"phone_number_sid"`
	Status         string  `json:"status"`
	StartTime      string  `json:"start_time"`
	EndTime        string  `json:"end_time"`
	Duration       int     `json:"duration,string"`
	PriceUnit      string  `json:"price_unit"`
	Price          *string `json:"price,omitempty"`
	Direction      string  `json:"direction"`
	AnsweredBy     string  `json:"answered_by"`
	ApiVersion     string  `json:"api_version"`
	Annotation     string  `json:"annotation"`
	ForwardedFrom  string  `json:"forwarded_from"`
	GroupSid       string  `json:"group_sid"`
	CallerName     string  `json:"caller_name"`
	Uri            string  `json:"uri"`
	// TODO: handle SubresourceUris
	// TODO: handle annotation
}

// DateCreatedAsTime returns VoiceResponse.DateCreated as a time.Time object
// instead of a string.
func (vr *VoiceResponse) DateCreatedAsTime() (time.Time, error) {
	return time.Parse(time.RFC1123Z, vr.DateCreated)
}

// DateUpdatedAsTime returns VoiceResponse.DateUpdated as a time.Time object
// instead of a string.
func (vr *VoiceResponse) DateUpdatedAsTime() (time.Time, error) {
	return time.Parse(time.RFC1123Z, vr.DateUpdated)
}

// StartTimeAsTime returns VoiceResponse.StartTime as a time.Time object
// instead of a string.
func (vr *VoiceResponse) StartTimeAsTime() (time.Time, error) {
	return time.Parse(time.RFC1123Z, vr.StartTime)
}

// EndTimeAsTime returns VoiceResponse.EndTime as a time.Time object
// instead of a string.
func (vr *VoiceResponse) EndTimeAsTime() (time.Time, error) {
	return time.Parse(time.RFC1123Z, vr.EndTime)
}

// NewCallbackParameters returns a CallbackParameters type with the specified url and
// CallbackParameters.Timeout set to 60.
func NewCallbackParameters(url string) *CallbackParameters {
	return &CallbackParameters{Url: url, Timeout: 60}
}

// GetCall uses Twilio to get information about a voice call.
// See https://www.twilio.com/docs/voice/api/call
func (twilio *Twilio) GetCall(sid string) (*VoiceResponse, *Exception, error) {
	var voiceResponse *VoiceResponse
	var exception *Exception
	twilioUrl := twilio.BaseUrl + "/Accounts/" + twilio.AccountSid + "/Calls/" + sid + ".json"

	res, err := twilio.get(twilioUrl)
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close()

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, nil, err
	}
	if res.StatusCode != http.StatusOK {
		exception = new(Exception)
		err = json.Unmarshal(responseBody, exception)
		return nil, exception, err
	}

	voiceResponse = new(VoiceResponse)
	err = json.Unmarshal(responseBody, voiceResponse)
	return voiceResponse, nil, err
}

// Place a voice call with a list of callbacks specified.
func (twilio *Twilio) CallWithUrlCallbacks(from, to string, callbackParameters *CallbackParameters) (*VoiceResponse, *Exception, error) {
	formValues := url.Values{}
	formValues.Set("From", from)
	formValues.Set("To", to)
	formValues.Set("Url", callbackParameters.Url)

	// Optional values
	if callbackParameters.Method != "" {
		formValues.Set("Method", callbackParameters.Method)
	}
	if callbackParameters.FallbackUrl != "" {
		formValues.Set("FallbackUrl", callbackParameters.FallbackUrl)
	}
	if callbackParameters.FallbackMethod != "" {
		formValues.Set("FallbackMethod", callbackParameters.FallbackMethod)
	}
	if callbackParameters.StatusCallback != "" {
		formValues.Set("StatusCallback", callbackParameters.StatusCallback)
	}
	if callbackParameters.StatusCallbackMethod != "" {
		formValues.Set("StatusCallbackMethod", callbackParameters.StatusCallbackMethod)
	}
	for _, event := range callbackParameters.StatusCallbackEvent {
		formValues.Add("StatusCallbackEvent", event)
	}
	if callbackParameters.SendDigits != "" {
		formValues.Set("SendDigits", callbackParameters.SendDigits)
	}
	if callbackParameters.Timeout != 0 {
		formValues.Set("Timeout", strconv.Itoa(callbackParameters.Timeout))
	}
	if callbackParameters.MachineDetection != "" {
		formValues.Set("MachineDetection", callbackParameters.MachineDetection)
	}
	if callbackParameters.MachineDetectionTimeout != 0 {
		formValues.Set(
			"MachineDetectionTimeout",
			strconv.Itoa(callbackParameters.MachineDetectionTimeout),
		)
	}
	if callbackParameters.MachineDetectionSpeechThreshold != 0 {
		formValues.Set(
			"MachineDetectionSpeechThreshold",
			strconv.Itoa(callbackParameters.MachineDetectionSpeechThreshold),
		)
	}
	if callbackParameters.MachineDetectionSpeechEndThreshold != 0 {
		formValues.Set(
			"MachineDetectionSpeechEndThreshold",
			strconv.Itoa(callbackParameters.MachineDetectionSpeechEndThreshold),
		)
	}
	if callbackParameters.MachineDetectionSilenceTimeout != 0 {
		formValues.Set(
			"MachineDetectionSilenceTimeout",
			strconv.Itoa(callbackParameters.MachineDetectionSilenceTimeout),
		)
	}

	if callbackParameters.Record {
		formValues.Set("Record", "true")

		if callbackParameters.RecordingChannels != "" {
			formValues.Set("RecordingChannels", callbackParameters.RecordingChannels)
		}
		if callbackParameters.RecordingStatusCallback != "" {
			formValues.Set("RecordingStatusCallback", callbackParameters.RecordingStatusCallback)
		}
		if callbackParameters.RecordingStatusCallbackMethod != "" {
			formValues.Set("RecordingStatusCallbackMethod", callbackParameters.RecordingStatusCallbackMethod)
		}
	} else {
		formValues.Set("Record", "false")
	}

	if callbackParameters.AsyncAmd {
		formValues.Set("AsyncAmd", "true")

		if callbackParameters.AsyncAmdStatusCallback != "" {
			formValues.Set("AsyncAmdStatusCallback", callbackParameters.AsyncAmdStatusCallback)
		}
		if callbackParameters.AsyncAmdStatusCallbackMethod != "" {
			formValues.Set("AsyncAmdStatusCallbackMethod", callbackParameters.AsyncAmdStatusCallbackMethod)
		}
	}

	return twilio.voicePost("Calls.json", formValues)
}

// Place a voice call with an ApplicationSid specified.
func (twilio *Twilio) CallWithApplicationCallbacks(from, to, applicationSid string) (*VoiceResponse, *Exception, error) {
	formValues := url.Values{}
	formValues.Set("From", from)
	formValues.Set("To", to)
	formValues.Set("ApplicationSid", applicationSid)

	return twilio.voicePost("Calls.json", formValues)
}

// Update an existing call
func (twilio *Twilio) CallUpdate(callSid string, formValues url.Values) (*VoiceResponse, *Exception, error) {
	return twilio.voicePost("Calls/"+callSid+".json", formValues)
}

// This is a private method that has the common bits for placing or updating a voice call.
func (twilio *Twilio) voicePost(resourcePath string, formValues url.Values) (*VoiceResponse, *Exception, error) {
	var voiceResponse *VoiceResponse
	var exception *Exception

	twilioUrl := twilio.buildUrl(resourcePath)
	res, err := twilio.post(formValues, twilioUrl)
	if err != nil {
		return voiceResponse, exception, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
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
