package gotwilio

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/google/go-querystring/query"
)

// PhoneNumberType defines whether a phone number is local, toll-free, or mobile.
type PhoneNumberType int

const (
	PhoneNumberLocal PhoneNumberType = iota
	PhoneNumberTollFree
	PhoneNumberMobile
)

var numberTypeMapping = map[PhoneNumberType]string{
	PhoneNumberLocal:    "Local",
	PhoneNumberTollFree: "TollFree",
	PhoneNumberMobile:   "Mobile",
}

func (t PhoneNumberType) String() string {
	return numberTypeMapping[t]
}

// AvailablePhoneNumbersOptions are all of the options that can be passed to an GetAvailablePhoneNumber query.
type AvailablePhoneNumbersOptions struct {
	AreaCode                      string  `url:",omitempty"`
	Contains                      string  `url:",omitempty"`
	SMSEnabled                    Boolean `url:"SmsEnabled,omitempty"`
	MMSEnabled                    Boolean `url:"MmsEnabled,omitempty"`
	VoiceEnabled                  Boolean `url:",omitempty"`
	FaxEnabled                    Boolean `url:",omitempty"`
	ExcludeAllAddressRequired     Boolean `url:",omitempty"`
	ExcludeLocalAddressRequired   Boolean `url:",omitempty"`
	ExcludeForeignAddressRequired Boolean `url:",omitempty"`
	Beta                          Boolean `url:",omitempty"`
	NearNumber                    string  `url:",omitempty"`
	NearLatLong                   string  `url:",omitempty"`
	Distance                      int     `url:",omitempty"`
	InPostalCode                  string  `url:",omitempty"`
	InRegion                      string  `url:",omitempty"`
	InRateCenter                  string  `url:",omitempty"`
	InLATA                        string  `url:"InLata,omitempty"`
	InLocality                    string  `url:",omitempty"`
}

// ToQueryString converts the provided options to a query string to be used in the outbound HTTP request.
func (o AvailablePhoneNumbersOptions) ToQueryString() (url.Values, error) {
	return query.Values(o)
}

// AvailablePhoneNumber represents a Twilio phone number available for purchase
// https://www.twilio.com/docs/phone-numbers/api/availablephonenumber-resource
type AvailablePhoneNumber struct {
	FriendlyName string  `json:"friendly_name"`
	PhoneNumber  string  `json:"phone_number"`
	LATA         string  `json:"lata"`
	RateCenter   string  `json:"rate_center"`
	Region       string  `json:"region"`
	Locality     string  `json:"locality"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	PostalCode   string  `json:"postal_code"`
	Beta         bool    `json:"beta"`

	Capabilities struct {
		MMS   bool `json:"mms"`
		SMS   bool `json:"sms"`
		Voice bool `json:"voice"`
	} `json:"capabilities"`
}

// GetAvailablePhoneNumbers retrieves list of available phone numbers
func (twilio *Twilio) GetAvailablePhoneNumbers(numberType PhoneNumberType, country string, options AvailablePhoneNumbersOptions) ([]*AvailablePhoneNumber, *Exception, error) {
	// build initial request
	resourceName := country + "/" + numberType.String() + ".json"
	req, err := http.NewRequest(http.MethodGet, twilio.buildUrl("AvailablePhoneNumbers/"+resourceName), nil)
	if err != nil {
		return nil, nil, err
	}

	// authenticate
	req.SetBasicAuth(twilio.getBasicAuthCredentials())

	// set query string
	queryValues, err := query.Values(options)
	if err != nil {
		return nil, nil, err
	}
	req.URL.RawQuery = queryValues.Encode()

	// perform request
	res, err := twilio.do(req)
	if err != nil {
		return nil, nil, err
	}

	decoder := json.NewDecoder(res.Body)
	if res.StatusCode != http.StatusOK {
		exception := new(Exception)
		err = decoder.Decode(exception)
		return nil, exception, err
	}

	// decode response
	availablePhoneNumberResponse := new(struct {
		AvailablePhoneNumbers []*AvailablePhoneNumber `json:"available_phone_numbers"`
	})
	decoder.Decode(availablePhoneNumberResponse)
	return availablePhoneNumberResponse.AvailablePhoneNumbers, nil, nil
}

// IncomingPhoneNumber represents a phone number resource owned by the calling account in Twilio
type IncomingPhoneNumber struct {
	SID          string `json:"sid"`
	PhoneNumber  string `url:",omitempty" json:"phone_number"`
	AreaCode     string `url:",omitempty"`
	FriendlyName string `url:",omitempty" json:"friendly_name"`

	SMSApplicationSID string `url:"SmsApplicationSid,omitempty" json:"sms_application_sid"`
	SMSMethod         string `url:"SmsMethod,omitempty" json:"sms_method"`
	SMSURL            string `url:"SmsUrl,omitempty" json:"sms_url"`
	SMSFallbackMethod string `url:"SmsFallbackMethod,omitempty" json:"sms_fallback_method"`
	SMSFallbackURL    string `url:"SmsFallbackUrl,omitempty" json:"sms_fallback_url"`

	StatusCallback       string `url:",omitempty" json:"status_callback"`
	StatusCallbackMethod string `url:",omitempty" json:"status_callback_method"`

	VoiceApplicationSID string  `url:"VoiceApplicationSid,omitempty"`
	VoiceMethod         string  `url:",omitempty"`
	VoiceURL            string  `url:"VoiceUrl,omitempty"`
	VoiceFallbackMethod string  `url:",omitempty"`
	VoiceFallbackURL    string  `url:"VoiceFallbackUrl,omitempty"`
	VoiceCallerIDLookup Boolean `url:",omitempty"`

	// Either "Active" or "Inactive"
	EmergencyStatus    string `url:",omitempty"`
	EmergencyStatusSID string `url:"EmergencyStatusSid,omitempty"`

	TrunkSID    string `url:"TrunkSid,omitempty"`
	IdentitySID string `url:"IdentitySid,omitempty"`
	AddressSID  string `url:"AddressSid,omitempty"`

	// Either "fax" or "voice". Defaults to "voice"
	VoiceReceiveMode string `url:",omitempty"`
}

// CreateIncomingPhoneNumber creates an IncomingPhoneNumber resource via the Twilio REST API.
// https://www.twilio.com/docs/phone-numbers/api/incomingphonenumber-resource#create-an-incomingphonenumber-resource
func (twilio *Twilio) CreateIncomingPhoneNumber(options IncomingPhoneNumber) (*IncomingPhoneNumber, *Exception, error) {
	// convert options to HTTP form
	form, err := query.Values(options)
	if err != nil {
		return nil, nil, err
	}

	res, err := twilio.post(form, twilio.buildUrl("IncomingPhoneNumbers.json"))
	if err != nil {
		return nil, nil, err
	}

	decoder := json.NewDecoder(res.Body)

	// handle NULL response
	if res.StatusCode != http.StatusCreated {
		exception := new(Exception)
		err = decoder.Decode(exception)
		return nil, exception, err
	}

	incomingPhoneNumber := new(IncomingPhoneNumber)
	err = decoder.Decode(incomingPhoneNumber)
	return incomingPhoneNumber, nil, err
}
