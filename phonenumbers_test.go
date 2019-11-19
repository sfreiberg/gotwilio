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

func initTestTwilioClient() *Twilio {
	return NewTwilioClient(testTwilioAccountSID, testTwilioAuthToken)
}

func validateTwilioException(t *testing.T, e *Exception) {
	if e != nil {
		t.Errorf("twilio exception. status: %d, message: %s, code: %d, more_info: %s", e.Status, e.Message, e.Code, e.MoreInfo)
	}
}

func TestGetAvailablePhoneNumbers(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	options := AvailablePhoneNumbersOptions{
		AreaCode:   "925",
		SMSEnabled: NewBoolean(true),
	}

	// get available phone numbers
	client := initTestTwilioClient()
	res, exception, err := client.GetAvailablePhoneNumbers(PhoneNumberLocal, "US", options)
	validateTwilioException(t, exception)
	assert.NoError(t, err)
	assert.NotNil(t, res)

	for _, availablePhoneNumber := range res {
		log.Debug(availablePhoneNumber.PhoneNumber)
		assert.NotNil(t, availablePhoneNumber)
		assert.NotEmpty(t, availablePhoneNumber.PhoneNumber)
	}
}

func TestAvailablePhoneNumberOptionsToQueryString(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	options := AvailablePhoneNumbersOptions{
		AreaCode:     "415",
		SMSEnabled:   NewBoolean(true),
		VoiceEnabled: NewBoolean(false),
	}

	queryString, err := options.ToQueryString()
	assert.NoError(t, err)

	assert.Empty(t, queryString.Get("in_region"))

	// test our ternary boolean
	assert.Equal(t, "true", queryString.Get("SmsEnabled"))
	assert.Equal(t, "false", queryString.Get("VoiceEnabled"))
	assert.Empty(t, queryString.Get("MmsEnabled"))
}

func TestCreateUpdateDeleteIncomingPhoneNumber(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	client := initTestTwilioClient()
	phoneNumber := IncomingPhoneNumber{
		AreaCode: "925",
	}

	// create
	number, exception, err := client.CreateIncomingPhoneNumber(phoneNumber)
	validateTwilioException(t, exception)
	assert.NoError(t, err)
	assert.NotNil(t, number)

	// update
	number, exception, err = client.UpdateIncomingPhoneNumber(number.SID, IncomingPhoneNumber{
		FriendlyName: "test name",
	})
	validateTwilioException(t, exception)
	assert.NoError(t, err)
	assert.NotNil(t, number)
	assert.Equal(t, "test name", number.FriendlyName)

	// delete
	exception, err = client.DeleteIncomingPhoneNumber(number.SID)
	validateTwilioException(t, exception)
	assert.NoError(t, err)
}
