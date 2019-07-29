package gotwilio

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	log "github.com/sirupsen/logrus"
)

var (
	testTwilioAccountSID = os.Getenv("TWILIO_ACCOUNT_SID")
	testTwilioAuthToken  = os.Getenv("TWILIO_AUTH_TOKEN")
)

func TestGetAvailablePhoneNumbers(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	client := NewTwilioClient(testTwilioAccountSID, testTwilioAuthToken)

	options := AvailablePhoneNumbersOptions{
		AreaCode:   "415",
		SMSEnabled: True,
	}

	// get available phone numbers
	res, err := client.GetAvailablePhoneNumbers(PhoneNumberTollFree, "US", options)
	assert.NoError(t, err)
	assert.NotNil(t, res)

	for _, availablePhoneNumber := range res {
		assert.NotNil(t, availablePhoneNumber)
		assert.NotEmpty(t, availablePhoneNumber.PhoneNumber)
	}
}

func TestToQueryString(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	options := AvailablePhoneNumbersOptions{
		AreaCode:     "415",
		SMSEnabled:   True,
		VoiceEnabled: False,
	}

	queryString, err := options.ToQueryString()
	assert.NoError(t, err)

	assert.Empty(t, queryString.Get("in_region"))

	// test our ternary boolean
	assert.Equal(t, "true", queryString.Get("sms_enabled"))
	assert.Equal(t, "false", queryString.Get("voice_enabled"))
	assert.Empty(t, queryString.Get("mms_enabled"))
}
