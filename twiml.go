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
	Url    string `xml:",chardata"`
	Method string `xml:"method,attr,omitempty"`
}

type Say struct {
	Text     string `xml:",chardata"`
	Voice    string `xml:"voice,attr,omitempty"`
	Loop     string `xml:"loop,attr,omitempty"`
	Language string `xml:"language,attr,omitempty"`
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
