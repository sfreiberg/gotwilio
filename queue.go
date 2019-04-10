package gotwilio

import (
	"encoding/json"
	"net/http"
	"net/url"
)

const (
	ErrorQueueAlreadyExists ExceptionCode = 22003
)

type QueueResponse struct {
	Sid          string `json:"sid"`
	FriendlyName string `json:"friendly_name"`
	MaxSize      int    `json:"max_size"`
}

func (twilio *Twilio) CreateQueue(friendlyName string) (*QueueResponse, *Exception, error) {
	var queueResponse *QueueResponse
	var exception *Exception
	twilioUrl := twilio.buildUrl("Queues.json")

	formValues := url.Values{}
	formValues.Set("FriendlyName", friendlyName)

	res, err := twilio.post(formValues, twilioUrl)
	if err != nil {
		return queueResponse, exception, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	if res.StatusCode != http.StatusCreated {
		exception = new(Exception)
		err = decoder.Decode(exception)

		// We aren't checking the error because we don't actually care.
		// It's going to be passed to the client either way.
		return queueResponse, exception, err
	}

	queueResponse = new(QueueResponse)
	err = decoder.Decode(queueResponse)
	return queueResponse, exception, err
}
