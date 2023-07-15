package v1

import (
	"fmt"
	"net/http"

	"github.com/bob1118/fs/db"
	"github.com/bob1118/fs/ec"
	"github.com/gin-gonic/gin"
)

// GetBlacklists function.
//
// request: GET /api/v1/blacklists?uuid=xxx&caller=xxx&callee=xxx
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
func GetBlacklists(c *gin.Context) {
	rtmsg := ``
	rtcode := ec.SUCCESS
	condition := "true"
	data := make(map[string]interface{})

	if uuid := c.Query("uuid"); len(uuid) > 0 {
		condition += fmt.Sprintf(" and blacklist_uuid='%s'", uuid)
	}
	if caller := c.Query("caller"); len(caller) > 0 {
		condition += fmt.Sprintf(" and blacklist_caller='%s'", caller)
	}
	if callee := c.Query("callee"); len(callee) > 0 {
		condition += fmt.Sprintf(" and blacklist_callee='%s'", callee)
	}

	if rtcode == ec.SUCCESS {
		if blacklists, err := db.SelectBlacklistsWithCondition(condition); err != nil {
			rtcode = ec.ERROR_DATABSE_QUERY
			rtmsg = err.Error()
		} else {
			data["len"] = len(blacklists)
			data["lists"] = blacklists
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": gin.H{"rtcode": rtcode, "rtmsg": rtmsg},
		"data": data,
	})
}

// PostBlacklist function.
//
// request: POST /api/v1/blacklist, a Blacklist{} json.
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
func PostBlacklist(c *gin.Context) {
	rtmsg := ``
	rtcode := ec.SUCCESS
	blacklist := db.Blacklist{}
	blacklists := make([]db.Blacklist, 0)
	data := make(map[string]interface{})

	if c.ContentType() != gin.MIMEJSON {
		rtcode = ec.ERROR_HTTP_REQUEST_CONTENTTYPE
	} else {
		if err := c.BindJSON(&blacklist); err != nil {
			rtcode = ec.ERROR_HTTP_REQUEST_CONTEXTBINDJSON
			rtmsg = err.Error()
		} else {
			if len(blacklist.Bcaller) == 0 || len(blacklist.Bcallee) == 0 {
				rtcode = ec.ERROR_HTTP_REQUEST_JSONITEMNULL
			}
		}
	}

	if rtcode == ec.SUCCESS {
		blacklists = append(blacklists, blacklist)
		if rt, err := db.InsertBlacklists(blacklists); err != nil {
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

// PutBlacklist function.
//
// request: PUT /api/v1/blacklist/:uuid, a Blacklist{} json.
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
func PutBlacklist(c *gin.Context) {
	rtmsg := ``
	rtcode := ec.SUCCESS
	blacklist := db.Blacklist{}
	blacklists := make([]db.Blacklist, 0)
	data := make(map[string]interface{})

	if c.ContentType() != gin.MIMEJSON {
		rtcode = ec.ERROR_HTTP_REQUEST_CONTENTTYPE
	} else {
		if uuid := c.Param("uuid"); len(uuid) == 0 {
			rtcode = ec.ERROR_HTTP_REQUEST_URLPARAM
		} else {
			if err := c.BindJSON(&blacklist); err != nil {
				rtcode = ec.ERROR_HTTP_REQUEST_CONTEXTBINDJSON
				rtmsg = err.Error()
			} else {
				if rt, err := db.UpdateBlacklistsBlacklist(uuid, blacklist); err != nil {
					rtcode = ec.ERROR_DATABSE_UPDATE
					rtmsg = err.Error()
				} else {
					blacklists = append(blacklists, rt)
					data["len"] = len(blacklists)
					data["lists"] = blacklists
				}
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": gin.H{"rtcode": rtcode, "rtmsg": rtmsg},
		"data": data,
	})
}

// DeleteBlacklist function.
//
// request: DELETE /api/v1/blacklist/:uuid
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
func DeleteBlacklist(c *gin.Context) {
	rtmsg := ``
	rtcode := ec.SUCCESS
	blacklists := make([]db.Blacklist, 0)
	data := make(map[string]interface{})

	if c.ContentType() != gin.MIMEJSON {
		rtcode = ec.ERROR_HTTP_REQUEST_CONTENTTYPE
	} else {
		if uuid := c.Param("uuid"); len(uuid) == 0 {
			rtcode = ec.ERROR_HTTP_REQUEST_URLPARAM
		} else {
			{
				if rt, err := db.DeleteBlacklistsBlacklist(uuid); err != nil {
					rtcode = ec.ERROR_DATABSE_DELETE
					rtmsg = err.Error()
				} else {
					blacklists = append(blacklists, rt)
					data["len"] = len(blacklists)
					data["lists"] = blacklists
				}
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": gin.H{"rtcode": rtcode, "rtmsg": rtmsg},
		"data": data,
	})
}
