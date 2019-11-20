package gotwilio

import (
	"fmt"
	"net/url"

	"github.com/gorilla/schema"
)

// LookupReq Go-representation of Twilio REST API's lookup request.
// https://www.twilio.com/docs/api/rest/lookups#lookups-query-parameters
type LookupReq struct {
	PhoneNumber string
	Type        string
	CountryCode string
}

// Lookup Go-representation of Twilio REST API's lookup.
// https://www.twilio.com/docs/api/rest/lookups
type Lookup struct {
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

	url := fmt.Sprintf("%s/PhoneNumbers/%s?%s", twilio.LookupURL, req.PhoneNumber, values.Encode())
	res := Lookup{}
	err := twilio.getJSON(url, &res)
	return res, err
}

// LookupNoCarrier looks up a phone number's details without the carrier
func (twilio *Twilio) LookupNoCarrier(phoneNumber string) (Lookup, error) {
	req := LookupReq{PhoneNumber: phoneNumber}
	return twilio.SubmitLookup(req)
}
