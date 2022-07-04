package fsapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//PostFromXmlCurl
func PostFromXmlCurl(c *gin.Context) {
	var responseBody string
	value := c.PostForm("section")
	switch value {
	case "configuration":
		responseBody = doConfiguration(c)
	case "dialplan":
		responseBody = doDialplan(c)
	case "directory":
		responseBody = doDirectory(c)
	case "phrases":
		responseBody = doPhrases(c)
	}
	//for debug
	//log.Println("data:", c.Request.PostForm)
	c.String(http.StatusOK, responseBody)
}
