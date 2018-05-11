package gotwilio

import (
	"testing"
)

var params map[string]string

func init() {
	params = make(map[string]string)
	params["SID"] = "AC0f30491286ab4abb4a108abefbd05d8a"
	params["TOKEN"] = "1dcf52d7a1f3853ed78f0ee20d056dd0"
	params["FROM"] = "+15005550006"
	params["TO"] = "+19135551234"
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

func TestMMS(t *testing.T) {
	msg := "Welcome to gotwilio"
	twilio := NewTwilioClient(params["SID"], params["TOKEN"])
	_, exc, err := twilio.SendMMS(params["FROM"], params["TO"], msg, "http://www.google.com/images/logo.png", "", "")
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

func TestProxyServiceCRUD(t *testing.T) {
	t.Skip("Skipping as default test")

	twilio := NewTwilioClient(params["SID"], params["TOKEN"])

	req := ProxyServiceRequest{
		UniqueName:  "Test Service Name",
		CallbackURL: "https://www.example.com/",
	}

	resp, exc, err := twilio.NewProxyService(req)
	if err != nil {
		t.Fatal(err)
	}
	if exc != nil {
		t.Fatal(exc)
	}

	serviceID := resp.Sid

	resp, exc, err = twilio.GetProxyService(serviceID)
	if err != nil {
		t.Fatal(err)
	}
	if exc != nil {
		t.Fatal(exc)
	}

	req.OutOfSessionCallbackURL = "https://www.example.com/out"
	resp, exc, err = twilio.UpdateProxyService(serviceID, req)
	if err != nil {
		t.Fatal(err)
	}
	if exc != nil {
		t.Fatal(exc)
	}

	exc, err = twilio.DeleteProxyService(serviceID)
	if err != nil {
		t.Fatal(err)
	}

	if exc != nil {
		t.Fatal(exc)
	}

}

func TestProxySessionCRUD(t *testing.T) {
	t.Skip("Skipping as default test")

	twilio := NewTwilioClient(params["SID"], params["TOKEN"])
	// New Service to attach Session
	sreq := ProxyServiceRequest{
		UniqueName:  "Test Service Name",
		CallbackURL: "https://www.example.com/",
	}

	service, exc, err := twilio.NewProxyService(sreq)
	if err != nil {
		t.Fatal(err)
	}
	if exc != nil {
		t.Fatal(exc)
	}

	// Create Session
	req := ProxySessionRequest{
		UniqueName: "First Session Name",
		Mode:       "message-only",
	}

	session, exc, err := twilio.NewProxySession(service.Sid, req)
	if err != nil {
		t.Fatal(err)
	}
	if exc != nil {
		t.Fatal(exc)
	}

	// Get Session
	session, exc, err = twilio.GetProxySession(service.Sid, session.Sid)
	if err != nil {
		t.Fatal(err)
	}
	if exc != nil {
		t.Fatal(exc)
	}

	// Update Session
	req.Status = "closed"
	session, exc, err = twilio.UpdateProxySession(service.Sid, session.Sid, req)
	if err != nil {
		t.Fatal(err)
	}
	if exc != nil {
		t.Fatal(exc)
	}

	// Delete Session
	exc, err = twilio.DeleteProxySession(service.Sid, session.Sid)
	if err != nil {
		t.Fatal(err)
	}

	if exc != nil {
		t.Fatal(exc)
	}

	// Cleanup Service
	exc, err = twilio.DeleteProxyService(service.Sid)
	if err != nil {
		t.Fatal(err)
	}

	if exc != nil {
		t.Fatal(exc)
	}

}
