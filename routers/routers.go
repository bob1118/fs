package routers

import (
	"github.com/bob1118/fs/routers/fsapi"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	//receive mod_xml_curl request
	r.POST("/fsapi", fsapi.PostFromXmlCurl)
	return r
}
