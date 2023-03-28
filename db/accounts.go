// gateway table accounts
//
// CREATE TABLE IF NOT EXISTS %s (
// 	account_uuid uuid NOT NULL DEFAULT gen_random_uuid(),
// 	account_id varchar NOT NULL,
// 	account_name varchar NULL,
// 	account_auth varchar NULL,
// 	account_password varchar NULL,
// 	account_a1hash varchar NULL,
// 	account_group varchar NULL,
// 	account_domain varchar NULL,
// 	account_proxy varchar NULL,
// 	account_cacheable varchar NULL,
// 	CONSTRAINT accounts_pkey PRIMARY KEY (account_uuid),
// 	CONSTRAINT accounts_un UNIQUE (account_id, account_domain)
// );

package db

import (
	"fmt"
	"strings"
)

type Account struct {
	Auuid      string `db:"account_uuid" json:"uuid"`
	Aid        string `db:"account_id" json:"id"`
	Aname      string `db:"account_name" json:"name"`
	Aauth      string `db:"account_auth" json:"auth"`
	Apassword  string `db:"account_password" json:"password"`
	Aa1hash    string `db:"account_a1hash" json:"a1hash"`
	Agroup     string `db:"account_group" json:"group"`
	Adomain    string `db:"account_domain" json:"domain"`
	Aproxy     string `db:"account_proxy" json:"proxy"`
	Acacheable string `db:"account_cacheable" json:"cacheable"`
}

func SelectAccountsDistinctDomain() ([]string, error) {
	var domains []string
	q := fmt.Sprintf("select distinct(account_domain) from %saccounts", GetTablesGatewayPrifex())
	err := GetGatewaydb().Select(&domains, q)
	return domains, err
}

func GetAccountsAccount(id, domain string) (*Account, error) {
	var account Account
	q := fmt.Sprintf("select * from %saccounts where true and account_id='%s' and account_domain='%s'", GetTablesGatewayPrifex(), id, domain)
	err := GetGatewaydb().Get(&account, q)
	return &account, err
}

func SelectAccountsWithCondition(condition string) (rt []Account, e error) {
	var err error
	var accounts []Account
	q := fmt.Sprintf("select * from %saccounts where %s", GetTablesGatewayPrifex(), condition)
	if err = GetGatewaydb().Select(&accounts, q); err != nil {
		fmt.Println(err)
	}
	return accounts, err
}

func InsertAccounts(in []Account) (rt []Account, e error) {
	var ua Account
	var newuas []Account
	var q = fmt.Sprintf("insert into %saccounts(account_id,account_name,account_auth,account_password,account_a1hash,account_group,account_domain,account_proxy,account_cacheable)values", GetTablesGatewayPrifex())

	if len := len(in); len > 0 {
		for index := 0; index < len; index++ {
			ua = in[index]
			value := fmt.Sprintf("('%s','%s','%s','%s','%s','%s','%s','%s','%s'),", ua.Aid, ua.Aname, ua.Aauth, ua.Apassword, ua.Aa1hash, ua.Agroup, ua.Adomain, ua.Aproxy, ua.Acacheable)
			q += value
		}
		q = strings.TrimSuffix(q, ",")
		q += (" returning *;")
	}
	err := GetGatewaydb().Select(&newuas, q)
	return newuas, err
}

func UpdateAccountsAccount(uuid string, in Account) (out Account, e error) {
	var ua Account
	var q = fmt.Sprintf("update %saccounts set ", GetTablesGatewayPrifex())
	if len(in.Aid) > 0 {
		q += fmt.Sprintf("account_id='%s',", in.Aid)
	}
	if len(in.Aname) > 0 {
		q += fmt.Sprintf("account_name='%s',", in.Aname)
	}
	if len(in.Aauth) > 0 {
		q += fmt.Sprintf("account_auth='%s',", in.Aauth)
	}
	if len(in.Apassword) > 0 {
		q += fmt.Sprintf("account_password='%s',", in.Apassword)
	}
	if len(in.Aa1hash) > 0 {
		q += fmt.Sprintf("account_a1hash='%s',", in.Aa1hash)
	}
	if len(in.Agroup) > 0 {
		q += fmt.Sprintf("account_group='%s',", in.Agroup)
	}
	if len(in.Adomain) > 0 {
		q += fmt.Sprintf("account_domain='%s',", in.Adomain)
	}
	if len(in.Aproxy) > 0 {
		q += fmt.Sprintf("account_proxy='%s',", in.Aproxy)
	}
	if len(in.Acacheable) > 0 {
		q += fmt.Sprintf("account_cacheable='%s',", in.Acacheable)
	}
	q = strings.TrimSuffix(q, ",")
	q += fmt.Sprintf(" where account_uuid='%s'", uuid)
	q += (" return *;")

	err := GetGatewaydb().Select(&ua, q)
	return ua, err
}

func DeleteAccountsAccount(uuid string) (out Account, e error) {
	var ua = Account{}
	var q = fmt.Sprintf("delete from %saccounts ", GetTablesGatewayPrifex())
	q += fmt.Sprintf("where account_uuid='%s'", uuid)
	err := GetGatewaydb().Select(&ua, q)
	return ua, err
}
