package gotwilio

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// SmsResponse is returned after a text/sms message is posted to Twilio
type SmsResponse struct {
	Sid                 string  `json:"sid"`
	DateCreated         string  `json:"date_created"`
	DateUpdate          string  `json:"date_updated"`
	DateSent            string  `json:"date_sent"`
	AccountSid          string  `json:"account_sid"`
	To                  string  `json:"to"`
	From                string  `json:"from"`
	NumMedia            string  `json:"num_media"`
	Body                string  `json:"body"`
	NumSegments         string  `json:"num_segments"`
	Status              string  `json:"status"`
	MessagingServiceSid *string `json:"messaging_service_sid,omitempty"`
	Direction           string  `json:"direction"`
	ApiVersion          string  `json:"api_version"`
	Price               *string `json:"price,omitempty"`
	PriceUnit           *string `json:"price_unit,omitempty"`
	Url                 string  `json:"uri"`
}

// MessageResponse is returned when checking on the status of a message.
// See: https://www.twilio.com/docs/sms/api/message-resource#fetch-a-message-resource
type MessageResponse struct {
	Sid                 string            `json:"sid"`
	DateCreated         string            `json:"date_created"`
	DateUpdate          string            `json:"date_updated"`
	DateSent            string            `json:"date_sent"`
	AccountSid          string            `json:"account_sid"`
	To                  string            `json:"to"`
	From                string            `json:"from"`
	NumMedia            string            `json:"num_media"`
	Body                string            `json:"body"`
	NumSegments         string            `json:"num_segments"`
	Status              string            `json:"status"`
	MessagingServiceSid *string           `json:"messaging_service_sid,omitempty"`
	Direction           string            `json:"direction"`
	ApiVersion          string            `json:"api_version"`
	Price               *string           `json:"price,omitempty"`
	PriceUnit           *string           `json:"price_unit,omitempty"`
	Url                 string            `json:"uri"`
	ErrorCode           *int              `json:"error_code,omitempty"`
	ErrorMessage        *string           `json:"error_message,omitempty"`
	SubresourceURIs     map[string]string `json:"subresource_uris,omitempty"`
}

// SmsPriceResponse is returned price information base on country code
type SmsPriceResponse struct {
	Country           string             `json:"country"`
	ISOCountry        string             `json:"iso_country"`
	InboundSmsPrices  []InboundSmsPrice  `json:"inbound_sms_prices"`
	OutboundSmsPrices []OutboundSmsPrice `json:"outbound_sms_prices"`
	PriceUnit         string             `json:"price_unit"`
	Url               string             `json:"url"`
}

type InboundSmsPrice SMSPrice

type OutboundSmsPrice struct {
	Carrier string     `json:"carrier"`
	Mcc     string     `json:"mcc"`
	Mnc     string     `json:"mnc"`
	Prices  []SMSPrice `json:"prices"`
}

type SMSPrice struct {
	BasePrice    string `json:"base_price"`
	CurrentPrice string `json:"current_price"`
	NumberType   string `json:"number_type"`
}

// SmsCountryesponse is returned all countries about sms price.
type SmsCountryResponse struct {
	Meta      SmsCountryMeta `json:"meta"`
	Countries []SmsCountry   `json:"countries"`
}

type SmsCountryMeta struct {
	Page            int    `json:"page"`
	PageSize        int    `json:"page_size"`
	FirstPageUrl    string `json:"first_page_url"`
	PreviousPageUrl string `json:"previous_page_url"`
	Url             string `json:"url"`
	NextPageUrl     string `json:"next_page_url"`
	Key             string `json:"key"`
}

type SmsCountry struct {
	Country    string `json:"country"`
	ISOCountry string `json:"iso_country"`
	Url        string `json:"url"`
}

// Optional SMS parameters
var (
	// Settings for what Twilio should do with addresses in message logs
	SmsAddressRetentionObfuscate = &Option{"AddressRetention", "obfuscate"}
	SmsAddressRetentionRetain    = &Option{"AddressRetention", "retain"}
	// Settings for what Twilio should do with message content in message logs
	SmsContentRetentionDiscard = &Option{"ContentRetention", "discard"}
	SmsContentRetentionRetain  = &Option{"ContentRetention", "retain"}
)

