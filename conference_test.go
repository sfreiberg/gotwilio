package gotwilio

import (
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestTwilio_GetConference(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	client := initTestTwilioClient()

	res, exception, err := client.GetConference("CFfafdfecc9552c7362ba273260694c160")
	validateTwilioException(t, exception)
	assert.NoError(t, err)
	assert.NotNil(t, res)

	assert.Equal(t, "CFfafdfecc9552c7362ba273260694c160", res.Sid)
	assert.NotEmpty(t, res.FriendlyName)
}
