// const ACCE164S = `
// CREATE TABLE IF NOT EXISTS %s (
//
//	acce164_uuid uuid NOT NULL DEFAULT gen_random_uuid(),
//	account_id varchar NOT NULL,
//	account_domain varchar NOT NULL,
//	gateway_name varchar NOT NULL,
//	e164_number varchar NOT NULL,
//	acce164_isdefault bool NOT NULL DEFAULT false,
//	CONSTRAINT acce164_pkey PRIMARY KEY (acce164_uuid)
//
// );
// COMMENT ON TABLE %s IS 'account e164 number for outgoing call';
// ALTER TABLE %s ADD CONSTRAINT acce164_fk FOREIGN KEY (account_id,account_domain) REFERENCES %s(account_id,account_domain) ON DELETE CASCADE ON UPDATE CASCADE;
// ALTER TABLE %s ADD CONSTRAINT acce164_fk_1 FOREIGN KEY (gateway_name) REFERENCES %s(gateway_name) ON DELETE CASCADE ON UPDATE CASCADE;
// ALTER TABLE %s ADD CONSTRAINT acce164_fk_2 FOREIGN KEY (e164_number) REFERENCES %s(e164_number) ON DELETE CASCADE ON UPDATE CASCADE;
// `
package db

import (
	"fmt"
	"strings"
)

// acce164 struct
type ACCE164 struct {
	AEuuid    string `db:"acce164_uuid" json:"uuid"`
	Aid       string `db:"account_id" json:"id"`
	Adomain   string `db:"account_domain" json:"domain"`
	Gname     string `db:"gateway_name" json:"name"`
	Enumber   string `db:"e164_number" json:"number"`
	Isdefault bool   `db:"acce164_isdefault" json:"isdefault"`
}

// SelectAcce164sWithCondition
func SelectAcce164sWithCondition(condition string) ([]ACCE164, error) {
	var acce164s []ACCE164
	query := fmt.Sprintf("select * from %sacce164s where %s", GetTablesGatewayPrifex(), condition)
	err := GetGatewaydb().Select(&acce164s, query)
	return acce164s, err
}

// InsertAcce164s
func InsertAcce164s(in []ACCE164) ([]ACCE164, error) {
	var acce164 ACCE164
	var newacce164s []ACCE164
	var q = fmt.Sprintf("insert into %sacce164s(account_id,account_domain,e164_number,gateway_name) values", GetTablesGatewayPrifex())

	if len := len(in); len > 0 {
		for index := 0; index < len; index++ {
			acce164 = in[index]
			value := fmt.Sprintf("('%s','%s','%s','%s'),", acce164.Aid, acce164.Adomain, acce164.Enumber, acce164.Gname)
			q += value
		}
		q = strings.TrimSuffix(q, ",")
		q += (" returning *;")
	}
	err := GetGatewaydb().Select(&newacce164s, q)
	return newacce164s, err
}

// UpdateAcce164sAcce164
func UpdateAcce164sAcce164(uuid string, in ACCE164) (ACCE164, error) {
	var acce164 = ACCE164{}
	var q = fmt.Sprintf("update %sacce164s set ", GetTablesGatewayPrifex())
	if len(in.Aid) > 0 {
		q += fmt.Sprintf("account_id='%s',", in.Aid)
	}
	if len(in.Adomain) > 0 {
		q += fmt.Sprintf("account_domain='%s',", in.Adomain)
	}
	if len(in.Gname) > 0 {
		q += fmt.Sprintf("gateway_name='%s',", in.Gname)
	}
	if len(in.Enumber) > 0 {
		q += fmt.Sprintf("e164_number='%s',", in.Enumber)
	}
	if in.Isdefault {
		q += "acce164_isdefault=true,"
	} else {
		q += "acce164_isdefault=false,"
	}

	q = strings.TrimSuffix(q, ",")
	q += fmt.Sprintf(" where acce164_uuid='%s'", uuid)
	q += (" returning *;")

	err := GetGatewaydb().Get(&acce164, q)
	return acce164, err
}

// DeleteAcce164sAcce164
func DeleteAcce164sAcce164(uuid string) (ACCE164, error) {
	var acce164 = ACCE164{}
	var q = fmt.Sprintf("delete from %sacce164s ", GetTablesGatewayPrifex())
	q += fmt.Sprintf("where acce164_uuid='%s'", uuid)
	q += (" returning *;")
	err := GetGatewaydb().Get(&acce164, q)
	return acce164, err
}
