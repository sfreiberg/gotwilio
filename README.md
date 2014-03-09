## Overview
This is the start of a library for [Twilio](http://www.twilio.com/). Gotwilio supports making voice calls, sending text messages, validating requests, and creating TWiML responses.

## License
Gotwilio is licensed under a BSD license.

## Installation
To install gotwilio, simply run `go get github.com/Januzellij/gotwilio`, until the original author has merged my pull request.

## Todo
* Maybe make a TWiML verb interface type?
* Start implementing the account usage API
* Generate docs with http://godoc.org

## Getting Started
Just create a Twilio client with either `NewTwilioClient(accountSid, authToken)` or `NewTwilioClientFromEnvironment()`, and store the accountSid and authToken in `TWILIO_ACCOUNT_SID` and `TWILIO_AUTH_TOKEN` environment variables, respectively.

## SMS Example

	package main

	import "github.com/Januzellij/gotwilio"

	func main() {
		accountSid := "ABC123..........ABC123"
		authToken := "ABC123..........ABC123"
		twilio := gotwilio.NewTwilioClient(accountSid, authToken)

		from := "+15555555555"
		to := "+15555555555"
		message := "Welcome to gotwilio!"
		twilio.SendSMS(from, to, message, "", "")
	}
	
## Voice Example

	package main

	import "github.com/Januzellij/gotwilio"

	func main() {
		accountSid := "ABC123..........ABC123"
		authToken := "ABC123..........ABC123"
		twilio := gotwilio.NewTwilioClient(accountSid, authToken)

		from := "+15555555555"
		to := "+15555555555"
		callbackParams := gotwilio.NewCallbackParameters("http://example.com")
		twilio.CallWithUrlCallbacks(from, to, callbackParams)
	}

## Validate Example

	package main

	import (
		"github.com/Januzellij/gotwilio"
		"net/http"
	)

	func root(w http.ResponseWriter, r *http.Request) {
		url := "http://example.com/"
		authToken := "12345"
		err := gotwilio.Validate(r, url, authToken)
		if err != nil {
			// do something
		}
	}

	func main() {
		http.HandleFunc("/", root)
		http.ListenAndServe(":8080", nil)
	}

## Twiml Response Example

	package main

	import (
		"github.com/Januzellij/gotwilio"
		"os"
	)
	
	func main() {
		resp := gotwilio.NewTwimlResponse()
		newGather := gotwilio.Gather{Method: "POST"}
		newGather.Say = gotwilio.Say{Text: "test", Voice: "alice"}

		newDial := gotwilio.Dial{}
		firstNumber := gotwilio.Number{Text: "I am the first number"}
		secondNumber := gotwilio.Number{Text: "I am the second number"}
		newDial.Numbers = []gotwilio.Number{firstNumber, secondNumber}

		newPause := gotwilio.Pause{Length: "2"}

		// be careful of the order in which you add verbs, that's the order they'll appear in the XML
		resp.AddVerb(newPause)
		resp.AddVerbs([]interface{}{newDial, newGather})
		err := resp.SendTwimlResponse(os.Stdout) // when using Twiml in a real web app, this would actually be written to a http.ResponseWriter.
		if err != nil {
			// your verbs were invalid XML
		}
	}