package db

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

var gatewayTableprifex, serverTablePrefix string
var pgdb, switchdb, gatewaydb, serverdb *sqlx.DB

func init() {}

func Initdb() {

	//gatewayTablesPrefix
	g := viper.GetString(`gateway.db.tableprefix`)
	if len(g) > 0 {
		gatewayTableprifex = fmt.Sprintf(`%s_`, g)
	} else {
		gatewayTableprifex = g
	}

	//serverTablePrefix
	s := viper.GetString(`server.db.tableprefix`)
	if len(s) > 0 {
		serverTablePrefix = fmt.Sprintf(`%s_`, s)
	} else {
		serverTablePrefix = s
	}
	pgsqlInitdb()
}

func GetGatewaydb() *sqlx.DB { return gatewaydb }

func GetSwitchdb() *sqlx.DB { return switchdb }

func GetServerdb() *sqlx.DB { return serverdb }

func GetTablesGatewayPrifex() string { return gatewayTableprifex }

func GetTablesServerPrifex() string { return serverTablePrefix }

func pgsqlOpen(strcon string) (*sqlx.DB, error) { return sqlx.Connect(`postgres`, strcon) }

func pgsqlClose(db *sqlx.DB) { db.Close() }

func pgsqlInitdb() {
	var err error
	strcon := fmt.Sprintf(`user=%s password=%s host=%s dbname=%s sslmode=disable`,
		viper.GetString(`postgres.user`), viper.GetString(`postgres.password`), viper.GetString(`postgres.host`), viper.GetString(`postgres.name`))
	pgdb, err = pgsqlOpen(strcon)
	if err != nil {
		fmt.Println(err)
	} else {
		// for debug
		if true {
			pgsqlDestroySwitchdb(pgdb)
		}
		// initSwitchdb
		if switchdb, err = pgsqlInitSwitchdb(pgdb); err != nil {
			fmt.Println(err)
		} else {
			pgsqlInitSwitchMododbccdrTables(switchdb)
		}
		//initGatewaydb
		if gatewaydb, err = pgsqlInitGatewaydb(pgdb); err != nil {
			fmt.Println(err)
		} else {
			pgsqlInitGatewayTables(gatewaydb)
		}
		//initServerdb
		if serverdb, err = pgsqlInitServerdb(pgdb); err != nil {
			fmt.Println(err)
		} else {
			pgsqlInitServerTables(serverdb)
		}
		pgsqlClose(pgdb)
	}
}

func pgsqlInitSwitchdb(db *sqlx.DB) (*sqlx.DB, error) {
	var err error
	var isFound bool
	switchDbhost := viper.GetString(`switch.db.host`)
	switchDbUser := viper.GetString(`switch.db.user`)
	switchDbName := viper.GetString(`switch.db.name`)
	switchDbPassword := viper.GetString(`switch.db.password`)
	//init switch db
	if err = db.Get(&isFound, "select count(1)!=0 as isFound from pg_user where usename=$1", switchDbUser); err != nil {
		fmt.Println(err)
	} else {
		if !isFound { //create user.
			createuser := fmt.Sprintf(USER_CREATE, switchDbUser, switchDbPassword)
			db.MustExec(createuser)
		}
	}
	if err = db.Get(&isFound, "select count(1)!=0 as isFound from pg_database where datname=$1", switchDbName); err != nil {
		fmt.Println(err)
	} else {
		if !isFound { //create db.
			createdb := fmt.Sprintf(DB_CREATE, switchDbName, switchDbUser)
			auth := fmt.Sprintf(DBUSER_AUTH, switchDbName, switchDbUser)
			db.MustExec(createdb)
			db.MustExec(auth)
		}
	}
	//return switch db
	switchstr := fmt.Sprintf(`user=%s password=%s host=%s dbname=%s sslmode=disable`, switchDbUser, switchDbPassword, switchDbhost, switchDbName)
	return pgsqlOpen(switchstr)
}

