package gotwilio

import (
	"encoding/json"
	"net/url"

	"github.com/google/go-querystring/query"
)

// Boolean is a custom ternary nullable bool type
type Boolean *bool

var (
	// True represents "true" in our ternary optional bool
	True Boolean
	// False represents "false" in our ternary optional bool
	False Boolean
)

// initialize our Boolean type
func init() {
	a := true
	b := false
	True = &a
	False = &b
}

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
	AreaCode                      string  `url:"area_code,omitempty"`
	Contains                      string  `url:"contains,omitempty"`
	SMSEnabled                    Boolean `url:"sms_enabled,omitempty"`
	MMSEnabled                    Boolean `url:"mms_enabled,omitempty"`
	VoiceEnabled                  Boolean `url:"voice_enabled,omitempty"`
	FaxEnabled                    Boolean `url:"fax_enabled,omitempty"`
	ExcludeAllAddressRequired     Boolean `url:"exclude_all_address_required,omitempty"`
	ExcludeLocalAddressRequired   Boolean `url:"exclude_local_address_required,omitempty"`
	ExcludeForeignAddressRequired Boolean `url:"exclude_foreign_address_required,omitempty"`
	Beta                          Boolean `url:"beta,omitempty"`
	NearNumber                    string  `url:"near_number,omitempty"`
	NearLatLong                   string  `url:"near_lat_long,omitempty"`
	Distance                      int     `url:"distance,omitempty"`
	InPostalCode                  string  `url:"in_postal_code,omitempty"`
	InRegion                      string  `url:"in_region,omitempty"`
	InRateCenter                  string  `url:"in_rate_center,omitempty"`
	InLATA                        string  `url:"in_lata,omitempty"`
	InLocality                    string  `url:"in_locality,omitempty"`
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
func (twilio *Twilio) GetAvailablePhoneNumbers(numberType PhoneNumberType, country string, options AvailablePhoneNumbersOptions) ([]*AvailablePhoneNumber, error) {
	resourceName := country + "/" + numberType.String() + ".json"
	res, err := twilio.get(twilio.buildUrl("AvailablePhoneNumbers/" + resourceName))
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(res.Body)
	availablePhoneNumberResponse := new(struct {
		AvailablePhoneNumbers []*AvailablePhoneNumber `json:"available_phone_numbers"`
	})
	decoder.Decode(availablePhoneNumberResponse)
	return availablePhoneNumberResponse.AvailablePhoneNumbers, nil
}
