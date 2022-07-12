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

type Conf struct {
	Cuuid       string `db:"conf_uuid" json:"uuid"`
	Cfilename   string `db:"conf_filename" json:"filename"`
	Cfunction   string `db:"conf_function" json:"function"`
	Cprofile    string `db:"conf_profile" json:"profile"`
	Ccontent    string `db:"conf_content" json:"content"`
	Cnewcontent string `db:"conf_newcontent" json:"newcontent"`
}

func InsertGatewayConfsConf(c *Conf) error {
	var err error
	var conf = c
	if len(conf.Cfilename) > 0 && len(conf.Ccontent) > 0 {
		insertsql := fmt.Sprintf("insert into %s_confs(conf_filename, conf_function, conf_profile, conf_content, conf_newcontent) values(:conf_filename,:conf_function,:conf_profile,:conf_content,:conf_newcontent)",
			viper.GetString(`gateway.db.tableprefix`))
		_, err = GetGatewaydb().NamedExec(insertsql, conf)
	} else {
		err = errors.New("InsertGatewayConfsConf in param null")
	}
	return err
}

func GetGatewayConfsConf(filename, function, profile string) (*Conf, error) {
	var err error
	var conf = Conf{}
	query := fmt.Sprintf("select * from %s_confs where conf_filename = '%s' and conf_function='%s' and conf_profile='%s'", viper.GetString(`gateway.db.tableprefix`), filename, function, profile)
	if len(filename) > 0 {
		err = GetGatewaydb().Get(&conf, query)
	} else {
		err = errors.New("GetGatewayConfsConf filename null")
	}
	return &conf, err
}