func pgsqlGetDatabaseOpenConnections(pgdb *sqlx.DB, dbname string) int {
	var connections int
	query := fmt.Sprintf(`select count(1) as connections from pg_stat_activity where datname = '%s';`, dbname)
	if err := pgdb.Get(&connections, query); err != nil {
		fmt.Println(err)
	}
	return connections
}

func pgsqlDestroySwitchdb(db *sqlx.DB) {
	switchDbUser := viper.GetString(`switch.db.user`)
	switchDbName := viper.GetString(`switch.db.name`)
	connections := pgsqlGetDatabaseOpenConnections(db, switchDbName)
	if connections > 0 {
		log.Printf("pgsqlGetDatabaseOpenConnections database:%s connections:%d\n", switchDbName, connections)
	} else {
		dropdb := fmt.Sprintf(DATABASE_DROP, switchDbName)
		dropuser := fmt.Sprintf(USER_DROP, switchDbUser)
		db.MustExec(dropdb)
		db.MustExec(dropuser)
	}
}

func pgsqlInitSwitchMododbccdrTables(db *sqlx.DB) {
	var err error
	var isFound bool
	var tables = make([]string, 0, 10)
	modname := viper.GetString(`switch.cdr.modname`)
	alegname := viper.GetString(`switch.cdr.a-leg`)
	blegname := viper.GetString(`switch.cdr.b-leg`)
	bothname := viper.GetString(`switch.cdr.both`)

	if len(alegname) > 0 {
		tables = append(tables, alegname)
	}
	if len(blegname) > 0 {
		tables = append(tables, blegname)
	}
	if len(bothname) > 0 {
		tables = append(tables, bothname)
	}

	if strings.EqualFold(modname, `mod_odbc_cdr`) {
		for _, table := range tables {
			if err = db.Get(&isFound, "select count(1)!=0 as isFound from pg_tables where tablename =$1", table); err != nil {
				fmt.Println(err)
			} else {
				if !isFound {
					cdrleg := fmt.Sprintf(CDR_LEG, table)
					db.MustExec(cdrleg)
					//db.MustExec(`inserted by switch mod_odbc_cdr`)
				}
			}
		}
	}
}

func pgsqlInitGatewaydb(db *sqlx.DB) (*sqlx.DB, error) {
	var err error
	var isFound bool

	gdbhost := viper.GetString(`gateway.db.host`)
	gdbname := viper.GetString(`gateway.db.name`)
	gdbuser := viper.GetString(`gateway.db.user`)
	gdbpassword := viper.GetString(`gateway.db.password`)
	//init gateway db
	if err = db.Get(&isFound, "select count(1)!=0 as isFound from pg_user where usename=$1", gdbuser); err != nil {
		fmt.Println(err)
	} else {
		if !isFound { //create user.
			createuser := fmt.Sprintf(USER_CREATE, gdbuser, gdbpassword)
			db.MustExec(createuser)
		}
	}
	if err = db.Get(&isFound, "select count(1)!=0 as isFound from pg_database where datname=$1", gdbname); err != nil {
		fmt.Println(err)
	} else {
		if !isFound { //create db.
			createdb := fmt.Sprintf(DB_CREATE, gdbname, gdbuser)
			auth := fmt.Sprintf(DBUSER_AUTH, gdbname, gdbuser)
			db.MustExec(createdb)
			db.MustExec(auth)
		}
	}
	//return gateway db
	gatewaystr := fmt.Sprintf(`user=%s password=%s host=%s dbname=%s sslmode=disable`, gdbuser, gdbpassword, gdbhost, gdbname)
	return pgsqlOpen(gatewaystr)
}

