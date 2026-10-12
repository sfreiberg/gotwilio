> [!WARNING]
> **gotwilio is no longer maintained.** Please use the official Twilio Go library, [twilio-go](https://github.com/twilio/twilio-go), instead. It covers every Twilio API and is maintained by Twilio.
>
> This repository is archived. The final release, v1.1.0, remains available, but it will not receive further fixes or features.

## Overview
This is the start of a library for [Twilio](http://www.twilio.com/). Gotwilio supports making voice calls and sending text messages.

## License
Gotwilio is licensed under a BSD license.

## Installation
To install gotwilio, simply run `go get github.com/sfreiberg/gotwilio`.

## SMS Example

	package main

	import (
		"github.com/sfreiberg/gotwilio"
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

## MMS Example

	package main

	import (
		"github.com/sfreiberg/gotwilio"
	)

	func main() {
		accountSid := "ABC123..........ABC123"
		authToken := "ABC123..........ABC123"
		twilio := gotwilio.NewTwilioClient(accountSid, authToken)

		from := "+15555555555"
		to := "+15555555555"
		message := "Welcome to gotwilio!"
		mediaUrl := []string{"http://host/myimage.gif"}
		twilio.SendMMS(from, to, message, mediaUrl, "", "")
	}

## Voice Example

	package main

	import (
		"github.com/sfreiberg/gotwilio"
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

## Video example

	package main

	import (
		"github.com/sfreiberg/gotwilio"
	)

	func main() {
		accountSid := "ABC123..........ABC123"
		authToken := "ABC123..........ABC123"
		twilio := gotwilio.NewTwilioClient(accountSid, authToken)

		twilio.CreateVideoRoom(gotwilio.DefaultVideoRoomOptions)
	}
