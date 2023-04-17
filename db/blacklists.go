// const BLACKLISTS = `
// CREATE TABLE IF NOT EXISTS %s (
// 	blacklist_uuid uuid NOT NULL DEFAULT gen_random_uuid(),
// 	blacklist_caller varchar NOT NULL,
// 	blacklist_callee varchar NOT NULL,
// 	CONSTRAINT blacklists_pkey PRIMARY KEY (blacklist_uuid)
// );
// COMMENT ON TABLE %s IS 'call filter blacklist include caller and callee';
// `

package db

import (
	"fmt"
	"strings"
)

type Blacklist struct {
	Buuid   string `db:"blacklist_uuid" json:"uuid"`
	Bcaller string `db:"blacklist_caller" json:"caller"`
	Bcallee string `db:"blacklist_callee" json:"callee"`
}

// // IsExistBlacklistCaller
// func IsExistBlacklistCaller(caller, callee string) (b Blacklist, exist bool) {
// 	var is bool
// 	blacklist := Blacklist{}
// 	query := fmt.Sprintf("select * from %sblacklist where blacklist_caller='%s' and blacklist_callee='%s'", GetTablesServerPrifex(), caller, callee)
// 	if err := GetServerdb().Get(&blacklist, query); err != nil {
// 		if err == sql.ErrNoRows {
// 			is = false
// 		}
// 	} else {
// 		is = true
// 	}
// 	return blacklist, is
// }

// SelectBlacklistsWithCondition
func SelectBlacklistsWithCondition(condition string) ([]Blacklist, error) {
	var blacklists = []Blacklist{}
	q := fmt.Sprintf(`select * from %sblacklist where %s`, GetTablesServerPrifex(), condition)
	err := GetServerdb().Select(&blacklists, q)
	return blacklists, err
}

// InsertBlacklists
func InsertBlacklists(in []Blacklist) ([]Blacklist, error) {
	var blacklist Blacklist
	var blacklists []Blacklist
	var q = fmt.Sprintf("insert into %sblacklists(blacklist_caller,blacklist_callee) values ", GetTablesServerPrifex())

	if len := len(in); len > 0 {
		for index := 0; index < len; index++ {
			blacklist = in[index]
			value := fmt.Sprintf("('%s','%s'),", blacklist.Bcaller, blacklist.Bcallee)
			q += value
		}
		q = strings.TrimSuffix(q, ",")
		q += (" returning *;")
	}
	err := GetServerdb().Select(&blacklists, q)
	return blacklists, err
}

// UpdateBlacklistsBlacklist
func UpdateBlacklistsBlacklist(uuid string, in Blacklist) (Blacklist, error) {
	var blacklist = Blacklist{}
	var q = fmt.Sprintf("update %sblacklists set ", GetTablesServerPrifex())
	if len(in.Bcaller) > 0 {
		q += fmt.Sprintf("blacklist_caller='%s',", in.Bcaller)
	}
	if len(in.Bcallee) > 0 {
		q += fmt.Sprintf("blacklist_callee='%s',", in.Bcallee)
	}
	q = strings.TrimSuffix(q, ",")
	q += fmt.Sprintf(" where blacklist_uuid='%s'", uuid)
	q += (" return *;")

	err := GetServerdb().Select(&blacklist, q)
	return blacklist, err
}

// DeleteBlacklistsBlacklist
func DeleteBlacklistsBlacklist(uuid string) (Blacklist, error) {
	var blacklist = Blacklist{}
	var q = fmt.Sprintf("delete from %sblacklist ", GetTablesServerPrifex())
	q += fmt.Sprintf(" where blacklist_uuid='%s'", uuid)
	q += (" return *;")
	err := GetServerdb().Select(&blacklist, q)
	return blacklist, err
}
