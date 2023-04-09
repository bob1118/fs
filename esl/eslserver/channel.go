package eslserver

import (
	"errors"
	"fmt"

	"github.com/bob1118/fs/db"
	"github.com/bob1118/fs/esl/eventsocket"
	"github.com/bob1118/fs/utils"
)

////////////////////first event CHANNEL_DATA action///////////////////////////

// DefaultChannelAction
func ChannelDefaultAction(c *eventsocket.Connection, ev *eventsocket.Event) error {
	var myerr error
	if call, err := NewCall(ev); err != nil {
		myerr = err
		fmt.Println(err)
	} else {
		//send myevents
		c.SendCommand("myevents")
		//call route.
		if utils.IsEqual(call.direction, "inbound") { //incoming call
			switch call.profile {
			case "internal", "internal-ipv6": //internal ua incoming
				myerr = channelInternalIncomingProc(c, call)
			case "external", "external-ipv6": //external gateway incoming
				myerr = channelExternalIncomingProc(c, call)
			default:
				c.APPHangup("CALL_REJECT")
				myerr = errors.New("CHANNEL_DATA: unknown profile")
			}
		} else { //outgoing call hit socket. bgapi origination... &socket ?
			switch call.profile {
			case "internal", "internal-ipv6":
				myerr = channelInternaOutgoingProc(c, call)
			case "external", "external-ipv6":
				myerr = channelExternalOutgoingProc(c, call)
			default:
				c.APPHangup("CALL_REJECT")
				myerr = errors.New("CHANNEL_DATA: unknown profile")
			}
		}
	}
	return myerr
}

// channelInternalIncomingProc
func channelInternalIncomingProc(c *eventsocket.Connection, call *CALL) (err error) {
	var (
		uuid  string
		myerr error
	)

	if uuid, myerr = c.APICreateUUID(); err != nil {
		return myerr
	}

	//continue_on_fail=true/continue_on_fail=NORMAL_TEMPORARY_FAILURE,USER_BUSY,NO_ANSWER,NO_ROUTE_DESTINATION
	c.APPSet(`continue_on_fail=NORMAL_TEMPORARY_FAILURE,USER_BUSY,NO_ANSWER,NO_ROUTE_DESTINATION`, true)
	c.APPSet(`hangup_after_bridge=true`, true)

	if call.CallerIsUa() {
		if call.CalleeIsUa() { //ua dial ua
			appargv := fmt.Sprintf(`{origination_uuid=%s,origination_caller_id_name=%s,origination_caller_id_number=%s,ignore_early_media=true}sofia/%s/%s`, uuid, "local", call.ani, call.domain, call.distinationnumber)
			c.APPBridge(appargv, true)
		} else { //ua dial out through gateway.
			q := fmt.Sprintf(`account_id='%s' and account_domain='%s' and acce164_isdefault=true limit 1`, call.ani, call.domain)
			if acce164s, err := db.GetAcce164s(q); err != nil { //row not found.
				c.APPHangup("NO_ROUTE_DESTINATION")
				myerr = err
			} else {
				gatewayname := acce164s[0].Gname
				gatewaye164number := acce164s[0].Enumber
				appargv := fmt.Sprintf(`{origination_uuid=%s,origination_caller_id_number=%s,ignore_early_media=true,codec_string="PCMU,PCMA"}sofia/gateway/%s/%s`, uuid, gatewaye164number, gatewayname, call.distinationnumber)
				c.APPBridge(appargv, true)
			}
		}
	} else {
		c.APPHangup("USER_NOT_REGISTERED")
	}
	return myerr
}

// channelExternalIncomingProc
func channelExternalIncomingProc(c *eventsocket.Connection, call *CALL) (err error) {

	if !call.CallFilterPassed() {
		c.Hangup("CALL_REJECT")
		return errors.New("function CallFilterPassed fail, Call Reject")
	} else {
		return channelExternalExecuteFifo(c)
	}
}

func channelExternalExecuteFifo(c *eventsocket.Connection) error {
	//Put a caller into a FIFO queue
	//<action application="fifo" data="myqueue in /tmp/exit-message.wav /tmp/music-on-hold.wav"/>
	var err error
	argv := `fifomember@fifos in`
	c.APPSet(`hangup_after_bridge=true`, true)
	if err = c.APPFifo(argv, true); err != nil {
		fmt.Println(err)
	}
	return err
}

// channelInternaOutgoingProc function.
func channelInternaOutgoingProc(c *eventsocket.Connection, call *CALL) error { return nil }

// channelExternalOutgoingProc function.
func channelExternalOutgoingProc(c *eventsocket.Connection, call *CALL) error { return nil }

// //////////////////////channel event action////////////////////
// ChannelAction function
func ChannelAction(c *eventsocket.Connection, e *eventsocket.Event) {
	//e.LogPrint()
	eventName := e.Get("Event-Name")
	if len(eventName) > 0 {
		switch eventName {
		case "CHANNEL_STATE":
			channelstateAction(c, e)
		case "CHANNEL_CALLSTATE":
			channelcallstateAction(c, e)
		case "CHANNEL_HANGUP":
			channelhangupAction(c, e)
		case "CHANNEL_DESTROY":
			channelCDRAction(c, e)
		default:
			//nothing todo.
		}
	}
}

// channelstateAction function.
func channelstateAction(c *eventsocket.Connection, e *eventsocket.Event) {}

// channelcallstateAction function.
func channelcallstateAction(c *eventsocket.Connection, e *eventsocket.Event) {}

// channelhangupAction function.
func channelhangupAction(c *eventsocket.Connection, e *eventsocket.Event) {}

// channelCDRAction function. channel cdr.
func channelCDRAction(c *eventsocket.Connection, e *eventsocket.Event) {}
