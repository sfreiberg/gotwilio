package gotwilio

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

var LookupType = struct {
	Carrier    string
	CallerName string
	Fraud      string
}{
	Carrier:    "carrier",
	CallerName: "caller-name",
	Fraud:      "fraud",
}

type LookupResponse struct {
	CallerName struct {
		CallerName string      `json:"caller_name"`
		CallerType string      `json:"caller_type"`
		ErrorCode  interface{} `json:"error_code"`
	} `json:"caller_name"`
	Carrier struct {
		ErrorCode         interface{} `json:"error_code"`
		MobileCountryCode string      `json:"mobile_country_code"`
		MobileNetworkCode string      `json:"mobile_network_code"`
		Name              string      `json:"name"`
		Type              string      `json:"type"`
	} `json:"carrier"`
	Fraud struct {
		ErrorCode         interface{} `json:"error_code"`
		MobileCountryCode string      `json:"mobile_country_code"`
		MobileNetworkCode string      `json:"mobile_network_code"`
		AdvancedLineType  string      `json:"advanced_line_type"`
		CallerName        string      `json:"caller_name"`
		IsPorted          bool        `json:"is_ported"`
		LastPortedDate    string      `json:"last_ported_date"`
	} `json:"fraud"`
	CountryCode    string `json:"country_code"`
	NationalFormat string `json:"national_format"`
	PhoneNumber    string `json:"phone_number"`
	AddOns         struct {
		Status  string      `json:"status"`
		Message interface{} `json:"message"`
		Code    interface{} `json:"code"`
		Results struct {
		} `json:"results"`
	} `json:"add_ons"`
	URL string `json:"url"`
}

// LookupPhone uses the Twilio lookup API to request additional data about a phone number
// See https://www.twilio.com/docs/lookup/api
func (twilio *Twilio) LookupPhone(phone string, types ...string) (lookupResponse *LookupResponse, err error) {
	twilioUrl := twilio.LookupUrl + "/" + phone

	if len(types) != 0 {
		twilioUrl += "?Type=" + types[0]
		for _, t := range types[1:] {
			twilioUrl += "&Type=" + t
		}
	}

	res, err := twilio.get(twilioUrl)
	if err != nil {
		return lookupResponse, err
	}
	defer res.Body.Close()

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return lookupResponse, err
	}

	if res.StatusCode != http.StatusCreated {
		exception := new(Exception)
		err = json.Unmarshal(responseBody, exception)

		// We aren't checking the error because we don't actually care.
		// It's going to be passed to the client either way.
		return lookupResponse, err
	}

	lookupResponse = new(LookupResponse)
	err = json.Unmarshal(responseBody, lookupResponse)
	return lookupResponse, err
}
