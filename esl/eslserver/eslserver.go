//eslserver is a tcp server for dialplan application socket
//while mod_sofia receive a incoming call, dialplan execute socket(ip:port async full) and socket connect to eslserver.

package eslserver

import (
	"fmt"

	"github.com/bob1118/fs/esl/eventsocket"
	"github.com/spf13/viper"
)

// ServerRun function.
func ServerRun() {
	addr := viper.GetString(`server.eventsocket.serveraddr`)
	if err := eventsocket.ListenAndServe(addr, handler); err != nil {
		fmt.Println(err)
	}
}

// ServerRestart function.
func ServerRestart() {}

// handler function.
func handler(c *eventsocket.Connection) {
	fmt.Println("new client:", c, "from:", c.RemoteAddr())
	if ev, err := c.SendCommandEx("connect"); err != nil {
		//if err := c.SendCommand("connect"); err != nil {
		fmt.Println(err)
	} else {
		if err := c.SendCommand("myevents"); err != nil {
			fmt.Println(err)
		} else {
			//send connect return default event CHANNEL_DATA.
			eventChannelDefaultAction(c, ev)
			eventChannelReadLoop(c)
		}
	}
}

// eventChannelDefaultAction function
func eventChannelDefaultAction(c *eventsocket.Connection, e *eventsocket.Event) (err error) {
	return ChannelDefaultAction(c, e)
}

// eventChannelReadLoop function
func eventChannelReadLoop(c *eventsocket.Connection) error {
	isLoop := true
	for isLoop {
		if e, err := c.ReadEvent(); err != nil {
			return err
		} else {
			ChannelEventAction(c, e)
		}
	}
	return nil
}
