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

	"github.com/spf13/viper"
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

func SelectGateways(condition string) ([]Gateway, error) {
	gateways := []Gateway{}
	var realtableprefix string
	tableprefix := viper.GetString(`gateway.db.tableprefix`)
	if len(tableprefix) > 0 {
		realtableprefix = fmt.Sprintf(`%s_`, tableprefix)
	} else {
		realtableprefix = tableprefix
	}
	query := fmt.Sprintf("select * from %sgateways where true %s", realtableprefix, condition)
	err := GetGatewaydb().Select(&gateways, query)
	return gateways, err
}
