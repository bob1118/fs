package routers

import (
	"github.com/bob1118/fs/routers/fsapi"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	//for debug
	if true {gin.SetMode(gin.DebugMode)}
	r := gin.Default()

	//receive mod_xml_curl request
	r.POST("/fsapi", fsapi.PostFromXmlCurl)
	return r
}
