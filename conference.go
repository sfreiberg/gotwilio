package gotwilio

import (
	"encoding/json"
	"fmt"
	"github.com/google/go-querystring/query"
	"net/http"
)

// Conference represents a Twilio Voice conference call
type Conference struct {
	Sid          string `json:"sid"`
	FriendlyName string `json:"friendly_name"`
	Status       string `json:"status"`
	Region       string `json:"region"`
}

// ConferenceOptions are used for updating Conferences
type ConferenceOptions struct {
	Status         string `url:"status,omitempty"`
	AnnounceURL    string `url:"announceURL,omitempty"`
	AnnounceMethod string `url:"announceMethod,omitempty"`
}

// ConferenceParticipant represents a Participant in responses from the Twilio API
type ConferenceParticipant struct {
	CallSid                string `json:"call_sid"`
	ConferenceSid          string `json:"conference_sid"`
	Muted                  bool   `json:"muted"`
	Hold                   bool   `json:"hold"`
	Status                 string `json:"status"`
	StartConferenceOnEnter bool   `json:"start_conference_on_enter"`
	EndConferenceOnExit    bool   `json:"end_conference_on_exit"`
	Coaching               bool   `json:"coaching"`
	CallSidToCoach         string `json:"call_sid_to_coach"`
}

// ConferenceParticipantOptions are used for creating and updating Conference Participants.
type ConferenceParticipantOptions struct {
	From                                    string `url:"From,omitempty"`
	To                                      string `url:"To,omitempty"`
	StatusCallback                          string `url:"StatusCallback,omitempty"`
	StatusCallbackMethod                    string `url:"StatusCallbackMethod,omitempty"`
	StatusCallbackEvent                     string `url:"statusCallbackEvent,omitempty"`
	Timeout                                 int    `url:"Timeout"`
	Record                                  *bool  `url:"Record,omitempty"`
	Muted                                   *bool  `url:"Muted,omitempty"`
	Beep                                    *bool  `url:"Beep,omitempty"`
	StartConferenceOnEnter                  *bool  `url:"StartConferenceOnEnter,omitempty"`
	EndConferenceOnExit                     *bool  `url:"EndConferenceOnExit,omitempty"`
	WaitURL                                 string `url:"WaitURL,omitempty"`
	WaitMethod                              string `url:"WaitMethod,omitempty"`
	EarlyMedia                              *bool  `url:"EarlyMedia,omitempty"`
	MaxParticipants                         int    `url:"MaxParticipants"`
	ConferenceRecord                        string `url:"ConferenceRecord,omitempty"`
	ConferenceTrim                          string `url:"ConferenceTrim,omitempty"`
	ConferenceStatusCallback                string `url:"ConferenceStatusCallback,omitempty"`
	ConferenceStatusCallbackMethod          string `url:"ConferenceStatusCallbackMethod,omitempty"`
	ConferenceStatusCallbackEvent           string `url:"ConferenceStatusCallbackEvent,omitempty"`
	RecordingChannels                       string `url:"RecordingChannels,omitempty"`
	RecordingStatusCallback                 string `url:"RecordingStatusCallback,omitempty"`
	RecordingStatusCallbackMethod           string `url:"RecordingStatusCallbackMethod,omitempty"`
	RecordingStatusCallbackEvent            string `url:"RecordingStatusCallbackEvent,omitempty"`
	SipAuthUsername                         string `url:"SipAuthUsername,omitempty"`
	SipAuthPassword                         string `url:"SipAuthPassword,omitempty"`
	Region                                  string `url:"Region,omitempty"`
	ConferenceRecordingStatusCallback       string `url:"ConferenceRecordingStatusCallback,omitempty"`
	ConferenceRecordingStatusCallbackMethod string `url:"ConferenceRecordingStatusCallbackMethod,omitempty"`
	Coaching                                *bool  `url:"Coaching,omitempty"`
	CallSidToCoach                          string `url:"CallSidToCoach,omitempty"`
}

// GetConference fetches details for a single conference instance
// https://www.twilio.com/docs/voice/api/conference-resource#fetch-a-conference-resource
func (twilio *Twilio) GetConference(conferenceSid string) (*Conference, *Exception, error) {
	res, err := twilio.get(twilio.buildUrl(fmt.Sprintf("Conferences/%s.json", conferenceSid)))
	if err != nil {
		return nil, nil, err
	}

	decoder := json.NewDecoder(res.Body)

	// handle NULL response
	if res.StatusCode != http.StatusOK {
		exception := new(Exception)
		err = decoder.Decode(exception)
		return nil, exception, err
	}

	conf := new(Conference)
	err = decoder.Decode(conf)
	return conf, nil, err
}

