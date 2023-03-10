package eslserver

import (
	"errors"

	"github.com/bob1118/fs/esl/eventsocket"
	"github.com/bob1118/fs/esl/run_time"
)

type CALL struct {
	//CALL BASIC INFO
	coreuuid          string //Core-Uuid
	fsipv4            string //Freeswitch-Ipv4
	eventname         string //Event-Name
	uuid              string //Variable_uuid
	callid            string //Variable_call_uuid
	direction         string //Variable_direction
	profile           string //Variable_sofia_profile_name
	domain            string //Variable_domain_name
	gateway           string //Variable_sip_gateway
	ani               string //Caller-Ani
	distinationnumber string //Caller-Destination-Number
}

func NewCall(ev *eventsocket.Event) (c *CALL, e error) {
	var err error

	call := &CALL{
		coreuuid:          ev.Get("Core-Uuid"),
		fsipv4:            ev.Get("Freeswitch-Ipv4"),
		eventname:         ev.Get("Event-Name"),
		uuid:              ev.Get("Variable_uuid"),
		callid:            ev.Get("Variable_call_uuid"),
		direction:         ev.Get("Variable_direction"),
		profile:           ev.Get("Variable_sofia_profile_name"),
		domain:            ev.Get("Variable_domain_name"),
		gateway:           ev.Get("Variable_sip_gateway"),
		ani:               ev.Get("Caller-Ani"),
		distinationnumber: ev.Get("Caller-Destination-Number"),
	}
	if len(call.coreuuid) > 0 &&
		len(call.fsipv4) > 0 &&
		len(call.eventname) > 0 &&
		len(call.uuid) > 0 &&
		len(call.callid) > 0 &&
		len(call.direction) > 0 &&
		len(call.profile) > 0 &&
		(len(call.domain) > 0 || len(call.gateway) > 0) &&
		len(call.ani) > 0 &&
		len(call.distinationnumber) > 0 {
		err = nil
	} else {
		err = errors.New("newcall param check fail")
	}
	return call, err
}

func (c *CALL) CallerIsUa() bool {
	mkey := c.ani + "@" + c.domain
	return run_time.IsUa(mkey)
}
func (c *CALL) CalleeIsUa() bool {
	mkey := c.distinationnumber + "@" + c.domain
	return run_time.IsUa(mkey)
}

func (c *CALL) CallFilterPassed() bool {
	passed := false
	// if c.e164CalleeExist() {
	// 	if !c.blacklistCallerExist() {
	// 		passed = true
	// 	}
	// }
	return passed
}

// func (c *CALL) e164CalleeExist() bool {
// 	is := false
// 	if len(c.distinationnumber) > 0 && len(c.gateway) > 0 {
// 		if e164, exist := db.IsExistE164BynumberEx(c.gateway, c.distinationnumber); !exist {
// 			is = false
// 		} else {
// 			if e164.Eenable {
// 				if !e164.Elockin {
// 					is = true
// 				}
// 			}
// 		}
// 	}
// 	return is
// }

// func (c *CALL) blacklistCallerExist() bool {
// 	is := false
// 	if len(c.ani) > 0 {
// 		if _, exist := db.IsExistBlacklistCaller(c.ani, c.distinationnumber); exist {
// 			is = true
// 		}
// 	}
// 	return is
// }
