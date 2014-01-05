## Overview
This is the start of a library for [Twilio](http://www.twilio.com/). Gotwilio supports making voice calls and sending text messages.

## License
Gotwilio is licensed under a BSD license.

## Installation
To install gotwilio, simply run `go get github.com/Januzellij/gotwilio`, until the original author has merged my pull request.

## SMS Example

	package main

	import (
		"github.com/Januzellij/gotwilio"
	)

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

	import (
		"github.com/Januzellij/gotwilio"
	)

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
		resp.Message(gotwilio.Message{Body: "test", To: "+1234567890"})
		resp.SendTwimlResponse(os.Stdout) // when using Twiml in a real project, this would actually be written to a http.ResponseWriter.
	}