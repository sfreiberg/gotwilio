package gotwilio

import (
	"context"
	"net/http"
	"net/url"

	json "github.com/bytedance/sonic"
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
	AreaCode                      string `url:",omitempty"`
	Contains                      string `url:",omitempty"`
	SMSEnabled                    *bool  `url:"SmsEnabled,omitempty"`
	MMSEnabled                    *bool  `url:"MmsEnabled,omitempty"`
	VoiceEnabled                  *bool  `url:",omitempty"`
	FaxEnabled                    *bool  `url:",omitempty"`
	ExcludeAllAddressRequired     *bool  `url:",omitempty"`
	ExcludeLocalAddressRequired   *bool  `url:",omitempty"`
	ExcludeForeignAddressRequired *bool  `url:",omitempty"`
	Beta                          *bool  `url:",omitempty"`
	NearNumber                    string `url:",omitempty"`
	NearLatLong                   string `url:",omitempty"`
	Distance                      int    `url:",omitempty"`
	InPostalCode                  string `url:",omitempty"`
	InRegion                      string `url:",omitempty"`
	InRateCenter                  string `url:",omitempty"`
	InLATA                        string `url:"InLata,omitempty"`
	InLocality                    string `url:",omitempty"`
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

	decoder := json.ConfigStd.NewDecoder(res.Body)
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

	VoiceApplicationSID string `url:"VoiceApplicationSid,omitempty"`
	VoiceMethod         string `url:",omitempty"`
	VoiceURL            string `url:"VoiceUrl,omitempty"`
	VoiceFallbackMethod string `url:",omitempty"`
	VoiceFallbackURL    string `url:"VoiceFallbackUrl,omitempty"`
	VoiceCallerIDLookup *bool  `url:",omitempty"`

	// Either "Active" or "Inactive"
	EmergencyStatus    string `url:",omitempty"`
	EmergencyStatusSID string `url:"EmergencyStatusSid,omitempty"`

	TrunkSID    string `url:"TrunkSid,omitempty"`
	IdentitySID string `url:"IdentitySid,omitempty"`
	AddressSID  string `url:"AddressSid,omitempty"`

	// Either "fax" or "voice". Defaults to "voice"
	VoiceReceiveMode string `url:",omitempty"`
}

type GetIncomingPhoneNumbersRequest struct {
	Beta         *bool  `url:"Beta,omitempty"`
	FriendlyName string `url:"FriendlyName,omitempty"`
	PhoneNumber  string `url:"PhoneNumber,omitempty"`
	Origin       string `url:"Origin,omitempty"`
}

type getIncomingPhoneNumbersResponse struct {
	IncomingPhoneNumbers []*IncomingPhoneNumber `json:"incoming_phone_numbers"`
}

// GetIncomingPhoneNumbers reads multiple IncomingPhoneNumbers from the Twilio REST API, with optional filtering
// https://www.twilio.com/docs/phone-numbers/api/incomingphonenumber-resource#read-multiple-incomingphonenumber-resources
func (twilio *Twilio) GetIncomingPhoneNumbers(request GetIncomingPhoneNumbersRequest) ([]*IncomingPhoneNumber, *Exception, error) {
	return twilio.GetIncomingPhoneNumbersWithContext(context.Background(), request)
}

func (twilio *Twilio) GetIncomingPhoneNumbersWithContext(ctx context.Context, request GetIncomingPhoneNumbersRequest) ([]*IncomingPhoneNumber, *Exception, error) {
	// convert request to url.Values for encoding into querystring
	form, err := query.Values(request)
	if err != nil {
		return nil, nil, err
	}

	// build URL with query string
	endpoint := twilio.buildUrl("IncomingPhoneNumbers.json")
	reqURL, err := url.Parse(endpoint)
	if err != nil {
		return nil, nil, err
	}
	reqURL.RawQuery = form.Encode()

	res, err := twilio.get(ctx, reqURL.String())
	if err != nil {
		return nil, nil, err
	}

	decoder := json.ConfigStd.NewDecoder(res.Body)

	// handle NULL response
	if res.StatusCode != http.StatusOK {
		exception := new(Exception)
		err = decoder.Decode(exception)
		return nil, exception, err
	}

	response := new(getIncomingPhoneNumbersResponse)
	err = decoder.Decode(&response)
	return response.IncomingPhoneNumbers, nil, err
}

// CreateIncomingPhoneNumber creates an IncomingPhoneNumber resource via the Twilio REST API.
// https://www.twilio.com/docs/phone-numbers/api/incomingphonenumber-resource#create-an-incomingphonenumber-resource
func (twilio *Twilio) CreateIncomingPhoneNumber(options IncomingPhoneNumber) (*IncomingPhoneNumber, *Exception, error) {
	return twilio.CreateIncomingPhoneNumberWithContext(context.Background(), options)
}

func (twilio *Twilio) CreateIncomingPhoneNumberWithContext(ctx context.Context, options IncomingPhoneNumber) (*IncomingPhoneNumber, *Exception, error) {
	// convert options to HTTP form
	form, err := query.Values(options)
	if err != nil {
		return nil, nil, err
	}

	res, err := twilio.post(ctx, form, twilio.buildUrl("IncomingPhoneNumbers.json"))
	if err != nil {
		return nil, nil, err
	}

	decoder := json.ConfigStd.NewDecoder(res.Body)

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

// UpdateIncomingPhoneNumber updates an IncomingPhoneNumber resource via the Twilio REST API.
// https://www.twilio.com/docs/phone-numbers/api/incomingphonenumber-resource#update-an-incomingphonenumber-resource
func (twilio *Twilio) UpdateIncomingPhoneNumber(sid string, options IncomingPhoneNumber) (*IncomingPhoneNumber, *Exception, error) {
	return twilio.UpdateIncomingPhoneNumberWithContext(context.Background(), sid, options)
}

func (twilio *Twilio) UpdateIncomingPhoneNumberWithContext(ctx context.Context, sid string, options IncomingPhoneNumber) (*IncomingPhoneNumber, *Exception, error) {
	// convert options to HTTP form
	form, err := query.Values(options)
	if err != nil {
		return nil, nil, err
	}

	res, err := twilio.post(ctx, form, twilio.buildUrl("IncomingPhoneNumbers/"+sid+".json"))
	if err != nil {
		return nil, nil, err
	}

	decoder := json.ConfigStd.NewDecoder(res.Body)

	// handle NULL response
	if res.StatusCode != http.StatusOK {
		exception := new(Exception)
		err = decoder.Decode(exception)
		return nil, exception, err
	}

	incomingPhoneNumber := new(IncomingPhoneNumber)
	err = decoder.Decode(incomingPhoneNumber)
	return incomingPhoneNumber, nil, err
}

// DeleteIncomingPhoneNumber deletes an IncomingPhoneNumber resource via the Twilio REST API.
// https://www.twilio.com/docs/phone-numbers/api/incomingphonenumber-resource#delete-an-incomingphonenumber-resource
func (twilio *Twilio) DeleteIncomingPhoneNumber(sid string) (*Exception, error) {
	return twilio.DeleteIncomingPhoneNumberWithContext(context.Background(), sid)
}

func (twilio *Twilio) DeleteIncomingPhoneNumberWithContext(ctx context.Context, sid string) (*Exception, error) {
	resourceName := sid + ".json"
	res, err := twilio.delete(ctx, twilio.buildUrl("IncomingPhoneNumbers/"+resourceName))
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusNoContent {
		exception := new(Exception)
		decoder := json.ConfigStd.NewDecoder(res.Body)
		err = decoder.Decode(exception)
		return exception, err
	}

	return nil, nil
}
