// gateway table confs
//
// CREATE TABLE IF NOT EXISTS %s (
// 	conf_uuid uuid NOT NULL DEFAULT gen_random_uuid(),
// 	conf_filename varchar NOT NULL,
// 	conf_content varchar NOT NULL,
// 	conf_newcontent varchar NULL,
// 	CONSTRAINT confs_pkey PRIMARY KEY (conf_uuid),
// 	CONSTRAINT confs_un UNIQUE (conf_filename)
// );
// COMMENT ON TABLE %s IS 'switch config files which mod_xml_curl requested';

package db

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

var db = GetGatewaydb()
var tablename = fmt.Sprintf(`%s_confs`, viper.GetString(`gateway.db.tableprefix`))

type Conf struct {
	Cuuid       string `db:"conf_uuid" json:"uuid"`
	Cfilename   string `db:"conf_filename" json:"filename"`
	Ccontent    string `db:"conf_content" json:"content"`
	Cnewcontent string `db:"conf_newcontent" json:"newcontent"`
}

func InsertGatewayConfsConf(c *Conf) error {
	var err error
	var conf = c
	if len(conf.Cfilename) > 0 && len(conf.Ccontent) > 0 {
		insertsql := fmt.Sprintf("insert into %s(conf_filenam, conf_content, conf_newcontent) values(:Cfilename,:Ccontent,:Cnewcontent)", tablename)
		_, err = db.NamedExec(insertsql, conf)
	} else {
		err = errors.New("InsertGatewayConfsConf in param null")
	}
	return err
}

//
func GetGatewayConfsConf(filename string) (*Conf, error) {
	var err error
	var conf = Conf{}
	query := fmt.Sprintf("select * from %s where conf_filename = %s", tablename, filename)
	if len(filename) > 0 {
		err = db.Get(&conf, query)
	} else {
		err = errors.New("GetGatewayConfsContent in param null")
	}
	return &conf, err
}
