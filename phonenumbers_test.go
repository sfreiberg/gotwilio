package gotwilio

import (
	"testing"

	"github.com/stretchr/testify/assert"

	log "github.com/sirupsen/logrus"
)

const (
	testAccountSID      = ""
	testTwilioAuthToken = ""
)

func TestGetAvailablePhoneNumbers(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	client := NewTwilioClient(testAccountSID, testTwilioAuthToken)

	options := AvailablePhoneNumbersOptions{
		AreaCode:   "415",
		SMSEnabled: True,
	}

	// get available phone numbers
	res, err := client.GetAvailablePhoneNumbers(PhoneNumberTollFree, "US", options)
	assert.Nil(t, err)
	assert.NotNil(t, res)

	for _, availablePhoneNumber := range res {
		log.Debug(availablePhoneNumber)

		assert.NotNil(t, availablePhoneNumber)
		assert.NotEmpty(t, availablePhoneNumber.PhoneNumber)
	}
}