// DateCreatedAsTime returns SmsResponse.DateCreated as a time.Time object
// instead of a string.
func (sms *SmsResponse) DateCreatedAsTime() (time.Time, error) {
	return time.Parse(time.RFC1123Z, sms.DateCreated)
}

// DateUpdateAsTime returns SmsResponse.DateUpdate as a time.Time object
// instead of a string.
func (sms *SmsResponse) DateUpdateAsTime() (time.Time, error) {
	return time.Parse(time.RFC1123Z, sms.DateUpdate)
}

// DateSentAsTime returns SmsResponse.DateSent as a time.Time object
// instead of a string.
func (sms *SmsResponse) DateSentAsTime() (time.Time, error) {
	return time.Parse(time.RFC1123Z, sms.DateSent)
}

func whatsapp(phone string) string {
	return "whatsapp:" + phone
}

// SendWhatsApp uses Twilio to send a WhatsApp message.
// See https://www.twilio.com/docs/sms/whatsapp/tutorial/send-and-receive-media-messages-whatsapp-python
func (twilio *Twilio) SendWhatsApp(from, to, body, statusCallback, applicationSid string) (smsResponse *SmsResponse, exception *Exception, err error) {
	return twilio.SendSMS(whatsapp(from), whatsapp(to), body, statusCallback, applicationSid)
}

// SendWhatsAppMedia uses Twilio to send a WhatsApp message with Media enabled.
// See https://www.twilio.com/docs/sms/whatsapp/tutorial/send-and-receive-media-messages-whatsapp-python
func (twilio *Twilio) SendWhatsAppMedia(from, to, body string, mediaURL []string, statusCallback, applicationSid string) (smsResponse *SmsResponse, exception *Exception, err error) {
	return twilio.SendWhatsAppMediaWithContext(context.Background(), from, to, body, mediaURL, statusCallback, applicationSid)
}

func (twilio *Twilio) SendWhatsAppMediaWithContext(ctx context.Context, from, to, body string, mediaURL []string, statusCallback, applicationSid string) (smsResponse *SmsResponse, exception *Exception, err error) {
	formValues := initFormValues(whatsapp(to), body, mediaURL, statusCallback, applicationSid)
	formValues.Set("From", whatsapp(from))

	return twilio.sendMessage(ctx, formValues)
}

// SendSMS uses Twilio to send a text message.
// See http://www.twilio.com/docs/api/rest/sending-sms for more information.
func (twilio *Twilio) SendSMS(from, to, body, statusCallback, applicationSid string, opts ...*Option) (smsResponse *SmsResponse, exception *Exception, err error) {
	return twilio.SendSMSWithContext(context.Background(), from, to, body, statusCallback, applicationSid, opts...)
}

func (twilio *Twilio) SendSMSWithContext(ctx context.Context, from, to, body, statusCallback, applicationSid string, opts ...*Option) (smsResponse *SmsResponse, exception *Exception, err error) {
	formValues := initFormValues(to, body, nil, statusCallback, applicationSid)
	formValues.Set("From", from)

	for _, opt := range opts {
		if opt != nil {
			formValues.Set(opt.Key, opt.Value)
		}
	}

	smsResponse, exception, err = twilio.sendMessage(ctx, formValues)
	return
}

// GetSMS uses Twilio to get information about a text message.
// See https://www.twilio.com/docs/api/rest/sms for more information.
func (twilio *Twilio) GetSMS(sid string) (smsResponse *SmsResponse, exception *Exception, err error) {
	return twilio.GetSMSWithContext(context.Background(), sid)
}

