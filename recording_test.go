package gotwilio

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTwilio_CreateRecording(t *testing.T) {
	client := initTestTwilioClient()

	param := &CreateRecordingParameters{
		RecordingStatusCallbackEvent: []RecordingEventType{RecordingEventCompleted, RecordingEventInProgress},
		RecordingStatusCallback:      `http://macbook-huazhouliu-outreach.ngrok.io/call/status`,
		RecordingChannels:            ChannelDual,
	}

	// need to get acct sid and call sid of the call
	accSid := "get your own account sid"
	callSid := "get your own call sid"
	res, _, err := client.CreateRecording(accSid, callSid, param)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	updateParam := &UpdateRecordingParameters{
		Status:        RecordingStatusPaused,
		PauseBehavior: RecordingPauseSkip,
	}
	res1, _, err := client.UpdateRecording(res.AccountSid, res.CallSid, res.Sid, updateParam)
	assert.NoError(t, err)
	assert.NotNil(t, res1)
	updateParam = &UpdateRecordingParameters{
		Status: RecordingStatusInProgress,
	}
	res1, _, err = client.UpdateRecording(res.AccountSid, res.CallSid, res.Sid, updateParam)
	assert.NoError(t, err)
	assert.NotNil(t, res1)
	updateParam = &UpdateRecordingParameters{
		Status: RecordingStatusStopped,
	}
	res1, _, err = client.UpdateRecording(res.AccountSid, res.CallSid, res.Sid, updateParam)
	assert.NoError(t, err)
	assert.NotNil(t, res1)
}
