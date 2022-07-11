// mod_sofia
// switch.conf/autoload_configs/sofia.conf.xml

// request:
// hostname=D1130&section=configuration&tag_name=configuration&key_name=name&key_value=sofia.conf&Event-Name=REQUEST_PARAMS&Core-UUID=297f80ae-fee7-4a80-87a4-32625cafb18d&FreeSWITCH-Hostname=D1130&FreeSWITCH-Switchname=D1130&FreeSWITCH-IPv4=10.10.10.10&FreeSWITCH-IPv6=%3A%3A1&Event-Date-Local=2022-07-06%2018%3A25%3A12&Event-Date-GMT=Wed,%2006%20Jul%202022%2010%3A25%3A12%20GMT&Event-Date-Timestamp=1657103112943844&Event-Calling-File=sofia.c&Event-Calling-Function=config_sofia&Event-Calling-Line-Number=4489&Event-Sequence=30
// response:
// <document type="freeswitch/xml">
//   <section name="configuration">
//      <!-- sofia.conf.xml -->
//   </section>
// </document>

// sofia.conf.xml
// <configuration name="sofia.conf" description="sofia Endpoint">
// <global_settings>
//   <param name="log-level" value="0"/>
//   <!-- <param name="abort-on-empty-external-ip" value="true"/> -->
//   <!-- <param name="auto-restart" value="false"/> -->
//   <param name="debug-presence" value="0"/>
//   <!-- <param name="capture-server" value="udp:homer.domain.com:5060"/> -->
//   <!--
//   the new format for HEPv2/v3 and capture ID
//   protocol:host:port;hep=2;capture_id=200;
//   -->
//   <!-- <param name="capture-server" value="udp:homer.domain.com:5060;hep=3;capture_id=100"/> -->
// </global_settings>
// <!--
// 	The rabbit hole goes deep.  This includes all the
// 	profiles in the sip_profiles directory that is up
// 	one level from this directory.
// -->
// <profiles>
//   <X-PRE-PROCESS cmd="include" data="../sip_profiles/*.xml"/>
// </profiles>
// </configuration>

// <X-PRE-PROCESS cmd="include" data="../sip_profiles/*.xml"/>
// <profile name="external">
// ...
// </profile>
// <profile name="external-ipv6">
// ...
// </profile>
// <profile name="internal">
// ...
// </profile>
// <profile name="internal-ipv6">
// ...
// </profile>

package sofia

import (
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func init() {}

func Default() (string, error) { return SOFIA_CONF_XML, nil }

func Read(c *gin.Context) (string, error) {
	var file string
	function := c.PostForm(`Event-Calling-Function`)
	profile := c.PostForm(`profile`)
	switch function {
	case "config_sofia":
		switch profile {
		case "": //sofia.conf.xml
			file = fmt.Sprintf("%s/autoload_configs/sofia.conf.xml", viper.GetString(`switch.conf`))
		case "internal": //sofia profile internal xxx
		case "internal-ipv6": //sofia profile internal-ipv6 xxx
		case "external": //sofia profile external xxx
		case "external-ipv6": //sofia profile external-ipv6 xxx
		}
	case "launch_sofia_worker_thread":
		switch profile {
		case "internal":
			file = fmt.Sprintf("%s/sip_profiles/internal.xml", viper.GetString(`switch.conf`))
		case "internal-ipv6":
			file = fmt.Sprintf("%s/sip_profiles/internal-ipv6.xml", viper.GetString(`switch.conf`))
		case "external":
			file = fmt.Sprintf("%s/sip_profiles/external.xml", viper.GetString(`switch.conf`))
		case "external-ipv6":
			file = fmt.Sprintf("%s/sip_profiles/external-ipv6.xml", viper.GetString(`switch.conf`))
		}
	}
	content, err := os.ReadFile(file)
	return string(content), err
}

func Build(c *gin.Context, content string) (string, error) {
	//<X-PRE-PROCESS cmd="include" data="../sip_profiles/*.xml"/>
	newprofiles := ``
	oldprofiles := `    <X-PRE-PROCESS cmd="include" data="../sip_profiles/*.xml"/>`
	newcontent := strings.ReplaceAll(content, oldprofiles, newprofiles)

	//external
	//
	return newcontent, nil
}
