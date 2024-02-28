package v1

import (
	"fmt"
	"net/http"

	"github.com/bob1118/fs/db"
	"github.com/bob1118/fs/ec"
	"github.com/bob1118/fs/utils"
	"github.com/gin-gonic/gin"
)

// GetE164accs function.
//
// request: GET /api/v1/e164accs?uuid=xxx&gname=xxx&enumber=xxx&aid=xxx&adomain=xxx&fname=xxx&isfifo=true
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
func GetE164accs(c *gin.Context) {
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
	if fname := c.Query("fname"); len(fname) > 0 {
		condition += fmt.Sprintf(" and fifo_name ='%s'", fname)
	}
	if isdefault := c.Query("isfifo"); len(isdefault) > 0 {
		if utils.IsEqual(isdefault, `true`) {
			condition += " and acce164_isdefault is true"
		} else {
			condition += " and acce164_isdefault is false"
		}
	}

	if rtcode == ec.SUCCESS {
		if e164accs, err := db.SelectE164accsWithCondition(condition); err != nil {
			rtcode = ec.ERROR_DATABSE_QUERY
			rtmsg = err.Error()
		} else {
			data["len"] = len(e164accs)
			data["lists"] = e164accs
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": gin.H{"rtcode": rtcode, "rtmsg": rtmsg},
		"data": data,
	})
}

// PostE164acc function.
//
// request: POST /api/v1/acce164, a E164ACC{} json.
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
func PostE164acc(c *gin.Context) {
	rtmsg := ``
	rtcode := ec.SUCCESS
	e164acc := db.E164ACC{}
	e164accs := make([]db.E164ACC, 0)
	data := make(map[string]interface{})

	if c.ContentType() != gin.MIMEJSON {
		rtcode = ec.ERROR_HTTP_REQUEST_CONTENTTYPE
	} else {
		if err := c.BindJSON(&e164acc); err != nil {
			rtcode = ec.ERROR_HTTP_REQUEST_CONTEXTBINDJSON
			rtmsg = err.Error()
		} else {
			if len(e164acc.Aid) == 0 || len(e164acc.Adomain) == 0 || len(e164acc.Enumber) == 0 || len(e164acc.Gname) == 0 {
				rtcode = ec.ERROR_HTTP_REQUEST_JSONITEMNULL
			}
		}
	}

	if rtcode == ec.SUCCESS {
		e164accs = append(e164accs, e164acc)
		if rte164accs, err := db.InsertE164accs(e164accs); err != nil {
			rtcode = ec.ERROR_DATABSE_INSERT
			rtmsg = err.Error()
		} else {
			data["len"] = len(rte164accs)
			data["lists"] = rte164accs
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": gin.H{"rtcode": rtcode, "rtmsg": rtmsg},
		"data": data,
	})
}

// PutE164acc function.
//
// request: PUT /api/v1/acce164/:uuid, a E164ACC{} json.
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
func PutE164acc(c *gin.Context) {
	rtmsg := ``
	rtcode := ec.SUCCESS
	e164acc := db.E164ACC{}
	e164accs := make([]db.E164ACC, 0)
	data := make(map[string]interface{})

	if c.ContentType() != gin.MIMEJSON {
		rtcode = ec.ERROR_HTTP_REQUEST_CONTENTTYPE
	} else {
		if uuid := c.Param("uuid"); len(uuid) == 0 {
			rtcode = ec.ERROR_HTTP_REQUEST_URLPARAM
		} else {
			if err := c.BindJSON(&e164acc); err != nil {
				rtcode = ec.ERROR_HTTP_REQUEST_CONTEXTBINDJSON
				rtmsg = err.Error()
			} else {
				if rte164acc, err := db.UpdateE164accsE164acc(uuid, e164acc); err != nil {
					rtcode = ec.ERROR_DATABSE_UPDATE
					rtmsg = err.Error()
				} else {
					e164accs = append(e164accs, rte164acc)
					data["len"] = len(e164accs)
					data["lists"] = e164accs
				}
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": gin.H{"rtcode": rtcode, "rtmsg": rtmsg},
		"data": data,
	})
}

// DeleteE164acc function.
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
func DeleteE164acc(c *gin.Context) {
	rtmsg := ``
	rtcode := ec.SUCCESS
	e164accs := make([]db.ACCE164, 0)
	data := make(map[string]interface{})

	if c.ContentType() != gin.MIMEJSON {
		rtcode = ec.ERROR_HTTP_REQUEST_CONTENTTYPE
	} else {
		if uuid := c.Param("uuid"); len(uuid) == 0 {
			rtcode = ec.ERROR_HTTP_REQUEST_URLPARAM
		} else {
			{
				if rte164acc, err := db.DeleteE164accsE164acc(uuid); err != nil {
					rtcode = ec.ERROR_DATABSE_DELETE
					rtmsg = err.Error()
				} else {
					e164accs = append(e164accs, rte164acc)
					data["len"] = len(e164accs)
					data["lists"] = e164accs
				}
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": gin.H{"rtcode": rtcode, "rtmsg": rtmsg},
		"data": data,
	})
}
