// gateway table gateways
//
// CREATE TABLE IF NOT EXISTS %s (
// 	gateway_uuid uuid NOT NULL DEFAULT gen_random_uuid(),
// 	profile_name varchar NOT NULL DEFAULT 'external',
// 	gateway_name varchar NOT NULL,
// 	gateway_username varchar NULL,
// 	gateway_realm varchar NULL,
// 	gateway_fromuser varchar NULL,
// 	gateway_fromdomain varchar NULL,
// 	gateway_password varchar NULL,
// 	gateway_extension varchar NULL,
// 	gateway_proxy varchar NULL,
// 	gateway_registerproxy varchar NULL,
// 	gateway_expire varchar NULL,
// 	gateway_register varchar NULL,
// 	gateway_calleridinfrom varchar NULL,
// 	gateway_extensionincontact varchar NULL,
// 	gateway_optionping varchar NULL,
// 	CONSTRAINT gateways_pkey PRIMARY KEY (gateway_uuid),
// 	CONSTRAINT gateways_un UNIQUE (gateway_name)
// );

package db

import (
	"fmt"
	"strings"
)

type Gateway struct {
	Guuid               string `db:"gateway_uuid" json:"uuid"`
	Pname               string `db:"profile_name" json:"pname"`
	Gname               string `db:"gateway_name" json:"gname"`
	Gusername           string `db:"gateway_username" json:"username"`
	Grealm              string `db:"gateway_realm" json:"realm"`
	Gfromuser           string `db:"gateway_fromuser" json:"fromuser"`
	Gfromdomain         string `db:"gateway_fromdomain" json:"fromdomain"`
	Gpassword           string `db:"gateway_password" json:"password"`
	Gextension          string `db:"gateway_extension" json:"extension"`
	Gproxy              string `db:"gateway_proxy" json:"proxy"`
	Gregisterproxy      string `db:"gateway_registerproxy" json:"registerproxy"`
	Gexpire             string `db:"gateway_expire" json:"expire"`
	Gregister           string `db:"gateway_register" json:"register"`
	Gcalleridinfrom     string `db:"gateway_calleridinfrom" json:"calleridinfrom"`
	Gextensionincontact string `db:"gateway_extensionincontact" json:"extensionincontact"`
	Goptionping         string `db:"gateway_optionping" json:"optionping"`
}

func SelectGatewaysWithCondition(condition string) ([]Gateway, error) {
	gateways := []Gateway{}
	query := fmt.Sprintf("select * from %sgateways where %s", GetTablesGatewayPrifex(), condition)
	err := GetGatewaydb().Select(&gateways, query)
	return gateways, err
}

func InsertGateways(in []Gateway) (rt []Gateway, e error) {
	var gw Gateway
	var newgws []Gateway
	var q = fmt.Sprintf("insert into %sgateways(gateway_name,gateway_username,gateway_realm,gateway_fromuser,gateway_fromdomain,gateway_password,gateway_extension,gateway_proxy,gateway_registerproxy,gateway_expire,gateway_register,gateway_calleridinfrom,gateway_extensionincontact,gateway_optionping)values", GetTablesGatewayPrifex())

	if len := len(in); len > 0 {
		for index := 0; index < len; index++ {
			gw = in[index]
			value := fmt.Sprintf("('%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s'),",
				gw.Gname, gw.Gusername, gw.Grealm, gw.Gfromuser, gw.Gfromdomain, gw.Gpassword, gw.Gextension, gw.Gproxy, gw.Gregisterproxy, gw.Gexpire, gw.Gregister, gw.Gcalleridinfrom, gw.Gextensionincontact, gw.Goptionping)
			q += value
		}
		q = strings.TrimSuffix(q, ",")
		q += (" returning *;")
	}
	err := GetGatewaydb().Select(&newgws, q)
	return newgws, err
}

func UpdateGatewaysGateway(uuid string, in Gateway) (out Gateway, e error) {
	var gw Gateway
	var q = fmt.Sprintf("update %sgateways set ", GetTablesGatewayPrifex())

	if len(in.Pname) > 0 {
		q += fmt.Sprintf("profile_name='%s',", in.Pname)
	}
	if len(in.Gname) > 0 {
		q += fmt.Sprintf("gateway_name='%s',", in.Gname)
	}
	if len(in.Gusername) > 0 {
		q += fmt.Sprintf("gateway_username='%s',", in.Gusername)
	}
	if len(in.Grealm) > 0 {
		q += fmt.Sprintf("gateway_realm='%s',", in.Grealm)
	}
	if len(in.Gfromuser) > 0 {
		q += fmt.Sprintf("gateway_fromuser='%s',", in.Gfromuser)
	}
	if len(in.Gfromdomain) > 0 {
		q += fmt.Sprintf("gateway_fromdomain='%s',", in.Gfromdomain)
	}
	if len(in.Gpassword) > 0 {
		q += fmt.Sprintf("gateway_password='%s',", in.Gpassword)
	}
	if len(in.Gextension) > 0 {
		q += fmt.Sprintf("gateway_extension='%s',", in.Gextension)
	}
	if len(in.Gproxy) > 0 {
		q += fmt.Sprintf("gateway_proxy='%s',", in.Gproxy)
	}
	if len(in.Gregisterproxy) > 0 {
		q += fmt.Sprintf("gateway_registerproxy='%s',", in.Gregisterproxy)
	}
	if len(in.Gexpire) > 0 {
		q += fmt.Sprintf("gateway_expire='%s',", in.Gexpire)
	}
	if len(in.Gregister) > 0 {
		q += fmt.Sprintf("gateway_register='%s',", in.Gregister)
	}
	if len(in.Gcalleridinfrom) > 0 {
		q += fmt.Sprintf("gateway_calleridinfrom='%s',", in.Gcalleridinfrom)
	}
	if len(in.Gextensionincontact) > 0 {
		q += fmt.Sprintf("gateway_extensionincontact='%s',", in.Gextensionincontact)
	}
	if len(in.Goptionping) > 0 {
		q += fmt.Sprintf("gateway_optionping='%s',", in.Goptionping)
	}
	q = strings.TrimSuffix(q, ",")
	q += fmt.Sprintf(" where gateway_uuid='%s'", uuid)
	q += (" return *;")

	err := GetGatewaydb().Select(&gw, q)
	return gw, err
}

func DeleteGatewaysGateway(uuid string) (out Gateway, e error) {
	var gw = Gateway{}
	var q = fmt.Sprintf("delete from %sgateways ", GetTablesGatewayPrifex())
	q += fmt.Sprintf("where gateway_uuid='%s'", uuid)
	q += (" return *;")
	err := GetGatewaydb().Select(&gw, q)
	return gw, err
}
