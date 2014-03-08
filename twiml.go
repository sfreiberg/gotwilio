package gotwilio

import (
	"bytes"
	"encoding/xml"
	"io"
)

// A response has a single field, a slice of all of its Verbs
type Response struct {
	Verbs []interface{}
}

// these structs define the XML encoding of Twiml verbs
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
	Text   string `xml:",chardata"`
	Method string `xml:"method,attr,omitempty"`
}

type Say struct {
	Text     string `xml:",chardata"`
	Voice    string `xml:"voice,attr,omitempty"`
	Loop     string `xml:"loop,attr,omitempty"`
	Language string `xml:"language,attr,omitempty"`
}

type Play struct {
	Text   string `xml:",chardata"`
	Loop   string `xml:"loop,attr,omitempty"`
	Digits string `xml:"digits,attr,omitempty"`
}

type Pause struct {
	Length string `xml:"length,attr"`
}

type Gather struct {
	XMLName     xml.Name `xml:"Gather"`
	Action      string   `xml:"action,attr,omitempty"`
	Method      string   `xml:"method,attr,omitempty"`
	Timeout     string   `xml:"timeout,attr,omitempty"`
	FinishOnKey string   `xml:"finishOnKey,attr,omitempty"`
	NumDigits   string   `xml:"numDigits,attr,omitempty"`
	Say         Say      `xml:"Say"`
	Play        Play     `xml:"Play"`
	Pause       Pause    `xml:"Pause"`
}

type Record struct {
	Action             string `xml:"action,attr,omitempty"`
	Method             string `xml:"method,attr,omitempty"`
	Timeout            string `xml:"timeout,attr,omitempty"`
	FinishOnKey        string `xml:"finishOnKey,attr,omitempty"`
	MaxLength          string `xml:"maxLength,attr,omitempty"`
	Transcribe         string `xml:"transcribe,attr,omitempty"`
	TranscribeCallback string `xml:"transcribeCallback,attr,omitempty"`
	PlayBeep           string `xml:"playBeep,attr,omitempty"`
}

type Sms struct {
	Text           string `xml:",chardata"`
	To             string `xml:"to,attr,omitempty"`
	From           string `xml:"from,attr,omitempty"`
	Action         string `xml:"action,attr,omitempty"`
	Method         string `xml:"method,attr,omitempty"`
	StatusCallback string `xml:"statusCallback,attr,omitempty"`
}

type Number struct {
	Text       string `xml:",chardata"`
	SendDigits string `xml:"sendDigits,attr,omitempty"`
	Url        string `xml:"url,attr,omitempty"`
	Method     string `xml:"method,attr,omitempty"`
}

type Sip struct {
	Text     string `xml:",chardata"`
	Username string `xml:"username,attr,omitempty"`
	Password string `xml:"password,attr,omitempty"`
	Url      string `xml:"url,attr,omitempty"`
	Method   string `xml:"method,attr,omitempty"`
}

type Client struct {
	Text   string `xml:",chardata"`
	Url    string `xml:"url,attr,omitempty"`
	Method string `xml:"method,attr,omitempty"`
}

type Conference struct {
	Text                   string `xml:",chardata"`
	Muted                  string `xml:"muted,attr,omitempty"`
	Beep                   string `xml:"beep,attr,omitempty"`
	StartConferenceOnEnter string `xml:"startConferenceOnEnter,attr,omitempty"`
	EndConferenceOnExit    string `xml:"endConferenceOnExit,attr,omitempty"`
	WaitUrl                string `xml:"waitUrl,attr,omitempty"`
	WaitMethod             string `xml:"waitMethod,attr,omitempty"`
	MaxParticipants        string `xml:"maxParticipants,attr,omitempty"`
}

type Queue struct {
	Text   string `xml:",chardata"`
	Url    string `xml:"url,attr,omitempty"`
	Method string `xml:"method,attr,omitempty"`
}

type Dial struct {
	Text         string   `xml:",chardata"`
	Action       string   `xml:"action,attr,omitempty"`
	Method       string   `xml:"method,attr,omitempty"`
	Timeout      string   `xml:"timeout,attr,omitempty"`
	HangupOnStar string   `xml:"hangupOnStar,attr,omitempty"`
	TimeLimit    string   `xml:"timeLimit,attr,omitempty"`
	CallerId     string   `xml:"callerId,attr,omitempty"`
	Record       string   `xml:"record,attr,omitempty"`
	Numbers      []Number `xml:"number"`
	Clients      []Client `xml:"client"`
}

// constructor method to make a Response
func NewTwimlResponse() *Response {
	return &Response{}
}

// method to add verbs to a response
func (resp *Response) AddVerb(verb interface{}) {
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
