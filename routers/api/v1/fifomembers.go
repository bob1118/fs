package v1

import (
	"fmt"
	"net/http"

	"github.com/bob1118/fs/db"
	"github.com/bob1118/fs/ec"
	"github.com/gin-gonic/gin"
)

// GetFifomembers function.
//
// request: GET /api/v1/fifomembers?uuid=xxx&fname=xxx&mstring=xxx
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
func GetFifomembers(c *gin.Context) {
	rtmsg := ``
	rtcode := ec.SUCCESS
	condition := "true"
	data := make(map[string]interface{})

	if uuid := c.Query("uuid"); len(uuid) > 0 {
		condition += fmt.Sprintf(" and fifomember_uuid='%s'", uuid)
	}
	if fname := c.Query("fname"); len(fname) > 0 {
		condition += fmt.Sprintf(" and fifo_name='%s'", fname)
	}
	if mstring := c.Query("mstring"); len(mstring) > 0 {
		condition += fmt.Sprintf(" and member_string='%s'", mstring)
	}

	if rtcode == ec.SUCCESS {
		if fifomembers, err := db.SelectFifomembersWithCondition(condition); err != nil {
			rtcode = ec.ERROR_DATABSE_QUERY
			rtmsg = err.Error()
		} else {
			data["len"] = len(fifomembers)
			data["lists"] = fifomembers
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": gin.H{"rtcode": rtcode, "rtmsg": rtmsg},
		"data": data,
	})
}

// PostFifomember function.
//
// request: POST /api/v1/fifo, a FifoMember{} json.
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
func PostFifomember(c *gin.Context) {
	rtmsg := ``
	rtcode := ec.SUCCESS
	fifomember := db.FifoMember{}
	fifomembers := make([]db.FifoMember, 0)
	data := make(map[string]interface{})

	if c.ContentType() != gin.MIMEJSON {
		rtcode = ec.ERROR_HTTP_REQUEST_CONTENTTYPE
	} else {
		if err := c.BindJSON(&fifomember); err != nil {
			rtcode = ec.ERROR_HTTP_REQUEST_CONTEXTBINDJSON
			rtmsg = err.Error()
		} else {
			if len(fifomember.Fname) == 0 || len(fifomember.Mstring) == 0 {
				rtcode = ec.ERROR_HTTP_REQUEST_JSONITEMNULL
			}
		}
	}

	if rtcode == ec.SUCCESS {
		fifomembers = append(fifomembers, fifomember)
		if rt, err := db.InsertFifoMembers(fifomembers); err != nil {
			rtcode = ec.ERROR_DATABSE_INSERT
			rtmsg = err.Error()
		} else {
			data["len"] = len(rt)
			data["lists"] = rt
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": gin.H{"rtcode": rtcode, "rtmsg": rtmsg},
		"data": data,
	})
}

// PutFifomember function.
//
// request: PUT /api/v1/fifomember/:uuid, a FifoMember{} json.
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
func PutFifomember(c *gin.Context) {
	rtmsg := ``
	rtcode := ec.SUCCESS
	fifomember := db.FifoMember{}
	fifomembers := make([]db.FifoMember, 0)
	data := make(map[string]interface{})

	if c.ContentType() != gin.MIMEJSON {
		rtcode = ec.ERROR_HTTP_REQUEST_CONTENTTYPE
	} else {
		if uuid := c.Param("uuid"); len(uuid) == 0 {
			rtcode = ec.ERROR_HTTP_REQUEST_URLPARAM
		} else {
			if err := c.BindJSON(&fifomember); err != nil {
				rtcode = ec.ERROR_HTTP_REQUEST_CONTEXTBINDJSON
				rtmsg = err.Error()
			} else {
				if rtfifomember, err := db.UpdateFifomembersFifomember(uuid, fifomember); err != nil {
					rtcode = ec.ERROR_DATABSE_UPDATE
					rtmsg = err.Error()
				} else {
					fifomembers = append(fifomembers, rtfifomember)
					data["len"] = len(fifomembers)
					data["lists"] = fifomembers
				}
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": gin.H{"rtcode": rtcode, "rtmsg": rtmsg},
		"data": data,
	})
}

// DeleteFifomember function.
//
// request: DELETE /api/v1/fifomember/:uuid
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
func DeleteFifomember(c *gin.Context) {
	rtmsg := ``
	rtcode := ec.SUCCESS
	fifomembers := make([]db.FifoMember, 0)
	data := make(map[string]interface{})

	if c.ContentType() != gin.MIMEJSON {
		rtcode = ec.ERROR_HTTP_REQUEST_CONTENTTYPE
	} else {
		if uuid := c.Param("uuid"); len(uuid) == 0 {
			rtcode = ec.ERROR_HTTP_REQUEST_URLPARAM
		} else {
			{
				if rt, err := db.DeleteFifomembersFifomember(uuid); err != nil {
					rtcode = ec.ERROR_DATABSE_DELETE
					rtmsg = err.Error()
				} else {
					fifomembers = append(fifomembers, rt)
					data["len"] = len(fifomembers)
					data["lists"] = fifomembers
				}
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": gin.H{"rtcode": rtcode, "rtmsg": rtmsg},
		"data": data,
	})
}
