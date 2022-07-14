package fsapi

import (
	"fmt"

	"github.com/bob1118/fs/fsconf"
	"github.com/bob1118/fs/fsconf/modules"
	"github.com/gin-gonic/gin"
)

//doConfiguration function return mod_xxx xml config.
func doConfiguration(c *gin.Context) (b string) {

	body := fsconf.NOT_FOUND
	value := c.PostForm(`key_value`)
	switch value {
	//switch boot order.
	// case "console.conf": //?
	// case "logfile.conf": //?
	// case "enum.conf": //?
	// case "xml_curl.conf": //?
	case "odbc_cdr.conf": //1th request.
		if odbc_cdr, err := modules.GetConfiguration(c); err == nil {
			body = fmt.Sprintf(fsconf.CONFIGURATION, odbc_cdr)
		}
	case "sofia.conf": //2th request(a request per profile).
		if sofia, err := modules.GetConfiguration(c); err == nil {
			body = fmt.Sprintf(fsconf.CONFIGURATION, sofia)
		}
	case "loopback.conf": //3th
	case "verto.conf": //4th
	case "conference.conf": //5th
	case "db.conf": //6th
		if switchdb, err := modules.GetConfiguration(c); err == nil {
			body = fmt.Sprintf(fsconf.CONFIGURATION, switchdb)
		}
	case "fifo.conf": //7th
	case "hash.conf": //8th
	case "voicemail.conf": //9th
	case "httapi.conf": //10th
	case "spandsp.conf": //11th
	case "amr.conf": //12th
	case "opus.conf": //13th
	case "avformat.conf": //14th
	case "avcodec.conf": //15th
	case "sndfile.conf": //16th
	case "local_stream.conf": //17th
	case "lua.conf": //18th
	case "post_load_modules.conf": //19th
	case "event_socket.conf": //20th
	case "acl.conf": //21th
	case "post_load_switch.conf": //22th
	case "switch.conf": //23th? switch version v1.10.7 missing.
	}
	return body
}
