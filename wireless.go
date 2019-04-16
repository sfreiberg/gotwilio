package gotwilio

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const (
	SimStatus_New         = "new"
	SimStatus_Ready       = "ready"
	SimStatus_Active      = "active"
	SimStatus_Suspended   = "suspended"
	SimStatus_Deactivated = "deactivated"
	SimStatus_Canceled    = "canceled"
	SimStatus_Scheduled   = "scheduled"
	SimStatus_Updating    = "updating"
)

type GetSimCardByIdResponse struct {
	Sims []SimCard
}

type GetRatePlansResponse struct {
	Plans []RatePlan `json:"rate_plans"`
}

type SimCard struct {
	AccountSId             string    `json:"account_sid"`
	UniqueName             string    `json:"unique_name"`
	CommandsCallbackMethod string    `json:"commands_callback_method"`
	CommandsCallbackURL    string    `json:"commands_callback_url"`
	Created                time.Time `json:"date_created"`
	Updated                time.Time `json:"date_updated"`
	FriendlyName           string    `json:"friendly_name"`
	SMSFallbackMethod      string    `json:"sms_fallback_method"`
	SMSFallbackURL         string    `json:"sms_fallback_url"`
	SMSMethod              string    `json:"sms_method"`
	SMSURL                 string    `json:"sms_url"`
	VoiceFallbackMethod    string    `json:"voice_fallback_method"`
	VoiceFallbackURL       string    `json:"voice_fallback_url"`
	VoiceMethod            string    `json:"voice_method"`
	VoiceURL               string    `json:"voice_url"`
	RatePlanSId            string    `json:"rate_plan_sid"`
	SId                    string    `json:"sid"`
	IccId                  string    `json:"iccid"`
	EId                    string    `json:"e_id"`
	Status                 string    `json:"status"`
	ResetStatus            *string   `json:"reset_status"`
	Url                    string    `json:"url"`
	IpAddress              string    `json:"ip_address"`
}

type RatePlan struct {
	AccountSId                    string    `json:"account_sid"`
	UniqueName                    string    `json:"unique_name"`
	DataEnabled                   bool      `json:"data_enabled"`
	DataLimit                     uint      `json:"data_limit"`
	DataMetering                  string    `json:"data_metering"`
	DateCreated                   time.Time `json:"date_created"`
	DateUpdated                   time.Time `json:"date_updated"`
	FriendlyName                  string    `json:"friendly_name"`
	MessagingEnabled              bool      `json:"messaging_enabled"`
	VoiceEnabled                  bool      `json:"voice_enabled"`
	NationalRoamingEnabled        bool      `json:"national_roaming_enabled"`
	NationalRoamingDataLimit      uint      `json:"national_roaming_data_limit"`
	InternationalRoamingDataLimit uint      `json:"international_roaming_data_limit"`
	SId                           string    `json:"sid"`
	Url                           string    `json:"url"`
}

func (t *Twilio) UpdateSimStatus(status, sid string) (*Exception, error) {
	formValues := url.Values{}
	formValues.Set("Status", status)

	resp, err := t.post(formValues, fmt.Sprintf("https://wireless.twilio.com/v1/Sims/%v", sid))
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(resp.Body)
	if resp.StatusCode != http.StatusAccepted {
		exc := new(Exception)
		err = decoder.Decode(exc)
		return exc, err
	}

	return nil, nil
}

func (t *Twilio) UpdateSimRatePlan(planSId, simSId string) (*Exception, error) {
	formValues := url.Values{}
	formValues.Set("RatePlan", planSId)

	resp, err := t.post(formValues, fmt.Sprintf("https://wireless.twilio.com/v1/Sims/%v", simSId))
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(resp.Body)
	if resp.StatusCode != http.StatusOK {
		exc := new(Exception)
		err = decoder.Decode(exc)
		return exc, err
	}

	return nil, nil
}

func (t *Twilio) GetSimCardByIccId(sid string) (*SimCard, *Exception, error) {

	resp, err := t.get(fmt.Sprintf("https://wireless.twilio.com/v1/Sims?Iccid=%v", sid))
	if err != nil {
		return nil, nil, err
	}

	decoder := json.NewDecoder(resp.Body)

	if resp.StatusCode != http.StatusOK {
		exc := new(Exception)
		err = decoder.Decode(exc)
		return nil, exc, err
	}

	var cards GetSimCardByIdResponse
	err = decoder.Decode(&cards)
	if err != nil {
		return nil, nil, err
	}

	if len(cards.Sims) != 1 {
		return nil, nil, errors.New("Invalid number of sim cards returned")
	}

	return &cards.Sims[0], nil, nil
}

func (t *Twilio) GetAllRatePlans() ([]RatePlan, *Exception, error) {
	resp, err := t.get("https://wireless.twilio.com/v1/RatePlans")
	if err != nil {
		return nil, nil, err
	}

	decoder := json.NewDecoder(resp.Body)

	if resp.StatusCode != http.StatusOK {
		exc := new(Exception)
		err = decoder.Decode(exc)
		return nil, exc, err
	}

	var ratePlans GetRatePlansResponse
	err = decoder.Decode(&ratePlans)
	if err != nil {
		return nil, nil, err
	}

	return ratePlans.Plans, nil, nil
}
