package gotwilio

import (
	"os"
	"testing"
)

var (
	testTwilioAccountSID = os.Getenv("TWILIO_ACCOUNT_SID")
	testTwilioAuthToken  = os.Getenv("TWILIO_AUTH_TOKEN")
)

func initTestTwilioClient() *Twilio {
	return NewTwilioClient(testTwilioAccountSID, testTwilioAuthToken)
}

func validateTwilioException(t *testing.T, e *Exception) {
	if e != nil {
		t.Errorf("twilio exception. status: %d, message: %s, code: %d, more_info: %s", e.Status, e.Message, e.Code, e.MoreInfo)
	}
}
