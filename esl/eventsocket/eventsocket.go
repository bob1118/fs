// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// FreeSWITCH Event Socket library for the Go programming language.
//
// eventsocket supports both inbound and outbound event socket connections,
// acting either as a client connecting to FreeSWITCH or as a server accepting
// connections from FreeSWITCH to control calls.
//
// Reference:
// http://wiki.freeswitch.org/wiki/Event_Socket
// http://wiki.freeswitch.org/wiki/Event_Socket_Outbound
//
// WORK IN PROGRESS, USE AT YOUR OWN RISK.
package eventsocket

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/textproto"
	"net/url"
	"strconv"
	"time"
)

const bufferSize = 1024 << 6 // For the socket reader
const eventsBuffer = 16      // For the events channel (memory eater!)
const timeoutPeriod = 16 * time.Second

var errMissingAuthRequest = errors.New("Missing auth request")
var errInvalidPassword = errors.New("Invalid password")
var errInvalidCommand = errors.New("Invalid command contains \\r or \\n")
var errTimeout = errors.New("Timeout")

// Connection is the event socket connection handler.
type Connection struct {
	conn                net.Conn
	reader              *bufio.Reader
	textreader          *textproto.Reader
	cmderr, apierr, err chan error
	cmd, api, evt       chan *Event
}

// newConnection allocates a new Connection and initialize its buffers.
func newConnection(c net.Conn) *Connection {
	h := Connection{
		conn:   c,
		reader: bufio.NewReaderSize(c, bufferSize),
		cmd:    make(chan *Event),
		cmderr: make(chan error, 1),
		api:    make(chan *Event),
		apierr: make(chan error, 1),
		evt:    make(chan *Event, eventsBuffer),
		err:    make(chan error, 1),
	}
	h.textreader = textproto.NewReader(h.reader)
	return &h
}

// HandleFunc is the function called on new incoming connections.
type HandleFunc func(*Connection)

// ListenAndServe listens for incoming connections from FreeSWITCH and calls
// HandleFunc in a new goroutine for each client.
//
// Example:
//
//	func main() {
//		eventsocket.ListenAndServe(":9090", handler)
//	}
//
//	func handler(c *eventsocket.Connection) {
//		ev, err := c.Send("connect") // must always start with this
//		ev.PrettyPrint()             // print event to the console
//		...
//		c.Send("myevents")
//		for {
//			ev, err = c.ReadEvent()
//			...
//		}
//	}
func ListenAndServe(addr string, fn HandleFunc) error {
	srv, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	for {
		c, err := srv.Accept()
		if err != nil {
			return err
		}
		h := newConnection(c)
		go h.readLoop()
		go fn(h)
	}
}

// Dial attemps to connect to FreeSWITCH and authenticate.
//
// Example:
//
//	c, _ := eventsocket.Dial("localhost:8021", "ClueCon")
//	ev, _ := c.Send("events plain ALL") // or events json ALL
//	for {
//		ev, _ = c.ReadEvent()
//		ev.PrettyPrint()
//		...
//	}
func Dial(addr, passwd string) (*Connection, error) {
	c, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	h := newConnection(c)
	m, err := h.textreader.ReadMIMEHeader()
	if err != nil {
		c.Close()
		return nil, err
	}
	if m.Get("Content-Type") != "auth/request" {
		c.Close()
		return nil, errMissingAuthRequest
	}
	fmt.Fprintf(c, "auth %s\r\n\r\n", passwd)
	m, err = h.textreader.ReadMIMEHeader()
	if err != nil {
		c.Close()
		return nil, err
	}
	if m.Get("Reply-Text") != "+OK accepted" {
		c.Close()
		return nil, errInvalidPassword
	}
	go h.readLoop()
	return h, err
}

// readLoop calls readOne until a fatal error occurs, then close the socket.
func (h *Connection) readLoop() {
	for h.readOne() {
	}
	h.Close()
}

