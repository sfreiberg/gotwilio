package gotwilio

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/gorilla/schema"
)

const (
	// LookupTypeString format for single request type
	LookupTypeString = "Type"

	// LookupTypesString format for multiple request types
	LookupTypesString = "Types"
)

// LookupReq Go-representation of Twilio REST API's lookup request.
// https://www.twilio.com/docs/api/rest/lookups#lookups-query-parameters
type LookupReq struct {
	PhoneNumber string
	Type        string
	Types       []string
	CountryCode string
}

// Lookup Go-representation of Twilio REST API's lookup.
// https://www.twilio.com/docs/api/rest/lookups
type Lookup struct {
	CallerName *struct {
		ErrorCode  *int   `json:"error_code"`
		CallerName string `json:"caller_name"`
		CallerType string `json:"caller_type"`
	} `json:"caller_name"`
	Carrier *struct {
		ErrorCode         *int   `json:"error_code"`
		MobileCountryCode string `json:"mobile_country_code"`
		MobileNetworkCode string `json:"mobile_network_code"`
		Name              string `json:"name"`
		Type              string `json:"type"`
	} `json:"carrier"`
	CountryCode    string `json:"country_code"`
	NationalFormat string `json:"national_format"`
	PhoneNumber    string `json:"phone_number"`
	URL            string `json:"url"`
}

// SubmitLookup sends a lookup request populating form fields only if they
// contain a non-zero value.
func (twilio *Twilio) SubmitLookup(req LookupReq) (Lookup, error) {
	encoder := schema.NewEncoder()
	values := url.Values{}

	if err := encoder.Encode(req, values); err != nil {
		return Lookup{}, err
	}

	// check for multiple types
	var types string
	if len(req.Types) > 0 {
		types = fmt.Sprintf("%s=%s", LookupTypeString, strings.Join(req.Types, "&Type="))
	} else {
		types = fmt.Sprintf("%s=%s", LookupTypeString, values.Get(LookupTypeString))
	}

	// remove req.Type value
	values.Del(LookupTypeString)
	values.Del(LookupTypesString)

	url := fmt.Sprintf("%s/PhoneNumbers/%s?%s&%s", twilio.LookupURL, req.PhoneNumber, values.Encode(), types)
	res := Lookup{}
	err := twilio.getJSON(url, &res)
	return res, err
}

// LookupNoCarrier looks up a phone number's details without the carrier
func (twilio *Twilio) LookupNoCarrier(phoneNumber string) (Lookup, error) {
	req := LookupReq{PhoneNumber: phoneNumber}
	return twilio.SubmitLookup(req)
}
