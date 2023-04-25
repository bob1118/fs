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
	"fmt"
	"strings"
)

type Conf struct {
	Cuuid       string `db:"conf_uuid" json:"uuid"`
	Cfilename   string `db:"conf_filename" json:"filename"`
	Cprofile    string `db:"conf_profile" json:"profile"`
	Ccontent    string `db:"conf_content" json:"content"`
	Cnewcontent string `db:"conf_newcontent" json:"newcontent"`
}

// GetConfsConfWithCondition
func GetConfsConfWithCondition(condition string) (Conf, error) {
	var conf Conf
	query := fmt.Sprintf("select * from %sconfs where %s", GetTablesGatewayPrifex(), condition)
	err := GetGatewaydb().Get(&conf, query)
	return conf, err
}

// SelectConfsWithCondition
func SelectConfsWithCondition(condition string) ([]Conf, error) {
	var confs []Conf
	query := fmt.Sprintf("select * from %sconfs where %s", GetTablesGatewayPrifex(), condition)
	err := GetGatewaydb().Select(&confs, query)
	return confs, err
}

// InsertConfsConf
func InsertConfsConf(in Conf) (Conf, error) {
	var conf Conf
	var q = fmt.Sprintf("insert into %sconfs(conf_filename,conf_profile,conf_content,conf_newcontent) values", GetTablesGatewayPrifex())
	q += fmt.Sprintf("('%s','%s','%s','%s')", in.Cfilename, in.Cprofile, in.Ccontent, in.Cnewcontent)
	q += (" returning *;")
	err := GetGatewaydb().Get(&conf, q)
	return conf, err
}

// InsertConfs
func InsertConfs(in []Conf) ([]Conf, error) {
	var conf Conf
	var newconfs []Conf
	var q = fmt.Sprintf("insert into %sconfs(conf_filename,conf_profile,conf_content,conf_newcontent) values", GetTablesGatewayPrifex())

	if len := len(in); len > 0 {
		for index := 0; index < len; index++ {
			conf = in[index]
			value := fmt.Sprintf("('%s','%s','%s','%s'),", conf.Cfilename, conf.Cprofile, conf.Ccontent, conf.Cnewcontent)
			q += value
		}
		q = strings.TrimSuffix(q, ",")
		q += (" returning *;")
	}
	err := GetGatewaydb().Select(&newconfs, q)
	return newconfs, err
}

// UpdateConfsConf
func UpdateConfsConf(uuid string, in Conf) (Conf, error) {
	var conf = Conf{}
	var q = fmt.Sprintf("update %sconfs set ", GetTablesGatewayPrifex())
	if len(in.Cfilename) > 0 {
		q += fmt.Sprintf("account_id='%s',", in.Cfilename)
	}
	if len(in.Cprofile) > 0 {
		q += fmt.Sprintf("account_domain='%s',", in.Cprofile)
	}
	if len(in.Ccontent) > 0 {
		q += fmt.Sprintf("gateway_name='%s',", in.Ccontent)
	}
	if len(in.Cnewcontent) > 0 {
		q += fmt.Sprintf("e164_number='%s',", in.Cnewcontent)
	}
	q = strings.TrimSuffix(q, ",")
	q += fmt.Sprintf(" where conf_uuid='%s'", uuid)
	q += (" returning *;")

	err := GetGatewaydb().Select(&conf, q)
	return conf, err
}

// DeleteConfsConf
func DeleteConfsConf(uuid string) (Conf, error) {
	var conf = Conf{}
	var q = fmt.Sprintf("delete from %sconfs ", GetTablesGatewayPrifex())
	q += fmt.Sprintf("where conf_uuid='%s'", uuid)
	q += (" returning *;")
	err := GetGatewaydb().Select(&conf, q)
	return conf, err
}
