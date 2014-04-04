package gotwilio

// https://twilio-python.readthedocs.org/en/latest/usage/usage.html
// and https://www.twilio.com/docs/api/rest/usage-records

import (
	"encoding/json"
	"net/http"
)

type UsageRecords struct {
	FirstPageUrl    string        `json:"first_page_uri"`
	End             int           `json:"end"`
	PreviousPageUrl string        `json:"previous_page_uri"`
	Url             string        `json:"uri"`
	PageSize        int           `json:"page_size"`
	Start           int           `json:"start"`
	UsageRecords    []UsageRecord `json:"usage_records"`
}

type UsageRecord struct {
	Category    string   `json:"category"`
	Description string   `json:"description"`
	AccountSid  string   `json:"account_sid"`
	StartDate   string   `json:"start_date"`
	EndDate     string   `json:"end_date"`
	Count       int64    `json:"count"`
	CountUnit   string   `json:"count_unit"`
	Usage       int64    `json:"usage"`
	UsageUnit   string   `json:"usage_unit"`
	Price       *float32 `json:"price,omitempty"`
	PriceUnit   string   `json:"price_unit"`
	ApiVersion  string   `json:"api_version"`
	Url         string   `json:"uri"`
}

func (twilio *Twilio) UsageRecords() (*UsageRecords, *Exception, error) {
	var usageRecords *UsageRecords
	var exception *Exception
	twilioUrl := twilio.BaseUrl + "/Accounts/" + twilio.AccountSid + "/Usage/Records.json"
	res, err := twilio.get(twilioUrl)
	if err != nil {
		return usageRecords, exception, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	if res.StatusCode != http.StatusOK {
		exception = new(Exception)
		err = decoder.Decode(exception)

		// We aren't checking the error because we don't actually care.
		// It's going to be passed to the client either way.
		return usageRecords, exception, err
	}

	usageRecords = new(UsageRecords)
	err = decoder.Decode(usageRecords)
	return usageRecords, exception, err
}
