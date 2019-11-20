package gotwilio

import (
	"encoding/json"
	"net/http"
)

// Conference represents a Twilio Voice conference call
type Conference struct {
	Sid          string `json:"sid"`
	FriendlyName string `json:"friendly_name"`
	Status       string `json:"status"`
	Region       string `json:"region"`
}

// ConferenceParticipant refers to a participant in a Twilio Voice Conference.
type ConferenceParticipant struct {
	From                 string `json:"from"`
	To                   string `json:"to"`
	StatusCallback       string `json:"status_callback"`
	StatusCallbackMethod string `json:"status_callback_method"`
	StatusCallbackEvent  string `json:"status_callback_event,omitempty"`
	Timeout              int    `json:"timeout"`
	Record               bool   `json:"record"`
	Muted                bool   `json:"muted"`
}

// GetConference fetches details for a single conference instance
// https://www.twilio.com/docs/voice/api/conference-resource#fetch-a-conference-resource
func (twilio *Twilio) GetConference(conferenceSid string) (*Conference, *Exception, error) {
	res, err := twilio.get(twilio.buildUrl("Conferences/" + conferenceSid + ".json"))
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
