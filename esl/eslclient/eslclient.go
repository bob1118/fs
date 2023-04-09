//eslclient is a tcp client connect to mod_evnet_socket.
//while mod_sofia receive a incoming call, dialplan execute app park.
//now, do what you want to before received park execute complete event.

package eslclient

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"syscall"
	"time"

	"github.com/bob1118/fs/esl/eventsocket"
	"github.com/bob1118/fs/utils"
	"github.com/spf13/viper"
)

// var CHfsisrun chan bool
var is_mod_odbc_cdr bool
var ClientCon *eventsocket.Connection

func init() {
	//CHfsisrun = make(chan bool)
}

// clientRun
func ClientRun() {
	modname := viper.GetString(`switch.cdr.modname`)
	is_mod_odbc_cdr = utils.IsEqual(modname, "mod_odbc_cdr")
	if err := clientReconnect(); err != nil {
		fmt.Println(err)
	}
}

// clientReconnect
func clientReconnect() error {
	var e error
	var addr, password string
	addr = fmt.Sprintf("%s:%s", viper.GetString(`switch.eventsocket.ipaddr`), viper.GetString(`switch.eventsocket.port`))
	password = viper.GetString(`switch.eventsocket.password`)
	alwaysrun := true
	//	if isrun, ok := <-CHfsisrun; ok {
	//		if isrun {
	for alwaysrun {
		fmt.Println("->start reconnect.")
		c, err := eventsocket.Dial(addr, password)
		if err != nil {
			//if errors.Is(err, syscall.WSAECONNRESET+7) { //syscall.Errno=10061 (No connection could be made because the target machine actively refused it)
			fmt.Println(err)
			e = err
			//}
		} else {
			ClientCon = c
			fmt.Println("->connect successful.")
			if eventSubscribe("plain") &&
				eventUnsubscribe("plain", "RE_SCHEDULE", "HEARTBEAT", "MESSAGE_WAITING", "MESSAGE_QUERY") {
				//if eventSubscribe("plain", "API", "BACKGROUND_JOB", "CUSTOM sofia::register sofia::unregister sofia::expire sofia::gateway_state") {
				if err := eventReadLoop(); err != nil {
					e = err
					if errors.Is(err, io.EOF) {
						fmt.Println(err)
					}
					//if errors.Is(err, syscall.WSAECONNRESET) { //windows
					//	fmt.Println(err)
					//}
					if errors.Is(err, syscall.ECONNRESET) { //linux
						fmt.Println(err)
					}
				}
			}
		}
		//		}
		//	}
		<-time.After(8 * time.Second)
	}
	return e
}

// EventLoop function.
func eventReadLoop() error {
	isLoop := true
	for isLoop {
		if e, err := ClientCon.ReadEvent(); err != nil {
			return err
		} else {
			eventAction(e)
		}
	}
	return nil
}

// eventSubscribe function.
func eventSubscribe(format string, enames ...string) bool {
	var isOK bool
	var command string

	command = fmt.Sprintf("event %s", format)
	if len(enames) == 0 {
		command += " all"
	} else {
		for _, ename := range enames {
			command += fmt.Sprintf(" %s", ename)
		}
	}

	if event, err := ClientCon.Send(command); err != nil {
		isOK = false
		fmt.Println(err)
	} else {
		reply := event.Header["Reply-Text"]
		if strings.Contains(reply.(string), "+OK") {
			isOK = true
		}
	}
	return isOK
}

// eventUnsubscribe function.
func eventUnsubscribe(format string, enames ...string) bool {
	var isOK bool
	var command string

	command = fmt.Sprintf("nixevent %s", format)
	if len(enames) == 0 {
		command = "noevents"
	} else {
		for _, ename := range enames {
			command += fmt.Sprintf(" %s", ename)
		}
	}

	if event, err := ClientCon.Send(command); err != nil {
		isOK = false
		fmt.Println(err.Error())
	} else {
		reply := event.Header["Reply-Text"]
		if strings.Contains(reply.(string), "+OK") {
			isOK = true
		}
	}
	return isOK
}
