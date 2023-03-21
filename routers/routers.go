package routers

import (
	"github.com/bob1118/fs/routers/api"
	v1 "github.com/bob1118/fs/routers/api/v1"
	"github.com/bob1118/fs/routers/api/v1/apicmd"
	"github.com/bob1118/fs/routers/api/v1/bgapicmd"
	"github.com/bob1118/fs/routers/fsapi"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

const T_GATEWAY = 1
const T_SERVER = 2

func NewRouter(mytype int) *gin.Engine {
	//for debug
	if true {
		gin.SetMode(gin.DebugMode)
	}

	r := gin.Default()
	switch mytype {
	case T_GATEWAY:
		gatewayRouter(r)
	case T_SERVER:
		serverRouter(r)
	}
	return r
}

func gatewayRouter(r *gin.Engine) {
	//receive mod_xml_curl request
	r.POST("/fsapi", fsapi.PostFromXmlCurl)
}

func serverRouter(r *gin.Engine) {
	switchRecorddingDir := viper.GetString(`switch.record.dir`)

	r.GET("/api", api.DefaultOK)
	apiv1 := r.Group("/api/v1")
	{
		//api/v1/recorddings
		apiv1.Static("/recorddings", switchRecorddingDir)
		//api/v1/api?cmd=xxx
		apiv1.GET("/api", apicmd.Get)
		//api/v1/bgapi?cmd=xxx
		apiv1.GET("/bgapi", bgapicmd.Get)

		////gateway tables
		//table confs
		apiv1.GET("/confs")
		apiv1.POST("/confs")
		apiv1.PUT("/confs/:uuid")
		//apiv1.PATCH("/confs/:uuid")
		apiv1.DELETE("/confs/:uuid")
		//table accounts
		apiv1.GET("/accounts", v1.GetAccounts)
		apiv1.POST("/account", v1.PostAccount)
		apiv1.POST("/accounts", v1.PostAccounts)
		apiv1.PUT("/account/:uuid", v1.PutAccount)
		apiv1.DELETE("/account/:uuid", v1.DeleteAccount)
		//table gateways
		//table e164s
		//table acce164s
		//table fifos
		//table fifomembers
		////server tables
		//table backlists
		//table bgjobs
		////cdr tables
		//table cdr
	}
	apiv2 := r.Group("/api/v2")
	{
		apiv2.POST("/default")
	}
}
