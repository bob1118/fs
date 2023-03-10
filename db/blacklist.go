package db

import (
	"database/sql"
	"fmt"
)

type Blacklist struct {
	Buuid   string `db:"blacklist_uuid" json:"uuid"`
	Bcaller string `db:"blacklist_caller" json:"caller"`
	Bcallee string `db:"blacklist_callee" json:"callee"`
}

func IsExistBlacklistCaller(caller, callee string) (b Blacklist, exist bool) {
	var is bool
	blacklist := Blacklist{}
	query := fmt.Sprintf("select * from cc_blacklist where blacklist_caller='%s' and blacklist_callee='%s'", caller, callee)
	if err := GetServerdb().Get(&blacklist, query); err != nil {
		if err == sql.ErrNoRows {
			is = false
		}
	} else {
		is = true
	}
	return blacklist, is
}
