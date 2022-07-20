// mod_voicemail
// `switch.conf`/autoload_configs/voicemail.conf.xml

// request:
// hostname=D1130&section=configuration&tag_name=configuration&key_name=name&key_value=voicemail.conf
// response:
// <document type="freeswitch/xml">
//   <section name="configuration">
//      <!-- voicemail.conf.xml -->
//   </section>
// </document>

package voicemail

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
	//<param name="odbc-dsn" value="$${pg_handle}"/>
	old := `<!--<param name="odbc-dsn" value="dsn:user:pass"/>-->`
	newcontent := strings.ReplaceAll(content, old, ODBC_DSN)
	return newcontent, nil
}
