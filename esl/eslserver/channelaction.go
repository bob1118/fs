// //////////////////////channel event action////////////////////

package eslserver

import "github.com/bob1118/fs/esl/eventsocket"

// ChannelAction function
func ChannelAction(c *eventsocket.Connection, e *eventsocket.Event) {
	//e.LogPrint()
	eventName := e.Get("Event-Name")
	if len(eventName) > 0 {
		switch eventName {
		case "CHANNEL_STATE":
			channelstateAction(c, e)
		case "CHANNEL_CALLSTATE":
			channelcallstateAction(c, e)
		case "CHANNEL_HANGUP":
			channelhangupAction(c, e)
		case "CHANNEL_DESTROY":
			channelCDRAction(c, e)
		default:
			//nothing todo.
		}
	}
}

// channelstateAction function.
func channelstateAction(c *eventsocket.Connection, e *eventsocket.Event) {}

// channelcallstateAction function.
func channelcallstateAction(c *eventsocket.Connection, e *eventsocket.Event) {}

// channelhangupAction function.
func channelhangupAction(c *eventsocket.Connection, e *eventsocket.Event) {}

// channelCDRAction function. channel cdr.
func channelCDRAction(c *eventsocket.Connection, e *eventsocket.Event) {}
