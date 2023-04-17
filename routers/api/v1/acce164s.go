package v1

import (
	"fmt"
	"net/http"

	"github.com/bob1118/fs/db"
	"github.com/bob1118/fs/ec"
	"github.com/bob1118/fs/utils"
	"github.com/gin-gonic/gin"
)

// GetAcce164s function.
//
// request: GET /api/v1/acce164s?uuid=xxx&aid=xxx&adomain=xxx&gname=xxx&enumber=xxx&isdefault=true
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
func GetAcce164s(c *gin.Context) {
	rtmsg := ``
	rtcode := ec.SUCCESS
	condition := "true"
	data := make(map[string]interface{})

	if uuid := c.Query("uuid"); len(uuid) > 0 {
		condition += fmt.Sprintf(" and acce164_uuid='%s'", uuid)
	}
	if id := c.Query("aid"); len(id) > 0 {
		condition += fmt.Sprintf(" and account_id='%s'", id)
	}
	if domain := c.Query("adomain"); len(domain) > 0 {
		condition += fmt.Sprintf(" and account_domain='%s'", domain)
	}
	if gname := c.Query("gname"); len(gname) > 0 {
		condition += fmt.Sprintf(" and gateway_name='%s'", gname)
	}
	if number := c.Query("enumber"); len(number) > 0 {
		condition += fmt.Sprintf(" and e164_number='%s'", number)
	}
	if isdefault := c.Query("isdefault"); len(isdefault) > 0 {
		if utils.IsEqual(isdefault, `true`) {
			condition += " and acce164_isdefault is true"
		} else {
			condition += " and acce164_isdefault is false"
		}
	}

	if rtcode == ec.SUCCESS {
		if acce164s, err := db.SelectAcce164sWithCondition(condition); err != nil {
			rtcode = ec.ERROR_DATABSE_QUERY
			rtmsg = err.Error()
		} else {
			data["len"] = len(acce164s)
			data["lists"] = acce164s
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": gin.H{"rtcode": rtcode, "rtmsg": rtmsg},
		"data": data,
	})
}

// PostAcce164 function.
//
// request: POST /api/v1/acce164, a ACCE164{} json.
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
func PostAcce164(c *gin.Context) {
	rtmsg := ``
	rtcode := ec.SUCCESS
	acce164 := db.ACCE164{}
	acce164s := make([]db.ACCE164, 0)
	data := make(map[string]interface{})

	if c.ContentType() != gin.MIMEJSON {
		rtcode = ec.ERROR_HTTP_REQUEST_CONTENTTYPE
	} else {
		if err := c.BindJSON(&acce164); err != nil {
			rtcode = ec.ERROR_HTTP_REQUEST_CONTEXTBINDJSON
			rtmsg = err.Error()
		} else {
			if len(acce164.Aid) == 0 || len(acce164.Adomain) == 0 || len(acce164.Enumber) == 0 || len(acce164.Gname) == 0 {
				rtcode = ec.ERROR_HTTP_REQUEST_JSONITEMNULL
			}
		}
	}

	if rtcode == ec.SUCCESS {
		acce164s = append(acce164s, acce164)
		if rte164s, err := db.InsertAcce164s(acce164s); err != nil {
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

// PutAcce164 function.
//
// request: PUT /api/v1/acce164/:uuid, a ACCE164{} json.
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
func PutAcce164(c *gin.Context) {
	rtmsg := ``
	rtcode := ec.SUCCESS
	acce164 := db.ACCE164{}
	acce164s := make([]db.ACCE164, 0)
	data := make(map[string]interface{})

	if c.ContentType() != gin.MIMEJSON {
		rtcode = ec.ERROR_HTTP_REQUEST_CONTENTTYPE
	} else {
		if uuid := c.Param("uuid"); len(uuid) == 0 {
			rtcode = ec.ERROR_HTTP_REQUEST_URLPARAM
		} else {
			if err := c.BindJSON(&acce164); err != nil {
				rtcode = ec.ERROR_HTTP_REQUEST_CONTEXTBINDJSON
				rtmsg = err.Error()
			} else {
				if rtacce164, err := db.UpdateAcce164sAcce164(uuid, acce164); err != nil {
					rtcode = ec.ERROR_DATABSE_UPDATE
					rtmsg = err.Error()
				} else {
					acce164s = append(acce164s, rtacce164)
					data["len"] = len(acce164s)
					data["lists"] = acce164s
				}
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": gin.H{"rtcode": rtcode, "rtmsg": rtmsg},
		"data": data,
	})
}

// DeleteAcce164 function.
//
// request: DELETE /api/v1/acce164/:uuid
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
func DeleteAcce164(c *gin.Context) {
	rtmsg := ``
	rtcode := ec.SUCCESS
	acce164s := make([]db.ACCE164, 0)
	data := make(map[string]interface{})

	if c.ContentType() != gin.MIMEJSON {
		rtcode = ec.ERROR_HTTP_REQUEST_CONTENTTYPE
	} else {
		if uuid := c.Param("uuid"); len(uuid) == 0 {
			rtcode = ec.ERROR_HTTP_REQUEST_URLPARAM
		} else {
			{
				if rtacce164, err := db.DeleteAcce164sAcce164(uuid); err != nil {
					rtcode = ec.ERROR_DATABSE_DELETE
					rtmsg = err.Error()
				} else {
					acce164s = append(acce164s, rtacce164)
					data["len"] = acce164s
					data["lists"] = acce164s
				}
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": gin.H{"rtcode": rtcode, "rtmsg": rtmsg},
		"data": data,
	})
}
