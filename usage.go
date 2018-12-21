package gotwilio

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

// These are the parameters to use when you are requesting account usage.
// See https://www.twilio.com/docs/usage/api/usage-record#read-multiple-usagerecord-resources
// for more info.
type UsageParameters struct {
	Category           string // Optional
	StartDate          string // Optional, in YYYY-MM-DD or as offset
	EndDate            string // Optional, in YYYY-MM-DD or as offset
	IncludeSubaccounts bool   // Optional
}

// UsageRecord specifies the usage for a particular usage category.
// See https://www.twilio.com/docs/usage/api/usage-record#usagerecord-properties
// for more info.
type UsageRecord struct {
	AccountSid  string `json:"account_sid"`
	Category    string `json:"category"`
	Description string `json:"description"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	Price       string `json:"price"`
	PriceUnit   string `json:"price_unit"`
	Count       int    `json:"count,string"`
	CountUnit   string `json:"count_unit"`
	Usage       int    `json:"usage,string"`
	UsageUnit   string `json:"usage_unit"`
	// TODO: handle SubresourceUris
}

// UsageResponse contains information about account usage.
type UsageResponse struct {
	PageSize     int           `json:"page_size"`
	Page         int           `json:"page"`
	UsageRecords []UsageRecord `json:"usage_records"`
}

func (twilio *Twilio) GetUsage(category, startDate, endDate string, includeSubaccounts bool) (*UsageResponse, *Exception, error) {
	formValues := url.Values{}
	formValues.Set("category", category)
	formValues.Set("start_date", startDate)
	formValues.Set("end_date", endDate)
	formValues.Set("include_subaccounts", strconv.FormatBool(includeSubaccounts))

	var usageResponse *UsageResponse
	var exception *Exception
	twilioUrl := twilio.BaseUrl + "/Accounts/" + twilio.AccountSid + "/Usage/Records.json"

	res, err := twilio.get(twilioUrl + "?" + formValues.Encode())
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close()

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, nil, err
	}
	if res.StatusCode != http.StatusOK {
		exception = new(Exception)
		err = json.Unmarshal(responseBody, exception)
		return nil, exception, err
	}

	usageResponse = new(UsageResponse)
	err = json.Unmarshal(responseBody, usageResponse)
	return usageResponse, nil, err
}
