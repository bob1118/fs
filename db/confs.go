// gateway table confs
//
// CREATE TABLE IF NOT EXISTS %s (
// 	conf_uuid uuid NOT NULL DEFAULT gen_random_uuid(),
// 	conf_filename varchar NOT NULL,
// 	conf_function varchar NULL,
// 	conf_profile varchar NULL,
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
	"log"

	"github.com/spf13/viper"
)

type Conf struct {
	Cuuid       string `db:"conf_uuid" json:"uuid"`
	Cfilename   string `db:"conf_filename" json:"filename"`
	Cprofile    string `db:"conf_profile" json:"profile"`
	Ccontent    string `db:"conf_content" json:"content"`
	Cnewcontent string `db:"conf_newcontent" json:"newcontent"`
}

func InsertConfsConf(c *Conf) error {
	var err error
	var conf = c
	var realtableprefix string
	tableprefix := viper.GetString(`gateway.db.tableprefix`)
	if len(tableprefix) > 0 {
		realtableprefix = fmt.Sprintf(`%s_`, tableprefix)
	} else {
		realtableprefix = tableprefix
	}
	if len(conf.Cfilename) > 0 && len(conf.Ccontent) > 0 {
		insertsql := fmt.Sprintf("insert into %sconfs(conf_filename, conf_profile, conf_content, conf_newcontent) values(:conf_filename,:conf_profile,:conf_content,:conf_newcontent)", realtableprefix)
		if _, err = GetGatewaydb().NamedExec(insertsql, conf); err != nil {
			log.Println(err)
		}
	} else {
		err = errors.New("InsertGatewayConfsConf in param null")
	}
	return err
}

func GetConfsConf(filename, function, profile string) (*Conf, error) {
	var err error
	var conf Conf
	var realtableprefix string
	tableprefix := viper.GetString(`gateway.db.tableprefix`)
	if len(tableprefix) > 0 {
		realtableprefix = fmt.Sprintf(`%s_`, tableprefix)
	} else {
		realtableprefix = tableprefix
	}
	query := fmt.Sprintf("select * from %sconfs where conf_filename = '%s'and conf_profile='%s'", realtableprefix, filename, profile)
	if len(filename) > 0 {
		if err = GetGatewaydb().Get(&conf, query); err != nil {
			log.Println(err)
		}
	} else {
		err = errors.New("GetGatewayConfsConf filename null")
	}
	return &conf, err
}

func SelectConfs() {}
