package v1

import (
	"fmt"
	"net/http"

	"github.com/bob1118/fs/db"
	"github.com/bob1118/fs/ec"
	"github.com/gin-gonic/gin"
)

// GetCallDetailRecords function return cdrs by condition.
//
// request: GET /api/v1/cdrs?uuid=xxx
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
func GetCallDetailRecords(c *gin.Context) {
	rtmsg := ``
	rtcode := ec.SUCCESS
	condition := "true"
	data := make(map[string]interface{})

	if uuid := c.Query("uuid"); len(uuid) > 0 {
		condition += fmt.Sprintf(" and uuid='%s'", uuid)
	}
	if rtcode == ec.SUCCESS {
		if cdrs, err := db.SelectCallDetailRecordsWithCondition(condition); err != nil {
			rtcode = ec.ERROR_DATABSE_QUERY
			rtmsg = err.Error()
		} else {
			data["len"] = len(cdrs)
			data["lists"] = cdrs
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": gin.H{"rtcode": rtcode, "rtmsg": rtmsg},
		"data": data,
	})
}
