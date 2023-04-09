package run_time

import (
	"fmt"
	"log"
	"sync"

	"github.com/bob1118/fs/esl/eventsocket"
)

var uamap sync.Map
var gwmap sync.Map
var chmap sync.Map

// runtime useragent
type rtua struct {
	coreuuid    string //Core-Uuid
	hostipv4    string //Freeswitch-Ipv4
	hostipv6    string //Freeswitch-Ipv6
	profilename string //Profile-Name
	callid      string //Call-Id
	expires     string //Expires
	user        string //User_Name
	domain      string //Domain_Name

	//from db
	defaultgawayname  string
	defaulte164number string
}

// runtime gateway
type rtgw struct {
	coreuuid     string //Core-Uuid
	hostipv4     string //Freeswitch-Ipv4
	hostipv6     string //Freeswitch-Ipv6
	registerip   string //Register-Network-Ip
	registerport string //Register-Network-Port
	gateway      string //Gateway:myfsgateway
	pingstatus   string //Ping-Status:DOWN/UP
	state        string //State:DOWN/TRYING/REGISTER/REGED
	//phrase       string //Phrase:OK
	//status       string //Status:200
}

// runtime channel
type rtch struct{}

func SetUaOnline(e *eventsocket.Event) {
	ua := rtua{
		coreuuid:    e.Get("Core-Uuid"),
		hostipv4:    e.Get("Freeswitch-Ipv4"),
		hostipv6:    e.Get("Freeswitch-Ipv6"),
		profilename: e.Get("Profile-Name"),
		callid:      e.Get("Call-Id"),
		expires:     e.Get("Expires"),
		user:        e.Get("User_Name"),
		domain:      e.Get("Domain_Name"),
	}

	k := ua.user + "@" + ua.domain
	if _, isloaded := uamap.LoadOrStore(k, &ua); isloaded {
		uamap.Delete(k)
		uamap.Store(k, &ua)
		fmt.Println("SetUaOnline: update value")
	} else {
		fmt.Println(ua)
	}
}

func SetUaOffline(e *eventsocket.Event) {
	k := e.Get("User_Name") + "@" + e.Get("Domain_Name")
	if ua, isloaded := uamap.LoadAndDelete(k); isloaded {
		fmt.Println(ua)
	}
}

func IsUa(k interface{}) bool {
	_, ok := uamap.Load(k)
	return ok
}

func GetUaDefaultGatewayName(k interface{}) (n string) {
	var name string
	if value, ok := uamap.Load(k); ok {
		name = value.(rtua).defaultgawayname
	} else {
		name = ``
	}
	return name
}
func GetUaDefaultE164Number(k interface{}) (n string) {
	var number string
	if value, ok := uamap.Load(k); ok {
		number = value.(rtua).defaulte164number
	} else {
		number = ``
	}
	return number
}

func SetGatewayState(e *eventsocket.Event) {
	gw := rtgw{
		coreuuid:     e.Get("Core-Uuid"),
		hostipv4:     e.Get("Freeswitch-Ipv4"),
		hostipv6:     e.Get("Freeswitch-Ipv6"),
		registerip:   e.Get("Register-Network-Ip"),
		registerport: e.Get("Register-Network-Port"),
		gateway:      e.Get("Gateway"),
		pingstatus:   e.Get("Ping-Status"),
		state:        e.Get("State"),
	}
	if _, isloaded := gwmap.LoadOrStore(gw.gateway, &gw); isloaded {
		gwmap.Delete(gw.gateway)
		gwmap.Store(gw.gateway, &gw)
		log.Printf("SetGatewayState:%v\n", gw)
	} else {
		fmt.Println(gw)
	}
}

// func mapClean(m sync.Map) {
// 	m.Range(func(k, v interface{}) bool { m.Delete(k); return true })
// }
