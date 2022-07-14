// mod_db
// `switch.conf`/autoload_configs/db.conf.xml

// request:
// hostname=D1130&section=configuration&tag_name=configuration&key_name=name&key_value=db.conf
// response:
// <document type="freeswitch/xml">
//   <section name="configuration">
//      <!-- db.conf.xml -->
//   </section>
// </document>

package switch_db

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

func Default() (string, error) { return DB_CONF_XML, nil }

func Build(c *gin.Context, content string) (string, error) {
	//<param name="odbc-dsn" value="$${pg_handle}"/>
	newOdbcdsn := `<param name="odbc-dsn" value="$${pg_handle}"/>`
	newcontent := strings.ReplaceAll(content, ODBC_DSN, newOdbcdsn)
	return newcontent, nil
}
