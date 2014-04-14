package gotwilio

import (
	"bytes"
	"os"
	"testing"
)

var params map[string]string

func init() {
	params = make(map[string]string)
	params["SID"] = os.Getenv("TWILIO_ACCOUNT_SID")
	params["TOKEN"] = os.Getenv("TWILIO_AUTH_TOKEN")
	params["FROM"] = os.Getenv("TWILIO_FROM")
	params["TO"] = os.Getenv("TWILIO_TO")
}

func TestSMS(t *testing.T) {
	msg := "Welcome to gotwilio"
	twilio := NewTwilioClient(params["SID"], params["TOKEN"])
	_, exc, err := twilio.SendSMS(params["FROM"], params["TO"], msg, "", "")
	if err != nil {
		t.Fatal(err)
	}

	if exc != nil {
		t.Fatal(exc)
	}
}

func TestVoice(t *testing.T) {
	callback := NewCallbackParameters("http://example.com")
	twilio := NewTwilioClient(params["SID"], params["TOKEN"])
	_, exc, err := twilio.CallWithUrlCallbacks(params["FROM"], params["TO"], callback)
	if err != nil {
		t.Fatal(err)
	}

	if exc != nil {
		t.Fatal(exc)
	}
}

func TestTwiml(t *testing.T) {
	var b bytes.Buffer
	const properResponse = `<?xml version="1.0" encoding="UTF-8"?>` + "\n" +
		`<Response><Say voice="alice">test</Say><Pause length="2"></Pause></Response>`
	newSay := Say{Text: "test", Voice: "alice"}
	newPause := Pause{Length: "2"}
	resp := NewTwimlResponse(newSay, newPause)
	err := resp.SendTwimlResponse(&b)
	if err != nil {
		t.Fatal(err)
	}

	if b.String() != properResponse {
		t.Fatalf("Expected: %s, Got: %s", properResponse, b.String())
	}
}
