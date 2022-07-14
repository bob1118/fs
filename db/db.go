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

var pgdb, switchdb, gatewaydb, serverdb *sqlx.DB

func init() {}

func Initdb() { pgsqlInitdb() }

func GetGatewaydb() *sqlx.DB { return gatewaydb }

func GetSwitchdb() *sqlx.DB { return switchdb }

func GetServerdb() *sqlx.DB { return serverdb }

func pgsqlOpen(strcon string) (*sqlx.DB, error) { return sqlx.Connect(`postgres`, strcon) }

func pgsqlClose(db *sqlx.DB) { db.Close() }

func pgsqlInitdb() {
	var err error
	strcon := fmt.Sprintf(`user=%s password=%s host=%s dbname=%s`,
		viper.GetString(`postgres.user`), viper.GetString(`postgres.password`), viper.GetString(`postgres.host`), viper.GetString(`postgres.name`))
	pgdb, err = pgsqlOpen(strcon)
	if err != nil {
		log.Println(err)
	} else {
		// for debug
		if true {
			pgsqlDestroySwitchdb(pgdb)
		}
		// initSwitchdb
		if switchdb, err = pgsqlInitSwitchdb(pgdb); err != nil {
			log.Println(err)
		} else {
			pgsqlInitSwitchMododbccdrTables(switchdb)
		}
		//initGatewaydb
		if gatewaydb, err = pgsqlInitGatewaydb(pgdb); err != nil {
			log.Println(err)
		} else {
			pgsqlInitGatewayTables(gatewaydb)
		}
		//initServerdb
		if serverdb, err = pgsqlInitServerdb(pgdb); err != nil {
			log.Println(err)
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
		log.Println(err)
	} else {
		if !isFound { //create user.
			createuser := fmt.Sprintf(USER_CREATE, switchDbUser, switchDbPassword)
			db.MustExec(createuser)
		}
	}
	if err = db.Get(&isFound, "select count(1)!=0 as isFound from pg_database where datname=$1", switchDbName); err != nil {
		log.Println(err)
	} else {
		if !isFound { //create db.
			createdb := fmt.Sprintf(DB_CREATE, switchDbName, switchDbUser)
			auth := fmt.Sprintf(DBUSER_AUTH, switchDbName, switchDbUser)
			db.MustExec(createdb)
			db.MustExec(auth)
		}
	}
	//return switch db
	switchstr := fmt.Sprintf(`user=%s password=%s host=%s dbname=%s`, switchDbUser, switchDbPassword, switchDbhost, switchDbName)
	return pgsqlOpen(switchstr)
}

func pgsqlGetDatabaseOpenConnections(pgdb *sqlx.DB, dbname string) int {
	var connections int
	query := fmt.Sprintf(`select count(*) as connections from pg_stat_activity where datname = '%s';`, dbname)
	if err := pgdb.Get(&connections, query); err != nil {
		log.Println(err)
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
				log.Println(err)
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
		log.Println(err)
	} else {
		if !isFound { //create user.
			createuser := fmt.Sprintf(USER_CREATE, gdbuser, gdbpassword)
			db.MustExec(createuser)
		}
	}
	if err = db.Get(&isFound, "select count(1)!=0 as isFound from pg_database where datname=$1", gdbname); err != nil {
		log.Println(err)
	} else {
		if !isFound { //create db.
			createdb := fmt.Sprintf(DB_CREATE, gdbname, gdbuser)
			auth := fmt.Sprintf(DBUSER_AUTH, gdbname, gdbuser)
			db.MustExec(createdb)
			db.MustExec(auth)
		}
	}
	//return gateway db
	gatewaystr := fmt.Sprintf(`user=%s password=%s host=%s dbname=%s`, gdbuser, gdbpassword, gdbhost, gdbname)
	return pgsqlOpen(gatewaystr)
}

func pgsqlInitGatewayTables(db *sqlx.DB) {
	var err error
	var isFound bool
	ip := viper.GetString(`switch.vars.ipv4`)
	tablePrefix := viper.GetString(`gateway.db.tableprefix`)

	//table confs define gateway response content.
	tableConfs := fmt.Sprintf(`%s_confs`, tablePrefix)
	if err = db.Get(&isFound, "select count(1)!=0 as isFound from pg_tables where tablename =$1", tableConfs); err != nil {
		log.Println(err)
	} else {
		if !isFound {
			sql := fmt.Sprintf(CONFS, tableConfs, tableConfs)
			db.MustExec(sql)
			//insert DEFAULT_CONFS 
		}
	}
	//table accounts
	tableAccounts := fmt.Sprintf(`%s_accounts`, tablePrefix)
	if err = db.Get(&isFound, "select count(1)!=0 as isFound from pg_tables where tablename =$1", tableAccounts); err != nil {
		log.Println(err)
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
	tableGateways := fmt.Sprintf(`%s_gateways`, tablePrefix)
	if err = db.Get(&isFound, "select count(1)!=0 as isFound from pg_tables where tablename =$1", tableGateways); err != nil {
		log.Println(err)
	} else {
		if !isFound {
			sql := fmt.Sprintf(GATEWAYS, tableGateways, tableGateways)
			gatewaysql := fmt.Sprintf(DEFAULT_GATEWAYS, tableGateways)
			db.MustExec(sql)
			db.MustExec(gatewaysql)
		}
	}
	//table e164s
	tableE164s := fmt.Sprintf(`%s_e164s`, tablePrefix)
	if err = db.Get(&isFound, "select count(1)!=0 as isFound from pg_tables where tablename =$1", tableE164s); err != nil {
		log.Println(err)
	} else {
		if !isFound {
			sql := fmt.Sprintf(E164S, tableE164s, tableE164s)
			e164sql := fmt.Sprintf(DEFAULT_E164S, tableE164s)
			db.MustExec(sql)
			db.MustExec(e164sql)
		}
	}
	//table acce164
	tableAcce164 := fmt.Sprintf(`%s_acce164`, tablePrefix)
	if err = db.Get(&isFound, "select count(1)!=0 as isFound from pg_tables where tablename =$1", tableAcce164); err != nil {
		log.Println(err)
	} else {
		if !isFound {
			sql := fmt.Sprintf(ACCE164, tableAcce164, tableAcce164, tableAcce164, tableAccounts, tableAcce164, tableGateways, tableAcce164, tableE164s)
			acce164sql := fmt.Sprintf(DEFAULT_ACCE164, tableAcce164)
			db.MustExec(sql)
			db.MustExec(acce164sql)
		}
	}
	//table fifos
	tableFifo := fmt.Sprintf(`%s_fifos`, tablePrefix)
	if err = db.Get(&isFound, "select count(1)!=0 as isFound from pg_tables where tablename =$1", tableFifo); err != nil {
		log.Println(err)
	} else {
		if !isFound {
			sql := fmt.Sprintf(FIFOS, tableFifo, tableFifo)
			fifosql := fmt.Sprintf(DEFAULT_FIFOS, tableFifo)
			db.MustExec(sql)
			db.MustExec(fifosql)
		}
	}
	//table fifomembers
	tableFifomembers := fmt.Sprintf(`%s_fifomembers`, tablePrefix)
	if err = db.Get(&isFound, "select count(1)!=0 as isFound from pg_tables where tablename =$1", tableFifomembers); err != nil {
		log.Println(err)
	} else {
		if !isFound {
			sql := fmt.Sprintf(FIFOMEMBERS, tableFifomembers, tableFifomembers, tableFifomembers, tableFifo)
			fifomembersql := fmt.Sprintf(DEFAULT_FIFOMEMBERS, tableFifomembers)
			db.MustExec(sql)
			db.MustExec(fifomembersql)
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
		log.Println(err)
	} else {
		if !isFound { //create user.
			createuser := fmt.Sprintf(USER_CREATE, sdbuser, sdbpassword)
			db.MustExec(createuser)
		}
	}
	if err = db.Get(&isFound, "select count(1)!=0 as isFound from pg_database where datname=$1", sdbname); err != nil {
		log.Println(err)
	} else {
		if !isFound { //create db.
			createdb := fmt.Sprintf(DB_CREATE, sdbname, sdbuser)
			auth := fmt.Sprintf(DBUSER_AUTH, sdbname, sdbuser)
			db.MustExec(createdb)
			db.MustExec(auth)
		}
	}
	//return server db
	serverstr := fmt.Sprintf(`user=%s password=%s host=%s dbname=%s`, sdbuser, sdbpassword, sdbhost, sdbname)
	return pgsqlOpen(serverstr)
}

func pgsqlInitServerTables(db *sqlx.DB) {
	var err error
	var isFound bool
	tablePrefix := viper.GetString(`server.db.tableprefix`)

	//table backlists
	tableBlacklists := fmt.Sprintf(`%s_blacklists`, tablePrefix)
	if err = db.Get(&isFound, "select count(1)!=0 as isFound from pg_tables where tablename =$1", tableBlacklists); err != nil {
		log.Println(err)
	} else {
		if !isFound {
			sql := fmt.Sprintf(BLACKLISTS, tableBlacklists, tableBlacklists)
			insersql := fmt.Sprintf(DEFAULT_BLACKLISTS, tableBlacklists)
			db.MustExec(sql)
			db.MustExec(insersql)
		}
	}
	//table bgjobs
	tableBgjobs := fmt.Sprintf(`%s_bgjobs`, tablePrefix)
	if err = db.Get(&isFound, "select count(1)!=0 as isFound from pg_tables where tablename =$1", tableBgjobs); err != nil {
		log.Println(err)
	} else {
		if !isFound {
			sql := fmt.Sprintf(BGJOBS, tableBgjobs, tableBgjobs)
			db.MustExec(sql)
			//db.MustExec(`inserted by fs server --run`)
		}
	}
}
