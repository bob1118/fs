package apicmd

import (
	"github.com/bob1118/fs/esl/eslclient"
	"github.com/gin-gonic/gin"
)

// some freeswitch api/cmd response.
// request: Get /api?cmd=xxx
// response: string
func Get(c *gin.Context) {
	var err error
	var cmd, result string
	cmd = c.Query("cmd")
	if result, err = eslclient.ClientCon.SendApiCommand(cmd); err != nil {
		result = err.Error()
	}
	c.String(200, result)
}
