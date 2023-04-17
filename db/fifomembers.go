// const FIFOMEMBERS = `
// CREATE TABLE IF NOT EXISTS %s (
// 	fifomember_uuid uuid NOT NULL DEFAULT gen_random_uuid(),
// 	fifo_name varchar NOT NULL,
// 	member_string varchar NOT NULL,
// 	member_simo varchar NULL DEFAULT 1,
// 	member_timeout varchar NULL DEFAULT 10,
// 	member_lag varchar NULL DEFAULT 10,
// 	CONSTRAINT fifomembers_pkey PRIMARY KEY (fifomember_uuid)
// );
// COMMENT ON TABLE %s IS 'mod_fifo fifo members';
// ALTER TABLE %s ADD CONSTRAINT fifomember_fk FOREIGN KEY (fifo_name) REFERENCES %s(fifo_name) ON DELETE CASCADE ON UPDATE CASCADE;
// `

package db

import (
	"fmt"
	"strings"
)

type FifoMember struct {
	Muuid    string `db:"fifomember_uuid" json:"uuid"`
	Fname    string `db:"fifo_name" json:"name"`
	Mstring  string `db:"member_string" json:"string"`
	Msimo    string `db:"member_simo" json:"simo"`
	Mtimeout string `db:"member_timeout" json:"timeout"`
	Mlag     string `db:"member_lag" json:"lag"`
}

// SelectFifomembersWithCondition
func SelectFifomembersWithCondition(condition string) ([]FifoMember, error) {
	var fifomembers = []FifoMember{}
	q := fmt.Sprintf(`select * from %sfifomembers where %s`, GetTablesGatewayPrifex(), condition)
	err := GetGatewaydb().Select(&fifomembers, q)
	return fifomembers, err
}

// InsertFifoMembers
func InsertFifoMembers(in []FifoMember) ([]FifoMember, error) {
	var fifomember FifoMember
	var fifomembers []FifoMember
	var q = fmt.Sprintf("insert into %sfifomembers(fifo_name,member_string,member_simo,member_timeout,member_lag) values", GetTablesGatewayPrifex())

	if len := len(in); len > 0 {
		for index := 0; index < len; index++ {
			fifomember = in[index]
			value := fmt.Sprintf("('%s','%s','%s','%s','%s'),", fifomember.Fname, fifomember.Mstring, fifomember.Msimo, fifomember.Mtimeout, fifomember.Mlag)
			q += value
		}
		q = strings.TrimSuffix(q, ",")
		q += (" returning *;")
	}
	err := GetGatewaydb().Select(&fifomembers, q)
	return fifomembers, err
}

// UpdateFifomembersFifomember
func UpdateFifomembersFifomember(uuid string, in FifoMember) (FifoMember, error) {
	var fifomember = FifoMember{}
	var q = fmt.Sprintf("update %sfifomembers set ", GetTablesGatewayPrifex())
	if len(in.Fname) > 0 {
		q += fmt.Sprintf("fifo_name='%s',", in.Fname)
	}
	if len(in.Mstring) > 0 {
		q += fmt.Sprintf("member_string='%s',", in.Mstring)
	}
	if len(in.Msimo) > 0 {
		q += fmt.Sprintf("member_simo='%s',", in.Msimo)
	}
	if len(in.Mtimeout) > 0 {
		q += fmt.Sprintf("member_timeout='%s',", in.Mtimeout)
	}
	if len(in.Mlag) > 0 {
		q += fmt.Sprintf("member_lag='%s',", in.Mlag)
	}
	q = strings.TrimSuffix(q, ",")
	q += fmt.Sprintf(" where fifomember_uuid='%s'", uuid)
	q += (" return *;")

	err := GetGatewaydb().Select(&fifomember, q)
	return fifomember, err
}

// DeleteFifomembersFifomember
func DeleteFifomembersFifomember(uuid string) (FifoMember, error) {
	var fifomember = FifoMember{}
	var q = fmt.Sprintf("delete from %sfifomembers ", GetTablesGatewayPrifex())
	q += fmt.Sprintf(" where fifomember_uuid='%s'", uuid)
	q += (" return *;")
	err := GetGatewaydb().Select(&fifomember, q)
	return fifomember, err
}
