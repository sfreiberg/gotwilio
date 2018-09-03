package gotwilio

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

type PricingResponse struct {
	Country           string             `json:"country"`
	ISOCountry        string             `json:"iso_country"`
	PriceUnit         string             `json:"price_unit"`
	OutboundSMSPrices []OutboundSMSPrice `json:"outbound_sms_prices"`
	InboundSMSPrice   []Price            `json:"inbound_sms_prices"`
	Url               string             `json:"uri"`
}

type OutboundSMSPrice struct {
	MCC     string  `json:"mcc"`
	MNC     string  `json:"mnc"`
	Carrier string  `json:"carrier"`
	Prices  []Price `json:"prices"`
}

type Price struct {
	NumberType   string `json:"number_type"`
	BasePrice    string `json:"base_price"`
	CurrentPrice string `json:"current_price"`
}

func (twilio *Twilio) GetPricing(countryISO string) (pricingResponse *PricingResponse, exception *Exception, err error) {
	pricingUrl := twilio.PricingUrl + "/Messaging/Countries/" + strings.ToUpper(countryISO)

	res, err := twilio.get(pricingUrl)
	if err != nil {
		return pricingResponse, exception, err
	}
	defer res.Body.Close()

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return pricingResponse, exception, err
	}

	if res.StatusCode != http.StatusOK {
		exception = new(Exception)
		err = json.Unmarshal(responseBody, exception)

		// We aren't checking the error because we don't actually care.
		// It's going to be passed to the client either way.
		return pricingResponse, exception, err
	}

	pricingResponse = new(PricingResponse)
	err = json.Unmarshal(responseBody, pricingResponse)

	return pricingResponse, exception, err
}
