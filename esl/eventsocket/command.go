// Package eventsocket send command and wait return.
package eventsocket

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/bob1118/fs/utils"
)

// Send sends a single command to the server and returns a Event.
//
// See https://developer.signalwire.com/freeswitch/FreeSWITCH-Explained/Modules/mod_event_socket_1048924#3-command-documentation for
// details.
func (h *Connection) Send(command string) (*Event, error) {
	// Sanity check to avoid breaking the parser
	//if strings.IndexAny(command, "\r\n") > 0 {
	//	return nil, errInvalidCommand
	//}
	if _, err := fmt.Fprintf(h.conn, "%s\r\n\r\n", command); err != nil {
		return nil, err
	}
	var (
		ev  *Event
		err error
	)
	select {
	case err = <-h.cmderr:
		return nil, err
	case err = <-h.apierr:
		return nil, err
	case ev = <-h.cmd:
		return ev, nil
	case ev = <-h.api:
		return ev, nil
	case <-time.After(timeoutPeriod):
		return nil, errTimeout
	}
}

// SendCommand function send command and return send result.
//
// https://developer.signalwire.com/freeswitch/FreeSWITCH-Explained/Modules/mod_event_socket_1048924/#3-command-documentation
func (h *Connection) SendCommand(cmd string) error {

	// auth ClueCon

	// Content-Type: command/reply
	// Reply-Text: +OK accepted

	// hello

	// Content-Type: command/reply
	// Reply-Text: -ERR command not found

	_, err := h.Send(cmd)
	return err
}

// SendCommandEx function,return event.
//
// https://developer.signalwire.com/freeswitch/FreeSWITCH-Explained/Modules/mod_event_socket_1048924/#3-command-documentation
func (h *Connection) SendCommandEx(cmd string) (*Event, error) {
	return h.Send(cmd)
}

// eventsocket send api command return api response body.
//
// 3.1 api, blocking mode.
//
// command api: https://developer.signalwire.com/freeswitch/FreeSWITCH-Explained/Modules/mod_event_socket_1048924#31-api
//
// mod_command: https://developer.signalwire.com/freeswitch/FreeSWITCH-Explained/Modules/mod_commands_1966741/#--
func (h *Connection) SendApiCommandSync(cmd string) (string, error) {
	// api version

	// Content-Type: api/response
	// Content-Length: 85

	// FreeSWITCH Version 1.10.6-release-18-1ff9d0a60e~64bit (-release-18-1ff9d0a60e 64bit)

	// 	api ver

	// Content-Type: api/response
	// Content-Length: 28

	// -ERR ver Command not found!

	var (
		response string
		myerr    error
	)
	if len(cmd) > 0 {
		apicommand := fmt.Sprintf("api %s", cmd)
		if ev, err := h.Send(apicommand); err != nil {
			myerr = err
		} else {
			response = ev.Body
		}
	}
	return response, myerr
}

// eventsocket send bgapi command return command reply Job-UUID.
//
// 3.2 bgapi, non-blocking mode.
//
// https://developer.signalwire.com/freeswitch/FreeSWITCH-Explained/Modules/mod_event_socket_1048924#32-bgapi
//
// https://developer.signalwire.com/freeswitch/FreeSWITCH-Explained/Modules/mod_commands_1966741/#--
func (h *Connection) SendBgapiCommandAsync(cmd string) (uuid string, e error) {
	// 	bgapi version

	// Content-Type: command/reply
	// Reply-Text: +OK Job-UUID: 2b5563b2-d465-4d90-8abd-52a032d0933f
	// ...
	// job working...
	// ...
	// Content-Type: command/reply
	// Reply-Text: +OK Job-UUID: 2b5563b2-d465-4d90-8abd-52a032d0933f
	// Job-UUID: 2b5563b2-d465-4d90-8abd-52a032d0933f

	var (
		jobuuid string
		myerr   error
	)
	if len(cmd) > 0 {
		bgapicommand := fmt.Sprintf("bgapi %s", cmd)
		if ev, err := h.Send(bgapicommand); err != nil {
			myerr = err
		} else {
			jobuuid = ev.Get("Job-Uuid")
		}
	}
	return jobuuid, myerr
}

// eventsocket send linger.
//
// 3.3 linger
//
// https://developer.signalwire.com/freeswitch/FreeSWITCH-Explained/Modules/mod_event_socket_1048924#33-linger
func (h *Connection) Linger() error { return h.SendCommand(`linger`) }

// eventsocket send nolinger.
//
// 3.4 nolinger
//
// https://developer.signalwire.com/freeswitch/FreeSWITCH-Explained/Modules/mod_event_socket_1048924#34-nolinger
func (h *Connection) NoLinger() error { return h.SendCommand(`nolinger`) }

