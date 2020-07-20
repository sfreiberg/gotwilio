package gotwilio

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLookup(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, testLookupResponse)
	}))
	defer srv.Close()

	twilio := NewTwilioClient("", "")
	twilio.LookupURL = srv.URL
	req := &LookupReq{PhoneNumber: "+11231231234"}
	lookup, err := twilio.SubmitLookup(*req)
	if err != nil {
		t.Fatalf("Failed: %s", err.Error())
	}
	bs, err := json.MarshalIndent(lookup, "", "  ")
	if err != nil {
		t.Fatalf("Failed: %s", err.Error())
	}
	t.Logf("Lookup Result:\n%s\n", string(bs))
}

// Example from https://www.twilio.com/docs/usage/api/usage-record:
const testLookupResponse = `
{
	"caller_name": {
		"error_code": null,
		"caller_name": "My Business, INC",
		"caller_type": "CONSUMER",
	},
	"carrier": {
	  "error_code": null,
	  "mobile_country_code": "310",
	  "mobile_network_code": "456",
	  "name": "verizon",
	  "type": "mobile"
	},
	"country_code": "US",
	"national_format": "(510) 867-5310",
	"phone_number": "+15108675310",
	"fraud": null,
	"add_ons": null,
	"url": "https://lookups.twilio.com/v1/PhoneNumbers/phone_number"
  }
`
