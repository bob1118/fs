// switch define
// pre_load_switch.conf.xml	`switch.conf`/autoload_configs/pre_load_switch.conf.xml
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

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func Read(c *gin.Context) (string, error) {
	filename := c.PostForm(`key_value`)
	file := fmt.Sprintf("%s/autoload_configs/%s%s", viper.GetString(`switch.conf`), filename, `.xml`)
	content, err := os.ReadFile(file)
	return string(content), err
}

func Default() (string, error) { return MOD_CONF_XML, nil }

func Build(c *gin.Context, content string) (string, error) {
	// //<param name="odbc-dsn" value="$${pg_handle}"/>
	// old := `<!-- <param name="core-db-dsn" value="dsn:username:password" /> -->`
	// newcontent := strings.ReplaceAll(content, old, ODBC_DSN)
	//<param name="odbc-dsn" value="$${pg_handle}"/>
	return ``, errors.New("switch_main.Build() nothing")
}