// EVENT function.
// suscribe event.
//
// 3.5 EVENT
// https://developer.signalwire.com/freeswitch/FreeSWITCH-Explained/Modules/mod_event_socket_1048924#35-event
func (h *Connection) EVENT() {}

// Filter function.
// filter event.
//
// 3.6 Filter
// https://developer.signalwire.com/freeswitch/FreeSWITCH-Explained/Modules/mod_event_socket_1048924#36-filter
func (h *Connection) Filter() {}

// FilterDelete function.
// filter delete event.
//
// 3.7 FilterDelete
// https://developer.signalwire.com/freeswitch/FreeSWITCH-Explained/Modules/mod_event_socket_1048924#36-filter
func (h *Connection) FilterDelete() {}

// SendEvent function.
// send event.
//
// 3.8 SendEvent
//
// 3.8.1.1 Switch phone MWI led (tested on yealink)
//
// 3.8.1.2 Make Snom phones reread their settings from the settings server
//
// 3.8.1.3 sendevent examples with a message body
//
// 3.8.1.4 SIP Proxy Example
//
// 3.8.1.5 Sipura/Linksys/Cisco phone or ATA to resync config with a specified profile URL
//
// 3.8.1.6 Example usage for CSTA event:
//
// 3.8.1.7 Display a text message on a Snom 370 or Snom 820
//
// https://developer.signalwire.com/freeswitch/FreeSWITCH-Explained/Modules/mod_event_socket_1048924/#38-sendevent
func (h *Connection) SendEvent() {}

type MSG map[string]string

// SendMSG function.
//
// 3.9 sendmsg
//
// MSG is the container used by SendMsg to store messages sent to FreeSWITCH.
// It's supposed to be populated with directives supported by the sendmsg
// command only, like "call-command: execute".
//
// See https://developer.signalwire.com/freeswitch/FreeSWITCH-Explained/Modules/mod_event_socket_1048924#39-sendmsg for details.
//
// Keys with empty values are ignored; uuid and appData are optional.
// If appData is set, a "content-length" header is expected (lower case!).
//
// See https://developer.signalwire.com/freeswitch/FreeSWITCH-Explained/Modules/mod_event_socket_1048924#39-sendmsg for details.
func (h *Connection) SendMsg(m MSG, uuid, appData string) (*Event, error) {
	// sendmsg <UUID>
	// <headers>
	//
	// <body>
	//
	// Example:
	// sendmsg <uuid>
	// call-command: execute
	// execute-app-name: playback
	// execute-app-arg: /tmp/test.wav

	b := bytes.NewBufferString("sendmsg")
	if uuid != "" {
		// Make sure there's no \r or \n in the UUID.
		if strings.IndexAny(uuid, "\r\n") > 0 {
			return nil, errInvalidCommand
		}
		b.WriteString(" " + uuid)
	}
	b.WriteString("\n")
	for k, v := range m {
		// Make sure there's no \r or \n in the key, and value.
		if strings.IndexAny(k, "\r\n") > 0 {
			return nil, errInvalidCommand
		}
		if v != "" {
			if strings.IndexAny(v, "\r\n") > 0 {
				return nil, errInvalidCommand
			}
			b.WriteString(fmt.Sprintf("%s: %s\n", k, v))
		}
	}
	b.WriteString("\n")
	if m["content-length"] != "" && appData != "" {
		b.WriteString(appData)
	}
	if _, err := b.WriteTo(h.conn); err != nil {
		return nil, err
	}
	var (
		ev  *Event
		err error
	)
	select {
	case err = <-h.cmderr:
		return nil, err
	case ev = <-h.cmd:
		return ev, nil
	case <-time.After(timeoutPeriod):
		return nil, errTimeout
	}
}

// SendMSG function.
// send msg.
//
// https://freeswitch.org/confluence/display/FREESWITCH/mod_event_socket#mod_event_socket-3.9sendmsg
func (h *Connection) SendMSG(uuid string, m MSG, content string) (*Event, error) {
	return h.SendMsg(m, uuid, content)
}