func pgsqlInitGatewayTables(db *sqlx.DB) {
	var err error
	var isFound bool
	ip := viper.GetString(`switch.vars.ipv4`)

	//table confs
	tableConfs := fmt.Sprintf(`%sconfs`, gatewayTableprifex)
	if err = db.Get(&isFound, "select count(1)!=0 as isFound from pg_tables where tablename =$1", tableConfs); err != nil {
		fmt.Println(err)
	} else {
		if !isFound {
			sql := fmt.Sprintf(CONFS, tableConfs, tableConfs)
			db.MustExec(sql)
			//insert DEFAULT_CONFS
		}
	}
	//table accounts
	tableAccounts := fmt.Sprintf(`%saccounts`, gatewayTableprifex)
	if err = db.Get(&isFound, "select count(1)!=0 as isFound from pg_tables where tablename =$1", tableAccounts); err != nil {
		fmt.Println(err)
	} else {
		if !isFound {
			sql := fmt.Sprintf(ACCOUNTS, tableAccounts, tableAccounts)
			db.MustExec(sql)
			//insert DEFAULT_ACCOUNTS
			for i := 8000; i < 8010; i++ {
				s := strconv.Itoa(i)
				ip_domain_sql := fmt.Sprintf(DEFAULT_ACCOUNTS, tableAccounts, s, s, s, s, ip, ip)
				my_domain_sql := fmt.Sprintf(DEFAULT_ACCOUNTS, tableAccounts, s, s, s, s, `mydomain`, ip)
				db.MustExec(ip_domain_sql) //8000-8009 @ip
				db.MustExec(my_domain_sql) //8000-8009 @mydomain
			}
		}
	}
	//table gateways
	tableGateways := fmt.Sprintf(`%sgateways`, gatewayTableprifex)
	if err = db.Get(&isFound, "select count(1)!=0 as isFound from pg_tables where tablename =$1", tableGateways); err != nil {
		fmt.Println(err)
	} else {
		if !isFound {
			sql := fmt.Sprintf(GATEWAYS, tableGateways, tableGateways)
			gatewaysql := fmt.Sprintf(DEFAULT_GATEWAYS, tableGateways)
			db.MustExec(sql)
			db.MustExec(gatewaysql)
		}
	}
	//table e164s
	tableE164s := fmt.Sprintf(`%se164s`, gatewayTableprifex)
	if err = db.Get(&isFound, "select count(1)!=0 as isFound from pg_tables where tablename =$1", tableE164s); err != nil {
		fmt.Println(err)
	} else {
		if !isFound {
			sql := fmt.Sprintf(E164S, tableE164s, tableE164s)
			e164sql := fmt.Sprintf(DEFAULT_E164S, tableE164s)
			db.MustExec(sql)
			db.MustExec(e164sql)
		}
	}
	//table fifos
	tableFifos := fmt.Sprintf(`%sfifos`, gatewayTableprifex)
	if err = db.Get(&isFound, "select count(1)!=0 as isFound from pg_tables where tablename =$1", tableFifos); err != nil {
		fmt.Println(err)
	} else {
		if !isFound {
			sql := fmt.Sprintf(FIFOS, tableFifos, tableFifos)
			fifosql := fmt.Sprintf(DEFAULT_FIFOS, tableFifos)
			db.MustExec(sql)
			db.MustExec(fifosql)
		}
	}
	//table fifomembers
	tableFifomembers := fmt.Sprintf(`%sfifomembers`, gatewayTableprifex)
	if err = db.Get(&isFound, "select count(1)!=0 as isFound from pg_tables where tablename =$1", tableFifomembers); err != nil {
		fmt.Println(err)
	} else {
		if !isFound {
			sql := fmt.Sprintf(FIFOMEMBERS, tableFifomembers, tableFifomembers, tableFifomembers, tableFifos)
			fifomembersql := fmt.Sprintf(DEFAULT_FIFOMEMBERS, tableFifomembers)
			db.MustExec(sql)
			db.MustExec(fifomembersql)
		}
	}
	//table acce164s
	tableAcce164s := fmt.Sprintf(`%sacce164s`, gatewayTableprifex)
	if err = db.Get(&isFound, "select count(1)!=0 as isFound from pg_tables where tablename =$1", tableAcce164s); err != nil {
		fmt.Println(err)
	} else {
		if !isFound {
			sql := fmt.Sprintf(ACCE164S, tableAcce164s, tableAcce164s, tableAcce164s, tableAccounts, tableAcce164s, tableGateways, tableAcce164s, tableE164s)
			acce164sql := fmt.Sprintf(DEFAULT_ACCE164S, tableAcce164s)
			db.MustExec(sql)
			db.MustExec(acce164sql)
		}
	}
	//table e164accs
	tableE164accs := fmt.Sprintf(`%se164accs`, gatewayTableprifex)
	if err = db.Get(&isFound, "select count(1)!=0 as isFound from pg_tables where tablename =$1", tableE164accs); err != nil {
		fmt.Println(err)
	} else {
		if !isFound {
			sql := fmt.Sprintf(E164ACCS, tableE164accs, tableE164accs, tableE164accs, tableAccounts, tableE164accs, tableGateways, tableE164accs, tableE164s, tableE164accs, tableFifos)
			e164accsql := fmt.Sprintf(DEFAULT_E164ACCES, tableE164accs)
			db.MustExec(sql)
			db.MustExec(e164accsql)
		}
	}
}

