// const E164S = `
// CREATE TABLE IF NOT EXISTS %s (
//
//	e164_uuid uuid NOT NULL DEFAULT gen_random_uuid(),
//	gateway_name varchar NOT NULL DEFAULT '',
//	e164_number varchar NOT NULL,
//	e164_enable bool NULL DEFAULT true,
//	e164_lockin bool NULL DEFAULT false,
//	e164_lockout bool NULL DEFAULT false,
//	CONSTRAINT e164s_pkey PRIMARY KEY (e164_uuid),
//	CONSTRAINT e164s_un UNIQUE (e164_number)
//
// );
// COMMENT ON TABLE %s IS 'phone numbers of external gateway';
// `
package db

import (
	"fmt"
	"strings"
)

type E164 struct {
	Euuid    string `db:"e164_uuid" json:"uuid"`
	Gname    string `db:"gateway_name" json:"gname"`
	Enumber  string `db:"e164_number" json:"number"`
	Eenable  bool   `db:"e164_enable" json:"enable"`
	Elockin  bool   `db:"e164_lockin" json:"lockin"`
	Elockout bool   `db:"e164_lockout" json:"lockout"`
}

func SelectE164sWithCondition(condition string) ([]E164, error) {
	e164s := []E164{}
	q := fmt.Sprintf("select * from %se164s where %s", GetTablesGatewayPrifex(), condition)
	err := GetGatewaydb().Select(&e164s, q)
	return e164s, err
}

func IsExistE164(gateway string, number string) (E164, error) {
	var err error
	var e164 E164
	var e164s []E164
	condition := fmt.Sprintf("gateway_name='%s'and e164_number='%s'", gateway, number)
	if e164s, err = SelectE164sWithCondition(condition); err == nil {
		if len(e164s) == 0 {
			err = fmt.Errorf("IsExistE164: no row")
		} else {
			e164 = e164s[0]
		}
	}
	return e164, err
}

func InsertE164s(in []E164) ([]E164, error) {
	var e164 E164
	var newe164s []E164
	var q = fmt.Sprintf("insert into %se164s(gateway_name,e164_number)values", GetTablesGatewayPrifex())

	if len := len(in); len > 0 {
		for index := 0; index < len; index++ {
			e164 = in[index]
			value := fmt.Sprintf("('%s','%s'),", e164.Gname, e164.Enumber)
			q += value
		}
		q = strings.TrimSuffix(q, ",")
		q += (" returning *;")
	}
	err := GetGatewaydb().Select(&newe164s, q)
	return newe164s, err
}

func UpdateE164sE164(uuid string, in E164) (E164, error) {
	var e164 = E164{}
	var q = fmt.Sprintf("update %se164s set ", GetTablesGatewayPrifex())
	if len(in.Gname) > 0 {
		q += fmt.Sprintf("gateway_name='%s',", in.Gname)
	}
	if len(in.Enumber) > 0 {
		q += fmt.Sprintf("e164_number='%s',", in.Enumber)
	}
	if in.Eenable {
		q += "e164_enable=true,"
	} else {
		q += "e164_enable=false,"
	}
	if in.Elockin {
		q += "e164_lockin=true,"
	} else {
		q += "e164_lockin=false,"
	}
	if in.Elockout {
		q += "e164_lockout=true,"
	} else {
		q += "e164_lockout=false,"
	}
	q = strings.TrimSuffix(q, ",")
	q += fmt.Sprintf(" where e164_uuid='%s'", uuid)
	q += (" return *;")

	err := GetGatewaydb().Select(&e164, q)
	return e164, err
}

func DeleteE164sE164(uuid string) (E164, error) {
	var e164 = E164{}
	var q = fmt.Sprintf("delete from %se164s ", GetTablesGatewayPrifex())
	q += fmt.Sprintf("where e164_uuid='%s'", uuid)
	q += (" return *;")
	err := GetGatewaydb().Select(&e164, q)
	return e164, err
}
