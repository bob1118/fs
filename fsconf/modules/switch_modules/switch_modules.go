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
	if strings.EqualFold(filename, `pre_load_modules.conf`) { content = PRE_LOAD_MODULES_CONF }
	if strings.EqualFold(filename, `post_load_modules.conf`) { content = POST_LOAD_MODULES_CONF }
	return content, nil
}

func Build(c *gin.Context, content string) (string, error) {
	return ``, errors.New("switch_modules Build() nothing")
}
