// mod_event_socket
// `switch.conf`/autoload_configs/event_socket.conf.xml

// request:
// hostname=D1130&section=configuration&tag_name=configuration&key_name=name&key_value=event_socket.conf
// response:
// <document type="freeswitch/xml">
//   <section name="configuration">
//      <!-- event_socket.conf.xml -->
//   </section>
// </document>

package event_socket

import (
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func Read(c *gin.Context) (string, error) {
	file := fmt.Sprintf("%s/autoload_configs/%s", viper.GetString(`switch.conf`), MOD_CONF_NAME)
	content, err := os.ReadFile(file)
	return string(content), err
}

func Default() (string, error) { return MOD_CONF_XML, nil }

func Build(c *gin.Context, content string) (string, error) {
	var newcontent = content
	ipaddr := viper.GetString(`switch.eventsocket.ipaddr`)
	port := viper.GetString(`switch.eventsocket.port`)
	password := viper.GetString(`switch.eventsocket.password`)
	//<param name="listen-ip" value="::"/>
	if len(ipaddr) > 0 {
		listenip := fmt.Sprintf(`<param name="listen-ip" value="%s"/>`, ipaddr)
		newcontent = strings.ReplaceAll(newcontent, `<param name="listen-ip" value="::"/>`, listenip)
	}
	//<param name="listen-port" value="8021"/>
	if len(port) > 0 {
		listenport := fmt.Sprintf(`<param name="listen-port" value="%s"/>`, port)
		newcontent = strings.ReplaceAll(newcontent, `<param name="listen-port" value="8021"/>`, listenport)
	}
	//<param name="password" value="ClueCon"/>
	if len(password) > 0 {
		authpassword := fmt.Sprintf(`<param name="password" value="%s"/>`, password)
		newcontent = strings.ReplaceAll(newcontent, `<param name="password" value="ClueCon"/>`, authpassword)
	}
	return newcontent, nil
}
