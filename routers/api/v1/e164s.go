package v1

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/bob1118/fs/db"
	"github.com/bob1118/fs/ec"
	"github.com/bob1118/fs/utils"
	"github.com/gin-gonic/gin"
)

// GetE164s function return e164s by condition.
//
// request: GET /api/v1/e164s?uuid=xxx&gname=xxx&number=xxx&enable=true
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
func GetE164s(c *gin.Context) {
	rtmsg := ``
	rtcode := ec.SUCCESS
	condition := "true"
	data := make(map[string]interface{})

	if uuid := c.Query("uuid"); len(uuid) > 0 {
		condition += fmt.Sprintf(" and e164_uuid='%s'", uuid)
	}
	if gname := c.Query("gname"); len(gname) > 0 {
		condition += fmt.Sprintf(" and gateway_name='%s'", gname)
	}
	if number := c.Query("number"); len(number) > 0 {
		condition += fmt.Sprintf(" and e164_number='%s'", number)
	}
	if enable := c.Query("enable"); len(enable) > 0 {
		if utils.IsEqual(enable, `true`) {
			condition += " and e164_enable is true"
		} else {
			condition += " and e164_enable is false"
		}
	}

	if rtcode == ec.SUCCESS {
		if gateways, err := db.SelectE164sWithCondition(condition); err != nil {
			rtcode = ec.ERROR_DATABSE_QUERY
			rtmsg = err.Error()
		} else {
			data["len"] = len(gateways)
			data["lists"] = gateways
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": gin.H{"rtcode": rtcode, "rtmsg": rtmsg},
		"data": data,
	})
}

// PostE164 function.
//
// request: POST /api/v1/e164, a E164{} json.
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
func PostE164(c *gin.Context) {
	rtmsg := ``
	rtcode := ec.SUCCESS
	e164 := db.E164{}
	e164s := make([]db.E164, 0)
	data := make(map[string]interface{})

	if c.ContentType() != gin.MIMEJSON {
		rtcode = ec.ERROR_HTTP_REQUEST_CONTENTTYPE
	} else {
		if err := c.BindJSON(&e164); err != nil {
			rtcode = ec.ERROR_HTTP_REQUEST_CONTEXTBINDJSON
			rtmsg = err.Error()
		} else {
			if len(e164.Enumber) == 0 {
				rtcode = ec.ERROR_HTTP_REQUEST_JSONITEMNULL
			}
		}
	}

	if rtcode == ec.SUCCESS {
		e164s = append(e164s, e164)
		if rte164s, err := db.InsertE164s(e164s); err != nil {
			rtcode = ec.ERROR_DATABSE_INSERT
			rtmsg = err.Error()
		} else {
			data["len"] = len(rte164s)
			data["lists"] = rte164s
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": gin.H{"rtcode": rtcode, "rtmsg": rtmsg},
		"data": data,
	})
}

// PostE164 function.
//
// request: POST /api/v1/e164s?gname=xxx&numberprefix=xxx&numberstart=xxx&numberend=xxx
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
func PostE164s(c *gin.Context) {
	rtmsg := ``
	rtcode := ec.SUCCESS
	e164 := db.E164{}
	e164s := make([]db.E164, 0)
	data := make(map[string]interface{})

	if c.ContentType() != gin.MIMEJSON {
		rtcode = ec.ERROR_HTTP_REQUEST_CONTENTTYPE
	} else {
		gname := c.Query("gname")
		numberPrefix := c.Query("numberprefix")
		if numberStart := c.Query("numberstart"); len(numberStart) > 0 {
			if start, err := strconv.Atoi(numberStart); err != nil {
				rtcode = ec.ERROR_HTTP_REQUEST_URLQUERYATOI
				rtmsg = err.Error()
			} else {
				if numberEnd := c.Query("numberend"); len(numberEnd) > 0 {
					if end, err := strconv.Atoi(numberEnd); err != nil {
						rtcode = ec.ERROR_HTTP_REQUEST_URLQUERYATOI
						rtmsg = err.Error()
					} else {
						for index := start; index <= end; index++ {
							e164.Gname = gname
							e164.Enumber = fmt.Sprintf("%s%d", numberPrefix, index)
							e164s = append(e164s, e164)
						}
					}
				} else {
					rtcode = ec.ERROR_HTTP_REQUEST_URLQUERYKEY
				}
			}
		} else {
			rtcode = ec.ERROR_HTTP_REQUEST_URLQUERYKEY
		}
	}

	if rtcode == ec.SUCCESS {
		if rte164s, err := db.InsertE164s(e164s); err != nil {
			rtcode = ec.ERROR_DATABSE_INSERT
			rtmsg = err.Error()
		} else {
			data["len"] = len(rte164s)
			data["lists"] = rte164s
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": gin.H{"rtcode": rtcode, "rtmsg": rtmsg},
		"data": data,
	})
}

// PutE164 function.
//
// request: PUT /api/v1/e164/:uuid, a E164{} json.
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
func PutE164(c *gin.Context) {
	rtmsg := ``
	rtcode := ec.SUCCESS
	e164 := db.E164{}
	e164s := make([]db.E164, 0)
	data := make(map[string]interface{})

	if c.ContentType() != gin.MIMEJSON {
		rtcode = ec.ERROR_HTTP_REQUEST_CONTENTTYPE
	} else {
		if uuid := c.Param("uuid"); len(uuid) == 0 {
			rtcode = ec.ERROR_HTTP_REQUEST_URLPARAM
		} else {
			if err := c.BindJSON(&e164); err != nil {
				rtcode = ec.ERROR_HTTP_REQUEST_CONTEXTBINDJSON
				rtmsg = err.Error()
			} else {
				if rte164, err := db.UpdateE164sE164(uuid, e164); err != nil {
					rtcode = ec.ERROR_DATABSE_UPDATE
					rtmsg = err.Error()
				} else {
					e164s = append(e164s, rte164)
					data["len"] = len(e164s)
					data["lists"] = e164s
				}
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": gin.H{"rtcode": rtcode, "rtmsg": rtmsg},
		"data": data,
	})
}

// DeleteE164 function.
//
// request: DELETE /api/v1/e164/:uuid
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
func DeleteE164(c *gin.Context) {
	rtmsg := ``
	rtcode := ec.SUCCESS
	e164s := make([]db.E164, 0)
	data := make(map[string]interface{})

	if c.ContentType() != gin.MIMEJSON {
		rtcode = ec.ERROR_HTTP_REQUEST_CONTENTTYPE
	} else {
		if uuid := c.Param("uuid"); len(uuid) == 0 {
			rtcode = ec.ERROR_HTTP_REQUEST_URLPARAM
		} else {
			{
				if rte164, err := db.DeleteE164sE164(uuid); err != nil {
					rtcode = ec.ERROR_DATABSE_DELETE
					rtmsg = err.Error()
				} else {
					e164s = append(e164s, rte164)
					data["len"] = e164s
					data["lists"] = e164s
				}
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": gin.H{"rtcode": rtcode, "rtmsg": rtmsg},
		"data": data,
	})
}
