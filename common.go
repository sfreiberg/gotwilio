package gotwilio

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// ListResources contains Twilio paging information.
// See https://www.twilio.com/docs/usage/twilios-response#response-formats-list-paging-information
type ListResources struct {
	Uri             string `json:"uri"`
	FirstPageUri    string `json:"first_page_uri"`
	NextPageUri     string `json:"next_page_uri"`
	PreviousPageUri string `json:"previous_page_uri"`
	Page            uint   `json:"page"`
	PageSize        uint   `json:"page_size"`

	Faxes    []*FaxResource `json:"faxes"`
	Messages []*SmsResponse `json:"messages"`

	t *Twilio
}

func (t *Twilio) newListResources() *ListResources {
	lr := new(ListResources)
	lr.t = t
	return lr
}

func (l *ListResources) hasNext() bool {
	return l.NextPageUri != ""
}

func (l *ListResources) next() (*Exception, error) {
	resp, err := l.t.get(l.NextPageUri)
	if err != nil {
		return nil, err
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		exc := new(Exception)
		err = json.Unmarshal(respBody, exc)
		return exc, err
	}
	return nil, json.Unmarshal(respBody, l)
}
