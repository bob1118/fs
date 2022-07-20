// switch configuration modules
// p.Dir()+`autoload_configs`+`xxx.conf.xml`

// request:
// hostname=D1130&section=configuration&tag_name=configuration&key_name=name&key_value=xxx.conf ...
// response:
// <document type="freeswitch/xml">
//   <section name="configuration">
//      <!-- xxx.conf.xml -->
//   </section>
// </document>

package modules

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/bob1118/fs/db"
	"github.com/bob1118/fs/fsconf/modules/acl"
	"github.com/bob1118/fs/fsconf/modules/av"
	"github.com/bob1118/fs/fsconf/modules/event_socket"
	"github.com/bob1118/fs/fsconf/modules/fifo"
	"github.com/bob1118/fs/fsconf/modules/odbc_cdr"
	"github.com/bob1118/fs/fsconf/modules/sofia"
	"github.com/bob1118/fs/fsconf/modules/switch_db"
	"github.com/bob1118/fs/fsconf/modules/switch_main"
	"github.com/bob1118/fs/fsconf/modules/switch_modules"
	"github.com/bob1118/fs/fsconf/modules/voicemail"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func GetConfiguration(c *gin.Context) (string, error) {
	var err error
	var content, newcontent string
	if content, err = readConfFromDatabase(c); err != nil {
		if content, err = readConfFromFile(c); err != nil {
			if content, err = constConfiguration(c); err != nil {
				log.Println(err)
				return "", err
			} else {
				writeConfToFile(c, content)
			}
		}
		if newcontent, err = buildConf(c, content); err != nil {
			writeConfToDatabase(c, content, ``)
		} else {
			writeConfToDatabase(c, content, newcontent)
			content = newcontent
		}
	}
	//for debug
	if true {
		filename := c.PostForm(`key_value`)
		function := c.PostForm(`Event-Calling-Function`)
		profile := c.PostForm(`profile`)
		log.Println("Request:", filename, function, profile)
		log.Println("Response:", content)
	}
	return content, err
}

func readConfFromDatabase(c *gin.Context) (string, error) {
	var err error
	var content string
	filename := c.PostForm(`key_value`)
	function := c.PostForm(`Event-Calling-Function`)
	profile := c.PostForm(`profile`)
	if conf, e := db.GetConfsConf(filename, function, profile); e != nil {
		err = e
	} else {
		if len(conf.Ccontent) > 0 {
			if len(conf.Cnewcontent) > 0 {
				content = conf.Cnewcontent
			} else {
				content = conf.Ccontent
			}
		} else {
			err = errors.New("db confs conf.Ccontent null")
		}
	}
	return content, err
}

func readConfFromFile(c *gin.Context) (string, error) {
	var err error
	var content string
	filename := c.PostForm(`key_value`)
	switch filename {
	case "odbc_cdr.conf":
		content, err = odbc_cdr.Read(c)
	case "sofia.conf":
		content, err = sofia.Read(c)
	case "db.conf":
		content, err = switch_db.Read(c)
	case "fifo.conf":
		content, err = fifo.Read(c)
	case "voicemail.conf":
		content, err = voicemail.Read(c)
	case "avformat.conf", "avcodec.conf":
		content, err = av.Read(c)
	case "pre_load_modules.conf", "post_load_modules.conf":
		content, err = switch_modules.Read(c)
	case "event_socket.conf":
		content, err = event_socket.Read(c)
	case "acl.conf":
		content, err = acl.Read(c)
	case "pre_load_switch.conf", "switch.conf", "post_load_switch.conf":
		content, err = switch_main.Read(c)
	//readModule DefaultConf
	case "verto.conf",
		"conference.conf",
		"hash.conf",
		"httapi.conf",
		"spandsp.conf",
		"amr.conf",
		"opus.conf",
		"sndfile.conf",
		"local_stream.conf",
		"lua.conf":
		content, err = readModuleDefaultConf(filename)
	default:
		errtext := fmt.Sprintf(`readConfFromFile filename:%s unsupport!`, filename)
		err = errors.New(errtext)
	}
	return content, err
}

func constConfiguration(c *gin.Context) (string, error) {
	var err error
	var content string
	filename := c.PostForm(`key_value`)
	switch filename {
	case "odbc_cdr.conf":
		content, err = odbc_cdr.Default()
	case "sofia.conf":
		content, err = sofia.Default()
	case "db.conf":
		content, err = switch_db.Default()
	case "fifo.conf":
		content, err = fifo.Default()
	case "pre_load_modules.conf", "post_load_modules.conf":
		content, err = switch_modules.Default(filename)
	case "eventsocket.conf":
		content, err = event_socket.Default()
	case "acl.conf":
		content, err = acl.Default()
	case "pre_load_switch.conf", "swtich.conf", "post_load_switch.conf":
		content, err = switch_main.Default(filename)
	default:
		errtext := fmt.Sprintf(`constConfiguration, filename:%s unsupport!`, filename)
		err = errors.New(errtext)
	}
	return content, err
}

func writeConfToFile(c *gin.Context, content string) {} //nothing todo.

func buildConf(c *gin.Context, old string) (string, error) {
	var err error
	var new string
	filename := c.PostForm(`key_value`)
	switch filename {
	case "odbc_cdr.conf":
		new, err = odbc_cdr.Build(c, old)
	case "sofia.conf":
		new, err = sofia.Build(c, old)
	case "db.conf":
		new, err = switch_db.Build(c, old)
	case "fifo.conf":
		new, err = fifo.Build(c, old)
	case "voicemail.conf":
		new, err = voicemail.Build(c, old)
	case "pre_load_modules.conf", "post_load_modules.conf":
		new, err = switch_modules.Build(c, old)
	case "event_socket.conf":
		new, err = event_socket.Build(c, old)
	case "acl.conf":
		new, err = acl.Build(c, old)
	case "pre_load_switch.conf", "switch.conf", "post_load_switch.conf":
		new, err = switch_main.Build(c, old)
	default:
		errtext := fmt.Sprintf(`buildConf, filename:%s unsupport!`, filename)
		err = errors.New(errtext)
	}
	return new, err
}

func writeConfToDatabase(c *gin.Context, content string, newcontent string) error {
	filename := c.PostForm(`key_value`)
	//function := c.PostForm(`Event-Calling-Function`)
	profile := c.PostForm(`profile`)
	conf := &db.Conf{
		Cfilename:   filename,
		Cprofile:    profile,
		Ccontent:    content,
		Cnewcontent: newcontent,
	}
	return db.InsertConfsConf(conf)
}

func readModuleDefaultConf(filename string) (string, error) {
	file := fmt.Sprintf("%s/autoload_configs/%s%s", viper.GetString(`switch.conf`), filename, `.xml`)
	data, err := os.ReadFile(file)
	return string(data), err
}
