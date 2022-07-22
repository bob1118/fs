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

	"github.com/spf13/viper"
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
	var realtableprefix string
	tableprefix := viper.GetString(`gateway.db.tableprefix`)
	if len(tableprefix) > 0 {
		realtableprefix = fmt.Sprintf(`%s_`, tableprefix)
	} else {
		realtableprefix = tableprefix
	}
	query := fmt.Sprintf("select distinct(account_domain) from %saccounts", realtableprefix)
	err := GetGatewaydb().Select(&domains, query)
	return domains, err
}

func GetAccountsAccount(id, domain string) (*Account, error) {
	var account Account
	var realtableprefix string
	tableprefix := viper.GetString(`gateway.db.tableprefix`)
	if len(tableprefix) > 0 {
		realtableprefix = fmt.Sprintf(`%s_`, tableprefix)
	} else {
		realtableprefix = tableprefix
	}
	query := fmt.Sprintf("select * from %saccounts where true and account_id='%s' and account_domain='%s'", realtableprefix, id, domain)
	err := GetGatewaydb().Get(&account, query)
	return &account, err
}
