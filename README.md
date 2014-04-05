## Overview
This is the start of a library for [Twilio](http://www.twilio.com/). Gotwilio supports making voice calls, sending text messages, validating requests, and creating TWiML responses.

## License
Gotwilio is licensed under a BSD license.

## Installation
To install gotwilio, simply run `go get github.com/Januzellij/gotwilio`, until the original author has merged my pull request.

## Getting Started
Just create a Twilio client with either `NewTwilioClient(accountSid, authToken)` or `NewTwilioClientFromEnvironment()`, and store the accountSid and authToken in `TWILIO_ACCOUNT_SID` and `TWILIO_AUTH_TOKEN` environment variables, respectively.

## Docs
The documentation can be found at http://godoc.org/github.com/Januzellij/gotwilio

## SMS Example

```go
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
```
	
## Voice Example

```go
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
```

## Validate Example

```go
package main

import (
	"github.com/Januzellij/gotwilio"
	"net/http"
	"log"
)

func root(w http.ResponseWriter, r *http.Request) {
	twilio, err := NewTwilioClientFromEnvironment()
	if err != nil {
		panic(err)
	}
	url := "http://example.com/"
	err = gotwilio.Validate(r, url, twilio.authToken)
	if err == nil {
		// proceed as normal, the request is from Twilio
	}
}

func main() {
	http.HandleFunc("/", root)
	http.ListenAndServe(":8080", nil)
}
```

## Twiml Response Example

```go
package main

import (
	"github.com/Januzellij/gotwilio"
	"os"
)

func main() {
	newSay := gotwilio.Say{Text: "test", Voice: "alice"}
	newPause := gotwilio.Pause{Length: "2"}
	resp := gotwilio.NewTwimlResponse([]interface{}{newSay, newPause})
	err := resp.SendTwimlResponse(os.Stdout) // when using Twiml in a real web app, this would actually be written to a http.ResponseWriter.
	if err != nil {
		// your verbs were invalid XML
	}
}
```