// gateway table fifomembers
//
// CREATE TABLE public.g_fifomembers (
// 	fifomember_uuid uuid NOT NULL DEFAULT gen_random_uuid(),
// 	fifo_name varchar NOT NULL,
// 	member_string varchar NOT NULL,
// 	member_simo varchar NULL DEFAULT 1,
// 	member_timeout varchar NULL DEFAULT 10,
// 	member_lag varchar NULL DEFAULT 10,
// 	CONSTRAINT fifomembers_pkey PRIMARY KEY (fifomember_uuid),
// 	CONSTRAINT fifomember_fk FOREIGN KEY (fifo_name) REFERENCES public.g_fifos(fifo_name) ON DELETE CASCADE ON UPDATE CASCADE
// );

package db

import (
	"fmt"

	"github.com/spf13/viper"
)

type FifoMember struct {
	Muuid    string `db:"fifomember_uuid" json:"uuid"`
	Fname    string `db:"fifo_name" json:"name"`
	Mstring  string `db:"member_string" json:"string"`
	Msimo    string `db:"member_simo" json:"simo"`
	Mtimeout string `db:"member_timeout" json:"timeout"`
	Mlag     string `db:"member_lag" json:"lag"`
}

func GetFifomembers(s string) ([]FifoMember, error) {
	return nil, nil
}

func SelectFifoMembers() {}

func SelectFifoMembersByFifoname(fifoname string) ([]FifoMember, error) {
	var fifomembers = []FifoMember{}
	var realtableprefix string
	tableprefix := viper.GetString(`gateway.db.tableprefix`)
	if len(tableprefix) > 0 {
		realtableprefix = fmt.Sprintf(`%s_`, tableprefix)
	} else {
		realtableprefix = tableprefix
	}
	query := fmt.Sprintf(`select * from %sfifomembers where fifo_name='%s'`, realtableprefix, fifoname)
	err := GetGatewaydb().Select(&fifomembers, query)
	return fifomembers, err
}
