package v1

import (
	"fmt"
	"net/http"

	"github.com/bob1118/fs/db"
	"github.com/bob1118/fs/ec"
	"github.com/gin-gonic/gin"
)

// GetBgjobs function.
//
// request: GET /api/v1/bgjobs?uuid=xxx&cmd=xxx&cmdarg=xxx
//
// response: json
//
//	{
//		code:{
//			rtcode: rtcode, //return ec.SUCCESS or not.
//			rtmsg: rtmsg	//return error message while some error occured.
//		},
//		data:{
//			len: len(slice),
//			lists:{slice[0],slice[1], ...}
//		}
//	}
func GetBgjobs(c *gin.Context) {
	rtmsg := ``
	rtcode := ec.SUCCESS
	condition := "true"
	data := make(map[string]interface{})

	if uuid := c.Query("uuid"); len(uuid) > 0 {
		condition += fmt.Sprintf(" and job_uuid='%s'", uuid)
	}
	if cmd := c.Query("cmd"); len(cmd) > 0 {
		condition += fmt.Sprintf(" and job_cmd='%s'", cmd)
	}
	if cmdarg := c.Query("cmdarg"); len(cmdarg) > 0 {
		condition += fmt.Sprintf(" and job_cmdarg='%s'", cmdarg)
	}

	if rtcode == ec.SUCCESS {
		if bgjobs, err := db.SelectBgjobsWithCondition(condition); err != nil {
			rtcode = ec.ERROR_DATABSE_QUERY
			rtmsg = err.Error()
		} else {
			data["len"] = len(bgjobs)
			data["lists"] = bgjobs
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": gin.H{"rtcode": rtcode, "rtmsg": rtmsg},
		"data": data,
	})
}
