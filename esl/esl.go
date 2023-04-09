package esl

import (
	"fmt"

	"github.com/bob1118/fs/esl/eslclient"
	"github.com/bob1118/fs/esl/eslserver"
)

func init() {}

// Run
// esl inbound connector execute system api,channel application, and more.
// esl outbound connector execute some channel application.
func Run(eslmode string) {
	switch eslmode {
	case "inbound", "Inbound", "INBOUND":
		eslclient.ClientRun()
	case "outbound", "Outbound", "OUTBOUND":
		eslserver.ServerRun()
	default:
		fmt.Println("known esl mode")
	}
}