// Execute is a shortcut to SendMsg with call-command: execute without UUID,
// suitable for use on outbound event socket connections (acting as server).
// https://developer.signalwire.com/freeswitch/FreeSWITCH-Explained/Modules/mod_event_socket_1048924/#3911-execute
// https://developer.signalwire.com/freeswitch/FreeSWITCH-Explained/Modules/mod_dptools_1970333/#-c
//
// 3.9.1.1 execute
//
//	Execute("playback", "/tmp/test.wav", false)
func (h *Connection) Execute(appName, appArg string, lock bool) (*Event, error) {
	// sendmsg <uuid>
	// call-command: execute
	// execute-app-name: <one of the applications>
	// loops: <number of times to invoke the command, default: 1>
	// content-type: text/plain
	// content-length: <content length>
	//
	// <application data>
	//
	// Example:
	// sendmsg
	// call-command: execute
	// execute-app-name: set
	// execute-app-arg: foo=bar\n\n
	// event-lock: true
	var evlock string
	if lock {
		// Could be strconv.FormatBool(lock), but we don't want to
		// send event-lock when it's set to false.
		evlock = "true"
	}
	return h.SendMsg(MSG{
		"call-command":     "execute",
		"execute-app-name": appName,
		"execute-app-arg":  appArg,
		"event-lock":       evlock,
	}, "", "")
}

// ExecuteUUID is similar to Execute, but takes a UUID and no lock. Suitable
// for use on inbound event socket connections (acting as client).
func (h *Connection) ExecuteUUID(uuid, appName, appArg string) (*Event, error) {
	return h.SendMsg(MSG{
		"call-command":     "execute",
		"execute-app-name": appName,
		"execute-app-arg":  appArg,
	}, uuid, "")
}

// ExecuteDptools Execute is a shortcut to SendMsg with call-command: execute without UUID,
//
// suitable for use on outbound event socket connections (acting as server).
//
// Execute("set", "a=b", true)
//
// https://developer.signalwire.com/freeswitch/FreeSWITCH-Explained/Modules/mod_event_socket_1048924/#3911-execute
func (h *Connection) ExecuteDptools(appName, appArg string, lock bool) (*Event, error) {
	var msg MSG
	if lock {
		msg = MSG{
			"call-command":     "execute",
			"execute-app-name": appName,
			"execute-app-arg":  appArg,
			"event-lock":       "true",
		}
	} else {
		msg = MSG{
			"call-command":     "execute",
			"execute-app-name": appName,
			"execute-app-arg":  appArg,
		}
	}
	return h.SendMSG("", msg, "")
}

// ExecuteDptoolsEx is similar to Execute, but takes a UUID and no lock. Suitable
// for use on inbound event socket connections (acting as client).
func (h *Connection) ExecuteDptoolsEx(uuid, appName, appArg string) (*Event, error) {
	//var msg MSG
	msg := MSG{
		"call-command":     "execute",
		"execute-app-name": appName,
		"execute-app-arg":  appArg,
	}
	return h.SendMSG(uuid, msg, "")
}

// Hangup function.
//
// 3.9.1.2 hangup
//
// sendmsg <uuid>
//
// call-command: hangup
//
// hangup-cause: <one of the causes listed below>
// https://developer.signalwire.com/freeswitch/FreeSWITCH-Explained/Troubleshooting-Debugging/Hangup-Cause-Code-Table_3964945/#q850-to-sip-code-table
func (h *Connection) Hangup(cause string) error {
	msg := MSG{
		"call-command": "hangup",
		"hangup-cause": cause,
	}
	_, err := h.SendMSG("", msg, "")
	return err
}

// HangupEx function, hangup with uuid
func (h *Connection) HangupEx(uuid, cause string) (*Event, error) {
	msg := MSG{
		"call-command": "hangup",
		"hangup-cause": cause,
	}
	return h.SendMSG(uuid, msg, "")
}

// Unicast function.
// unicast is used to hook up mod_spandsp for faxing over a socket.
//
// 3.9.1.3 unicast
func (h *Connection) Unicast() error { return nil }

// Nomedia function.
//
// sendmsg <uuid>
// call-command: nomedia
// nomedia-uuid: <noinfo>
//
// 3.9.1.4 nomedia
func (h *Connection) Nomedia() error { return nil }

// Xferext function.
//
// 3.9.1.5 xferext

func (h *Connection) Xferext() error { return nil }

// Exit function.
//
// 3.10 exit
func (h *Connection) Exit() error { return h.SendCommand(`exit`) }

// Auth function.
//
// 3.11 auth
func (h *Connection) Auth(password string) error {
	var myerr error
	authcmd := fmt.Sprintf("auth %s", password)
	if ev, err := h.SendCommandEx(authcmd); err != nil {
		myerr = err
	} else {
		if reply := ev.Get("Reply-Text"); !utils.IsEqual(reply, "+OK accepted") {
			myerr = fmt.Errorf(reply)
		}
	}
	return myerr
}

// Log function.
//
// 3.12 log

func (h *Connection) Log(level string) error {
	return h.SendCommand("log" + level)
}

// NoLog function.
//
// 3.13 NoLog
func (h *Connection) NoLog() error { return h.SendCommand("nolog") }