func (twilio *Twilio) GetSMSWithContext(ctx context.Context, sid string) (smsResponse *SmsResponse, exception *Exception, err error) {
	twilioUrl := twilio.BaseUrl + "/Accounts/" + twilio.AccountSid + "/SMS/Messages/" + sid + ".json"

	res, err := twilio.get(ctx, twilioUrl)
	if err != nil {
		return smsResponse, exception, err
	}
	defer res.Body.Close()

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return smsResponse, exception, err
	}

	if res.StatusCode != http.StatusOK {
		exception = new(Exception)
		err = json.Unmarshal(responseBody, exception)

		// We aren't checking the error because we don't actually care.
		// It's going to be passed to the client either way.
		return smsResponse, exception, err
	}

	smsResponse = new(SmsResponse)
	err = json.Unmarshal(responseBody, smsResponse)
	return smsResponse, exception, err
}

// GetMessage uses Twilio to get information about a text message.
//
// This can be used to check to see if a message has been successfully delivered, or if there was an error delivering the message.
//
// See https://www.twilio.com/docs/api/rest/sms for more information.
func (twilio *Twilio) GetMessage(ctx context.Context, sid string) (messageResponse *MessageResponse, exception *Exception, err error) {
	twilioUrl := twilio.BaseUrl + "/Accounts/" + twilio.AccountSid + "/Messages/" + sid + ".json"

	res, err := twilio.get(ctx, twilioUrl)
	if err != nil {
		return messageResponse, exception, err
	}
	defer res.Body.Close()

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return messageResponse, exception, err
	}

	if res.StatusCode != http.StatusOK {
		exception = new(Exception)
		err = json.Unmarshal(responseBody, exception)

		// We aren't checking the error because we don't actually care.
		// It's going to be passed to the client either way.
		return messageResponse, exception, err
	}

	messageResponse = new(MessageResponse)
	err = json.Unmarshal(responseBody, messageResponse)
	return messageResponse, exception, err
}

// SendSMSWithCopilot uses Twilio Copilot to send a text message.
// See https://www.twilio.com/docs/api/rest/sending-messages-copilot
func (twilio *Twilio) SendSMSWithCopilot(messagingServiceSid, to, body, statusCallback, applicationSid string) (smsResponse *SmsResponse, exception *Exception, err error) {
	return twilio.SendSMSWithCopilotWithContext(context.Background(), messagingServiceSid, to, body, statusCallback, applicationSid)
}

func (twilio *Twilio) SendSMSWithCopilotWithContext(ctx context.Context, messagingServiceSid, to, body, statusCallback, applicationSid string) (smsResponse *SmsResponse, exception *Exception, err error) {
	formValues := initFormValues(to, body, nil, statusCallback, applicationSid)
	formValues.Set("MessagingServiceSid", messagingServiceSid)

	smsResponse, exception, err = twilio.sendMessage(ctx, formValues)
	return
}

// GetSMSPrice uses Twilio to get price information base on country.
// See https://www.twilio.com/docs/sms/api/pricing for more information.
func (twilio *Twilio) GetSMSPrice(ctx context.Context, countryCode string) (smsPriceResponse *SmsPriceResponse, exception *Exception, err error) {
	twilioUrl := twilio.PriceUrl + "/Messaging/Countries/" + countryCode
	res, err := twilio.get(ctx, twilioUrl)
	if err != nil {
		return smsPriceResponse, exception, err
	}
	defer res.Body.Close()

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return smsPriceResponse, exception, err
	}

	if res.StatusCode != http.StatusOK {
		exception = new(Exception)
		err = json.Unmarshal(responseBody, exception)

		// We aren't checking the error because we don't actually care.
		// It's going to be passed to the client either way.
		return smsPriceResponse, exception, err
	}

	smsPriceResponse = new(SmsPriceResponse)
	err = json.Unmarshal(responseBody, smsPriceResponse)
	return smsPriceResponse, exception, err
}

