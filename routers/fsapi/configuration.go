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
		if verto, err := modules.GetConfiguration(c); err == nil {
			body = fmt.Sprintf(fsconf.CONFIGURATION, verto)
		}
	case "conference.conf": //5th
		if conference, err := modules.GetConfiguration(c); err == nil {
			body = fmt.Sprintf(fsconf.CONFIGURATION, conference)
		}
	case "db.conf": //6th
		if switchdb, err := modules.GetConfiguration(c); err == nil {
			body = fmt.Sprintf(fsconf.CONFIGURATION, switchdb)
		}
	case "fifo.conf": //7th
		if fifo, err := modules.GetConfiguration(c); err == nil {
			body = fmt.Sprintf(fsconf.CONFIGURATION, fifo)
		}
	case "hash.conf": //8th
		if hash, err := modules.GetConfiguration(c); err == nil {
			body = fmt.Sprintf(fsconf.CONFIGURATION, hash)
		}
	case "voicemail.conf": //9th
		if voicemail, err := modules.GetConfiguration(c); err == nil {
			body = fmt.Sprintf(fsconf.CONFIGURATION, voicemail)
		}
	case "httapi.conf": //10th
		if httapi, err := modules.GetConfiguration(c); err == nil {
			body = fmt.Sprintf(fsconf.CONFIGURATION, httapi)
		}
	case "spandsp.conf": //11th
		if spandsp, err := modules.GetConfiguration(c); err == nil {
			body = fmt.Sprintf(fsconf.CONFIGURATION, spandsp)
		}
	case "amr.conf": //12th
		if amr, err := modules.GetConfiguration(c); err == nil {
			body = fmt.Sprintf(fsconf.CONFIGURATION, amr)
		}
	case "opus.conf": //13th
		if opus, err := modules.GetConfiguration(c); err == nil {
			body = fmt.Sprintf(fsconf.CONFIGURATION, opus)
		}
	case "avformat.conf": //14th
		if av, err := modules.GetConfiguration(c); err == nil {
			body = fmt.Sprintf(fsconf.CONFIGURATION, av)
		}
	case "avcodec.conf": //15th
		if av, err := modules.GetConfiguration(c); err == nil {
			body = fmt.Sprintf(fsconf.CONFIGURATION, av)
		}
	case "sndfile.conf": //16th
		if sndfile, err := modules.GetConfiguration(c); err == nil {
			body = fmt.Sprintf(fsconf.CONFIGURATION, sndfile)
		}
	case "local_stream.conf": //17th
		if localstream, err := modules.GetConfiguration(c); err == nil {
			body = fmt.Sprintf(fsconf.CONFIGURATION, localstream)
		}
	case "lua.conf": //18th
		if lua, err := modules.GetConfiguration(c); err == nil {
			body = fmt.Sprintf(fsconf.CONFIGURATION, lua)
		}
	case "post_load_modules.conf": //19th
		if post_load_modules, err := modules.GetConfiguration(c); err == nil {
			body = fmt.Sprintf(fsconf.CONFIGURATION, post_load_modules)
		}
	case "event_socket.conf": //20th
		if event_socket, err := modules.GetConfiguration(c); err == nil {
			body = fmt.Sprintf(fsconf.CONFIGURATION, event_socket)
		}
	case "acl.conf": //21th
		if acl, err := modules.GetConfiguration(c); err == nil {
			body = fmt.Sprintf(fsconf.CONFIGURATION, acl)
		}
	case "post_load_switch.conf": //22th
		if post_load_switch, err := modules.GetConfiguration(c); err == nil {
			body = fmt.Sprintf(fsconf.CONFIGURATION, post_load_switch)
		}
	case "switch.conf": //23th.
		if switch_main, err := modules.GetConfiguration(c); err == nil {
			body = fmt.Sprintf(fsconf.CONFIGURATION, switch_main)
		}
	}
	return body
}
