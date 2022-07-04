package db

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

var pgdb, switchdb, gatewaydb, serverdb *sqlx.DB

func init() {}

func pgsqlOpen(strcon string) (*sqlx.DB, error) { return sqlx.Connect(`postgres`, strcon) }

func pgsqlClose(db *sqlx.DB) { db.Close() }

func Initdb() {
	var err error
	strcon := fmt.Sprintf(`user=%s password=%s host=%s dbname=%s`,
		viper.GetString(`postgres.user`), viper.GetString(`postgres.password`), viper.GetString(`postgres.host`), viper.GetString(`postgres.name`))
	pgdb, err = pgsqlOpen(strcon)
	if err != nil {
		log.Println(err)
	} else {
		// initSwitchdb
		if switchdb, err = initSwitchdb(pgdb); err != nil {
			log.Println(err)
		} else {
			initSwitchMododbccdrTables(switchdb)
		}
		//initGatewaydb
		if gatewaydb, err = initGatewaydb(pgdb); err != nil {
			log.Println(err)
		} else {
			initGatewaydbData(gatewaydb)
		}
		//initServerdb
		if serverdb, err = initServerdb(pgdb); err != nil {
			log.Println(err)
		} else {
			initServerdbData(serverdb)
		}
		pgsqlClose(pgdb)
	}
}

func initSwitchdb(db *sqlx.DB) (*sqlx.DB, error) {
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
	}
	//return switch db
	switchstr := fmt.Sprintf(`user=%s password=%s host=%s dbname=%s`, switchDbUser, switchDbPassword, switchDbhost, switchDbName)
	return pgsqlOpen(switchstr)
}

func destroySwitchdb(db *sqlx.DB) {
	switchDbUser := viper.GetString(`switch.db.user`)
	switchDbName := viper.GetString(`switch.db.name`)
	if db.Stats().InUse > 0 {
		log.Println(db.Stats())
	} else {
		dropsql := fmt.Sprintf(DATABASE_USER_DROP, switchDbName, switchDbUser)
		db.MustExec(dropsql)
	}
}

func initSwitchMododbccdrTables(db *sqlx.DB) {
	var err error
	var isFound bool
	tables := viper.GetStringSlice(`switch.cdr.tables`)
	for _, table := range tables {
		if err = db.Get(&isFound, "select count(1)!=0 as isFound from pg_tables where tablename =$1", table); err != nil {
			log.Println(err)
		} else {
			if !isFound {
				cdrleg := fmt.Sprintf(CDR_LEG, table)
				db.MustExec(cdrleg)
			}
		}
	}
}

func initGatewaydb(db *sqlx.DB) (*sqlx.DB, error) {

	return pgsqlOpen(``)
}

func initGatewaydbData(db *sqlx.DB) {}

func initServerdb(db *sqlx.DB) (*sqlx.DB, error) {
	return pgsqlOpen(``)
}

func initServerdbData(db *sqlx.DB) {}
