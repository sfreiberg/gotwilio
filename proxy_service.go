// Package gotwilio is a library for interacting with http://www.twilio.com/ API.
package gotwilio

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	json "github.com/bytedance/sonic"
)

// https://www.twilio.com/docs/proxy/api/pb-proxy-service
const ProxyBaseUrl = "https://proxy.twilio.com/v1"

type ProxyServiceRequest struct {
	UniqueName              string // optional
	CallbackURL             string // optional
	OutOfSessionCallbackURL string // optional
	InterceptCallbackURL    string // optional
	GeoMatchLevel           string // Default: country Options: country, area-code, extended-area-code	optional
	NumberSelectionBehavior string // Default: prefer-sticky Options: prefer-sticky, avoid-sticky	optional
	DefaultTtl              int    // (seconds) Default: 0	optional
}

type ProxyService struct {
	Sid                     string      `json:"sid"`
	UniqueName              string      `json:"unique_name"`
	AccountSid              string      `json:"account_sid"`
	CallbackURL             string      `json:"callback_url"`
	DefaultTTL              int         `json:"default_ttl"`
	NumberSelectionBehavior string      `json:"number_selection_behavior"`
	GeoMatchLevel           string      `json:"geo_match_level"`
	InterceptCallbackURL    interface{} `json:"intercept_callback_url"`
	OutOfSessionCallbackURL interface{} `json:"out_of_session_callback_url"`
	DateCreated             time.Time   `json:"date_created"`
	DateUpdated             time.Time   `json:"date_updated"`
	URL                     string      `json:"url"`
	Links                   struct {
		PhoneNumbers string `json:"phone_numbers"`
		ShortCodes   string `json:"short_codes"`
		Sessions     string `json:"sessions"`
	} `json:"links"`
}

// Create a new Twilio Service
func (twilio *Twilio) NewProxyService(service ProxyServiceRequest) (response *ProxyService, exception *Exception, err error) {
	return twilio.NewProxyServiceWithContext(context.Background(), service)
}

func (twilio *Twilio) NewProxyServiceWithContext(ctx context.Context, service ProxyServiceRequest) (response *ProxyService, exception *Exception, err error) {
	twilioUrl := ProxyBaseUrl + "/Services"

	res, err := twilio.post(ctx, proxyServiceFormValues(service), twilioUrl)
	if err != nil {
		return response, exception, err
	}
	defer res.Body.Close()

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return response, exception, err
	}

	if res.StatusCode != http.StatusCreated {
		exception = new(Exception)
		err = json.Unmarshal(responseBody, exception)

		// We aren't checking the error because we don't actually care.
		// It's going to be passed to the client either way.
		return response, exception, err
	}

	response = new(ProxyService)
	err = json.Unmarshal(responseBody, response)
	return response, exception, err

}

func (twilio *Twilio) GetProxyService(sid string) (response *ProxyService, exception *Exception, err error) {
	return twilio.GetProxyServiceWithContext(context.Background(), sid)
}

func (twilio *Twilio) GetProxyServiceWithContext(ctx context.Context, sid string) (response *ProxyService, exception *Exception, err error) {

	twilioUrl := fmt.Sprintf("%s/%s/%s", ProxyBaseUrl, "Services", sid)

	res, err := twilio.get(ctx, twilioUrl)
	if err != nil {
		return response, exception, err
	}
	defer res.Body.Close()

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return response, exception, err
	}

	if res.StatusCode != http.StatusOK {
		exception = new(Exception)
		err = json.Unmarshal(responseBody, exception)

		// We aren't checking the error because we don't actually care.
		// It's going to be passed to the client either way.
		return response, exception, err
	}

	response = new(ProxyService)
	err = json.Unmarshal(responseBody, response)
	return response, exception, err

}

func (twilio *Twilio) UpdateProxyService(sid string, service ProxyServiceRequest) (response *ProxyService, exception *Exception, err error) {
	return twilio.UpdateProxyServiceWithContext(context.Background(), sid, service)
}

func (twilio *Twilio) UpdateProxyServiceWithContext(ctx context.Context, sid string, service ProxyServiceRequest) (response *ProxyService, exception *Exception, err error) {

	twilioUrl := fmt.Sprintf("%s/%s/%s", ProxyBaseUrl, "Services", sid)

	res, err := twilio.post(ctx, proxyServiceFormValues(service), twilioUrl)
	if err != nil {
		return response, exception, err
	}
	defer res.Body.Close()

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return response, exception, err
	}

	if res.StatusCode != http.StatusOK {
		exception = new(Exception)
		err = json.Unmarshal(responseBody, exception)

		// We aren't checking the error because we don't actually care.
		// It's going to be passed to the client either way.
		return response, exception, err
	}

	response = new(ProxyService)
	err = json.Unmarshal(responseBody, response)
	return response, exception, err

}

func (twilio *Twilio) DeleteProxyService(sid string) (exception *Exception, err error) {
	return twilio.DeleteProxyServiceWithContext(context.Background(), sid)
}

func (twilio *Twilio) DeleteProxyServiceWithContext(ctx context.Context, sid string) (exception *Exception, err error) {

	twilioUrl := fmt.Sprintf("%s/%s/%s", ProxyBaseUrl, "Services", sid)

	res, err := twilio.delete(ctx, twilioUrl)
	if err != nil {
		return exception, err
	}

	respBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusNoContent {
		exc := new(Exception)
		err = json.Unmarshal(respBody, exc)
		return exc, err
	}
	return nil, nil
}

// Form values initialization
func proxyServiceFormValues(service ProxyServiceRequest) url.Values {
	formValues := url.Values{}

	formValues.Set("UniqueName", service.UniqueName)
	formValues.Set("CallbackUrl", service.CallbackURL)
	formValues.Set("OutOfSessionCallbackUrl", service.OutOfSessionCallbackURL)
	formValues.Set("InterceptCallbackUrl", service.InterceptCallbackURL)

	if service.GeoMatchLevel != "" {
		formValues.Set("GeoMatchLevel", service.GeoMatchLevel)
	}
	if service.NumberSelectionBehavior != "" {
		formValues.Set("NumberSelectionBehavior", service.NumberSelectionBehavior)
	}
	formValues.Set("DefaultTtl", fmt.Sprintf("%d", service.DefaultTtl))

	return formValues
}
