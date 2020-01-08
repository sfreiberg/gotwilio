// Package gotwilio is a library for interacting with http://www.twilio.com/ API.
package gotwilio

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// https://www.twilio.com/docs/proxy/api/pb-proxy-session

type ProxySessionRequest struct {
	Status       string                // default: open, other values: closed	optional
	UniqueName   string                // optional
	TTL          int                   // default: 0 (No TTL)	optional
	DateExpiry   time.Time             // optional
	Mode         string                // default: voice-and-message, other values: voice-only, message-only	optional
	Participants []*ParticipantRequest // optional, The Participant objects to include in the new session.
}

type ProxySession struct {
	Sid                 string      `json:"sid"`
	Status              string      `json:"status"`
	UniqueName          string      `json:"unique_name"`
	ClosedReason        interface{} `json:"closed_reason"`
	DateEnded           interface{} `json:"date_ended"`
	TTL                 int         `json:"ttl"`
	DateExpiry          interface{} `json:"date_expiry"`
	AccountSid          string      `json:"account_sid"`
	DateUpdated         time.Time   `json:"date_updated"`
	Mode                string      `json:"mode"`
	DateLastInteraction interface{} `json:"date_last_interaction"`
	URL                 string      `json:"url"`
	DateCreated         time.Time   `json:"date_created"`
	DateStarted         time.Time   `json:"date_started"`
	ServiceSid          string      `json:"service_sid"`
	Links               struct {
		Participants string `json:"participants"`
		Interactions string `json:"interactions"`
	} `json:"links"`

	// internal attribute
	twilio *Twilio
}

// Create a new Twilio Service
func (twilio *Twilio) NewProxySession(serviceID string, req ProxySessionRequest) (response *ProxySession, exception *Exception, err error) {

	twilioUrl := fmt.Sprintf("%s/%s/%s/%s", ProxyBaseUrl, "Services", serviceID, "Sessions")

	res, err := twilio.post(proxySessionFormValues(req), twilioUrl)
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

	response = new(ProxySession)
	err = json.Unmarshal(responseBody, response)

	// Pass connection
	response.twilio = twilio

	return response, exception, err
}

func (twilio *Twilio) GetProxySession(serviceID, sessionID string) (response *ProxySession, exception *Exception, err error) {

	twilioUrl := fmt.Sprintf("%s/%s/%s/%s/%s", ProxyBaseUrl, "Services", serviceID, "Sessions", sessionID)

	res, err := twilio.get(twilioUrl)
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

	response = new(ProxySession)
	err = json.Unmarshal(responseBody, response)

	// Pass connection
	response.twilio = twilio

	return response, exception, err
}

func (twilio *Twilio) UpdateProxySession(serviceID, sessionID string, req ProxySessionRequest) (response *ProxySession, exception *Exception, err error) {

	twilioUrl := fmt.Sprintf("%s/%s/%s/%s/%s", ProxyBaseUrl, "Services", serviceID, "Sessions", sessionID)

	formValue, err := proxySessionFormValues(req)
	if err != nil {
		return response, exception, err
	}
	res, err := twilio.post(formValue, twilioUrl)
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

	response = new(ProxySession)
	err = json.Unmarshal(responseBody, response)

	// Pass connection
	response.twilio = twilio

	return response, exception, err
}

func (twilio *Twilio) DeleteProxySession(serviceID, sessionID string) (exception *Exception, err error) {

	twilioUrl := fmt.Sprintf("%s/%s/%s/%s/%s", ProxyBaseUrl, "Services", serviceID, "Sessions", sessionID)

	res, err := twilio.delete(twilioUrl)
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

// Form values initialization
func proxySessionFormValues(req ProxySessionRequest) (url.Values, error) {
	formValues := url.Values{}

	formValues.Set("UniqueName", req.UniqueName)

	if req.Status != "" {
		formValues.Set("Status", req.Status)
	}
	formValues.Set("Ttl", fmt.Sprintf("%d", req.TTL))

	if !req.DateExpiry.IsZero() {
		formValues.Set("DateExpiry", req.DateExpiry.Format("20060101"))
	}
	if req.Mode != "" {
		formValues.Set("Mode", req.Mode)
	}
	for _, participant := range req.Participants {
		jsonByte, err := json.Marshal(participant)
		if err != nil {
			return nil, err
		}
		formValues.Set("Participants", string(jsonByte))
	}

	return formValues, nil
}
