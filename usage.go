package gotwilio

// https://twilio-python.readthedocs.org/en/latest/usage/usage.html
// and https://www.twilio.com/docs/api/rest/usage-records

type UsageRecord struct {
	Category    string
	Description string
	AccountSid  string
	StartDate   string
	EndDate     string
	Usage       string
	UsageUnit   string
	Count       int64
	CountUnit   string
	Price       float32
	PriceUnit   string
	Url         string
}
