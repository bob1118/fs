// mod_odbc_cdr
// apt-get install freeswitch-mod-odbc-cdr
// default install without configuration file "odbc_cdr.conf.xml"

// request:
// hostname=D1130&section=configuration&tag_name=configuration&key_name=name&key_value=odbc_cdr.conf
// response:
// <document type="freeswitch/xml">
//   <section name="configuration">
//     <configuration name="odbc_cdr.conf" description="ODBC CDR Configuration">
//       <settings>
//          <!--ADD your parameters here-->
//       </settings>
//     </configuration>
//   </section>
// </document>

package odbc_cdr

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
)

func init() {}

func Default() (string, error) { return ODBC_CDR_CONF_XML, nil }

func Read(c *gin.Context) (string, error) {
	//file := fmt.Sprintf("%s/autoload_configs/odbc_cdr.conf.xml", viper.GetString(`switch.conf`))
	//content, err := os.ReadFile(file)
	//return string(content), err
	return "", errors.New("odbc_cdr.Read(c) nothing")
}

func Build(c *gin.Context, content string) (string, error) {
	//<param name="odbc-dsn" value="$${pg_handle}"/>
	newOdbcdsn := `<param name="odbc-dsn" value="$${pg_handle}"/>`
	newcontent := strings.ReplaceAll(content, ODBC_DSN, newOdbcdsn)
	return newcontent, nil
}
