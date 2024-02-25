// const E164ACCS = `
// CREATE TABLE IF NOT EXISTS %s (
//
//	e164acc_uuid uuid NOT NULL DEFAULT gen_random_uuid(),
//	gateway_name varchar NOT NULL,
//	e164_number varchar NOT NULL,
//	account_id varchar NULL DEFAULT '',
//	account_domain varchar NULL DEFAULT '',
//	fifo_name varchar NULL DEFAULT '',
//	e164acc_isfifo bool NOT NULL DEFAULT false,
//	CONSTRAINT e164acc_pkey PRIMARY KEY (e164acc_uuid)
//
// );
// COMMENT ON TABLE %s IS 'gateway e164 receive incoming call, bridge account@domain or fifo fifo@fifoname in';
// ALTER TABLE %s ADD CONSTRAINT e164acc_fk FOREIGN KEY (account_id,account_domain) REFERENCES %s(account_id,account_domain) ON DELETE CASCADE ON UPDATE CASCADE;
// ALTER TABLE %s ADD CONSTRAINT e164acc_fk_1 FOREIGN KEY (gateway_name) REFERENCES %s(gateway_name) ON DELETE CASCADE ON UPDATE CASCADE;
// ALTER TABLE %s ADD CONSTRAINT e164acc_fk_2 FOREIGN KEY (e164_number) REFERENCES %s(e164_number) ON DELETE CASCADE ON UPDATE CASCADE;
// ALTER TABLE %s ADD CONSTRAINT e164acc_fk_3 FOREIGN KEY (fifo_name) REFERENCES %s(fifo_name) ON DELETE CASCADE ON UPDATE CASCADE;
// `
package db

import (
	"fmt"
	"strings"
)

// e164acc struct
type E164ACC struct {
	EAuuid  string `db:"e164acc_uuid" json:"uuid"`
	Gname   string `db:"gateway_name" json:"gname"`
	Enumber string `db:"e164_number" json:"number"`
	Aid     string `db:"account_id" json:"id"`
	Adomain string `db:"account_domain" json:"domain"`
	Fname   string `db:"fifo_name" json:"fname"`
	Isfifo  bool   `db:"e164acc_isfifo" json:"isfifo"`
}

// SelectE164accsWithCondition
func SelectE164accsWithCondition(condition string) ([]E164ACC, error) {
	var e164accs []E164ACC
	query := fmt.Sprintf("select * from %se164accs where %s", GetTablesGatewayPrifex(), condition)
	err := GetGatewaydb().Select(&e164accs, query)
	return e164accs, err
}

// InsertE164accs
func InsertE164accs(in []E164ACC) ([]E164ACC, error) {
	var e164acc E164ACC
	var newe164accs []E164ACC
	var q = fmt.Sprintf("insert into %se164accs(gateway_name,e164_number,account_id,account_domain,fifo_name,e164acc_isfifo) values", GetTablesGatewayPrifex())

	if len := len(in); len > 0 {
		for index := 0; index < len; index++ {
			e164acc = in[index]
			value := fmt.Sprintf("('%s','%s','%s','%s','%s',%t),", e164acc.Gname, e164acc.Enumber, e164acc.Aid, e164acc.Adomain, e164acc.Fname, e164acc.Isfifo)
			q += value
		}
		q = strings.TrimSuffix(q, ",")
		q += (" returning *;")
	}
	err := GetGatewaydb().Select(&newe164accs, q)
	return newe164accs, err
}

// UpdateE164accsE164acc
func UpdateE164accsE164acc(uuid string, in E164ACC) (E164ACC, error) {
	var e164acc = E164ACC{}
	var q = fmt.Sprintf("update %se164accs set ", GetTablesGatewayPrifex())
	if len(in.Gname) > 0 {
		q += fmt.Sprintf("gateway_name='%s',", in.Gname)
	}
	if len(in.Enumber) > 0 {
		q += fmt.Sprintf("e164_number='%s',", in.Enumber)
	}
	if len(in.Aid) > 0 {
		q += fmt.Sprintf("account_id='%s',", in.Aid)
	}
	if len(in.Adomain) > 0 {
		q += fmt.Sprintf("account_domain='%s',", in.Adomain)
	}
	if (len(in.Fname)) > 0 {
		q += fmt.Sprintf("fifo_name='%s',", in.Fname)
	}
	if in.Isfifo {
		q += "e164acc_isfifo=true,"
	} else {
		q += "e164acc_isfifo=false,"
	}

	q = strings.TrimSuffix(q, ",")
	q += fmt.Sprintf(" where e164acc_uuid='%s'", uuid)
	q += (" returning *;")

	err := GetGatewaydb().Get(&e164acc, q)
	return e164acc, err
}

// DeleteAcce164sAcce164
func DeleteAcce164sAcce164(uuid string) (ACCE164, error) {
	var acce164 = ACCE164{}
	var q = fmt.Sprintf("delete from %se164accs ", GetTablesGatewayPrifex())
	q += fmt.Sprintf("where e164acc_uuid='%s'", uuid)
	q += (" returning *;")
	err := GetGatewaydb().Get(&acce164, q)
	return acce164, err
}
