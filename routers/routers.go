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
		///
		apiv1.GET("/", v1.DefaultOK)
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
		apiv1.GET("/gateways", v1.GetGateways)
		apiv1.POST("/gateway", v1.PostGateway)
		apiv1.PUT("/gateway/:uuid", v1.PutGateway)
		apiv1.DELETE("gateway/:uuid", v1.DeleteGateway)
		//table e164s
		apiv1.GET("/e164s", v1.GetE164s)
		apiv1.POST("/e164", v1.PostE164)
		apiv1.POST("/e164s", v1.PostE164s)
		apiv1.PUT("/e164/:uuid", v1.PutE164)
		apiv1.DELETE("/e164/:uuid", v1.DeleteE164)
		//table acce164s
		apiv1.GET("/acce164s", v1.GetAcce164s)
		apiv1.POST("/acce164", v1.PostAcce164)
		apiv1.PUT("/acce164/:uuid", v1.PutAcce164)
		apiv1.DELETE("/acce164/:uuid", v1.DeleteAcce164)
		//table fifos
		apiv1.GET("/fifos", v1.GetFifos)
		apiv1.POST("/fifo", v1.PostFifo)
		apiv1.PUT("/fifo/:uuid", v1.PutFifo)
		apiv1.DELETE("/fifo/:uuid", v1.DeleteFifo)
		//table fifomembers
		apiv1.GET("/fifomembers", v1.GetFifomembers)
		apiv1.POST("/fifomember", v1.PostFifomember)
		apiv1.PUT("/fifomember/:uuid", v1.PutFifomember)
		apiv1.DELETE("/fifomember/:uuid", v1.DeleteFifomember)

		////server tables
		//table jobs for outgoing call
		apiv1.GET("/jobs")
		apiv1.POST("/job")
		apiv1.PUT("/job/:uuid")
		apiv1.DELETE("/job/:uuid")
		//table backlists
		apiv1.GET("/backlists", v1.GetBlacklists)
		apiv1.POST("/backlist", v1.PostBlacklist)
		apiv1.PUT("/backlist/:uuid", v1.PutBlacklist)
		apiv1.DELETE("/backlist/:uuid", v1.DeleteBlacklist)
		//table bgjobs
		apiv1.GET("/bgjobs", v1.GetBgjobs)
		apiv1.POST("/bgjob")
		apiv1.PUT("/bgjob/:uuid")
		apiv1.DELETE("/bgjob/:uuid")
		////cdr tables
		//table cdr
		apiv1.GET("/cdr")
		apiv1.POST("/cdr")
		apiv1.PUT("/cdr/:uuid")
		apiv1.DELETE("/cdr/:uuid")
	}
	apiv2 := r.Group("/api/v2")
	{
		apiv2.POST("/default")
	}
}
