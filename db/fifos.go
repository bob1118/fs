// const FIFOS = `
// CREATE TABLE IF NOT EXISTS %s (
//
//	fifo_uuid uuid NOT NULL DEFAULT gen_random_uuid(),
//	fifo_name varchar NOT NULL,
//	fifo_importance varchar NULL DEFAULT 0,
//	fifo_announce varchar NULL DEFAULT '',
//	fifo_holdmusic varchar NULL DEFAULT '',
//	CONSTRAINT fifos_pkey PRIMARY KEY (fifo_uuid),
//	CONSTRAINT fifos_un UNIQUE (fifo_name)
//
// );
// COMMENT ON TABLE %s IS 'mod_fifo fifos';
// `
package db

import (
	"fmt"
	"strings"
)

type Fifo struct {
	Fuuid       string `db:"fifo_uuid" json:"uuid"`
	Fname       string `db:"fifo_name" json:"name"`
	Fimportance string `db:"fifo_importance" json:"importance"`
	Fannounce   string `db:"fifo_announce" json:"announce"`
	Fholdmusic  string `db:"fifo_holdmusic" json:"holdmusic"`
}

// SelectFifos
func SelectFifos() ([]Fifo, error) {
	fifos := []Fifo{}
	q := fmt.Sprintf("select * from %sfifos", GetTablesGatewayPrifex())
	err := GetGatewaydb().Select(&fifos, q)
	return fifos, err
}

// SelectFifosWithCondition
func SelectFifosWithCondition(condition string) ([]Fifo, error) {
	var fifos []Fifo
	q := fmt.Sprintf("select * from %sacce164s where %s", GetTablesGatewayPrifex(), condition)
	err := GetGatewaydb().Select(&fifos, q)
	return fifos, err
}

// InsertFifos
func InsertFifos(in []Fifo) ([]Fifo, error) {
	var fifo Fifo
	var fifos []Fifo
	var q = fmt.Sprintf("insert into %sfifos(fifo_name,fifo_importance,fifo_announce,fifo_holdmusic) values", GetTablesGatewayPrifex())

	if len := len(in); len > 0 {
		for index := 0; index < len; index++ {
			fifo = in[index]
			value := fmt.Sprintf("('%s','%s','%s','%s'),", fifo.Fname, fifo.Fimportance, fifo.Fannounce, fifo.Fholdmusic)
			q += value
		}
		q = strings.TrimSuffix(q, ",")
		q += (" returning *;")
	}
	err := GetGatewaydb().Select(&fifos, q)
	return fifos, err
}

// UpdateFifosFifo
func UpdateFifosFifo(uuid string, in Fifo) (Fifo, error) {
	var fifo = Fifo{}
	var q = fmt.Sprintf("update %sfifos set ", GetTablesGatewayPrifex())
	if len(in.Fname) > 0 {
		q += fmt.Sprintf("fifo_name='%s',", in.Fname)
	}
	if len(in.Fimportance) > 0 {
		q += fmt.Sprintf("fifo_importance='%s',", in.Fimportance)
	}
	if len(in.Fannounce) > 0 {
		q += fmt.Sprintf("fifo_announce='%s',", in.Fannounce)
	}
	if len(in.Fholdmusic) > 0 {
		q += fmt.Sprintf("fifo_holdmusic='%s',", in.Fholdmusic)
	}
	q = strings.TrimSuffix(q, ",")
	q += fmt.Sprintf(" where fifo_uuid='%s'", uuid)
	q += (" return *;")

	err := GetGatewaydb().Select(&fifo, q)
	return fifo, err
}

// DeleteFifosFifo
func DeleteFifosFifo(uuid string) (Fifo, error) {
	var fifo = Fifo{}
	var q = fmt.Sprintf("delete from %sfifos ", GetTablesGatewayPrifex())
	q += fmt.Sprintf(" where fifo_uuid='%s'", uuid)
	q += (" return *;")
	err := GetGatewaydb().Select(&fifo, q)
	return fifo, err
}
