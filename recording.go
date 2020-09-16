package gotwilio

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/gorilla/schema"
)

const (
	EmptyEnum                 = 0
	RecordingEventInProgress  = 1
	RecordingEventCompleted   = 2
	RecordingEventAbsent      = 3
	CallbackMethodGet         = 1
	CallbackMethodPost        = 2
	TrimTypeSilence           = 1
	TrimTypeNoTrim            = 2
	ChannelMono               = 1
	ChannelDual               = 2
	RecordingStatusStopped    = 1
	RecordingStatusPaused     = 2
	RecordingStatusInProgress = 3
	RecordingPauseSkip        = 1
	RecordingPauseSilence     = 2
)

type RecordingEventType int32
type CallbackMethodType int32
type TrimType int32
type ChannelType int32
type RecordingStatus int32
type RecordingPauseBehavior int32

var recordingEventTypeEncoder = func(value reflect.Value) string {
	switch value.Interface().(RecordingEventType) {
	case RecordingEventInProgress:
		return "in-progress"
	case RecordingEventCompleted:
		return "completed"
	case RecordingEventAbsent:
		return "absent"
	case EmptyEnum:
		return ""
	default:
		return "unknown"
	}
}

var callbackMethodTypeEncoder = func(value reflect.Value) string {
	switch value.Interface().(CallbackMethodType) {
	case CallbackMethodGet:
		return "GET"
	case CallbackMethodPost:
		return "POST"
	case EmptyEnum:
		return ""
	default:
		return "unknown"
	}
}

var trimTypeEncoder = func(value reflect.Value) string {
	switch value.Interface().(TrimType) {
	case TrimTypeSilence:
		return "trim-silence"
	case TrimTypeNoTrim:
		return "do-not-trim"
	case EmptyEnum:
		return ""
	default:
		return "unknown"
	}
}

var channelTypeEncoder = func(value reflect.Value) string {
	switch value.Interface().(ChannelType) {
	case ChannelMono:
		return "mono"
	case ChannelDual:
		return "dual"
	case EmptyEnum:
		return ""
	default:
		return "unknown"
	}
}

var recordingStatusEncoder = func(value reflect.Value) string {
	switch value.Interface().(RecordingStatus) {
	case RecordingStatusStopped:
		return "stopped"
	case RecordingStatusPaused:
		return "paused"
	case RecordingStatusInProgress:
		return "in-progress"
	case EmptyEnum:
		return ""
	default:
		return "unknown"
	}
}

var recordingPauseBheaviorEncoder = func(value reflect.Value) string {
	switch value.Interface().(RecordingPauseBehavior) {
	case RecordingPauseSilence:
		return "silence"
	case RecordingPauseSkip:
		return "skip"
	case EmptyEnum:
		return ""
	default:
		return "unknown"
	}
}

// CreateRecordingParameters defined the parameters required to create recording resource
// https://www.twilio.com/docs/voice/api/recording#create-a-recording-resource
type CreateRecordingParameters struct {
	RecordingStatusCallbackEvent  []RecordingEventType `url:"recording_status_callback_event,omitempty"`
	RecordingStatusCallback       string               `url:"recording_status_callback,omitempty"`
	RecordingStatusCallbackMethod CallbackMethodType   `url:"recording_status_callback_method,omitempty"`
	Trim                          TrimType             `url:"trim,omitempty"`
	RecordingChannels             ChannelType          `url:"recording_channels,omitempty"`
}

// UpdateRecordingParameters defined the parameters required to update recording resource
// https://www.twilio.com/docs/voice/api/recording#update-a-recording-resource
type UpdateRecordingParameters struct {
	Status        RecordingStatus        `url:"Status,omitempty"`
	PauseBehavior RecordingPauseBehavior `url:"PauseBehavior,omitempty"`
}

// RecordingProperties defined the properties returned from Twilio for create or update recording resources
// Currently only define subset of properties the full list can be found the link below
// https://www.twilio.com/docs/voice/api/recording#recording-properties
type RecordingProperties struct {
	Sid           string `json:"sid"`
	AccountSid    string `json:"account_sid"`
	CallSid       string `json:"call_sid"`
	ConferenceSid string `json:"conference_sid"`
	Status        string `json:"status"`
}

// CreateRecording is to create recording resource at twilio for a given call
func (twilio *Twilio) CreateRecording(acctSid, callSid string, params *CreateRecordingParameters) (*RecordingProperties, *Exception, error) {
	encoder := schema.NewEncoder()
	if len(params.RecordingStatusCallbackEvent) > 0 {
		encoder.RegisterEncoder(params.RecordingStatusCallbackEvent[0], recordingEventTypeEncoder)
	}
	encoder.RegisterEncoder(params.RecordingChannels, channelTypeEncoder)
	encoder.RegisterEncoder(params.RecordingStatusCallbackMethod, callbackMethodTypeEncoder)
	encoder.RegisterEncoder(params.Trim, trimTypeEncoder)
	encoder.SetAliasTag("url")

	values := make(map[string][]string)
	err := encoder.Encode(params, values)
	if err != nil {
		return nil, nil, err
	}

	res, err := twilio.post(values, twilio.buildUrlWithAcct(acctSid, fmt.Sprintf("Calls/%s/Recordings.json", callSid)))
	if err != nil {
		return nil, nil, err
	}

	decoder := json.NewDecoder(res.Body)

	if res.StatusCode != http.StatusOK {
		exception := new(Exception)
		err = decoder.Decode(exception)
		return nil, exception, err
	}

	r := new(RecordingProperties)
	err = decoder.Decode(r)
	return r, nil, err
}

// UpdateRecording is to update recording resource at twilio for a given recording resource
func (twilio *Twilio) UpdateRecording(acctSid, callSid, sid string, params *UpdateRecordingParameters) (*RecordingProperties, *Exception, error) {
	encoder := schema.NewEncoder()
	encoder.RegisterEncoder(params.PauseBehavior, recordingPauseBheaviorEncoder)
	encoder.RegisterEncoder(params.Status, recordingStatusEncoder)
	encoder.SetAliasTag("url")

	values := make(map[string][]string)
	err := encoder.Encode(params, values)
	if err != nil {
		return nil, nil, err
	}

	res, err := twilio.post(values, twilio.buildUrlWithAcct(acctSid, fmt.Sprintf("Calls/%s/Recordings/%s.json", callSid, sid)))
	if err != nil {
		return nil, nil, err
	}

	decoder := json.NewDecoder(res.Body)

	if res.StatusCode != http.StatusOK {
		exception := new(Exception)
		err = decoder.Decode(exception)
		return nil, exception, err
	}

	r := new(RecordingProperties)
	err = decoder.Decode(r)
	return r, nil, err
}
