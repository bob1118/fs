package v1

import (
	"fmt"
	"net/http"

	"github.com/bob1118/fs/db"
	"github.com/bob1118/fs/ec"
	"github.com/gin-gonic/gin"
)

// GetGateways function.
// request: GET /api/v1/gateways?uuid=xxx&pname=xxx&gname=xxx&username=xxx&realm=xxx&proxy=xxx&registerproxy=xxx&register=true
// response: json
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
		if accounts, err := db.SelectAccountsWithCondition(condition); err != nil {
			rtcode = ec.ERROR_DATABSE_QUERY
			rtmsg = err.Error()
		} else {
			data["len"] = len(accounts)
			data["lists"] = accounts
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": gin.H{"rtcode": rtcode, "rtmsg": rtmsg},
		"data": data,
	})
}

// PostGateway function.
func PostGateway(c *gin.Context) {}

// PutGateway function.
func PutGateway(c *gin.Context) {}

// DeleteGateway function.
func DeleteGateway(c *gin.Context) {}
