package gotwilio

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetUsage(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, testUsageResponse)
	}))
	defer srv.Close()

	twilio := NewTwilioClient("", "")
	twilio.BaseUrl = srv.URL

	usage, exc, err := twilio.GetUsage("", "", "", false)
	if err != nil {
		t.Fatal(err)
	}
	if exc != nil {
		t.Fatal(exc)
	}

	if len(usage.UsageRecords) != 2 {
		t.Errorf("Expected 2 records, got %d", len(usage.UsageRecords))
	}
}

// Example from https://www.twilio.com/docs/usage/api/usage-record:
const testUsageResponse = `
{
   "first_page_uri": "/2010-04-01/Accounts/ACXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX/Usage/Records.json?Page=0&PageSize=50",
   "previous_page_uri": null, 
   "uri": "/2010-04-01/Accounts/ACXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX/Usage/Records.json", 
   "page_size": 50,
   "usage_records": [
      {
         "category": "calleridlookups", 
         "count": "0", 
         "price_unit": "usd", 
         "subresource_uris": {
            "yearly": "/2010-04-01/Accounts/ACXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX/Usage/Records/Yearly.json?Category=calleridlookups", 
            "last_month": "/2010-04-01/Accounts/ACXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX/Usage/Records/LastMonth.json?Category=calleridlookups", 
            "monthly": "/2010-04-01/Accounts/ACXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX/Usage/Records/Monthly.json?Category=calleridlookups", 
            "yesterday": "/2010-04-01/Accounts/ACXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX/Usage/Records/Yesterday.json?Category=calleridlookups", 
            "daily": "/2010-04-01/Accounts/ACXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX/Usage/Records/Daily.json?Category=calleridlookups", 
            "all_time": "/2010-04-01/Accounts/ACXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX/Usage/Records/AllTime.json?Category=calleridlookups", 
            "this_month": "/2010-04-01/Accounts/ACXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX/Usage/Records/ThisMonth.json?Category=calleridlookups", 
            "today": "/2010-04-01/Accounts/ACXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX/Usage/Records/Today.json?Category=calleridlookups"
         }, 
         "description": "Caller ID Lookups", 
         "end_date": "2012-10-13", 
         "usage_unit": "lookups", 
         "price": "0", 
         "uri": "/2010-04-01/Accounts/ACXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX/Usage/Records.json?Category=calleridlookups&StartDate=2012-08-15&EndDate=2012-10-13", 
         "account_sid": "ACXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX", 
         "usage": "0", 
         "start_date": "2012-08-15", 
         "count_unit": "lookups"
      }, 
      {
         "category": "calls", 
         "count": "21", 
         "price_unit": "usd", 
         "subresource_uris": {
            "yearly": "/2010-04-01/Accounts/ACXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX/Usage/Records/Yearly.json?Category=calls", 
            "last_month": "/2010-04-01/Accounts/ACXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX/Usage/Records/LastMonth.json?Category=calls", 
            "monthly": "/2010-04-01/Accounts/ACXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX/Usage/Records/Monthly.json?Category=calls", 
            "yesterday": "/2010-04-01/Accounts/ACXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX/Usage/Records/Yesterday.json?Category=calls", 
            "daily": "/2010-04-01/Accounts/ACXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX/Usage/Records/Daily.json?Category=calls", 
            "all_time": "/2010-04-01/Accounts/ACXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX/Usage/Records/AllTime.json?Category=calls", 
            "this_month": "/2010-04-01/Accounts/ACXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX/Usage/Records/ThisMonth.json?Category=calls", 
            "today": "/2010-04-01/Accounts/ACXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX/Usage/Records/Today.json?Category=calls"
         }, 
         "description": "Total Calls", 
         "end_date": "2012-10-13", 
         "usage_unit": "minutes", 
         "price": "0.38", 
         "uri": "/2010-04-01/Accounts/ACXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX/Usage/Records.json?Category=calls&StartDate=2012-08-15&EndDate=2012-10-13", 
         "account_sid": "ACXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX", 
         "usage": "21", 
         "start_date": "2012-08-15", 
         "count_unit": "calls"
      }
   ], 
   "next_page_uri": null,
   "page": 0
}
`