// GetSMSCountries uses Twilio to get all countries about sms price.
// See https://www.twilio.com/docs/sms/api/pricing for more information.
func (twilio *Twilio) GetSMSCountries(ctx context.Context, nextPageUrl string, opts ...*Option) (smsCountryResponse *SmsCountryResponse, exception *Exception, err error) {
	var twilioUrl string
	if nextPageUrl == "" {
		queryValues := url.Values{}
		for _, opt := range opts {
			if opt != nil {
				queryValues.Set(opt.Key, opt.Value)
			}
		}

		twilioUrl = twilio.PriceUrl + "/Messaging/Countries?" + queryValues.Encode()
	} else {
		twilioUrl = nextPageUrl
	}

	res, err := twilio.get(ctx, twilioUrl)
	if err != nil {
		return smsCountryResponse, exception, err
	}
	defer res.Body.Close()

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return smsCountryResponse, exception, err
	}

	if res.StatusCode != http.StatusOK {
		exception = new(Exception)
		err = json.Unmarshal(responseBody, exception)

		// We aren't checking the error because we don't actually care.
		// It's going to be passed to the client either way.
		return smsCountryResponse, exception, err
	}

	smsCountryResponse = new(SmsCountryResponse)
	err = json.Unmarshal(responseBody, smsCountryResponse)
	return smsCountryResponse, exception, err
}

// SendMMS uses Twilio to send a multimedia message.
func (twilio *Twilio) SendMMS(from, to, body string, mediaUrl []string, statusCallback, applicationSid string) (smsResponse *SmsResponse, exception *Exception, err error) {
	return twilio.SendMMSWithContext(context.Background(), from, to, body, mediaUrl, statusCallback, applicationSid)
}

func (twilio *Twilio) SendMMSWithContext(ctx context.Context, from, to, body string, mediaUrl []string, statusCallback, applicationSid string) (smsResponse *SmsResponse, exception *Exception, err error) {
	formValues := initFormValues(to, body, mediaUrl, statusCallback, applicationSid)
	formValues.Set("From", from)

	smsResponse, exception, err = twilio.sendMessage(ctx, formValues)
	return
}

// SendMMSWithCopilot uses Twilio Copilot to send a multimedia message.
// See https://www.twilio.com/docs/api/rest/sending-messages-copilot
func (twilio *Twilio) SendMMSWithCopilot(messagingServiceSid, to, body string, mediaUrl []string, statusCallback, applicationSid string) (smsResponse *SmsResponse, exception *Exception, err error) {
	return twilio.SendMMSWithCopilotWithContext(context.Background(), messagingServiceSid, to, body, mediaUrl, statusCallback, applicationSid)
}

func (twilio *Twilio) SendMMSWithCopilotWithContext(ctx context.Context, messagingServiceSid, to, body string, mediaUrl []string, statusCallback, applicationSid string) (smsResponse *SmsResponse, exception *Exception, err error) {
	formValues := initFormValues(to, body, mediaUrl, statusCallback, applicationSid)
	formValues.Set("MessagingServiceSid", messagingServiceSid)

	smsResponse, exception, err = twilio.sendMessage(ctx, formValues)
	return
}

// Core method to send message
func (twilio *Twilio) sendMessage(ctx context.Context, formValues url.Values) (smsResponse *SmsResponse, exception *Exception, err error) {
	twilioUrl := twilio.BaseUrl + "/Accounts/" + twilio.AccountSid + "/Messages.json"

	res, err := twilio.post(ctx, formValues, twilioUrl)
	if err != nil {
		return smsResponse, exception, err
	}
	defer res.Body.Close()

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return smsResponse, exception, err
	}

	if res.StatusCode != http.StatusCreated {
		exception = new(Exception)
		err = json.Unmarshal(responseBody, exception)

		// We aren't checking the error because we don't actually care.
		// It's going to be passed to the client either way.
		return smsResponse, exception, err
	}

	smsResponse = new(SmsResponse)
	err = json.Unmarshal(responseBody, smsResponse)
	return smsResponse, exception, err
}

// Form values initialization
func initFormValues(to, body string, mediaUrl []string, statusCallback, applicationSid string) url.Values {
	formValues := url.Values{}

	formValues.Set("To", to)
	formValues.Set("Body", body)

	if len(mediaUrl) > 0 {
		for _, value := range mediaUrl {
			formValues.Add("MediaUrl", value)
		}
	}

	if statusCallback != "" {
		formValues.Set("StatusCallback", statusCallback)
	}

	if applicationSid != "" {
		formValues.Set("ApplicationSid", applicationSid)
	}

	return formValues
}
