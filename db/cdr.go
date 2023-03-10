package db

import (
	"fmt"
)

type CDRLEG struct {
	UUID           string `db:"uuid" json:"uuid"`
	OtherUUID      string `db:"otheruuid" json:"otheruuid"`
	OtherType      string `db:"othertype" json:"othertype"`
	Name           string `db:"name" json:"name"`
	Profile        string `db:"profile" json:"profile"`
	Direction      string `db:"direction" json:"direction"`
	Domain         string `db:"domain" json:"domain"`
	Gateway        string `db:"gateway" json:"gateway"`
	Calleridname   string `db:"calleridname" json:"calleridname"`
	Calleridnumber string `db:"calleridnumber" json:"calleridnumber"`
	Calleeidname   string `db:"calleeidname" json:"calleeidname"`
	Calleeidnumber string `db:"calleeidnumber" json:"calleeidnumber"`
	Destination    string `db:"destination" json:"destination"`
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
	var err error

	a := in
	q := fmt.Sprintf(`INSERT INTO cdr_aleg 
	(uuid, otheruuid, othertype, "name", profile, direction, "domain", gateway, 
		calleridname, calleridnumber, calleeidname, calleeidnumber, destination, app, appdata, dialstatus, 
		cause, q850, disposition, protocause, phrase, startepoch, answerepoch, endepoch, waitsec, billsec, duration) 
		VALUES('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s');`,
		a.UUID, a.OtherUUID, a.OtherType, a.Name, a.Profile, a.Direction, a.Domain, a.Gateway,
		a.Calleridname, a.Calleridnumber, a.Calleeidname, a.Calleeidnumber, a.Destination, a.App, a.Appdata, a.Appdialstatus,
		a.Cause, a.Q850, a.Disposition, a.Protocause, a.Phrase, a.Startepoch, a.Answerepoch, a.Endepoch, a.Waitsec, a.Billsec, a.Duration)
	_, err = GetServerdb().Exec(q)
	return err
}

// CreateCdrBleg
func CreateCdrBleg(in *CDRLEG) (e error) {
	var err error
	b := in
	q := fmt.Sprintf(`INSERT INTO cdr_bleg 
	(uuid, otheruuid, othertype, "name", profile, direction, "domain", gateway, 
		calleridname, calleridnumber, calleeidname, calleeidnumber, destination, app, appdata, dialstatus, 
		cause, q850, disposition, protocause, phrase, startepoch, answerepoch, endepoch, waitsec, billsec, duration) 
		VALUES('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s');`,
		b.UUID, b.OtherUUID, b.OtherType, b.Name, b.Profile, b.Direction, b.Domain, b.Gateway,
		b.Calleridname, b.Calleridnumber, b.Calleeidname, b.Calleeidnumber, b.Destination, b.App, b.Appdata, b.Appdialstatus,
		b.Cause, b.Q850, b.Disposition, b.Protocause, b.Phrase, b.Startepoch, b.Answerepoch, b.Endepoch, b.Waitsec, b.Billsec, b.Duration)
	_, err = GetServerdb().Exec(q)
	return err
}
