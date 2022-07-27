// switch define
// pre_load_switch.conf.xml		`switch.conf`/autoload_configs/pre_load_switch.conf.xml
// switch.conf.xml				`switch.conf`/autoload_configs/switch.conf.xml
// post_load_switch.conf.xml	`switch.conf`/autoload_configs/post_load_switch.conf.xml
//
// request:
// hostname=D1130&section=configuration&tag_name=configuration&key_name=name&key_value=switch.conf
// response:
// <document type="freeswitch/xml">
//   <section name="configuration">
//      <!-- switch.conf.xml -->
//   </section>
// </document>

package switch_main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func Read(c *gin.Context) (string, error) {
	filename := c.PostForm(`key_value`)
	file := fmt.Sprintf("%s/autoload_configs/%s%s", viper.GetString(`switch.conf`), filename, `.xml`)
	content, err := os.ReadFile(file)
	return string(content), err
}

func Default(filename string) (string, error) {
	var content string
	if strings.EqualFold(filename, `pre_load_switch.conf`) {
		content = PRE_LOAD_SWITCH_CONF
	}
	if strings.EqualFold(filename, `post_load_switch.conf`) {
		content = POST_LOAD_SWITCH_CONF
	}
	if strings.EqualFold(filename, `switch.conf`) {
		content = SWITCH_CONF
	}
	return content, nil
}

func Build(c *gin.Context, content string) (string, error) {
	////////NOTICE!!! change switch.conf param, but no effact!!!, set odbc-dsn before switch boot(fs config fsconfig --init)///////
	//<param name="sessions-per-second" value="30"/>
	//<param name="max-sessions" value="1000"/>
	//newcontent := strings.ReplaceAll(content, `<param name="max-sessions" value="1000"/>`, `<param name="max-sessions" value="999"/>`)
	//return newcontent, nil
	return ``, errors.New("switch_main.Build() nothing")
}
