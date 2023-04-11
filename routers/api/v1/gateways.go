package v1

import (
	"fmt"
	"net/http"

	"github.com/bob1118/fs/db"
	"github.com/bob1118/fs/ec"
	"github.com/gin-gonic/gin"
)

// GetGateways function return gateways by condition.
//
// request: GET /api/v1/gateways?uuid=xxx&pname=xxx&gname=xxx&username=xxx&realm=xxx&proxy=xxx&registerproxy=xxx&register=true
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
func GetGateways(c *gin.Context) {
	rtmsg := ``
	rtcode := ec.SUCCESS
	condition := "true"
	data := make(map[string]interface{})

	if uuid := c.Query("uuid"); len(uuid) > 0 {
		condition += fmt.Sprintf(" and gateway_uuid='%s'", uuid)
	}
	if pname := c.Query("pname"); len(pname) > 0 {
		condition += fmt.Sprintf(" and profile_name ='%s'", pname)
	}
	if gname := c.Query("gname"); len(gname) > 0 {
		condition += fmt.Sprintf(" and gateway_name='%s'", gname)
	}
	if username := c.Query("username"); len(username) > 0 {
		condition += fmt.Sprintf(" and gateway_username='%s'", username)
	}
	if realm := c.Query("realm"); len(realm) > 0 {
		condition += fmt.Sprintf(" and gateway_realm='%s'", realm)
	}
	if proxy := c.Query("proxy"); len(proxy) > 0 {
		condition += fmt.Sprintf(" and gateway_proxy='%s'", proxy)
	}
	if registerproxy := c.Query("registerproxy"); len(registerproxy) > 0 {
		condition += fmt.Sprintf(" and gateway_registerproxy='%s'", registerproxy)
	}
	if register := c.Query("register"); len(register) > 0 {
		condition += fmt.Sprintf(" and gateway_register='%s'", register)
	}

	if rtcode == ec.SUCCESS {
		if gateways, err := db.SelectGatewaysWithCondition(condition); err != nil {
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

// PostGateway function.
//
// request: POST /api/v1/gateway, a Gateway{} json.
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
func PostGateway(c *gin.Context) {
	rtmsg := ``
	rtcode := ec.SUCCESS
	gw := db.Gateway{}
	gws := make([]db.Gateway, 0)
	data := make(map[string]interface{})

	if c.ContentType() != gin.MIMEJSON {
		rtcode = ec.ERROR_HTTP_REQUEST_CONTENTTYPE
	} else {
		if err := c.BindJSON(&gw); err != nil {
			rtcode = ec.ERROR_HTTP_REQUEST_CONTEXTBINDJSON
			rtmsg = err.Error()
		} else {
			if len(gw.Gname) == 0 {
				rtcode = ec.ERROR_HTTP_REQUEST_JSONITEMNULL
			} else {
				if len(gw.Gusername) == 0 || len(gw.Gpassword) == 0 { /// account username password *required* ///
					//maybe sip trunk, if blank.
				}
				if len(gw.Grealm) == 0 { /// auth realm: *optional* same as gateway name, if blank ///
					gw.Grealm = gw.Gname
				}
				if len(gw.Gfromuser) == 0 { /// username to use in from: *optional* same as  username, if blank ///
					gw.Gfromuser = gw.Gusername
				}
				if len(gw.Gfromdomain) == 0 { /// domain to use in from: *optional* same as  realm, if blank ///
					gw.Gfromdomain = gw.Grealm
				}
				if len(gw.Gextension) == 0 { /// extension for inbound calls: *optional* same as username, if blank ///
					gw.Gextension = gw.Gusername
				}
				if len(gw.Gproxy) == 0 { /// proxy host: *optional* same as realm, if blank ///
					gw.Gproxy = gw.Grealm
				}
				if len(gw.Gregisterproxy) == 0 { ///// send register to this proxy: *optional* same as proxy, if blank ///
					gw.Gregisterproxy = gw.Gproxy
				}
				if len(gw.Gexpire) == 0 { /// expire in seconds: *optional* 3600, if blank ///
					gw.Gexpire = `3600`
				}
				if len(gw.Gregister) == 0 { /// do not register ///
					gw.Gregister = `false`
				}
				if len(gw.Gcalleridinfrom) == 0 { /// Use the callerid of an inbound call in the from field on outbound calls via this gateway
					gw.Gcalleridinfrom = `true`
				}
				if len(gw.Gextensionincontact) == 0 { /// Put the extension in the contact
					gw.Gextensionincontact = `false`
				}
				if len(gw.Goptionping) == 0 { /// send an options ping every x seconds, failure will unregister and/or mark it down
					gw.Goptionping = `0`
				}
			}
		}
	}

	if rtcode == ec.SUCCESS {
		gws = append(gws, gw)
		if rtgws, err := db.InsertGateways(gws); err != nil {
			rtcode = ec.ERROR_DATABSE_INSERT
			rtmsg = err.Error()
		} else {
			data["len"] = len(rtgws)
			data["lists"] = rtgws
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": gin.H{"rtcode": rtcode, "rtmsg": rtmsg},
		"data": data,
	})
}

// PutGateway function.
//
// request: PUT /api/v1/gateway/:uuid, a Gateway{} json.
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
func PutGateway(c *gin.Context) {}

// DeleteGateway function.
//
// request: DELETE /api/v1/gateway/:uuid
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
func DeleteGateway(c *gin.Context) {}
