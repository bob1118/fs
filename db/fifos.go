// gateway table fifos
//
// CREATE TABLE public.g_fifos (
// 	fifo_uuid uuid NOT NULL DEFAULT gen_random_uuid(),
// 	fifo_name varchar NOT NULL,
// 	fifo_importance varchar NULL DEFAULT 0,
// 	fifo_announce varchar NULL DEFAULT ''::character varying,
// 	fifo_holdmusic varchar NULL DEFAULT ''::character varying,
// 	CONSTRAINT fifos_pkey PRIMARY KEY (fifo_uuid),
// 	CONSTRAINT fifos_un UNIQUE (fifo_name)
// );
package db

import (
	"fmt"

	"github.com/spf13/viper"
)

type Fifo struct {
	Fuuid       string `db:"fifo_uuid" json:"uuid"`
	Fname       string `db:"fifo_name" json:"name"`
	Fimportance string `db:"fifo_importance" json:"importance"`
	Fannounce   string `db:"fifo_announce" json:"announce"`
	Fholdmusic  string `db:"fifo_holdmusic" json:"holdmusic"`
}

func SelectFifos() ([]Fifo, error) {
	fifos := []Fifo{}
	var realtableprefix string
	tableprefix := viper.GetString(`gateway.db.tableprefix`)
	if len(tableprefix) > 0 {
		realtableprefix = fmt.Sprintf(`%s_`, tableprefix)
	} else {
		realtableprefix = tableprefix
	}
	query := fmt.Sprintf("select * from %sfifos", realtableprefix)
	err := GetGatewaydb().Select(&fifos, query)
	return fifos, err
}
