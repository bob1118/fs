package fsapi

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

/////////////////////////////////// request when switch boot ...////////////////////////////////////
// http://localhost/fsapi
//
// HEADER:
// Content-Length [97]
// Content-Type [application/x-www-form-urlencoded]
// User-Agent [freeswitch-xml/1.0]
// Accept [*/*]
//
// BODY:
// hostname=D1130&section=configuration&tag_name=configuration&key_name=name&key_value=odbc_cdr.conf
//
//////////////////////////////////////request body list/////////////////////////////////////////////
//1 boot order.
// hostname=D1130&section=configuration&tag_name=configuration&key_name=name&key_value=odbc_cdr.conf
// hostname=D1130&section=configuration&tag_name=configuration&key_name=name&key_value=sofia.conf&Event-Name=REQUEST_PARAMS&Core-UUID=297f80ae-fee7-4a80-87a4-32625cafb18d&FreeSWITCH-Hostname=D1130&FreeSWITCH-Switchname=D1130&FreeSWITCH-IPv4=10.10.10.10&FreeSWITCH-IPv6=%3A%3A1&Event-Date-Local=2022-07-06%2018%3A25%3A12&Event-Date-GMT=Wed,%2006%20Jul%202022%2010%3A25%3A12%20GMT&Event-Date-Timestamp=1657103112943844&Event-Calling-File=sofia.c&Event-Calling-Function=config_sofia&Event-Calling-Line-Number=4489&Event-Sequence=30
// hostname=D1130&section=configuration&tag_name=configuration&key_name=name&key_value=sofia.conf&Event-Name=REQUEST_PARAMS&Core-UUID=297f80ae-fee7-4a80-87a4-32625cafb18d&FreeSWITCH-Hostname=D1130&FreeSWITCH-Switchname=D1130&FreeSWITCH-IPv4=10.10.10.10&FreeSWITCH-IPv6=%3A%3A1&Event-Date-Local=2022-07-06%2018%3A25%3A12&Event-Date-GMT=Wed,%2006%20Jul%202022%2010%3A25%3A12%20GMT&Event-Date-Timestamp=1657103112984742&Event-Calling-File=sofia.c&Event-Calling-Function=launch_sofia_worker_thread&Event-Calling-Line-Number=3049&Event-Sequence=39&profile=external
// hostname=D1130&section=configuration&tag_name=configuration&key_name=name&key_value=sofia.conf&Event-Name=REQUEST_PARAMS&Core-UUID=297f80ae-fee7-4a80-87a4-32625cafb18d&FreeSWITCH-Hostname=D1130&FreeSWITCH-Switchname=D1130&FreeSWITCH-IPv4=10.10.10.10&FreeSWITCH-IPv6=%3A%3A1&Event-Date-Local=2022-07-06%2018%3A25%3A12&Event-Date-GMT=Wed,%2006%20Jul%202022%2010%3A25%3A12%20GMT&Event-Date-Timestamp=1657103112984816&Event-Calling-File=sofia.c&Event-Calling-Function=launch_sofia_worker_thread&Event-Calling-Line-Number=3049&Event-Sequence=40&profile=internal-ipv6
// hostname=D1130&section=configuration&tag_name=configuration&key_name=name&key_value=sofia.conf&Event-Name=REQUEST_PARAMS&Core-UUID=297f80ae-fee7-4a80-87a4-32625cafb18d&FreeSWITCH-Hostname=D1130&FreeSWITCH-Switchname=D1130&FreeSWITCH-IPv4=10.10.10.10&FreeSWITCH-IPv6=%3A%3A1&Event-Date-Local=2022-07-06%2018%3A25%3A12&Event-Date-GMT=Wed,%2006%20Jul%202022%2010%3A25%3A12%20GMT&Event-Date-Timestamp=1657103112984816&Event-Calling-File=sofia.c&Event-Calling-Function=launch_sofia_worker_thread&Event-Calling-Line-Number=3049&Event-Sequence=41&profile=external-ipv6
// hostname=D1130&section=configuration&tag_name=configuration&key_name=name&key_value=sofia.conf&Event-Name=REQUEST_PARAMS&Core-UUID=297f80ae-fee7-4a80-87a4-32625cafb18d&FreeSWITCH-Hostname=D1130&FreeSWITCH-Switchname=D1130&FreeSWITCH-IPv4=10.10.10.10&FreeSWITCH-IPv6=%3A%3A1&Event-Date-Local=2022-07-06%2018%3A25%3A12&Event-Date-GMT=Wed,%2006%20Jul%202022%2010%3A25%3A12%20GMT&Event-Date-Timestamp=1657103112984846&Event-Calling-File=sofia.c&Event-Calling-Function=launch_sofia_worker_thread&Event-Calling-Line-Number=3049&Event-Sequence=42&profile=internal
// hostname=D1130&section=directory&tag_name=&key_name=&key_value=&Event-Name=REQUEST_PARAMS&Core-UUID=297f80ae-fee7-4a80-87a4-32625cafb18d&FreeSWITCH-Hostname=D1130&FreeSWITCH-Switchname=D1130&FreeSWITCH-IPv4=10.10.10.10&FreeSWITCH-IPv6=%3A%3A1&Event-Date-Local=2022-07-06%2018%3A25%3A13&Event-Date-GMT=Wed,%2006%20Jul%202022%2010%3A25%3A13%20GMT&Event-Date-Timestamp=1657103113026106&Event-Calling-File=sofia.c&Event-Calling-Function=launch_sofia_worker_thread&Event-Calling-Line-Number=3067&Event-Sequence=45&purpose=gateways&profile=internal
// hostname=D1130&section=directory&tag_name=&key_name=&key_value=&Event-Name=REQUEST_PARAMS&Core-UUID=297f80ae-fee7-4a80-87a4-32625cafb18d&FreeSWITCH-Hostname=D1130&FreeSWITCH-Switchname=D1130&FreeSWITCH-IPv4=10.10.10.10&FreeSWITCH-IPv6=%3A%3A1&Event-Date-Local=2022-07-06%2018%3A25%3A13&Event-Date-GMT=Wed,%2006%20Jul%202022%2010%3A25%3A13%20GMT&Event-Date-Timestamp=1657103113025527&Event-Calling-File=sofia.c&Event-Calling-Function=launch_sofia_worker_thread&Event-Calling-Line-Number=3067&Event-Sequence=44&purpose=gateways&profile=external
// hostname=D1130&section=configuration&tag_name=configuration&key_name=name&key_value=loopback.conf
// hostname=D1130&section=configuration&tag_name=configuration&key_name=name&key_value=verto.conf
// hostname=D1130&section=configuration&tag_name=configuration&key_name=name&key_value=conference.conf&Event-Name=COMMAND&Core-UUID=297f80ae-fee7-4a80-87a4-32625cafb18d&FreeSWITCH-Hostname=D1130&FreeSWITCH-Switchname=D1130&FreeSWITCH-IPv4=10.10.10.10&FreeSWITCH-IPv6=%3A%3A1&Event-Date-Local=2022-07-06%2018%3A25%3A14&Event-Date-GMT=Wed,%2006%20Jul%202022%2010%3A25%3A14%20GMT&Event-Date-Timestamp=1657103114557967&Event-Calling-File=mod_conference.c&Event-Calling-Function=send_presence&Event-Calling-Line-Number=3897&Event-Sequence=232&presence=true
// hostname=D1130&section=configuration&tag_name=configuration&key_name=name&key_value=db.conf
// hostname=D1130&section=configuration&tag_name=configuration&key_name=name&key_value=fifo.conf
// hostname=D1130&section=configuration&tag_name=configuration&key_name=name&key_value=hash.conf
// hostname=D1130&section=configuration&tag_name=configuration&key_name=name&key_value=voicemail.conf
// hostname=D1130&section=configuration&tag_name=configuration&key_name=name&key_value=voicemail.conf
// hostname=D1130&section=configuration&tag_name=configuration&key_name=name&key_value=httapi.conf
// hostname=D1130&section=configuration&tag_name=configuration&key_name=name&key_value=spandsp.conf
// hostname=D1130&section=configuration&tag_name=configuration&key_name=name&key_value=amr.conf
// hostname=D1130&section=configuration&tag_name=configuration&key_name=name&key_value=opus.conf
// hostname=D1130&section=configuration&tag_name=configuration&key_name=name&key_value=avformat.conf
// hostname=D1130&section=configuration&tag_name=configuration&key_name=name&key_value=avcodec.conf
// hostname=D1130&section=configuration&tag_name=configuration&key_name=name&key_value=sndfile.conf
// hostname=D1130&section=configuration&tag_name=configuration&key_name=name&key_value=local_stream.conf
// hostname=D1130&section=configuration&tag_name=configuration&key_name=name&key_value=lua.conf
// hostname=D1130&section=configuration&tag_name=configuration&key_name=name&key_value=post_load_modules.conf
// hostname=D1130&section=configuration&tag_name=configuration&key_name=name&key_value=event_socket.conf
// hostname=D1130&section=configuration&tag_name=configuration&key_name=name&key_value=acl.conf
// hostname=D1130&section=directory&tag_name=domain&key_name=name&key_value=10.10.10.10&Event-Name=GENERAL&Core-UUID=297f80ae-fee7-4a80-87a4-32625cafb18d&FreeSWITCH-Hostname=D1130&FreeSWITCH-Switchname=D1130&FreeSWITCH-IPv4=10.10.10.10&FreeSWITCH-IPv6=%3A%3A1&Event-Date-Local=2022-07-06%2018%3A25%3A15&Event-Date-GMT=Wed,%2006%20Jul%202022%2010%3A25%3A15%20GMT&Event-Date-Timestamp=1657103115428365&Event-Calling-File=switch_core.c&Event-Calling-Function=switch_load_network_lists&Event-Calling-Line-Number=1637&Event-Sequence=570&domain=10.10.10.10&purpose=network-list
// hostname=D1130&section=configuration&tag_name=configuration&key_name=name&key_value=post_load_switch.conf
// hostname=D1130&section=configuration&tag_name=configuration&key_name=name&key_value=switch.conf
// 2 reload mod_xxx
// hostname=D1130&section=configuration&tag_name=configuration&key_name=name&key_value=spandsp.conf
// hostname=D1130&section=configuration&tag_name=configuration&key_name=name&key_value=loopback.conf
// hostname=D1130&section=configuration&tag_name=configuration&key_name=name&key_value=enum.conf
// hostname=D1130&section=configuration&tag_name=configuration&key_name=name&key_value=timezones.conf
// ...
// xxx.conf request
// ...
// 3 reloadxml
// hostname=D1130&section=configuration&tag_name=configuration&key_name=name&key_value=spandsp.conf
// hostname=D1130&section=configuration&tag_name=configuration&key_name=name&key_value=loopback.conf
// hostname=D1130&section=configuration&tag_name=configuration&key_name=name&key_value=enum.conf
// hostname=D1130&section=configuration&tag_name=configuration&key_name=name&key_value=timezones.conf

// PostFromXmlCurl
func PostFromXmlCurl(c *gin.Context) {
	//for debug, notice!!! c.Request.Body can readed once only.
	if false {
		body, _ := io.ReadAll(c.Request.Body)
		fmt.Println(string(body))
	}

	var responseBody string
	value := c.PostForm("section")
	switch value {
	case "configuration":
		responseBody = doConfiguration(c)
	case "dialplan":
		responseBody = doDialplan(c)
	case "directory":
		responseBody = doDirectory(c)
	case "phrases":
		responseBody = doPhrases(c)
	default:
		responseBody = `default bad request was ignored !!!`
	}
	//fmt.Println(responseBody)
	c.String(http.StatusOK, responseBody)
}