// UpdateConference to end it or play an announcement
// https://www.twilio.com/docs/voice/api/conference-resource#update-a-conference-resource
func (twilio *Twilio) UpdateConference(conferenceSid string, options *ConferenceOptions) (*Conference, *Exception, error) {
	form, err := query.Values(options)
	if err != nil {
		return nil, nil, err
	}

	res, err := twilio.post(form, twilio.buildUrl(fmt.Sprintf("Conferences/%s.json", conferenceSid)))
	if err != nil {
		return nil, nil, err
	}

	decoder := json.NewDecoder(res.Body)

	if res.StatusCode != http.StatusOK {
		exception := new(Exception)
		err = decoder.Decode(exception)
		return nil, exception, err
	}

	c := new(Conference)
	err = decoder.Decode(c)
	return c, nil, err
}

// GetConferenceParticipant fetches details for a conference participant resource
// https://www.twilio.com/docs/voice/api/conference-participant-resource#fetch-a-participant-resource
func (twilio *Twilio) GetConferenceParticipant(conferenceSid, callSid string) (*ConferenceParticipant, *Exception, error) {
	res, err := twilio.get(twilio.buildUrl(fmt.Sprintf("Conferences/%s/Participants/%s.json", conferenceSid, callSid)))
	if err != nil {
		return nil, nil, err
	}

	decoder := json.NewDecoder(res.Body)

	// handle NULL response
	if res.StatusCode != http.StatusOK {
		exception := new(Exception)
		err = decoder.Decode(exception)
		return nil, exception, err
	}

	conf := new(ConferenceParticipant)
	err = decoder.Decode(conf)
	return conf, nil, err
}

// AddConferenceParticipant adds a Participant to a conference by dialing out a new call
// https://www.twilio.com/docs/voice/api/conference-participant-resource#create-a-participant-agent-conference-only
func (twilio *Twilio) AddConferenceParticipant(conferenceSid string, participant *ConferenceParticipantOptions) (*ConferenceParticipant, *Exception, error) {
	form, err := query.Values(participant)
	if err != nil {
		return nil, nil, err
	}

	res, err := twilio.post(form, twilio.buildUrl(fmt.Sprintf("Conferences/%s/Participants.json", conferenceSid)))
	if err != nil {
		return nil, nil, err
	}

	decoder := json.NewDecoder(res.Body)

	if res.StatusCode != http.StatusCreated {
		exception := new(Exception)
		err = decoder.Decode(exception)
		return nil, exception, err
	}

	conf := new(ConferenceParticipant)
	err = decoder.Decode(conf)
	return conf, nil, err
}

// UpdateConferenceParticipant
// https://www.twilio.com/docs/voice/api/conference-participant-resource#create-a-participant-agent-conference-only
func (twilio *Twilio) UpdateConferenceParticipant(conferenceSid string, callSid string, participant *ConferenceParticipantOptions) (*ConferenceParticipant, *Exception, error) {
	form, err := query.Values(participant)
	if err != nil {
		return nil, nil, err
	}

	res, err := twilio.post(form, twilio.buildUrl(fmt.Sprintf("Conferences/%s/Participants/%s.json", conferenceSid, callSid)))
	if err != nil {
		return nil, nil, err
	}

	decoder := json.NewDecoder(res.Body)

	if res.StatusCode != http.StatusOK {
		exception := new(Exception)
		err = decoder.Decode(exception)
		return nil, exception, err
	}

	p := new(ConferenceParticipant)
	err = decoder.Decode(p)
	return p, nil, err
}

// DeleteConferenceParticipant
func (twilio *Twilio) DeleteConferenceParticipant(conferenceSid string, callSid string) (*Exception, error) {
	res, err := twilio.delete(twilio.buildUrl(fmt.Sprintf("Conferences/%s/Participants/%s.json", conferenceSid, callSid)))
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		decoder := json.NewDecoder(res.Body)
		exception := new(Exception)
		err = decoder.Decode(exception)
		return exception, err
	}

	return nil, err
}
