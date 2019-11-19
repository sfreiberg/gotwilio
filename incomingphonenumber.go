package gotwilio

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// https://www.twilio.com/docs/phone-numbers/api/incomingphonenumber-resource

// IncomingPhoneNumber contains phone number settings from twilio
type IncomingPhoneNumber struct {
	AccountSid           string       `json:"account_sid"`
	AddressSid           string       `json:"address_sid"`
	AddressRequirements  string       `json:"address_requirements"`
	APIVersion           string       `json:"api_version"`
	Beta                 bool         `json:"beta"`
	Capabilities         Capabilities `json:"capabilities"`
	DateCreated          string       `json:"date_created"`
	DateUpdated          string       `json:"date_updated"`
	IdentitySid          string       `json:"identity_sid"`
	PhoneNumber          string       `json:"phone_number"`
	Origin               string       `json:"origin"`
	Sid                  string       `json:"sid"`
	SmsApplicationSid    string       `json:"sms_application_sid"`
	SmsFallbackMethod    string       `json:"sms_fallback_method"`
	SmsMethod            string       `json:"sms_method"`
	SmsURL               string       `json:"sms_url"`
	StatusCallback       string       `json:"status_callback"`
	StatusCallbackMethod string       `json:"status_callback_method"`
	TrunkSid             string       `json:"trunk_sid"`
	URI                  string       `json:"uri"`
	VoiceApplicationSid  string       `json:"voice_application_sid"`
	VoiceCallerIDLookup  bool         `json:"voice_caller_id_lookup"`
	VoiceFallbackMethod  string       `json:"voice_fallback_method"`
	VoiceFallbackURL     string       `json:"voice_fallback_url"`
	VoiceMethod          string       `json:"voice_method"`
	VoiceURL             string       `json:"voice_url"`
	EmergencyStatus      string       `json:"emergency_status"`
	EmergencyAddressSid  string       `json:"emergency_address_sid"`
}

// Capabilities define what phone number features are enabled
type Capabilities struct {
	MMS   bool `json:"mms"`
	SMS   bool `json:"sms"`
	Voice bool `json:"voice"`
}

// GetIncomingPhoneNumber retrieves phone number settings from twilio
func (twilio *Twilio) GetIncomingPhoneNumber(phoneNumberSid string) (incomingPhoneNumberResponse *IncomingPhoneNumber, exception *Exception, err error) {
	twilioURL := fmt.Sprintf("%v/Accounts/%v/IncomingPhoneNumbers/%v.json", twilio.BaseUrl, twilio.AccountSid, phoneNumberSid)

	response, err := twilio.get(twilioURL)
	if err != nil {
		return incomingPhoneNumberResponse, exception, err
	}
	defer response.Body.Close()

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return incomingPhoneNumberResponse, exception, err
	}

	if response.StatusCode != http.StatusOK {
		exception := new(Exception)
		err = json.Unmarshal(responseBody, exception)

		return incomingPhoneNumberResponse, exception, err

	}

	incomingPhoneNumberResponse = new(IncomingPhoneNumber)
	err = json.Unmarshal(responseBody, incomingPhoneNumberResponse)
	return incomingPhoneNumberResponse, exception, err
}

// UpdateIncomingPhoneNumber updates phone number settings from twilio
func (twilio *Twilio) UpdateIncomingPhoneNumber(phoneNumberSid string, req *IncomingPhoneNumber) (incomingPhoneNumberResponse *IncomingPhoneNumber, exception *Exception, err error) {
	twilioURL := fmt.Sprintf("%v/Accounts/%v/IncomingPhoneNumbers/%v.json", twilio.BaseUrl, twilio.AccountSid, phoneNumberSid)

	response, err := twilio.post(incomingPhoneNumberFormValues(req), twilioURL)
	if err != nil {
		return incomingPhoneNumberResponse, exception, err
	}
	defer response.Body.Close()

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return incomingPhoneNumberResponse, exception, err
	}

	if response.StatusCode != http.StatusOK {
		exception := new(Exception)
		err = json.Unmarshal(responseBody, exception)

		return incomingPhoneNumberResponse, exception, err

	}

	incomingPhoneNumberResponse = new(IncomingPhoneNumber)
	err = json.Unmarshal(responseBody, incomingPhoneNumberResponse)
	return incomingPhoneNumberResponse, exception, err
}

func incomingPhoneNumberFormValues(req *IncomingPhoneNumber) url.Values {
	formValues := url.Values{}

	formValues.Set("sid", req.Sid)

	if req.SmsMethod != "" {
		formValues.Set("SmsMethod", req.SmsMethod)
	}

	if req.SmsURL != "" {
		formValues.Set("SmsUrl", req.SmsURL)
	}

	if req.SmsApplicationSid != "" {
		formValues.Set("SmsApplicationSid", req.SmsApplicationSid)
	}

	// TODO: Handle remaining optional form values

	return formValues
}
