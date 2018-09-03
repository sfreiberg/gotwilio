package gotwilio

import (
	"testing"
)

var paramsLive map[string]string

func init() {
	paramsLive = make(map[string]string)

	// Only LIVE credentials possible, because of 20008 error
	paramsLive["SID"] = ""
	paramsLive["TOKEN"] = ""
}

func TestGetPricing(t *testing.T) {
	countryISO := "EE"
	twilio := NewTwilioClient(paramsLive["SID"], paramsLive["TOKEN"])
	_, exc, err := twilio.GetPricing(countryISO)
	if err != nil {
		t.Fatal(err)
	}

	if exc != nil {
		t.Fatal(exc)
	}
}