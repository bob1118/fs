// mod_sofia
// `switch.conf`/autoload_configs/sofia.conf.xml
// `switch.conf`/sip_profiles/*.xml

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

	"github.com/bob1118/fs/db"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func init() {}

func Default() (string, error) { return SOFIA_CONF_XML, nil }

func Read(c *gin.Context) (string, error) {
	var err error
	var file string
	var content string
	dir := viper.GetString(`switch.conf`)
	profile := c.PostForm(`profile`)
	//reconfig := c.PostForm(`reconfig`)
	//function := c.PostForm(`Event-Calling-Function`)

	// switch function {
	// case "config_sofia":
	// 	switch profile {
	// 	case "": //sofia.conf.xml
	// 		file = fmt.Sprintf("%s/autoload_configs/sofia.conf.xml", dir)
	// 	case "internal": //sofia profile internal xxx
	// 		file = fmt.Sprintf("%s/sip_profiles/internal.xml", dir)
	// 	case "internal-ipv6": //sofia profile internal-ipv6 xxx
	// 		file = fmt.Sprintf("%s/sip_profiles/internal-ipv6.xml", dir)
	// 	case "external": //sofia profile external xxx
	// 		file = fmt.Sprintf("%s/sip_profiles/external.xml", dir)
	// 	case "external-ipv6": //sofia profile external-ipv6 xxx
	// 		file = fmt.Sprintf("%s/sip_profiles/external-ipv6.xml", dir)
	// 	}
	// case "launch_sofia_worker_thread":
	// 	switch profile {
	// 	case "internal":
	// 		file = fmt.Sprintf("%s/sip_profiles/internal.xml", dir)
	// 	case "internal-ipv6":
	// 		file = fmt.Sprintf("%s/sip_profiles/internal-ipv6.xml", dir)
	// 	case "external":
	// 		file = fmt.Sprintf("%s/sip_profiles/external.xml", dir)
	// 	case "external-ipv6":
	// 		file = fmt.Sprintf("%s/sip_profiles/external-ipv6.xml", dir)
	// 	}

	switch profile {
	case "": //sofia.conf.xml
		file = fmt.Sprintf("%s/autoload_configs/%s", dir, MOD_CONF_NAME)
		if data, e := os.ReadFile(file); e != nil {
			err = e
			fmt.Println(e)
		} else {
			content = string(data)
		}
	case "internal": //sofia.conf.xml with profile internal
		file = fmt.Sprintf("%s/sip_profiles/internal.xml", dir)
		if data, e := os.ReadFile(file); e != nil {
			err = e
			fmt.Println(e)
		} else {
			content = fmt.Sprintf(SOFIA_CONF_XML_WITH_PROFILE, string(data))
		}
	case "internal-ipv6": //sofia.conf.xml with profile internal-ipv6
		file = fmt.Sprintf("%s/sip_profiles/internal-ipv6.xml", dir)
		if data, e := os.ReadFile(file); e != nil {
			err = e
			fmt.Println(e)
		} else {
			content = fmt.Sprintf(SOFIA_CONF_XML_WITH_PROFILE, string(data))
		}
	case "external": //sofia.conf.xml with profile external
		file = fmt.Sprintf("%s/sip_profiles/external.xml", dir)
		if data, e := os.ReadFile(file); e != nil {
			err = e
			fmt.Println(e)
		} else {
			content = fmt.Sprintf(SOFIA_CONF_XML_WITH_PROFILE, string(data))
		}
	case "external-ipv6": //sofia.conf.xml with profile external-ipv6
		file = fmt.Sprintf("%s/sip_profiles/external-ipv6.xml", dir)
		if data, e := os.ReadFile(file); e != nil {
			err = e
			fmt.Println(e)
		} else {
			content = fmt.Sprintf(SOFIA_CONF_XML_WITH_PROFILE, string(data))
		}
	}
	return content, err
}

func Build(c *gin.Context, content string) (string, error) {
	///////notice!!! param odbc-dsn not affect, set odbc-dsn before switch boot(fs config fsconfig --init)//////
	var err error
	var old, new, newcontent string
	profile := c.PostForm(`profile`)
	switch profile {
	case "": //sofia.conf.xml
		//<X-PRE-PROCESS cmd="include" data="../sip_profiles/*.xml"/>
		old = `<X-PRE-PROCESS cmd="include" data="../sip_profiles/*.xml"/>`
		new = `<X-PRE-PROCESS cmd="include" data="./sip_profiles/*.xml"/>`
		newcontent = strings.ReplaceAll(content, old, new)
	case "internal", "internal-ipv6":
		//<param name="force-register-domain" value="$${domain}"/>
		old = `<param name="force-register-domain" value="$${domain}"/>`
		new = `<!--<param name="force-register-domain" value="$${domain}"/>-->`
		newcontent = strings.ReplaceAll(newcontent, old, new)
		//<param name="force-subscription-domain" value="$${domain}"/>
		old = `<param name="force-subscription-domain" value="$${domain}"/>`
		new = `<!--<param name="force-subscription-domain" value="$${domain}"/>-->`
		newcontent = strings.ReplaceAll(newcontent, old, new)
		//<param name="force-register-db-domain" value="$${domain}"/>
		old = `<param name="force-register-db-domain" value="$${domain}"/>`
		new = `<!--<param name="force-register-db-domain" value="$${domain}"/>-->`
		newcontent = strings.ReplaceAll(newcontent, old, new)
	case "external", "external-ipv6":
		//<X-PRE-PROCESS cmd="include" data="external/*.xml"/>
		old = fmt.Sprintf(`<X-PRE-PROCESS cmd="include" data="%s/*.xml"/>`, profile)
		//new = fmt.Sprintf(`<X-PRE-PROCESS cmd="include" data="./sip_profiles/%s/*.xml"/>`, profile)
		if new, err = getProfileGateways(profile); err == nil {
			newcontent = strings.ReplaceAll(content, old, new)
		}
	}
	return newcontent, err
}

func getProfileGateways(s string) (string, error) {
	var e error
	var gatewaysConf string
	profile := s
	if len(profile) > 0 {
		condition := fmt.Sprintf(`profile_name='%s'`, profile)
		if gateways, err := db.SelectGatewaysWithCondition(condition); err != nil {
			e = err
			fmt.Println(err)
		} else {
			for _, gateway := range gateways {
				gatewayConf := fmt.Sprintf(SOFIA_PROFILE_GATEWAY,
					gateway.Gname,
					gateway.Gusername,
					gateway.Grealm,
					gateway.Gfromuser,
					gateway.Gfromdomain,
					gateway.Gpassword,
					gateway.Gextension,
					gateway.Gproxy,
					gateway.Gregisterproxy,
					gateway.Gexpire,
					gateway.Gregister,
					gateway.Gcalleridinfrom,
					gateway.Gextensionincontact,
					gateway.Goptionping)
				gatewaysConf = fmt.Sprintf("%s%s", gatewaysConf, gatewayConf)
			}
		}
	}
	return gatewaysConf, e
}
