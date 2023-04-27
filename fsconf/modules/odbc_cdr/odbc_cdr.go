// mod_odbc_cdr
// apt-get install freeswitch-mod-odbc-cdr
// default install without configuration file "odbc_cdr.conf.xml", so read return error.

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
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func init() {}

func Read(c *gin.Context) (string, error) {
	// file := fmt.Sprintf("%s/autoload_configs/%s", viper.GetString(`switch.conf`), MOD_CONF_NAME)
	// content, err := os.ReadFile(file)
	// return string(content), err
	return "", errors.New("odbc_cdr.Read(c) nothing, default odbc_cdr.conf.xml not found")
}

func Default() (string, error) {
	var err error
	var tables string

	alegname := viper.GetString(`switch.cdr.a-leg`)
	blegname := viper.GetString(`switch.cdr.b-leg`)
	bothname := viper.GetString(`switch.cdr.both`)
	if len(alegname) > 0 {
		table := fmt.Sprintf(ODBC_CDR_CONF_XML_TABLE, alegname, `a-leg`)
		tables = fmt.Sprintf("%s\n%s", tables, table)
	}
	if len(blegname) > 0 {
		table := fmt.Sprintf(ODBC_CDR_CONF_XML_TABLE, blegname, `b-leg`)
		tables = fmt.Sprintf("%s\n%s", tables, table)
	}
	if len(bothname) > 0 {
		table := fmt.Sprintf(ODBC_CDR_CONF_XML_TABLE, bothname, `both`)
		tables = fmt.Sprintf("%s\n%s", tables, table)
	}
	if len(tables) <= 0 {
		err = errors.New(`odbc_cdr.conf no table defined`)
	}
	return fmt.Sprintf(ODBC_CDR_CONF_XML, tables), err
}

func Build(c *gin.Context, content string) (string, error) {
	//<param name="odbc-dsn" value="$${pg_handle}"/>
	newodbcdsn := `<param name="odbc-dsn" value="$${pg_handle}"/>`
	newcontent := strings.ReplaceAll(content, ODBC_DSN, newodbcdsn)
	return newcontent, nil
}
