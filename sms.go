package gotwilio

import (
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
	formValues := initFormValues(whatsapp(to), body, mediaURL, statusCallback, applicationSid)
	formValues.Set("From", whatsapp(from))

	return twilio.sendMessage(formValues)
}

// SendSMS uses Twilio to send a text message.
// See http://www.twilio.com/docs/api/rest/sending-sms for more information.
func (twilio *Twilio) SendSMS(from, to, body, statusCallback, applicationSid string, opts ...*Option) (smsResponse *SmsResponse, exception *Exception, err error) {
	formValues := initFormValues(to, body, nil, statusCallback, applicationSid)
	formValues.Set("From", from)

	for _, opt := range opts {
		if opt != nil {
			formValues.Set(opt.Key, opt.Value)
		}
	}

	smsResponse, exception, err = twilio.sendMessage(formValues)
	return
}

// GetSMS uses Twilio to get information about a text message.
// See https://www.twilio.com/docs/api/rest/sms for more information.
func (twilio *Twilio) GetSMS(sid string) (smsResponse *SmsResponse, exception *Exception, err error) {
	twilioUrl := twilio.BaseUrl + "/Accounts/" + twilio.AccountSid + "/SMS/Messages/" + sid + ".json"

	res, err := twilio.get(twilioUrl)
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
func (twilio *Twilio) GetMessage(sid string) (messageResponse *MessageResponse, exception *Exception, err error) {
	twilioUrl := twilio.BaseUrl + "/Accounts/" + twilio.AccountSid + "/Messages/" + sid + ".json"

	res, err := twilio.get(twilioUrl)
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
	formValues := initFormValues(to, body, nil, statusCallback, applicationSid)
	formValues.Set("MessagingServiceSid", messagingServiceSid)

	smsResponse, exception, err = twilio.sendMessage(formValues)
	return
}

// SendMMS uses Twilio to send a multimedia message.
func (twilio *Twilio) SendMMS(from, to, body string, mediaUrl []string, statusCallback, applicationSid string) (smsResponse *SmsResponse, exception *Exception, err error) {
	formValues := initFormValues(to, body, mediaUrl, statusCallback, applicationSid)
	formValues.Set("From", from)

	smsResponse, exception, err = twilio.sendMessage(formValues)
	return
}

// SendMMSWithCopilot uses Twilio Copilot to send a multimedia message.
// See https://www.twilio.com/docs/api/rest/sending-messages-copilot
func (twilio *Twilio) SendMMSWithCopilot(messagingServiceSid, to, body string, mediaUrl []string, statusCallback, applicationSid string) (smsResponse *SmsResponse, exception *Exception, err error) {
	formValues := initFormValues(to, body, mediaUrl, statusCallback, applicationSid)
	formValues.Set("MessagingServiceSid", messagingServiceSid)

	smsResponse, exception, err = twilio.sendMessage(formValues)
	return
}

// Core method to send message
func (twilio *Twilio) sendMessage(formValues url.Values) (smsResponse *SmsResponse, exception *Exception, err error) {
	twilioUrl := twilio.BaseUrl + "/Accounts/" + twilio.AccountSid + "/Messages.json"

	res, err := twilio.post(formValues, twilioUrl)
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
