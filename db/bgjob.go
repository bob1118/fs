package db

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

type Bgjob struct {
	Juuid    string `db:"job_uuid" json:"uuid"`
	Jcmd     string `db:"job_cmd" json:"cmd"`
	Jcmdarg  string `db:"job_cmdarg" json:"cmdarg"`
	Jcontent string `db:"job_content" json:"content"`
}

func CreateBgjob(in *Bgjob) error {
	var err error
	var realtableprefix string
	tableprefix := viper.GetString(`server.db.tableprefix`)
	if len(tableprefix) > 0 {
		realtableprefix = fmt.Sprintf(`%s_`, tableprefix)
	} else {
		realtableprefix = tableprefix
	}
	job := in
	if len(job.Juuid) == 0 {
		err = errors.New("uuid not null")
	} else {
		query := fmt.Sprintf("insert into %sbgjobs(job_uuid,job_cmd,job_cmdarg,job_content)values('%s','%s','%s','%s')", realtableprefix, job.Juuid, job.Jcmd, job.Jcmdarg, job.Jcontent)
		GetServerdb().MustExec(query)
	}

	return err
}
