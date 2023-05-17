package v1

import (
	"fmt"
	"net/http"

	"github.com/bob1118/fs/db"
	"github.com/bob1118/fs/ec"
	"github.com/gin-gonic/gin"
)

// GetFifos function.
//
// request: GET /api/v1/fifos?uuid=xxx&name=xxx
//
// response: json
//
//	{
//		code:{
//			rtcode: rtcode, //return ec.SUCCESS or not.
//			rtmsg: rtmsg	//return error message while some error occured.
//		},
//		data:{
//			len: len(slice),
//			lists:{slice[0],slice[1], ...}
//		}
//	}
func GetFifos(c *gin.Context) {
	rtmsg := ``
	rtcode := ec.SUCCESS
	condition := "true"
	data := make(map[string]interface{})

	if uuid := c.Query("uuid"); len(uuid) > 0 {
		condition += fmt.Sprintf(" and fifo_uuid='%s'", uuid)
	}
	if name := c.Query("name"); len(name) > 0 {
		condition += fmt.Sprintf(" and fifo_name='%s'", name)
	}

	if rtcode == ec.SUCCESS {
		if fifos, err := db.SelectFifosWithCondition(condition); err != nil {
			rtcode = ec.ERROR_DATABSE_QUERY
			rtmsg = err.Error()
		} else {
			data["len"] = len(fifos)
			data["lists"] = fifos
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": gin.H{"rtcode": rtcode, "rtmsg": rtmsg},
		"data": data,
	})
}

// PostFifo function.
//
// request: POST /api/v1/fifo, a Fifo{} json.
//
// response: json
//
//	{
//		code:{
//			rtcode: rtcode, //return ec.SUCCESS or not.
//			rtmsg: rtmsg	//return error message while some error occured.
//		},
//		data:{
//			len: len(slice),
//			lists:{slice[0],slice[1], ...}
//		}
//	}
func PostFifo(c *gin.Context) {
	rtmsg := ``
	rtcode := ec.SUCCESS
	fifo := db.Fifo{}
	fifos := make([]db.Fifo, 0)
	data := make(map[string]interface{})

	if c.ContentType() != gin.MIMEJSON {
		rtcode = ec.ERROR_HTTP_REQUEST_CONTENTTYPE
	} else {
		if err := c.BindJSON(&fifo); err != nil {
			rtcode = ec.ERROR_HTTP_REQUEST_CONTEXTBINDJSON
			rtmsg = err.Error()
		} else {
			if len(fifo.Fname) == 0 {
				rtcode = ec.ERROR_HTTP_REQUEST_JSONITEMNULL
			}
		}
	}

	if rtcode == ec.SUCCESS {
		fifos = append(fifos, fifo)
		if rtfifos, err := db.InsertFifos(fifos); err != nil {
			rtcode = ec.ERROR_DATABSE_INSERT
			rtmsg = err.Error()
		} else {
			data["len"] = len(rtfifos)
			data["lists"] = rtfifos
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": gin.H{"rtcode": rtcode, "rtmsg": rtmsg},
		"data": data,
	})
}

// PutFifo function.
//
// request: PUT /api/v1/fifo/:uuid, a Fifo{} json.
//
// response: json
//
//	{
//		code:{
//			rtcode: rtcode, //return ec.SUCCESS or not.
//			rtmsg: rtmsg	//return error message while some error occured.
//		},
//		data:{
//			len: len(slice),
//			lists:{slice[0],slice[1], ...}
//		}
//	}
func PutFifo(c *gin.Context) {
	rtmsg := ``
	rtcode := ec.SUCCESS
	fifo := db.Fifo{}
	fifos := make([]db.Fifo, 0)
	data := make(map[string]interface{})

	if c.ContentType() != gin.MIMEJSON {
		rtcode = ec.ERROR_HTTP_REQUEST_CONTENTTYPE
	} else {
		if uuid := c.Param("uuid"); len(uuid) == 0 {
			rtcode = ec.ERROR_HTTP_REQUEST_URLPARAM
		} else {
			if err := c.BindJSON(&fifo); err != nil {
				rtcode = ec.ERROR_HTTP_REQUEST_CONTEXTBINDJSON
				rtmsg = err.Error()
			} else {
				if rtfifo, err := db.UpdateFifosFifo(uuid, fifo); err != nil {
					rtcode = ec.ERROR_DATABSE_UPDATE
					rtmsg = err.Error()
				} else {
					fifos = append(fifos, rtfifo)
					data["len"] = len(fifos)
					data["lists"] = fifos
				}
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": gin.H{"rtcode": rtcode, "rtmsg": rtmsg},
		"data": data,
	})
}

// DeleteFifo function.
//
// request: DELETE /api/v1/fifo/:uuid
//
// response: json
//
//	{
//		code:{
//			rtcode: rtcode, //return ec.SUCCESS or not.
//			rtmsg: rtmsg	//return error message while some error occured.
//		},
//		data:{
//			len: len(slice),
//			lists:{slice[0],slice[1], ...}
//		}
//	}
func DeleteFifo(c *gin.Context) {
	rtmsg := ``
	rtcode := ec.SUCCESS
	fifos := make([]db.Fifo, 0)
	data := make(map[string]interface{})

	if c.ContentType() != gin.MIMEJSON {
		rtcode = ec.ERROR_HTTP_REQUEST_CONTENTTYPE
	} else {
		if uuid := c.Param("uuid"); len(uuid) == 0 {
			rtcode = ec.ERROR_HTTP_REQUEST_URLPARAM
		} else {
			{
				if rtfifo, err := db.DeleteFifosFifo(uuid); err != nil {
					rtcode = ec.ERROR_DATABSE_DELETE
					rtmsg = err.Error()
				} else {
					fifos = append(fifos, rtfifo)
					data["len"] = len(fifos)
					data["lists"] = fifos
				}
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": gin.H{"rtcode": rtcode, "rtmsg": rtmsg},
		"data": data,
	})
}
