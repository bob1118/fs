// // CDR_LEG mod_odbc_cdr table define
// const CDR_LEG = `
// CREATE TABLE IF NOT EXISTS %s (
//
//	uuid varchar NOT NULL,
//	calluuid varchar NOT NULL,
//	otheruuid varchar NOT NULL DEFAULT '',
//	name varchar NOT NULL DEFAULT '',
//	direction varchar NOT NULL DEFAULT '',
//	sofiaprofile varchar NOT NULL DEFAULT '',
//	domain varchar NOT NULL DEFAULT '',
//	sipprofile varchar NOT NULL DEFAULT '',
//	gateway varchar NOT NULL DEFAULT '',
//	ani varchar NOT NULL DEFAULT '',
//	destination varchar NOT NULL DEFAULT '',
//	calleridname varchar NOT NULL DEFAULT '',
//	calleridnumber varchar NOT NULL DEFAULT '',
//	calleeidname varchar NOT NULL DEFAULT '',
//	calleeidnumber varchar NOT NULL DEFAULT '',
//	app varchar NOT NULL DEFAULT '',
//	appdata varchar NOT NULL DEFAULT '',
//	dialstatus varchar NOT NULL DEFAULT '',
//	cause varchar NOT NULL DEFAULT '',
//	q850 varchar NOT NULL DEFAULT '',
//	disposition varchar NOT NULL DEFAULT '',
//	protocause varchar NOT NULL DEFAULT '',
//	phrase varchar NOT NULL DEFAULT '',
//	startepoch varchar NOT NULL DEFAULT '',
//	answerepoch varchar NOT NULL DEFAULT '',
//	endepoch varchar NOT NULL DEFAULT '',
//	waitsec varchar NOT NULL DEFAULT '',
//	billsec varchar NOT NULL DEFAULT '',
//	duration varchar NOT NULL DEFAULT ''
//
// );
// `
package db

import (
	"fmt"

	"github.com/spf13/viper"
)

type CDRLEG struct {
	UUID           string `db:"uuid" json:"uuid"`
	OtherUUID      string `db:"otheruuid" json:"otheruuid"`
	BondUUID       string `db:"bonduuid" json:"bonduuid"`
	Name           string `db:"name" json:"name"`
	Direction      string `db:"direction" json:"direction"`
	Sofiaprofile   string `db:"sofiaprofile" json:"sofiaprofile"`
	Indomain       string `db:"indomain" json:"indomain"`
	Ingateway      string `db:"ingateway" json:"ingateway"`
	Outdomain      string `db:"outdomain" json:"outdomain"`
	Outgateway     string `db:"outgateway" json:"outgateway"`
	Ani            string `db:"ani" json:"ani"`
	Destination    string `db:"destination" json:"destination"`
	Calleridname   string `db:"calleridname" json:"calleridname"`
	Calleridnumber string `db:"calleridnumber" json:"calleridnumber"`
	Calleeidname   string `db:"calleeidname" json:"calleeidname"`
	Calleeidnumber string `db:"calleeidnumber" json:"calleeidnumber"`
	App            string `db:"app" json:"app"`
	Appdata        string `db:"appdata" json:"appdata"`
	Appdialstatus  string `db:"dialstatus" json:"dialstatus"`
	Cause          string `db:"cause" json:"cause"`
	Q850           string `db:"q850" json:"q850"`
	Disposition    string `db:"disposition" json:"disposition"`
	Protocause     string `db:"protocause" json:"protocause"`
	Phrase         string `db:"phrase" json:"phrase"`
	Startepoch     string `db:"startepoch" json:"startepoch"`
	Answerepoch    string `db:"answerepoch" json:"answerepoch"`
	Endepoch       string `db:"endepoch" json:"endepoch"`
	Waitsec        string `db:"waitsec" json:"waitsec"`
	Billsec        string `db:"billsec" json:"billsec"`
	Duration       string `db:"duration" json:"duration"`
}

// CreateCdrAleg
func CreateCdrAleg(in *CDRLEG) (e error) {
	return InsertCallDetailRecord(false, in)
}

// CreateCdrBleg
func CreateCdrBleg(in *CDRLEG) (e error) {
	return InsertCallDetailRecord(true, in)
}

// InsertCallDetailRecord
func InsertCallDetailRecord(isbleg bool, in *CDRLEG) error {
	var tablename string

	if isbleg {
		tablename = viper.GetString(`switch.cdr.a-leg`)
	} else {
		tablename = viper.GetString(`switch.cdr.b-leg`)
	}
	l := in
	q := fmt.Sprintf(`INSERT INTO %s(
		uuid,otheruuid,bonduuid,name,direction,sofiaprofile,indomain,ingateway,outdomain,outgateway,
		ani,destination,calleridname,calleridnumber,calleeidname,calleeidnumber,
		app,appdata,dialstatus,cause,q850,disposition,protocause,phrase,
		startepoch,answerepoch,endepoch,waitsec,billsec,duration)
		VALUES('%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s');`,
		tablename,
		l.UUID, l.OtherUUID, l.BondUUID, l.Name, l.Direction, l.Sofiaprofile, l.Indomain, l.Ingateway, l.Outdomain, l.Outgateway,
		l.Ani, l.Destination, l.Calleridname, l.Calleridnumber, l.Calleeidname, l.Calleeidnumber,
		l.App, l.Appdata, l.Appdialstatus, l.Cause, l.Q850, l.Disposition, l.Protocause, l.Phrase,
		l.Startepoch, l.Answerepoch, l.Endepoch, l.Waitsec, l.Billsec, l.Duration)
	_, err := GetSwitchdb().Exec(q)
	return err
}

// SelectCallDetailRecordsWithCondition
func SelectCallDetailRecordsWithCondition(condition string) ([]CDRLEG, error) {
	cdrs := []CDRLEG{}
	alegname := viper.GetString(`switch.cdr.a-leg`)
	//blegname := viper.GetString(`switch.cdr.b-leg`)

	q := fmt.Sprintf("select * from %s where %s", alegname, condition)
	err := GetSwitchdb().Select(&cdrs, q)
	return cdrs, err
}