// readOne reads a single event and send over the appropriate channel.
// It separates incoming events from api and command responses.
func (h *Connection) readOne() bool {
	hdr, err := h.textreader.ReadMIMEHeader()
	if err != nil {
		h.err <- err
		return false
	}
	resp := new(Event)
	resp.Header = make(EventHeader)
	if v := hdr.Get("Content-Length"); v != "" {
		length, err := strconv.Atoi(v)
		if err != nil {
			h.err <- err
			return false
		}
		b := make([]byte, length)
		if _, err := io.ReadFull(h.reader, b); err != nil {
			h.err <- err
			return false
		}
		resp.Body = string(b)
	}
	switch hdr.Get("Content-Type") {
	case "command/reply":
		reply := hdr.Get("Reply-Text")
		if reply[:2] == "-E" {
			h.cmderr <- errors.New(reply[5:])
			return true
		}
		if reply[0] == '%' {
			copyHeaders(&hdr, resp, true)
		} else {
			copyHeaders(&hdr, resp, false)
		}
		h.cmd <- resp
	case "api/response":
		if string(resp.Body[:2]) == "-E" {
			h.apierr <- errors.New(string(resp.Body)[5:])
			return true
		}
		copyHeaders(&hdr, resp, false)
		h.api <- resp
	case "text/event-plain":
		reader := bufio.NewReader(bytes.NewReader([]byte(resp.Body)))
		resp.Body = ""
		textreader := textproto.NewReader(reader)
		hdr, err = textreader.ReadMIMEHeader()
		if err != nil {
			h.err <- err
			return false
		}
		if v := hdr.Get("Content-Length"); v != "" {
			length, err := strconv.Atoi(v)
			if err != nil {
				h.err <- err
				return false
			}
			b := make([]byte, length)
			if _, err = io.ReadFull(reader, b); err != nil {
				h.err <- err
				return false
			}
			resp.Body = string(b)
		}
		copyHeaders(&hdr, resp, true)
		h.evt <- resp
	case "text/event-json":
		tmp := make(EventHeader)
		err := json.Unmarshal([]byte(resp.Body), &tmp)
		if err != nil {
			h.err <- err
			return false
		}
		// capitalize header keys for consistency.
		for k, v := range tmp {
			resp.Header[capitalize(k)] = v
		}
		if v, _ := resp.Header["_body"]; v != nil {
			resp.Body = v.(string)
			delete(resp.Header, "_body")
		} else {
			resp.Body = ""
		}
		h.evt <- resp
	case "text/disconnect-notice":
		copyHeaders(&hdr, resp, false)
		h.evt <- resp
	default:
		log.Fatal("Unsupported event:", hdr)
	}
	return true
}

// RemoteAddr returns the remote addr of the connection.
func (h *Connection) RemoteAddr() net.Addr {
	return h.conn.RemoteAddr()
}

// Close terminates the connection.
func (h *Connection) Close() {
	h.conn.Close()
}

// ReadEvent reads and returns events from the server. It supports both plain
// or json, but *not* XML.
//
// When subscribing to events (e.g. `Send("events json ALL")`) it makes no
// difference to use plain or json. ReadEvent will parse them and return
// all headers and the body (if any) in an Event struct.
func (h *Connection) ReadEvent() (*Event, error) {
	var (
		ev  *Event
		err error
	)
	select {
	case err = <-h.err:
		return nil, err
	case ev = <-h.evt:
		return ev, nil
	}
}

// copyHeaders copies all keys and values from the MIMEHeader to Event.Header,
// normalizing header keys to their capitalized version and values by
// unescaping them when decode is set to true.
//
// It's used after parsing plain text event headers, but not JSON.
func copyHeaders(src *textproto.MIMEHeader, dst *Event, decode bool) {
	var err error
	for k, v := range *src {
		k = capitalize(k)
		if decode {
			dst.Header[k], err = url.QueryUnescape(v[0])
			if err != nil {
				dst.Header[k] = v[0]
			}
		} else {
			dst.Header[k] = v[0]
		}
	}
}

// capitalize capitalizes strings in a very particular manner.
// Headers such as Job-UUID become Job-Uuid and so on. Headers starting with
// Variable_ only replace ^v with V, and headers staring with _ are ignored.
func capitalize(s string) string {
	if s[0] == '_' {
		return s
	}
	ns := bytes.ToLower([]byte(s))
	if len(s) > 9 && s[1:9] == "ariable_" {
		ns[0] = 'V'
		return string(ns)
	}
	toUpper := true
	for n, c := range ns {
		if toUpper {
			if 'a' <= c && c <= 'z' {
				c -= 'a' - 'A'
			}
			ns[n] = c
			toUpper = false
		} else if c == '-' || c == '_' {
			toUpper = true
		}
	}
	return string(ns)
}
