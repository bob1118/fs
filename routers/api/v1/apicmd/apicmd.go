package apicmd

import (
	"net/http"

	"github.com/bob1118/fs/esl/eslclient"
	"github.com/gin-gonic/gin"
)

// some freeswitch api/cmd response.
// request: Get api/v1/api?cmd=xxx
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
