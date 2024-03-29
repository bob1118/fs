package eslclient

import (
	"fmt"

	"github.com/bob1118/fs/db"
	"github.com/bob1118/fs/esl/eventsocket"
	"github.com/bob1118/fs/esl/run_time"
	"github.com/bob1118/fs/utils"
)

// eventAction function.
func eventAction(e *eventsocket.Event) {
	//e.LogPrint()
	eventName := e.Get("Event-Name")
	if len(eventName) > 0 {
		switch eventName {
		case "API":
			apiAction(e)
		case "CUSTOM":
			customAction(e)
		case "BACKGROUND_JOB":
			backgroundjobAction(e)
		case "CHANNEL_HANGUP_COMPLETE":
			// for debug cdr
			if false {
				e.LogPrint()
			}
			if !is_mod_odbc_cdr {
				channelCdrAction(e)
			}
		default:
			//nothing todo.
		}
	}
}

func apiAction(e *eventsocket.Event) {}

func customAction(e *eventsocket.Event) {
	user := e.Get("User_Name")
	domain := e.Get("Domain_Name")
	eventsubclass := e.Get("Event-Subclass")

	if len(eventsubclass) > 0 {
		switch eventsubclass {
		case "sofia::pre_register", "sofia::register_attempt", "sofia::register_failure": //sofia_reg_handle_register_token
		case "sofia::register": //sofia_reg_handle_register_token
			if len(user) > 0 && len(domain) > 0 {
				originate_string := fmt.Sprintf(`sofia/%s/%s`, domain, user)
				run_time.SetUaOnline(e)
				run_time.FifoMemberAdd(ClientCon, originate_string)
			}
		case "sofia::unregister": //sofia_reg_handle_register_token
			if len(user) > 0 && len(domain) > 0 {
				originate_string := fmt.Sprintf(`sofia/%s/%s`, domain, user)
				run_time.SetUaOffline(e)
				run_time.FifoMemberDel(ClientCon, originate_string)
			}
		case "sofia::expire": //sofia_reg_del_call_back
			if len(user) > 0 && len(domain) > 0 {
				originate_string := fmt.Sprintf(`sofia/%s/%s`, domain, user)
				run_time.SetUaOffline(e)
				run_time.FifoMemberDel(ClientCon, originate_string)
			}
		case "sofia::gateway_state": //sofia_reg_fire_custom_gateway_state_event
			run_time.SetGatewayState(e)
		default:
		}
	}
}

// backgroundjobAction function.
func backgroundjobAction(e *eventsocket.Event) {
	job := db.Bgjob{
		Juuid:    e.Get("Job-Uuid"),
		Jcmd:     e.Get("Job-Command"),
		Jcmdarg:  e.Get("Job-Command-Arg"),
		Jcontent: e.Body,
	}

	if false ||
		len(job.Juuid) == 0 ||
		len(job.Jcmd) == 0 {
		fmt.Println(e)
	} else {
		if err := db.CreateBgjob(job); err != nil {
			fmt.Println(err)
		}
	}
}

// channelCdrAction function
func channelCdrAction(e *eventsocket.Event) {
	//
	var isbleg bool
	var uuid, otherUUID, bondUUID, otherType string
	uuid = e.Get("Variable_uuid")
	bondUUID = e.Get("Variable_signal_bond")
	otherType = e.Get("Other-Type")
	isbleg = utils.IsEqual(otherType, "originator")

	if isbleg { //bleg
		otherUUID = utils.UUIDFormat(e.Get("Variable_originator"))
	} else { //if utils.IsEqual(otherType, "originatee"){}else{}
		otherUUID = utils.UUIDFormat(e.Get("Variable_originated_legs"))
	}
	leg := db.CDRLEG{
		UUID:           uuid,
		OtherUUID:      otherUUID,
		BondUUID:       bondUUID,
		Name:           e.Get("Variable_channel_name"),
		Direction:      e.Get("Variable_direction"),
		Sofiaprofile:   e.Get("Variable_sofia_profile_name"),
		Indomain:       e.Get("Variable_domain_name"),
		Ingateway:      e.Get("Variable_sip_gateway"),
		Outdomain:      e.Get("Variable_sip_profile_name"),
		Outgateway:     e.Get("Variable_sip_gateway_name"),
		Ani:            e.Get("Caller-Ani"),
		Destination:    e.Get("Caller-Destination-Number"),
		Calleridname:   e.Get("Caller-Caller-Id-Name"),
		Calleridnumber: e.Get("Caller-Caller-Id-Number"),
		Calleeidname:   e.Get("Caller-Callee-Id-Name"),
		Calleeidnumber: e.Get("Caller-Callee-Id-Number"),
		App:            e.Get("Variable_current_application"),
		Appdata:        e.Get("Variable_current_application_data"),
		Appdialstatus:  e.Get("Variable_dialstatus"),
		Cause:          e.Get("Variable_hangup_cause"),
		Q850:           e.Get("Variable_hangup_cause_q850"),
		Disposition:    e.Get("Variable_sip_hangup_disposition"),
		Protocause:     e.Get("Variable_proto_specific_hangup_cause"),
		Phrase:         e.Get("Variable_sip_hangup_phrase"),
		Startepoch:     e.Get("Variable_start_epoch"),
		Answerepoch:    e.Get("Variable_answer_epoch"),
		Endepoch:       e.Get("Variable_end_epoch"),
		Waitsec:        e.Get("Variable_waitsec"),
		Billsec:        e.Get("Variable_billsec"),
		Duration:       e.Get("Variable_duration"),
	}
	if isbleg {
		go db.CreateCdrBleg(&leg)
	} else {
		go db.CreateCdrAleg(&leg)
	}
}
