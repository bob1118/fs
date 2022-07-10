// switch configuration modules
// p.Dir()+`autoload_configs`+`xxx.conf.xml`
package modules

import (
	"errors"
	"fmt"
	"log"

	"github.com/bob1118/fs/db"
	"github.com/bob1118/fs/fsconf/modules/odbc_cdr"
	"github.com/spf13/viper"
)

func GetConfiguration(filename string) (string, error) {
	var content, newcontent string
	var err error

	if content, err = readConfFromDatabase(filename); err != nil {
		if content, err = readConfFromFile(filename); err != nil {
			if content, err = constConfiguration(filename); err != nil {
				log.Println(err)
			} else {
				writeConfToFile(filename, content)
			}
		}
		if newcontent, err = buildConf(filename, content); err != nil {
			writeConfToDatabase(filename, content, ``)
		} else {
			writeConfToDatabase(filename, content, newcontent)
		}
	}
	return content, err
}

func readConfFromDatabase(filename string) (string, error) {
	var err error
	var content string
	if conf, e := db.GetGatewayConfsConf(filename); e != nil {
		err = e
	} else {
		if len(conf.Ccontent) > 0 {
			if len(conf.Cnewcontent) > 0 {
				content = conf.Cnewcontent
			} else {
				content = conf.Ccontent
			}
		} else {
			err = errors.New("conf.Ccontent null")
		}
	}
	return content, err
}

func readConfFromFile(filename string) (string, error) {
	var err error
	var content string
	switch filename {
	case "odbc_cdr.conf":
		file := fmt.Sprintf("%s/autoload_configs/odbc_cdr.conf.xml", viper.GetString(`switch.conf`))
		content, err = odbc_cdr.Read(file)
	case "sofia.conf":
	}
	return content, err
}

func constConfiguration(filename string) (string, error) {
	var err error
	var content string
	switch filename {
	case "odbc_cdr.conf":
		content, err = odbc_cdr.Default()
	}
	return content, err
}

func writeConfToFile(filename, content string) {} //nothing todo.

func buildConf(filename, old string) (string, error) {
	var err error
	var content string
	switch filename {
	case "odbc_cdr.conf":
		content, err = odbc_cdr.Build(old)
	}
	return content, err
}

func writeConfToDatabase(filename string, content string, newcontent string) error {
	conf := &db.Conf{
		Cfilename:   filename,
		Ccontent:    content,
		Cnewcontent: newcontent,
	}
	return db.InsertGatewayConfsConf(conf)
}
