// Package gotwilio is a library for interacting with http://www.twilio.com/ API.
package gotwilio

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// https://www.twilio.com/docs/proxy/api/participants

type ParticipantRequest struct {
	Identifier      string // mandatory
	ProxyIdentifier string // optional
	FriendlyName    string // optional
}

type Participant struct {
	Sid                string      `json:"sid"`
	Identifier         string      `json:"identifier"`
	DateUpdated        time.Time   `json:"date_updated"`
	FriendlyName       interface{} `json:"friendly_name"`
	DateDeleted        interface{} `json:"date_deleted"`
	AccountSid         string      `json:"account_sid"`
	URL                string      `json:"url"`
	ProxyIdentifier    string      `json:"proxy_identifier"`
	ProxyIdentifierSid string      `json:"proxy_identifier_sid"`
	DateCreated        time.Time   `json:"date_created"`
	ParticipantSid     string      `json:"Participant_sid"`
	ServiceSid         string      `json:"service_sid"`
	Links              struct {
		MessageInteractions string `json:"message_interactions"`
	} `json:"links"`
}

type ParticipantList struct {
	Participants []Participant `json:"participants"`
	Meta         Meta          `json:"meta"`
}

type Interaction struct {
	InboundResourceStatus  string    `json:"inbound_resource_status"`
	InboundResourceSid     string    `json:"inbound_resource_sid"`
	OutboundResourceStatus string    `json:"outbound_resource_status"`
	OutboundResourceSid    string    `json:"outbound_resource_sid"`
	InboundResourceURL     string    `json:"inbound_resource_url"`
	Type                   string    `json:"type"`
	AccountSid             string    `json:"account_sid"`
	OutboundResourceType   string    `json:"outbound_resource_type"`
	DateCreated            time.Time `json:"date_created"`
	InboundResourceType    string    `json:"inbound_resource_type"`
	URL                    string    `json:"url"`
	DateUpdated            time.Time `json:"date_updated"`
	Sid                    string    `json:"sid"`
	OutboundParticipantSid string    `json:"outbound_participant_sid"`
	OutboundResourceURL    string    `json:"outbound_resource_url"`
	SessionSid             string    `json:"session_sid"`
	ServiceSid             string    `json:"service_sid"`
	Data                   string    `json:"data"`
	InboundParticipantSid  string    `json:"inbound_participant_sid"`
}

type InteractionList struct {
	Meta         Meta          `json:"meta"`
	Interactions []Interaction `json:"interactions"`
}

type Meta struct {
	Page            int         `json:"page"`
	PageSize        int         `json:"page_size"`
	FirstPageURL    string      `json:"first_page_url"`
	PreviousPageURL interface{} `json:"previous_page_url"`
	URL             string      `json:"url"`
	NextPageURL     interface{} `json:"next_page_url"`
	Key             string      `json:"key"`
}

type ProxyMessage struct {
	Body     string // The text to send to the participant
	MediaUrl string // Url to media file
	Callback string // Webhook url to handle processed record.
}

// Add Participant to Session
func (session *ProxySession) AddParticipant(req ParticipantRequest) (response Participant, exception *Exception, err error) {

	twilioUrl := fmt.Sprintf("%s/%s/%s/%s/%s/%s", ProxyBaseUrl, "Services", session.ServiceSid, "Sessions", session.Sid, "Participants")

	res, err := session.twilio.post(participantFormValues(req), twilioUrl)
	if err != nil {
		return response, exception, err
	}
	defer res.Body.Close()

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return response, exception, err
	}

	if res.StatusCode != http.StatusCreated {
		exception = new(Exception)
		err = json.Unmarshal(responseBody, exception)

		// We aren't checking the error because we don't actually care.
		// It's going to be passed to the client either way.
		return response, exception, err
	}

	err = json.Unmarshal(responseBody, &response)
	return response, exception, err

}

func (session *ProxySession) ListParticipants() (response []Participant, exception *Exception, err error) {

	twilioUrl := fmt.Sprintf("%s/%s/%s/%s/%s/%s", ProxyBaseUrl, "Services", session.ServiceSid, "Sessions", session.Sid, "Participants")

	res, err := session.twilio.get(twilioUrl)
	if err != nil {
		return response, exception, err
	}
	defer res.Body.Close()

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return response, exception, err
	}

	if res.StatusCode != http.StatusOK {
		exception = new(Exception)
		err = json.Unmarshal(responseBody, exception)

		// We aren't checking the error because we don't actually care.
		// It's going to be passed to the client either way.
		return response, exception, err
	}

	list := ParticipantList{}
	err = json.Unmarshal(responseBody, &list)
	return list.Participants, exception, err

}

