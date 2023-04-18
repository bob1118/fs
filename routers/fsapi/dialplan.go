package fsapi

import (
	"fmt"

	"github.com/bob1118/fs/fsconf"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func doDialplan(c *gin.Context) string {
	body := fsconf.NOT_FOUND
	addr := viper.GetString(`server.eventsocket.serveraddr`)
	if len(addr) > 0 { //DIALPLAN_APP_SOCKET, all incoming call hit app socket.
		body = fmt.Sprintf(fsconf.DIALPLAN_APP_SOCKET, addr, addr)
	}
	return body
}
