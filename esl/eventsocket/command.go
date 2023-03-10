//Package eventsocket send command and wait return.
package eventsocket

import (
	"fmt"
)

//SendCommand function return command reply result.
//
//https://freeswitch.org/confluence/display/FREESWITCH/mod_event_socket#mod_event_socket-3.CommandDocumentation
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

//SendCommandEx function,return event.
//
//https://freeswitch.org/confluence/display/FREESWITCH/Event+Socket+Outbound#EventSocketOutbound-ServerExamples
func (h *Connection) SendCommandEx(cmd string) (*Event, error) {
	return h.Send(cmd)
}

//eventsocket send api command return api response body.
//
//3.1 api, blocking mode.
//
//https://freeswitch.org/confluence/display/FREESWITCH/mod_event_socket#mod_event_socket-3.1api
//https://freeswitch.org/confluence/display/FREESWITCH/mod_commands
func (h *Connection) SendApiCommand(cmd string) (string, error) {
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

//eventsocket send bgapi command return command reply Job-UUID.
//
//3.2 bgapi, non-blocking mode.
//
//https://freeswitch.org/confluence/display/FREESWITCH/mod_event_socket#mod_event_socket-3.2bgapi
//https://freeswitch.org/confluence/display/FREESWITCH/mod_commands
func (h *Connection) SendBgapiCommand(cmd string) (uuid string, e error) {
	// 	bgapi version

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

//SendEVENT function.
// send event.
//
//3.8 sendevent
//https://freeswitch.org/confluence/display/FREESWITCH/mod_event_socket#mod_event_socket-3.8sendevent
func (h *Connection) SendEVENT() {}

//SendMSG function.
// send msg.
//
//3.9 sendmsg
//3.9.1 call-command
//3.9.1.1 execute
//3.9.1.2 hangup
//3.9.1.3 unicast
//3.9.1.4 nomedia
//3.9.1.5 xferxt
//
//https://freeswitch.org/confluence/display/FREESWITCH/mod_event_socket#mod_event_socket-3.9sendmsg
func (h *Connection) SendMSG(uuid string, m MSG, content string) (*Event, error) {
	return h.SendMsg(m, uuid, content)
}

//ExecuteDptools Execute is a shortcut to SendMsg with call-command: execute without UUID,
//
//suitable for use on outbound event socket connections (acting as server).
//
//Execute("set", "a=b", true)
//
//https://freeswitch.org/confluence/display/FREESWITCH/mod_event_socket#mod_event_socket-3.9.1.1execute
//https://freeswitch.org/confluence/display/FREESWITCH/Event+Socket+Outbound#EventSocketOutbound-FAQ
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

//ExecuteDptoolsEx is similar to Execute, but takes a UUID and no lock. Suitable
//for use on inbound event socket connections (acting as client).
func (h *Connection) ExecuteDptoolsEx(uuid, appName, appArg string) (*Event, error) {
	//var msg MSG
	msg := MSG{
		"call-command":     "execute",
		"execute-app-name": appName,
		"execute-app-arg":  appArg,
	}
	return h.SendMSG(uuid, msg, "")
}

//Hangup function.
//
// sendmsg <uuid>
//
// call-command: hangup
//
// hangup-cause: <one of the causes listed below>
//
func (h *Connection) Hangup(cause string) error {
	msg := MSG{
		"call-command": "hangup",
		"hangup-cause": cause,
	}
	_, err := h.SendMSG("", msg, "")
	return err
}

//HangupEx function, hangup with uuid
func (h *Connection) HangupEx(uuid, cause string) (*Event, error) {
	msg := MSG{
		"call-command": "hangup",
		"hangup-cause": cause,
	}
	return h.SendMSG(uuid, msg, "")
}

//
func (h *Connection) Unicast() {}

//
func (h *Connection) Xferxt() {}

////////////////////////////////////////////////////////////////////////////
/////////////////////////////api and app////////////////////////////////////
////////////////////////////////////////////////////////////////////////////

//APICreateUUID function.
// send api command "api create_uuid" and response uuid.
//
//https://freeswitch.org/confluence/display/FREESWITCH/mod_commands#mod_commands-create_uuid
func (h *Connection) APICreateUUID() (string, error) {

	return h.SendApiCommand("create_uuid")
}

//APIModuleExists function.
// send api command "api module_exists modulename" and response true or false.
//
//https://freeswitch.org/confluence/display/FREESWITCH/mod_commands#mod_commands-module_exists
func (h *Connection) APIModuleExists(modulename string) (bool, error) {
	var (
		exist bool
		myerr error
	)
	command := fmt.Sprintf("module_exists %s", modulename)
	if body, err := h.SendApiCommand(command); err != nil {
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

//APPSet function, dptools set.
//
//<action application="set" data="effective_caller_id_number=12345678"/>
//
//https://freeswitch.org/confluence/display/FREESWITCH/mod_dptools%3A+set
func (h *Connection) APPSet(data string, lock bool) error {
	_, err := h.ExecuteDptools("set", data, lock)
	return err
}

//APPBridge function, dptools bridge.
//
//<action application="bridge" data="endpoint/gateway/gateway_name/address"/>
//
//https://freeswitch.org/confluence/display/FREESWITCH/mod_dptools%3A+bridge
func (h *Connection) APPBridge(data string, lock bool) error {
	_, err := h.ExecuteDptools("bridge", data, lock)
	return err
}

//APPFifo function, dptools fifo.
//
//<action application="fifo" data="myqueue in /tmp/exit-message.wav /tmp/music-on-hold.wav"/>
//
//<action application="fifo" data="myqueue out nowait /tmp/caller-found.wav /tmp/agent-music-on-hold.wav"/>
//
//https://freeswitch.org/confluence/display/FREESWITCH/mod_fifo
func (h *Connection) APPFifo(data string, lock bool) error {
	_, err := h.ExecuteDptools("fifo", data, lock)
	return err
}

//APPHangup function
//
//<action application="hangup" data="USER_BUSY"/>
//
//https://freeswitch.org/confluence/display/FREESWITCH/mod_dptools%3A+hangup
func (h *Connection) APPHangup(data string) error {
	_, err := h.ExecuteDptools("hangup", data, false)
	return err
}

//APPAnswer function
//
//<action application="answer"/>
//
//https://freeswitch.org/confluence/display/FREESWITCH/mod_dptools%3A+answer
func (h *Connection) APPAnswer() error {
	_, err := h.ExecuteDptools("answer", "", true)
	return err
}

//APPAnswer function
//
//<action application="pre_answer"/>
//
//https://freeswitch.org/confluence/display/FREESWITCH/mod_dptools%3A+pre+answer
func (h *Connection) APPPreAnswer() error {
	_, err := h.ExecuteDptools("pre_answer", "", true)
	return err
}
