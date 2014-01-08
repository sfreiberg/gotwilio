package gotwilio

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strconv"
)

// A response has a single field, a slice of all of its Verbs
type Response struct {
	Verbs []interface{}
}

// constructor method to make it simpler to create Responses
func NewTwimlResponse() *Response {
	return &Response{}
}

// private method to easily add verbs to a response
func (resp *Response) addVerb(verb interface{}) {
	newVerbs := append(resp.Verbs, verb)
	resp.Verbs = newVerbs
}

// makes a buffer, writes the standard xml header and beginning response tag
// encodes all of the responses verbs as xml, and writes them to the buffer
// closes the response, and writes the buffer's contents to the provided writer
func (resp *Response) SendTwimlResponse(w io.Writer) error {
	var b bytes.Buffer
	b.WriteString(xml.Header)
	b.WriteString("<Response>")
	result, err := xml.Marshal(resp.Verbs)
	if err != nil {
		return err
	} else {
		b.Write(result)
		b.WriteString("</Response>")
		w.Write(b.Bytes())
		return nil
	}
}

// these next structs define the XML encoding of Twiml verbs
// e.g the Message struct defines the fields and attributes of the <Message> verb,
// and the proper XML encoding

type Message struct {
	XMLName        xml.Name `xml:"Message"`
	To             string   `xml:"to,attr,omitempty"`
	From           string   `xml:"from,attr,omitempty"`
	Action         string   `xml:"action,attr,omitempty"`
	Method         string   `xml:"method,attr,omitempty"`
	StatusCallback string   `xml:"statusCallback,attr,omitempty"`
	Body           string   `xml:"Body,omitempty"`
	Media          string   `xml:"Media,omitempty"`
}

type Redirect struct {
	Url    string `xml:",chardata"`
	Method string `xml:"method,attr,omitempty"`
}

type Say struct {
	Text     string `xml:",chardata"`
	Voice    string `xml:"voice,attr,omitempty"`
	Loop     string `xml:"loop,attr,omitempty"`
	Language string `xml:"language,attr,omitempty"`
}

// these next methods handle adding the respective verbs to the given Response
// e.g the Message method handles adding a Message verb to the given Response

// TODO: implement error checking in all verb adding methods

func (resp *Response) Message(params Message) {
	resp.addVerb(params)
}

func (resp *Response) Redirect(params Redirect) {
	resp.addVerb(params)
}

func (resp *Response) Say(params Say) error {

	// list of all valid choices for an attribute
	voices := map[string]bool{
		"":      true,
		"man":   true,
		"women": true,
		"alice": true,
	}
	regLangs := map[string]bool{
		"":      true,
		"en":    true,
		"en-gb": true,
		"es":    true,
		"fr":    true,
		"de":    true,
		"it":    true,
	}
	aliceLangs := map[string]bool{
		"":      true,
		"da-DK": true,
		"de-DE": true,
		"en-AU": true,
		"en-CA": true,
		"en-GB": true,
		"en-IN": true,
		"en-US": true,
		"ca-ES": true,
		"es-ES": true,
		"es-MX": true,
		"fi-FI": true,
		"fr-CA": true,
		"fr-FR": true,
		"it-IT": true,
		"ja-JP": true,
		"ko-KR": true,
		"nb-NO": true,
		"nl-NL": true,
		"pl-PL": true,
		"pt-BR": true,
		"pt-PT": true,
		"ru-RU": true,
		"sv-SE": true,
		"zh-CN": true,
		"zh-HK": true,
		"zh-TW": true,
	}

	invalidLangError := fmt.Errorf("The language you specified (%s) is not valid for your voice (%s).", params.Language, params.Voice)

	// checking if the chosen voice is valid for the chosen language
	if !voices[params.Voice] {
		return errors.New("Please select a valid voice: man, women, alice, or none.")
	} else if params.Voice != "alice" && params.Voice != "" {
		if !regLangs[params.Language] {
			return invalidLangError
		}
	} else if params.Voice == "alice" {
		if !aliceLangs[params.Language] {
			return invalidLangError
		}
	}

	// checking validity of loop attribute. If it passes, add the given Say verb to the response
	if loopNum, err := strconv.Atoi(params.Loop); err != nil {
		return err
	} else if loopNum < 0 {
		return errors.New("Please give a nonnegative loop")
	} else {
		resp.addVerb(params)
		return nil
	}
}
