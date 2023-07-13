package apicmd

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/bob1118/fs/db"
	"github.com/bob1118/fs/ec"
	"github.com/bob1118/fs/esl/eslclient"
	"github.com/gin-gonic/gin"
)

// some freeswitch api/cmd response.
// request: Get /api?cmd=xxx
// response: string
func Get(c *gin.Context) {
	var sc int
	var err error
	var cmd, result string
	cmd = c.Query("cmd")
	if len(cmd) > 0 {
		if result, err = eslclient.ClientCon.SendApiCommandSync(cmd); err != nil {
			result = err.Error()
		}
		sc = http.StatusOK
	} else {
		sc = http.StatusBadRequest
	}
	c.String(sc, result)
}

// LoadGateway function
// freeswitch rescan profile's gateways configuration, run LoadGateway after PostGateway.
// SendApiCommand(sofia profile external rescan)
// request: /gateway/:uuid/load
// response: json
//
//	{
//		code:{
//			rtcode: rtcode, //return ec.SUCCESS or not.
//			rtmsg: rtmsg	//return error message while some error occured.
//		},
//		data:{
//			apicmd: api command,
//			result: api result
//		}
//	}
func LoadGateway(c *gin.Context) {
	rtmsg := ``
	rtcode := ec.SUCCESS
	data := make(map[string]interface{})

	if uuid := c.Param("uuid"); len(uuid) == 0 {
		rtcode = ec.ERROR_HTTP_REQUEST_URLPARAM
	} else {
		condition := fmt.Sprintf("gateway_uuid='%s'", uuid)
		if gateways, err := db.SelectGatewaysWithCondition(condition); err != nil {
			rtcode = ec.ERROR_DATABSE_QUERY
			rtmsg = err.Error()
		} else {
			if len(gateways) == 0 {
				rtcode = ec.ERROR_DATABSE_QUERY_ITEMNOTEXIST
			} else {
				profilename := gateways[0].Pname
				cmd := fmt.Sprintf("sofia profile %s rescan", profilename)
				if result, err := eslclient.ClientCon.SendApiCommandSync(cmd); err != nil {
					rtcode = ec.ERROR_SWITCH_EXECUTE_API
					rtmsg = err.Error()
				} else {
					if !strings.Contains(result, `+OK`) {
						rtcode = ec.ERROR_SWITCH_EXECUTE_API_RESULT
						rtmsg = result
					} else {
						data["apicmd"] = cmd
						data["result"] = result
					}
				}
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": gin.H{"rtcode": rtcode, "rtmsg": rtmsg},
		"data": data,
	})
}

// KillGateway function
// freeswitch kill running gateway, run KillGateway after DeleteGateway.
// SendApiCommand(sofia profile external killgw gateway)
// request: /gateway/:uuid/kill
// response: json
//
//	{
//		code:{
//			rtcode: rtcode, //return ec.SUCCESS or not.
//			rtmsg: rtmsg	//return error message while some error occured.
//		},
//		data:{
//			apicmd: api command,
//			result: api result
//		}
//	}
func KillGateway(c *gin.Context) {
	rtmsg := ``
	rtcode := ec.SUCCESS
	data := make(map[string]interface{})

	if uuid := c.Param("uuid"); len(uuid) == 0 {
		rtcode = ec.ERROR_HTTP_REQUEST_URLPARAM
	} else {
		condition := fmt.Sprintf("gateway_uuid='%s'", uuid)
		if gateways, err := db.SelectGatewaysWithCondition(condition); err != nil {
			rtcode = ec.ERROR_DATABSE_QUERY
			rtmsg = err.Error()
		} else {
			if len(gateways) == 0 {
				rtcode = ec.ERROR_DATABSE_QUERY_ITEMNOTEXIST
			} else {
				profilename := gateways[0].Pname
				gatewayname := gateways[0].Gname
				cmd := fmt.Sprintf("sofia profile %s killgw %s", profilename, gatewayname)
				if result, err := eslclient.ClientCon.SendApiCommandSync(cmd); err != nil {
					rtcode = ec.ERROR_SWITCH_EXECUTE_API
					rtmsg = err.Error()
				} else {
					if !strings.Contains(result, `+OK`) {
						rtcode = ec.ERROR_SWITCH_EXECUTE_API_RESULT
						rtmsg = result
					} else {
						data["apicmd"] = cmd
						data["result"] = result
					}
				}
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": gin.H{"rtcode": rtcode, "rtmsg": rtmsg},
		"data": data,
	})
}

// ReloadGateway function
// freeswitch kill running gateway, and reload it. run KillGateway after PutGateway.
// sofia profile external killgw gateway
// sofia profile external rescan
// request: /gateway/:uuid/reload
// response: json
//
//	{
//		code:{
//			rtcode: rtcode, //return ec.SUCCESS or not.
//			rtmsg: rtmsg	//return error message while some error occured.
//		},
//		data:{
//			killcmd: sofia profile external killgw gateway,
//			killresult: api result,
//			loadcmd: sofia profile external rescan,
//			loadresult: api result
//		}
//	}
func ReloadGateway(c *gin.Context) {
	rtmsg := ``
	rtcode := ec.SUCCESS
	data := make(map[string]interface{})

	if uuid := c.Param("uuid"); len(uuid) == 0 {
		rtcode = ec.ERROR_HTTP_REQUEST_URLPARAM
	} else {
		condition := fmt.Sprintf("gateway_uuid='%s'", uuid)
		if gateways, err := db.SelectGatewaysWithCondition(condition); err != nil {
			rtcode = ec.ERROR_DATABSE_QUERY
			rtmsg = err.Error()
		} else {
			if len(gateways) == 0 {
				rtcode = ec.ERROR_DATABSE_QUERY_ITEMNOTEXIST
			} else {
				profilename := gateways[0].Pname
				gatewayname := gateways[0].Gname
				killcmd := fmt.Sprintf("sofia profile %s killgw %s", profilename, gatewayname)
				if killresult, err := eslclient.ClientCon.SendApiCommandSync(killcmd); err != nil {
					rtcode = ec.ERROR_SWITCH_EXECUTE_API
					rtmsg = err.Error()
				} else {
					if !strings.Contains(killresult, `+OK`) {
						rtcode = ec.ERROR_SWITCH_EXECUTE_API_RESULT
						rtmsg = killresult
					} else {
						data["killcmd"] = killcmd
						data["killresult"] = killresult
						//load
						loadcmd := fmt.Sprintf("sofia profile %s rescan", profilename)
						if loadresult, err := eslclient.ClientCon.SendApiCommandSync(loadcmd); err != nil {
							rtcode = ec.ERROR_SWITCH_EXECUTE_API
							rtmsg = err.Error()
						} else {
							if !strings.Contains(loadresult, `+OK`) {
								rtcode = ec.ERROR_SWITCH_EXECUTE_API_RESULT
								rtmsg = loadresult
							} else {
								data["loadcmd"] = loadcmd
								data["loadresult"] = loadresult
							}
						}
					}
				}
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": gin.H{"rtcode": rtcode, "rtmsg": rtmsg},
		"data": data,
	})
}

// RegisterGateway function
// gateway configuration gateway_register=true
// api command: sofia profile external register gateway
// request: /gateway/:uuid/register
// response: json
//
//	{
//		code:{
//			rtcode: rtcode, //return ec.SUCCESS or not.
//			rtmsg: rtmsg	//return error message while some error occured.
//		},
//		data:{
//			apicmd: api command,
//			result: api result
//		}
//	}
func RegisterGateway(c *gin.Context) {
	rtmsg := ``
	rtcode := ec.SUCCESS
	data := make(map[string]interface{})

	if uuid := c.Param("uuid"); len(uuid) == 0 {
		rtcode = ec.ERROR_HTTP_REQUEST_URLPARAM
	} else {
		condition := fmt.Sprintf("gateway_uuid='%s'", uuid)
		if gateways, err := db.SelectGatewaysWithCondition(condition); err != nil {
			rtcode = ec.ERROR_DATABSE_QUERY
			rtmsg = err.Error()
		} else {
			if len(gateways) == 0 {
				rtcode = ec.ERROR_DATABSE_QUERY_ITEMNOTEXIST
			} else {
				profilename := gateways[0].Pname
				gatewayname := gateways[0].Gname
				cmd := fmt.Sprintf("sofia profile %s register %s", profilename, gatewayname)
				if result, err := eslclient.ClientCon.SendApiCommandSync(cmd); err != nil {
					rtcode = ec.ERROR_SWITCH_EXECUTE_API
					rtmsg = err.Error()
				} else {
					if !strings.Contains(result, `+OK`) {
						rtcode = ec.ERROR_SWITCH_EXECUTE_API_RESULT
						rtmsg = result
					} else {
						data["apicmd"] = cmd
						data["result"] = result
					}
				}
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": gin.H{"rtcode": rtcode, "rtmsg": rtmsg},
		"data": data,
	})
}

// UnregisterGateway function
// gateway configuration gateway_register=true
// api command: sofia profile external unregister gateway
// request: /gateway/:uuid/unregister
// response: json
//
//	{
//		code:{
//			rtcode: rtcode, //return ec.SUCCESS or not.
//			rtmsg: rtmsg	//return error message while some error occured.
//		},
//		data:{
//			apicmd: api command,
//			result: api result
//		}
//	}
func UnregisterGateway(c *gin.Context) {
	rtmsg := ``
	rtcode := ec.SUCCESS
	data := make(map[string]interface{})

	if uuid := c.Param("uuid"); len(uuid) == 0 {
		rtcode = ec.ERROR_HTTP_REQUEST_URLPARAM
	} else {
		condition := fmt.Sprintf("gateway_uuid='%s'", uuid)
		if gateways, err := db.SelectGatewaysWithCondition(condition); err != nil {
			rtcode = ec.ERROR_DATABSE_QUERY
			rtmsg = err.Error()
		} else {
			if len(gateways) == 0 {
				rtcode = ec.ERROR_DATABSE_QUERY_ITEMNOTEXIST
			} else {
				profilename := gateways[0].Pname
				gatewayname := gateways[0].Gname
				cmd := fmt.Sprintf("sofia profile %s unregister %s", profilename, gatewayname)
				if result, err := eslclient.ClientCon.SendApiCommandSync(cmd); err != nil {
					rtcode = ec.ERROR_SWITCH_EXECUTE_API
					rtmsg = err.Error()
				} else {
					if !strings.Contains(result, `+OK`) {
						rtcode = ec.ERROR_SWITCH_EXECUTE_API_RESULT
						rtmsg = result
					} else {
						data["apicmd"] = cmd
						data["result"] = result
					}
				}
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": gin.H{"rtcode": rtcode, "rtmsg": rtmsg},
		"data": data,
	})
}
