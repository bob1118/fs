// mod_fifo
// `switch.conf`/autoload_configs/fifo.conf.xml

// request:
// hostname=D1130&section=configuration&tag_name=configuration&key_name=name&key_value=fifo.conf
// response:
// <document type="freeswitch/xml">
//   <section name="configuration">
//      <!-- fifo.conf.xml -->
//   </section>
// </document>

package fifo

import (
	"fmt"
	"os"
	"strings"

	"github.com/bob1118/fs/db"
	"github.com/bob1118/fs/utils"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func Read(c *gin.Context) (string, error) {
	file := fmt.Sprintf("%s/autoload_configs/%s", viper.GetString(`switch.conf`), MOD_CONF_NAME)
	content, err := os.ReadFile(file)
	return string(content), err
}

func Default() (string, error) { return FIFO_CONF_XML, nil }

func Build(c *gin.Context, content string) (string, error) {
	var err error
	var outbound_strategy string
	//<param name="outbound-strategy" value=""/>
	//<param name="odbc-dsn" value="$${pg_handle}"/>
	old := `<settings>`
	if utils.IsEqual(`enterprise`, viper.GetString(`switch.fifo.outbound-strategy`)) {
		outbound_strategy = OUTBOUND_STRATEGY_ENTERPRISE
	} else {
		outbound_strategy = OUTBOUND_STRATEGY_RINGALL
	}
	new := fmt.Sprintf("%s\n%s\n%s", old, outbound_strategy, ODBC_DSN)
	newcontent := strings.ReplaceAll(content, old, new)

	//maybe buildfifos()
	if fifos, e := buildFifos(); e != nil {
		err = e
	} else {
		newcontent = strings.ReplaceAll(newcontent, DEFAULT_FIFOS, fifos)
	}
	return newcontent, err
}

func buildFifos() (string, error) {
	var err error
	var myFifos string
	if fifos, e := db.SelectFifos(); err != nil {
		err = e
		fmt.Println(err)
	} else {
		for _, fifo := range fifos {
			myfifo := fmt.Sprintf(FIFO, fifo.Fname)
			myFifos = fmt.Sprintf("%s\n%s", myFifos, myfifo)
		}
	}
	return myFifos, err
}