func (session *ProxySession) GetParticipant(participantID string) (response Participant, exception *Exception, err error) {

	twilioUrl := fmt.Sprintf("%s/%s/%s/%s/%s/%s/%s", ProxyBaseUrl, "Services", session.ServiceSid, "Sessions", session.Sid, "Participants", participantID)

	res, err := session.twilio.get(twilioUrl)
	if err != nil {
		return response, exception, err
	}
	defer res.Body.Close()

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return response, exception, err
	}

	if res.StatusCode != http.StatusOK {
		exception = new(Exception)
		err = json.Unmarshal(responseBody, exception)

		// We aren't checking the error because we don't actually care.
		// It's going to be passed to the client either way.
		return response, exception, err
	}

	err = json.Unmarshal(responseBody, &response)
	return response, exception, err

}

// Participants cannot be changed once added. To add a new Participant, delete a Participant and add a new one.

func (session *ProxySession) DeleteParticipant(participantID string) (exception *Exception, err error) {

	twilioUrl := fmt.Sprintf("%s/%s/%s/%s/%s/%s/%s", ProxyBaseUrl, "Services", session.ServiceSid, "Sessions", session.Sid, "Participants", participantID)

	res, err := session.twilio.delete(twilioUrl)
	if err != nil {
		return exception, err
	}

	respBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusNoContent {
		exc := new(Exception)
		err = json.Unmarshal(respBody, exc)
		return exc, err
	}
	return nil, nil
}

//// INTERACTIONS

func (session *ProxySession) CreateInteraction(participantSid string, msg ProxyMessage) (response Interaction, exception *Exception, err error) {
	if msg.Body == "" {
		return response, exception, errors.New("Message Body Must exist")
	}

	twilioUrl := fmt.Sprintf("%s/%s/%s/%s/%s/%s/%s/%s", ProxyBaseUrl, "Services", session.ServiceSid, "Sessions", session.Sid, "Participants", participantSid, "MessageInteractions")

	formValues := url.Values{}
	formValues.Set("Body", msg.Body)

	if msg.MediaUrl != "" {
		formValues.Set("MediaUrl", msg.MediaUrl)
	}
	if msg.Callback != "" {
		formValues.Set("Callback", msg.Callback)
	}

	res, err := session.twilio.post(formValues, twilioUrl)
	if err != nil {
		return response, exception, err
	}
	defer res.Body.Close()

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return response, exception, err
	}

	if res.StatusCode != http.StatusOK {
		exception = new(Exception)
		err = json.Unmarshal(responseBody, exception)

		// We aren't checking the error because we don't actually care.
		// It's going to be passed to the client either way.
		return response, exception, err
	}

	err = json.Unmarshal(responseBody, &response)
	return response, exception, err
}

func (session *ProxySession) GetInteractions() (response InteractionList, exception *Exception, err error) {

	twilioUrl := fmt.Sprintf("%s/%s/%s/%s/%s/%s", ProxyBaseUrl, "Services", session.ServiceSid, "Sessions", session.Sid, "Interactions")

	res, err := session.twilio.get(twilioUrl)
	if err != nil {
		return response, exception, err
	}
	defer res.Body.Close()

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return response, exception, err
	}

	if res.StatusCode != http.StatusOK {
		exception = new(Exception)
		err = json.Unmarshal(responseBody, exception)

		// We aren't checking the error because we don't actually care.
		// It's going to be passed to the client either way.
		return response, exception, err
	}

	err = json.Unmarshal(responseBody, &response)
	return response, exception, err
}

// Form values initialization
func participantFormValues(req ParticipantRequest) url.Values {
	formValues := url.Values{}

	formValues.Set("Identifier", req.Identifier)

	if req.ProxyIdentifier != "" {
		formValues.Set("ProxyIdentifier", req.ProxyIdentifier)
	}

	if req.FriendlyName != "" {
		formValues.Set("FriendlyName", req.FriendlyName)
	}

	return formValues
}
