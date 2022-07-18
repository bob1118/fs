// mod_av
// `switch.conf`/autoload_configs/av.conf.xml

// request:
// hostname=D1130&section=configuration&tag_name=configuration&key_name=name&key_value=av.conf
// response:
// <document type="freeswitch/xml">
//   <section name="configuration">
//      <!-- av.conf.xml -->
//   </section>
// </document>

package av

import (
	"errors"
	"fmt"
	"os"

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
	return ``, errors.New(`av.Build() nothing`)
}
