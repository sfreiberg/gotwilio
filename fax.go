package gotwilio

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type FaxBase struct {
	Sid         string `json:"sid"`
	AccountSid  string `json:"account_sid"`
	DateCreated string `json:"date_created"`
	DateUpdated string `json:"date_updated"`
}

// DateCreatedAsTime returns FaxBase.DateCreated as a time.Time object
// instead of a string.
func (d *FaxBase) DateCreatedAsTime() (time.Time, error) {
	return time.Parse(time.RFC1123Z, d.DateCreated)
}

// DateUpdatesAsTime returns FaxBase.DateUpdated as a time.Time object
// instead of a string.
func (d *FaxBase) DateUpdatesAsTime() (time.Time, error) {
	return time.Parse(time.RFC1123Z, d.DateUpdated)
}

type FaxMediaResource struct {
	FaxBase
	FaxSid      string `json:"fax_sid"`
	ContentType string `json:"content_type"`
	Url         string `json:"url"`
}

type FaxResource struct {
	FaxBase
	From       string  `json:"from"`
	To         string  `json:"to"`
	Direction  string  `json:"direction"`
	NumPages   uint    `json:"num_pages,string"`
	Duration   uint    `json:"duration,string"`
	MediaSid   string  `json:"media_sid"`
	MediaUrl   string  `json:"media_url"`
	Status     string  `json:"status"`
	Quality    string  `json:"quality"`
	ApiVersion string  `json:"api_version"`
	Price      *string `json:"price,omitempty"`
	PriceUnit  *string `json:"price_unit,omitempty"`
}

type FaxResourcesList struct {
	ListResources
	Faxes []*FaxResource `json:"faxes"`
}

func (t *Twilio) CancelFax(faxSid string) (*Exception, error) {
	resp, err := t.post(url.Values{"Status": []string{"cancelled"}}, "https://fax.twilio.com/v1/Faxes/"+faxSid)
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
	return nil, nil
}

func (t *Twilio) DeleteFax(faxSid string) (*Exception, error) {
	resp, err := t.delete("https://fax.twilio.com/v1/Faxes/" + faxSid)
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
	return nil, nil
}

func (t *Twilio) GetFax(faxSid string) (*FaxResource, *Exception, error) {
	resp, err := t.get("https://fax.twilio.com/v1/Faxes/" + faxSid)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	if resp.StatusCode != http.StatusOK {
		exc := new(Exception)
		err = json.Unmarshal(respBody, exc)
		return nil, exc, err
	}

	fr := new(FaxResource)
	if err := json.Unmarshal(respBody, fr); err != nil {
		return nil, nil, err
	}
	return fr, nil, nil
}

// GetFaxes gets faxes for a Twilio account.
// See https://www.twilio.com/docs/fax/api/faxes#fax-list-resource
func (t *Twilio) GetFaxes(to, from, createdOnOrBefore, createdAfter string) ([]*FaxResource, *Exception, error) {
	values := url.Values{}
	if to != "" {
		values.Set("To", to)
	}
	if from != "" {
		values.Set("From", from)
	}
	if createdOnOrBefore != "" {
		values.Set("DateCreatedOnOrBefore", createdOnOrBefore)
	}
	if createdAfter != "" {
		values.Set("DateCreatedAfter", createdAfter)
	}

	resp, err := t.get("https://fax.twilio.com/v1/Faxes")
	if err != nil {
		return nil, nil, err
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	if resp.StatusCode != http.StatusOK {
		exc := new(Exception)
		err = json.Unmarshal(respBody, exc)
		return nil, exc, err
	}

	lr := t.newListResources()
	if err := json.Unmarshal(respBody, lr); err != nil {
		return nil, nil, err
	}
	frs := lr.Faxes
	for lr.hasNext() {
		if exc, err := lr.next(); exc != nil || err != nil {
			return nil, exc, err
		}
		frs = append(frs, lr.Faxes...)
	}
	return frs, nil, nil
}

// SendFax uses Twilio to send a fax.
// See https://www.twilio.com/docs/fax/api/faxes#list-post for more information.
func (t *Twilio) SendFax(to, from, mediaUrl, quality, statusCallback string, storeMedia bool) (*FaxResource, *Exception, error) {
	values := url.Values{}
	values.Set("To", to)
	values.Set("From", from)
	values.Set("MediaUrl", mediaUrl)
	if quality != "" {
		values.Set("Quality", quality)
	}
	if statusCallback != "" {
		values.Set("StatusCallback", statusCallback)
	}
	if storeMedia {
		values.Set("StoreMedia", "true")
	}

	resp, err := t.post(values, "https://fax.twilio.com/v1/Faxes")
	if err != nil {
		return nil, nil, err
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	if resp.StatusCode != http.StatusCreated {
		exc := new(Exception)
		err = json.Unmarshal(respBody, exc)
		return nil, exc, err
	}

	fr := new(FaxResource)
	if err := json.Unmarshal(respBody, fr); err != nil {
		return nil, nil, err
	}
	return fr, nil, nil
}
