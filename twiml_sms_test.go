package gotwilio

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

// TestTWiMLSmsRenderBasic
// example from: https://www.twilio.com/docs/sms/twiml#a-basic-twiml-sms-response-example
func TestTWiMLSmsRenderBasic(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, testLookupResponse)
	}))
	defer srv.Close()

	twilio := NewTwilioClient("", "")
	twilio.LookupURL = srv.URL

	// init
	space := regexp.MustCompile(`\s+`)
	var mr MessagingResponse

	// add message
	body := "hello world!"
	redirect := "https://demo.twilio.com/welcome/sms/"
	mr.Message(&TWiMLSmsMessage{
		Body:     &body,
		Redirect: &redirect,
	})

	xml, err := mr.TWiMLSmsRender()
	if err != nil {
		t.Fatalf("failed to render xml: %+v", err)
	}

	expected := `<?xml version="1.0" encoding="UTF-8"?> <Response> <Message> <Body>hello world!</Body> <Redirect>https://demo.twilio.com/welcome/sms/</Redirect> </Message> </Response>`
	if expected != space.ReplaceAllString(xml, " ") {
		t.Fatalf("TestTWiMLSmsRenderBasic - unexpected xml")
	}
}

func TestTWiMLSmsRenderSend2Messages(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, testLookupResponse)
	}))
	defer srv.Close()

	twilio := NewTwilioClient("", "")
	twilio.LookupURL = srv.URL

	// init
	var mr MessagingResponse

	// message 1
	mr.Message(&TWiMLSmsMessage{
		Message: "This is message 1 of 2.",
	})

	// message 2
	mr.Message(&TWiMLSmsMessage{
		Message: "This is message 2 of 2.",
	})

	xml, err := mr.TWiMLSmsRender()
	if err != nil {
		t.Fatalf("failed to render xml: %+v", err)
	}

	space := regexp.MustCompile(`\s+`)

	expected := `<?xmlversion="1.0"encoding="UTF-8"?><Response><Message>Thisismessage1of2.</Message><Message>Thisismessage2of2.</Message></Response>`
	if expected != space.ReplaceAllString(xml, "") {
		t.Fatalf("TestTWiMLSmsRenderSend2Messages - unexpected xml")
	}
}

func TestTWiMLSmsRenderSendingMMS(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, testLookupResponse)
	}))
	defer srv.Close()

	twilio := NewTwilioClient("", "")
	twilio.LookupURL = srv.URL

	// init
	var mr MessagingResponse

	// add message
	body := "Store Location: 123 Easy St."
	redirect := "https://demo.twilio.com/owl.png"
	mr.Message(&TWiMLSmsMessage{
		Body:     &body,
		Redirect: &redirect,
	})

	xml, err := mr.TWiMLSmsRender()
	if err != nil {
		t.Fatalf("failed to render message with MMS: %+v", err)
	}

	space := regexp.MustCompile(`\s+`)
	expected := `<?xmlversion="1.0"encoding="UTF-8"?><Response><Message><Body>StoreLocation:123EasySt.</Body><Redirect>https://demo.twilio.com/owl.png</Redirect></Message></Response>`
	if expected != space.ReplaceAllString(xml, "") {
		t.Fatalf("TestTWiMLSmsRenderSendingMMS - unexpected xml")
	}
}

func TestTWiMLSmsRenderMessageStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, testLookupResponse)
	}))
	defer srv.Close()

	twilio := NewTwilioClient("", "")
	twilio.LookupURL = srv.URL

	// init
	var mr MessagingResponse

	// add message
	action := "/SmsHandler.php"
	method := "POST"
	mr.Message(&TWiMLSmsMessage{
		Message: "Store Location: 123 Easy St.",
		Action:  &action,
		Method:  &method,
	})

	xml, err := mr.TWiMLSmsRender()
	if err != nil {
		t.Fatalf("failed to render message status: %+v", err)
	}

	space := regexp.MustCompile(`\s+`)
	expected := `<?xmlversion="1.0"encoding="UTF-8"?><Response><MessageAction="/SmsHandler.php"Method="POST">StoreLocation:123EasySt.</Message></Response>`
	if expected != space.ReplaceAllString(xml, "") {
		t.Fatalf("TestTWiMLSmsRenderMessageStatus - unexpected xml")
	}
}
