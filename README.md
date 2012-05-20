## Overview
This is the start of a library for [Twilio](http://www.twilio.com/). It currently only supports sending text messages at the moment. As time permits I'll work on adding more features.

## License
Gotwilio is licensed under a BSD license.

## Installation
To install gotwilio, simply run `go get github.com/sfreiberg/gotwilio`.

## Example

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
		twilio.SendTextMessage(from, to, message)
	}