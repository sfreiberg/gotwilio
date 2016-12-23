package gotwilio

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// LookupResponse is returned after a lookup
type LookupResponse struct {
	CountryCode    string `json:"country_code"`
	PhoneNumber    string `json:"phone_number"`
	NationalFormat string `json:"national_format"`
	Url            string `json:"url"`
}

// Lookup method
func (twilio *Twilio) Lookup(phone string) (lookupResponse *LookupResponse, exception *Exception, err error) {
	twilioUrl := twilio.LookupBaseUrl + phone

	res, err := twilio.get(twilioUrl)
	if err != nil {
		return lookupResponse, exception, err
	}
	defer res.Body.Close()

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return lookupResponse, exception, err
	}

	if res.StatusCode != http.StatusOK {
		exception = new(Exception)
		err = json.Unmarshal(responseBody, exception)

		// We aren't checking the error because we don't actually care.
		// It's going to be passed to the client either way.
		return lookupResponse, exception, err
	}

	lookupResponse = new(LookupResponse)
	err = json.Unmarshal(responseBody, lookupResponse)
	return lookupResponse, exception, err
}
