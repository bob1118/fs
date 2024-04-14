// //////////////////////channel event action////////////////////

package eslserver

import (
	"errors"
	"fmt"

	"github.com/bob1118/fs/db"
	"github.com/bob1118/fs/esl/eventsocket"
	"github.com/bob1118/fs/utils"
)

// ChannelAction function
func ChannelEventAction(c *eventsocket.Connection, e *eventsocket.Event) {
	//e.LogPrint()
	if utils.IsEqual(e.Get("Content-Type"), "text/disconnect-notice") {
		c.Close()
	}
	eventName := e.Get("Event-Name")
	if len(eventName) > 0 {
		switch eventName {
		case "CHANNEL_PARK":
			channelparkAction(c, e)
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

// channelparkAction function.
func channelparkAction(c *eventsocket.Connection, e *eventsocket.Event) {
	//ChannelDefaultAction(c, e)
	//fmt.Println("Event-Name: CHANNEL_PARK")
}

// DefaultChannelAction
func ChannelDefaultAction(c *eventsocket.Connection, ev *eventsocket.Event) error {
	var myerr error
	if call, err := NewCall(ev); err != nil {
		myerr = err
		fmt.Println(err)
	} else {
		//call route.
		if utils.IsEqual(call.direction, "inbound") { //incoming call hit DIALPLAN_APP_SOCKET
			switch call.profile {
			case "internal", "internal-ipv6": //internal ua incoming
				myerr = channelInternalIncomingProc(c, call)
			case "external", "external-ipv6": //external gateway incoming
				myerr = channelExternalIncomingProc(c, call)
			default:
				c.Hangup("CALL_REJECT")
				myerr = errors.New("CHANNEL_DATA: unknown profile")
			}
		} else { //outgoing call hit socket. bgapi origination... &socket ?
			switch call.profile {
			case "internal", "internal-ipv6":
				myerr = channelInternalOutgoingProc(c, call)
			case "external", "external-ipv6":
				myerr = channelExternalOutgoingProc(c, call)
			default:
				c.Hangup("CALL_REJECT")
				myerr = errors.New("CHANNEL_DATA: unknown profile")
			}
		}
	}
	return myerr
}

// channelInternalIncomingProc
// 1. ua ->ua;
// 2. ua -> gateway ->remote;
// 3. ua as upstream gateway.
func channelInternalIncomingProc(c *eventsocket.Connection, call *CALL) (err error) {
	var (
		uuid  string
		myerr error
	)

	if uuid, myerr = c.APICreateUUID(); err != nil {
		return myerr
	}

	//continue_on_fail=true/continue_on_fail=NORMAL_TEMPORARY_FAILURE,USER_BUSY,NO_ANSWER,NO_ROUTE_DESTINATION
	c.APPSet(`continue_on_fail=false`)
	c.APPSet(`hangup_after_bridge=true`)

	if call.CallerIsUa() {
		if call.CalleeIsUa() { //ua dial ua
			appargv := fmt.Sprintf(`{origination_uuid=%s,origination_caller_id_name=%s,origination_caller_id_number=%s,ignore_early_media=true}sofia/%s/%s`, uuid, "local", call.ani, call.domain, call.distinationnumber)
			c.APPBridge(appargv)
		} else { //ua dial out through gateway.
			q := fmt.Sprintf(`account_id='%s' and account_domain='%s' and acce164_isdefault=true limit 1`, call.ani, call.domain)
			if acce164s, err := db.SelectAcce164sWithCondition(q); err != nil {
				c.Hangup("NO_ROUTE_DESTINATION")
				myerr = err
			} else {
				if len(acce164s) == 0 { //no row.
					c.Hangup("NO_ROUTE_DESTINATION")
					myerr = fmt.Errorf("NO_ROUTE_DESTINATION")
				} else {
					gatewayname := acce164s[0].Gname
					gatewaye164number := acce164s[0].Enumber
					appargv := fmt.Sprintf(`{origination_uuid=%s,origination_caller_id_number=%s,ignore_early_media=true,codec_string=PCMU}sofia/gateway/%s/%s`, uuid, gatewaye164number, gatewayname, call.distinationnumber)
					c.APPBridge(appargv)
				}
			}
		}
	} else {
		//ua as upstream gateway, like channelExternalIncomingPro.
		c.Hangup("CALL_REJECT")
	}
	return myerr
}

// channelExternalIncomingProc
// remote -> gateway ->ua/fifo
func channelExternalIncomingProc(c *eventsocket.Connection, call *CALL) error {
	var myerr error

	//continue_on_fail=true/continue_on_fail=NORMAL_TEMPORARY_FAILURE,USER_BUSY,NO_ANSWER,NO_ROUTE_DESTINATION
	c.APPSet(`continue_on_fail=false`)
	c.APPSet(`hangup_after_bridge=true`)
	if !call.CallFilterPassed() {
		c.Hangup("CALL_REJECT")
		myerr = errors.New("function CallFilterPassed fail, Call Reject")
	} else {
		q := fmt.Sprintf("gateway_name='%s' and e164_number='%s'", call.gateway, call.distinationnumber)
		if e164accs, err := db.SelectE164accsWithCondition(q); err != nil {
			c.Hangup("NO_ROUTE_DESTINATION")
			myerr = err
		} else {
			if len(e164accs) == 0 {
				c.Hangup("NO_ROUTE_DESTINATION")
				myerr = err
			} else {
				e164acc := e164accs[0]
				if !e164acc.Isfifo { // do bridge sofia/mydomain/xxxx
					appargv := fmt.Sprintf(`{origination_caller_id_number=%s,ignore_early_media=true,codec_string=PCMU}sofia/%s/%s`, call.ani, e164acc.Adomain, e164acc.Aid)
					myerr = c.APPBridge(appargv)
				} else { //do fifo myfifo in
					appargv := fmt.Sprintf(`%s in`, e164acc.Fname)
					myerr = c.APPFifo(appargv)
				}
			}
		}
	}
	return myerr
}

// channelInternalOutgoingProc function.
// http client post data ->http server receive data ->bgapi originate sofia/profile/id@domain &socket ...
func channelInternalOutgoingProc(c *eventsocket.Connection, call *CALL) error { return nil }

// channelExternalOutgoingProc function.
// http client post data ->http server receive data ->bgapi originate sofia/gatewayname/reomte &socket...
func channelExternalOutgoingProc(c *eventsocket.Connection, call *CALL) error { return nil }

// channelstateAction function.
func channelstateAction(c *eventsocket.Connection, e *eventsocket.Event) {}

// channelcallstateAction function.
func channelcallstateAction(c *eventsocket.Connection, e *eventsocket.Event) {}

// channelhangupAction function.
func channelhangupAction(c *eventsocket.Connection, e *eventsocket.Event) {}

// channelCDRAction function. channel cdr.
func channelCDRAction(c *eventsocket.Connection, e *eventsocket.Event) {}