func pgsqlInitServerdb(db *sqlx.DB) (*sqlx.DB, error) {
	var err error
	var isFound bool

	sdbhost := viper.GetString(`server.db.host`)
	sdbname := viper.GetString(`server.db.name`)
	sdbuser := viper.GetString(`server.db.user`)
	sdbpassword := viper.GetString(`server.db.password`)
	//init server db
	if err = db.Get(&isFound, "select count(1)!=0 as isFound from pg_user where usename=$1", sdbuser); err != nil {
		fmt.Println(err)
	} else {
		if !isFound { //create user.
			createuser := fmt.Sprintf(USER_CREATE, sdbuser, sdbpassword)
			db.MustExec(createuser)
		}
	}
	if err = db.Get(&isFound, "select count(1)!=0 as isFound from pg_database where datname=$1", sdbname); err != nil {
		fmt.Println(err)
	} else {
		if !isFound { //create db.
			createdb := fmt.Sprintf(DB_CREATE, sdbname, sdbuser)
			auth := fmt.Sprintf(DBUSER_AUTH, sdbname, sdbuser)
			db.MustExec(createdb)
			db.MustExec(auth)
		}
	}
	//return server db
	serverstr := fmt.Sprintf(`user=%s password=%s host=%s dbname=%s sslmode=disable`, sdbuser, sdbpassword, sdbhost, sdbname)
	return pgsqlOpen(serverstr)
}

func pgsqlInitServerTables(db *sqlx.DB) {
	var err error
	var isFound bool

	//table backlists
	tableBlacklists := fmt.Sprintf(`%sblacklists`, serverTablePrefix)
	if err = db.Get(&isFound, "select count(1)!=0 as isFound from pg_tables where tablename =$1", tableBlacklists); err != nil {
		fmt.Println(err)
	} else {
		if !isFound {
			sql := fmt.Sprintf(BLACKLISTS, tableBlacklists, tableBlacklists)
			insersql := fmt.Sprintf(DEFAULT_BLACKLISTS, tableBlacklists)
			db.MustExec(sql)
			db.MustExec(insersql)
		}
	}

	//table bgjobs
	tableBgjobs := fmt.Sprintf(`%sbgjobs`, serverTablePrefix)
	if err = db.Get(&isFound, "select count(1)!=0 as isFound from pg_tables where tablename =$1", tableBgjobs); err != nil {
		fmt.Println(err)
	} else {
		if !isFound {
			sql := fmt.Sprintf(BGJOBS, tableBgjobs, tableBgjobs)
			db.MustExec(sql)
			//db.MustExec(`inserted by fs server --run`)
		}
	}

	//table outgoingcalls
	tableOutgoingcalls := fmt.Sprintf(`%soutgoingcalls`, serverTablePrefix)
	if err = db.Get(&isFound, "select count(1)!=0 as isFound from pg_tables where tablename=$1", tableOutgoingcalls); err != nil {
		fmt.Println(err)
	} else {
		if !isFound {
			sql := fmt.Sprintf(OUTGOINGCALLS, tableOutgoingcalls, tableOutgoingcalls)
			db.MustExec(sql)
			//db.MustExec(`inserted by post api/v1/outgoingcall`)
		}
	}
}
