package gotwilio

import (
	"encoding/json"
	"log"
	"net/url"
)

// UsageRecords contains a a list of requested UsageRecord's and metadata
type UsageRecords struct {
	FirstPageUri    string        `json:"first_page_uri"`
	End             int           `json:"end"`
	PreviousPageUri string        `json:"previous_page_uri"`
	Uri             string        `json:"uri"`
	PageSize        int           `json:"page_size"`
	Start           int           `json:"start"`
	UsageRecords    []UsageRecord `json:"usage_records"`
}

// UsageRecord contains all data for a Twilio Usage Record
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
	Uri         string   `json:"uri"`
}

// UsageFilter contains all UsageRecord filter query parameters
type UsageFilter struct {
	Category  string
	StartDate string
	EndDate   string
}

// UsageRecords retreives all UsageRecord's at a subresource if provided, defaulting to the list resource,
// with the given filter parameters, if provided.
func (twilio *Twilio) UsageRecords(subresource string, filter UsageFilter) (*UsageRecords, *Exception, error) {
	var usageRecords *UsageRecords
	var exception *Exception
	var twilioUrl string

	if subresource != "" {
		twilioUrl = twilio.BaseUrl + "/Accounts/" + twilio.AccountSid + "/Usage/Records/" + subresource
	} else {
		twilioUrl = twilio.BaseUrl + "/Accounts/" + twilio.AccountSid + "/Usage/Records"
	}
	if filter != nil {
		u, urlError := url.Parse(twilioUrl)
		if urlError != nil {
			log.Fatalln(urlError)
		}
		q := url.Values{}
		q.Set("Category", filter.Category)
		q.Set("StartDate", filter.StartDate)
		q.Set("EndDate", filter.EndDate)
		u.RawQuery = q.Encode()
		twilioUrl = u.String() + ".json"
	} else {
		twilioUrl = twilioUrl + ".json"
	}

	res, err := twilio.get(twilioUrl)
	if err != nil {
		return usageRecords, exception, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	if res.StatusCode != 200 {
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
