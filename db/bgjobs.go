// const BGJOBS = `
// CREATE TABLE IF NOT EXISTS %s (
// 	job_uuid uuid NOT NULL,
// 	job_cmd varchar,
// 	job_cmdarg varchar,
// 	job_content varchar,
// 	CONSTRAINT bgjobs_pkey PRIMARY KEY (job_uuid)
// );
// COMMENT ON TABLE %s IS 'eslclient execute bgapi command then receive EVENT BACKGROUND_JOB ';
// `

package db

import (
	"errors"
	"fmt"
)

type Bgjob struct {
	Juuid    string `db:"job_uuid" json:"uuid"`
	Jcmd     string `db:"job_cmd" json:"cmd"`
	Jcmdarg  string `db:"job_cmdarg" json:"cmdarg"`
	Jcontent string `db:"job_content" json:"content"`
}

func CreateBgjob(in *Bgjob) error {
	var err error
	job := in
	if len(job.Juuid) == 0 {
		err = errors.New("uuid not null")
	} else {
		q := fmt.Sprintf("insert into %sbgjobs(job_cmd,job_cmdarg,job_content)values('%s','%s','%s')", GetTablesServerPrifex(), job.Jcmd, job.Jcmdarg, job.Jcontent)
		GetServerdb().MustExec(q)
	}
	return err
}

// SelectBgjobsWithCondition
func SelectBgjobsWithCondition(condition string) ([]Bgjob, error) {
	var bgjobs = []Bgjob{}
	q := fmt.Sprintf(`select * from %sbgjobs where %s`, GetTablesServerPrifex(), condition)
	err := GetServerdb().Select(&bgjobs, q)
	return bgjobs, err
}
