package gotwilio

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"net/http"
	"net/url"
	"sort"
)

// takes in POST parameters and returns a string of the params concatenated like "keyvaluekeyvaluevalue"
func sortedFormString(f url.Values) string {
	keys := make([]string, len(f))
	values := make([][]string, len(f))

	i := 0
	for k, v := range f {
		keys[i] = k
		values[i] = v
		i++
	}

	// params must be sorted in alphabetical order
	sort.Strings(keys)

	// we use a buffer to concatenate the keys and values because it's a helluva lot faster and way easier than +ing everything
	var b bytes.Buffer
	for _, val := range keys {
		b.WriteString(val)
		for _, value := range f[val] {
			b.WriteString(value)
		}
	}

	return b.String()
}

// Validate checks if an http request actually came from Twilio, and is not faked.
// Validate uses directions from https://www.twilio.com/docs/security
func Validate(r *http.Request, url, authToken string) error {
	var urlString string

	// if the request is a POST request, get the string of the form
	if r.Method == "POST" {
		r.ParseForm()
		rawForm := r.PostForm
		formString := sortedFormString(rawForm)
		urlString = url + formString
	} else {
		urlString = url
	}

	mac := hmac.New(sha1.New, []byte(authToken))
	mac.Write([]byte(urlString))

	var b bytes.Buffer
	encoder := base64.NewEncoder(base64.StdEncoding, &b)
	macBytes := mac.Sum(nil)
	encoder.Write(macBytes)
	encoder.Close()

	// fetch the given twilio signature and compare against our signature
	twilioSig := r.Header.Get("X-Twilio-Signature")

	if bytes.Equal([]byte(twilioSig), b.Bytes()) == true {
		return nil
	} else {
		return errors.New("This request was spoofed")
	}
}
