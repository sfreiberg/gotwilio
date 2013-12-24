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

// Message defines the <Message> Twiml verb
type Message struct {
	XMLName xml.Name `xml:"Message"`
	Body    string   `xml:"Body"`
	Media   string   `xml:"Media"`
}

// constructor to make it simpler
func NewTwimlResponse() *Response {
	return &Response{}
}

// private method to easily add verbs to a response
func (resp *Response) addVerb(verb interface{}) {
	newVerbs := append(resp.Verbs, verb)
	resp.Verbs = newVerbs
}

// adds a message to the given response
func (resp *Response) Message(body, media string) {
	resp.addVerb(Message{Body: body, Media: media})
}

// makes a buffer, writes the standard xml header and beginning response tag
// encodes all of the responses verbs as xml, and writes them to the buffer
// closes the response, and writes the buffer's contents to the provided writer
func (resp *Response) SendTwimlResponse(w io.Writer) {
	var b bytes.Buffer
	b.WriteString(xml.Header)
	b.WriteString("<Response>")
	result, _ := xml.Marshal(resp.Verbs) // add error handling
	b.Write(result)
	b.WriteString("</Response>")
	w.Write(b.Bytes())
}