// NixEvent function.
//
// nixevent <event types | ALL  | CUSTOM custom event sub-class>
//
// 3.14 nixevent
func (h *Connection) NixEvent(format string, enames ...string) error {
	var eventnames string
	if len(enames) == 0 {
		eventnames = `all`
	} else {
		for _, ename := range enames {
			eventnames += fmt.Sprintf(" %s", ename)
		}
	}
	nixeventcmd := fmt.Sprintf("nixevent %s %s", format, eventnames)
	return h.SendCommand(nixeventcmd)
}

// NoEvents function.
//
// 3.15 noevents
func (h *Connection) NoEvents() error { return h.SendCommand("noevents") }

////////////////////////////////////////////////////////////////////////////
/////////////////////////////api and app////////////////////////////////////
////////////////////////////////////////////////////////////////////////////

// APICreateUUID function.
// send api command "api create_uuid" and response uuid.
//
// https://freeswitch.org/confluence/display/FREESWITCH/mod_commands#mod_commands-create_uuid
func (h *Connection) APICreateUUID() (string, error) {

	return h.SendApiCommandSync("create_uuid")
}

// APIModuleExists function.
// send api command "api module_exists modulename" and response true or false.
//
// https://freeswitch.org/confluence/display/FREESWITCH/mod_commands#mod_commands-module_exists
func (h *Connection) APIModuleExists(modulename string) (bool, error) {
	var (
		exist bool
		myerr error
	)
	command := fmt.Sprintf("module_exists %s", modulename)
	if body, err := h.SendApiCommandSync(command); err != nil {
		myerr = err
	} else {
		switch body {
		case "true", "TRUE":
			exist = true
		case "false", "FALSE":
			exist = false
		}
	}
	return exist, myerr
}

func (h *Connection) APILoadModule(modulename string) (bool, error) { return true, nil }

func (h *Connection) APIUnloadModule(modulename string) (bool, error) { return true, nil }

func (h *Connection) APIReloadModule(modulename string) (bool, error) { return true, nil }

// APPSet function, dptools set.
//
// <action application="set" data="effective_caller_id_number=12345678"/>
//
// https://freeswitch.org/confluence/display/FREESWITCH/mod_dptools%3A+set
func (h *Connection) APPSet(data string, lock bool) error {
	_, err := h.ExecuteDptools("set", data, lock)
	return err
}

// APPBridge function, dptools bridge.
//
// <action application="bridge" data="endpoint/gateway/gateway_name/address"/>
//
// https://freeswitch.org/confluence/display/FREESWITCH/mod_dptools%3A+bridge
func (h *Connection) APPBridge(data string, lock bool) error {
	_, err := h.ExecuteDptools("bridge", data, lock)
	return err
}

// APPFifo function, dptools fifo.
//
// <action application="fifo" data="myqueue in /tmp/exit-message.wav /tmp/music-on-hold.wav"/>
//
// <action application="fifo" data="myqueue out nowait /tmp/caller-found.wav /tmp/agent-music-on-hold.wav"/>
//
// https://freeswitch.org/confluence/display/FREESWITCH/mod_fifo
func (h *Connection) APPFifo(data string, lock bool) error {
	_, err := h.ExecuteDptools("fifo", data, lock)
	return err
}

// APPFifo function, dptools fifo.
//
// <action application="fifo" data="myqueue in /tmp/exit-message.wav /tmp/music-on-hold.wav"/>
//
// <action application="fifo" data="myqueue out nowait /tmp/caller-found.wav /tmp/agent-music-on-hold.wav"/>
//
// https://freeswitch.org/confluence/display/FREESWITCH/mod_fifo
func (h *Connection) APPAcd(data string, lock bool) error {
	_, err := h.ExecuteDptools("acd", data, lock)
	return err
}

// APPHangup function
//
// <action application="hangup" data="USER_BUSY"/>
//
// https://freeswitch.org/confluence/display/FREESWITCH/mod_dptools%3A+hangup
func (h *Connection) APPHangup(data string) error {
	_, err := h.ExecuteDptools("hangup", data, false)
	return err
}

// APPAnswer function
//
// <action application="answer"/>
//
// https://freeswitch.org/confluence/display/FREESWITCH/mod_dptools%3A+answer
func (h *Connection) APPAnswer() error {
	_, err := h.ExecuteDptools("answer", "", true)
	return err
}

// APPAnswer function
//
// <action application="pre_answer"/>
//
// https://freeswitch.org/confluence/display/FREESWITCH/mod_dptools%3A+pre+answer
func (h *Connection) APPPreAnswer() error {
	_, err := h.ExecuteDptools("pre_answer", "", true)
	return err
}
