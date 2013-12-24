package gotwilio

import (
	"encoding/xml"
	"bytes"
	"io"
)

// A response has a single field, a slice of all of its Verbs
type Response struct {
	Verbs []interface{}
}

// each of these structs denotes a different kind of Twiml <Message>; text, text and media, and only media
type Message struct {
	XMLName xml.Name `xml:"Message"`
	Body string `xml:"Body"`
}

type Media struct {
	XMLName xml.Name `xml:"Message"`
	Media string `xml:"Media"`
}

type MediaMessage struct {
	XMLName xml.Name `xml:"Message"`
	Body string `xml:"Body"`
	Media string `xml:"Media"`
}

// constructor to make it simpler
func NewTwimlResponse() *Response {
	return &Response{}
}

// private to easily add verbs to a response
func (resp *Response) addVerbs(verb interface{}) {
	newVerbs := append(resp.Verbs, verb)
	resp.Verbs = newVerbs
}

// adds a message containing only text
func (resp *Response) Message(body string) {
	resp.addVerbs(Message{Body: body})
}

// adds a message containing only media
func (resp *Response) Media(media string) {
	resp.addVerbs(Media{Media: media})
}

// adds a message containing media and text
func (resp *Response) MessageWithMedia(body, media string) {
	resp.addVerbs(MediaMessage{Body: body, Media: media})
}

// makes a buffer, writes the standard xml header and beginning response tag
// encodes all of the responses verbs as xml, and writes them to the buffer
// closes the response, and writes the buffer's contents to the provided writer
func (resp *Response) SendTwimlResponse(w io.Writer) {
	var b bytes.Buffer
	b.WriteString(xml.Header)
	b.WriteString("<Response>")
	result, _ := xml.Marshal(resp.Verbs)
	b.Write(result)
	b.WriteString("</Response>")
	w.Write(b.Bytes())
}