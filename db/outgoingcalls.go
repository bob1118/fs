package db

import (
	"fmt"
	"strings"
)

// acce164 struct
type OUTGOINGCALL struct {
	Auuid       string `db:"uuida" json:"uuida"`
	Buuid       string `db:"uuidb" json:"uuidb"`
	Id          string `db:"id" json:"id"`
	Domain      string `db:"domain" json:"domain"`
	E164        string `db:"e164" json:"e164"`
	Gateway     string `db:"gateway" json:"gateway"`
	Ani         string `db:"ani" json:"ani"`
	Destination string `db:"destination" json:"destination"`
}

// SelectOutgoingcallsWithCondition
func SelectOutgoingcallsWithCondition(condition string) ([]OUTGOINGCALL, error) {
	var outgoingcalls []OUTGOINGCALL
	query := fmt.Sprintf("select * from %soutgoingcalls where %s", GetTablesServerPrifex(), condition)
	err := GetGatewaydb().Select(&outgoingcalls, query)
	return outgoingcalls, err
}

// InsertOutgoingcalls
func InsertOutgoingcalls(in []OUTGOINGCALL) ([]OUTGOINGCALL, error) {
	var call OUTGOINGCALL
	var newoutgoingcalls []OUTGOINGCALL
	var q = fmt.Sprintf("insert into %soutgoingcalls(uuida,uuidb,id,domain,e164,gateway,ani,destination) values", GetTablesServerPrifex())

	if len := len(in); len > 0 {
		for index := 0; index < len; index++ {
			call = in[index]
			value := fmt.Sprintf("('%s','%s','%s','%s','%s','%s','%s','%s'),", call.Auuid, call.Id, call.Domain, call.Buuid, call.E164, call.Gateway, call.Ani, call.Destination)
			q += value
		}
		q = strings.TrimSuffix(q, ",")
		q += (" returning *;")
	}
	err := GetGatewaydb().Select(&newoutgoingcalls, q)
	return newoutgoingcalls, err
}
