// switch configuration modules
// p.Dir()+`autoload_configs`+`xxx.conf.xml`
package modules

import (
	"errors"
	"log"

	"github.com/bob1118/fs/db"
	"github.com/bob1118/fs/fsconf/modules/odbc_cdr"
	"github.com/bob1118/fs/fsconf/modules/sofia"
	"github.com/gin-gonic/gin"
)

func GetConfiguration(c *gin.Context) (string, error) {
	var err error
	var content, newcontent string
	if content, err = readConfFromDatabase(c); err != nil {
		if content, err = readConfFromFile(c); err != nil {
			if content, err = constConfiguration(c); err != nil {
				log.Println(err)
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
	return content, err
}

func readConfFromDatabase(c *gin.Context) (string, error) {
	var err error
	var content string
	filename := c.PostForm(`key_value`)
	function := c.PostForm(`Event-Calling-Function`)
	profile := c.PostForm(`profile`)
	if conf, e := db.GetGatewayConfsConf(filename, function, profile); e != nil {
		err = e
	} else {
		if len(conf.Ccontent) > 0 {
			if len(conf.Cnewcontent) > 0 {
				content = conf.Cnewcontent
			} else {
				content = conf.Ccontent
			}
		} else {
			err = errors.New("db.Conf.Ccontent null")
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
		//file := fmt.Sprintf("%s/autoload_configs/odbc_cdr.conf.xml", viper.GetString(`switch.conf`))
		content, err = odbc_cdr.Read(c)
	case "sofia.conf":
		//file := fmt.Sprintf("%s/autoload_configs/sofia.conf.xml", viper.GetString(`switch.conf`))
		content, err = sofia.Read(c)
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
	}
	return new, err
}

func writeConfToDatabase(c *gin.Context, content string, newcontent string) error {
	filename := c.PostForm(`key_value`)
	conf := &db.Conf{
		Cfilename:   filename,
		Ccontent:    content,
		Cnewcontent: newcontent,
	}
	return db.InsertGatewayConfsConf(conf)
}
