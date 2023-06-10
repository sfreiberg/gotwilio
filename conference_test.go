package gotwilio

import (
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var testConferenceSid = os.Getenv("TEST_CONFERENCE_SID")
var testPhoneNumberTo = os.Getenv("TEST_PHONE_NUMBER_TO")
var testPhoneNumberFrom = os.Getenv("TEST_PHONE_NUMBER_FROM")

func TestTwilio_GetConference(t *testing.T) {
	if testConferenceSid == "" {
		t.Skip("TEST_CONFERENCE_SID not set")
	}

	log.SetLevel(log.DebugLevel)
	client := initTestTwilioClient()

	res, exception, err := client.GetConference(testConferenceSid)
	validateTwilioException(t, exception)
	assert.NoError(t, err)
	assert.NotNil(t, res)

	assert.Equal(t, testConferenceSid, res.Sid)
	assert.NotEmpty(t, res.FriendlyName)
}

// Test Conference functionality end to end. Real Twilio Account SID and Auth Token are required.
// A real conference call must also be active.
func TestTwilio_Conference(t *testing.T) {
	if testConferenceSid == "" {
		t.Skip("TEST_CONFERENCE_SID not set")
	}

	log.SetLevel(log.DebugLevel)
	client := initTestTwilioClient()

	// cannot create a new conference in code

	conf, exception, err := client.GetConference(testConferenceSid)
	validateTwilioException(t, exception)
	assert.NoError(t, err)
	assert.NotNil(t, conf)
	assert.Equal(t, testConferenceSid, conf.Sid)
	assert.NotEmpty(t, conf.FriendlyName)

	// add participant to call
	participant, exception, err := client.AddConferenceParticipant(conf.Sid, &ConferenceParticipantOptions{
		From:    testPhoneNumberFrom,
		To:      testPhoneNumberTo,
		Timeout: 15,
		Record:  NewBoolean(false),
		Muted:   NewBoolean(false),
	})
	validateTwilioException(t, exception)
	assert.NoError(t, err)
	assert.NotNil(t, participant)
	assert.NotEmpty(t, participant.CallSid)

	// get same participant's data
	participant2, exception, err := client.GetConferenceParticipant(conf.Sid, participant.CallSid)
	validateTwilioException(t, exception)
	assert.NoError(t, err)
	assert.NotNil(t, participant2)
	assert.Equal(t, participant.CallSid, participant2.CallSid)

	// update the conference
	_, exception, err = client.UpdateConference(conf.Sid, &ConferenceOptions{
		AnnounceURL:    "https://google.com",
		AnnounceMethod: "GET",
	})
	validateTwilioException(t, exception)
	assert.NoError(t, err)

	// update the participant
	_, exception, err = client.UpdateConferenceParticipant(conf.Sid, participant.CallSid, &ConferenceParticipantOptions{
		Muted: NewBoolean(true),
	})
	validateTwilioException(t, exception)
	assert.NoError(t, err)

	// delete the participant
	exception, err = client.DeleteConferenceParticipant(conf.Sid, participant.CallSid)
	validateTwilioException(t, exception)
	assert.NoError(t, err)
}
