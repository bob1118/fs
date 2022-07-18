// switch modules define
// pre_load_modules.conf.xml	`switch.conf`/autoload_configs/pre_load_modules.conf.xml
// modules.conf.xml				`switch.conf`/autoload_configs/modules.conf.xml
// post_load_modules.conf.xml	`switch.conf`/autoload_configs/post_load_modules.conf.xml
//
// request:
// hostname=D1130&section=configuration&tag_name=configuration&key_name=name&key_value=post_load_modules.conf
// response:
// <document type="freeswitch/xml">
//   <section name="configuration">
//      <!-- post_load_modules.conf.xml -->
//   </section>
// </document>

package switch_modules

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

func Default() (string, error) { return MODULES_CONF_XML, nil }

func Build(c *gin.Context, content string) (string, error) {
	return ``, errors.New("switch_modules Build() nothing")
}
