// freeswitch run bgapi command do a backgroud job, execute bgapi command return job uuid.
// job in progress async, while receive BACKGROUND_JOB event with job uuid, job was down.
package bgapicmd

import (
	"net/http"

	"github.com/bob1118/fs/esl/eslclient"
	"github.com/gin-gonic/gin"
)

// some freeswitch bgapi/cmd reply.
// request: Get /bgapi?cmd=xxx
// response: job uuid.
func Get(c *gin.Context) {
	var cmd, result string
	cmd = c.Query("cmd")
	if jobuuid, err := eslclient.ClientCon.SendBgapiCommandAsync(cmd); err != nil {
		result = err.Error()
	} else {
		result = jobuuid
	}
	c.String(http.StatusOK, result)
}
