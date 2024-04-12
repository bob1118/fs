package v1

import (
	"fmt"
	"net/http"

	"github.com/bob1118/fs/db"
	"github.com/bob1118/fs/ec"
	"github.com/bob1118/fs/esl/eslclient"
	"github.com/bob1118/fs/utils"
	"github.com/gin-gonic/gin"
)

// GetOutgoingcalls function return Outgoingcalls by condition.
//
// request: GET /api/v1/outgoingcalls?uuida=xxx&id=xxx&domain=xxx&uuidb=xxx&e164=xxx&gateway=xxx&ani=xxx&destination=xxx
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
func GetOutgoingcalls(c *gin.Context) {
	rtmsg := ``
	rtcode := ec.SUCCESS
	condition := "true"
	data := make(map[string]interface{})

	if uuida := c.Query("uuida"); len(uuida) > 0 {
		condition += fmt.Sprintf(" and uuida='%s'", uuida)
	}
	if id := c.Query("id"); len(id) > 0 {
		condition += fmt.Sprintf(" and id='%s'", id)
	}
	if domain := c.Query("domain"); len(domain) > 0 {
		condition += fmt.Sprintf(" and domain='%s'", domain)
	}
	if uuidb := c.Query("uuidb"); len(uuidb) > 0 {
		condition += fmt.Sprintf(" and uuidb='%s'", uuidb)
	}
	if e164 := c.Query("e164"); len(e164) > 0 {
		condition += fmt.Sprintf(" and e164='%s'", e164)
	}
	if gateway := c.Query("gateway"); len(gateway) > 0 {
		condition += fmt.Sprintf(" and gateway='%s'", gateway)
	}
	if ani := c.Query("ani"); len(ani) > 0 {
		condition += fmt.Sprintf(" and ani='%s'", ani)
	}
	if destination := c.Query("destination"); len(destination) > 0 {
		condition += fmt.Sprintf(" and destination='%s'", destination)
	}

	if rtcode == ec.SUCCESS {
		if outgoingcalls, err := db.SelectOutgoingcallsWithCondition(condition); err != nil {
			rtcode = ec.ERROR_DATABSE_QUERY
			rtmsg = err.Error()
		} else {
			data["len"] = len(outgoingcalls)
			data["lists"] = outgoingcalls
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": gin.H{"rtcode": rtcode, "rtmsg": rtmsg},
		"data": data,
	})
}

// PostOutgoingcall function.
//
// request: POST /api/v1/outgoingcall, a OUTGOINGCALL{} json.
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
func PostOutgoingcall(c *gin.Context) {
	rtmsg := ``
	rtcode := ec.SUCCESS
	outgoingcall := db.OUTGOINGCALL{}
	outgoingcalls := make([]db.OUTGOINGCALL, 0)
	data := make(map[string]interface{})

	if c.ContentType() != gin.MIMEJSON {
		rtcode = ec.ERROR_HTTP_REQUEST_CONTENTTYPE
	} else {
		if err := c.BindJSON(&outgoingcall); err != nil {
			rtcode = ec.ERROR_HTTP_REQUEST_CONTEXTBINDJSON
			rtmsg = err.Error()
		} else {
			var uuida, uuidb, cmd, a, b string
			//make safe aleg uuid.
			if len(outgoingcall.Auuid) != 36 {
				if uuida, err = eslclient.ClientCon.APICreateUUID(); err != nil {
					rtcode = ec.ERROR_SWITCH_EXECUTE_API
					rtmsg = err.Error()
				} else {
					outgoingcall.Auuid = uuida
				}
			}
			for isExist := true; isExist; {
				if isExist, err = db.IsExistCdrsAlegUuid(outgoingcall.Auuid); err != nil {
					rtcode = ec.ERROR_DATABSE_QUERY
					rtmsg = err.Error()
					isExist = false
				} else {
					if isExist { //uuid duplicate, createuuid agin.
						if uuida, err = eslclient.ClientCon.APICreateUUID(); err != nil {
							rtcode = ec.ERROR_SWITCH_EXECUTE_API
							rtmsg = err.Error()
						} else {
							outgoingcall.Auuid = uuida
						}
					}
				}
			}
			//make safe bleg uuid.
			if len(outgoingcall.Buuid) != 36 {
				if uuidb, err = eslclient.ClientCon.APICreateUUID(); err != nil {
					rtcode = ec.ERROR_SWITCH_EXECUTE_API
					rtmsg = err.Error()
				} else {
					outgoingcall.Buuid = uuidb
				}
			}
			for isExist := true; isExist; {
				if isExist, err = db.IsExistCdrsBlegUuid(outgoingcall.Buuid); err != nil {
					rtcode = ec.ERROR_DATABSE_QUERY
					rtmsg = err.Error()
					isExist = false
				} else {
					if isExist { //uuid duplicate, createuuid agin.
						if uuidb, err = eslclient.ClientCon.APICreateUUID(); err != nil {
							rtcode = ec.ERROR_SWITCH_EXECUTE_API
							rtmsg = err.Error()
						} else {
							outgoingcall.Buuid = uuidb
						}
					}
				}
			}
			//originate ua bridge gateway or originate gateway bridge ua
			if utils.IsEqual(outgoingcall.Id, outgoingcall.Ani) {
				//aleg ua, bleg gateway.
				a = fmt.Sprintf("{origination_uuid=%s,origination_caller_id_name=%s,origination_caller_id_number=%s,ignore-early_media=true,codec_string=PCM,continue_on_fail=true,hangup_after_bridge=true}sofia/%s/%s", outgoingcall.Auuid, outgoingcall.Destination, outgoingcall.Destination, outgoingcall.Domain, outgoingcall.Id)
				b = fmt.Sprintf("{origination_uuid=%s,origination_caller_id_name=%s,origination_caller_id_number=%s,ignore-early_media=true,codec_string=PCM}sofia/gateway/%s/%s", outgoingcall.Buuid, outgoingcall.E164, outgoingcall.E164, outgoingcall.Gateway, outgoingcall.Destination)

			}
			if utils.IsEqual(outgoingcall.Id, outgoingcall.Destination) {
				//aleg gateway, bleg ua.
				a = fmt.Sprintf("{origination_uuid=%s,origination_caller_id_name=%s,origination_caller_id_number=%s,ignore-early_media=true,codec_string=PCM,continue_on_fail=true,hangup_after_bridge=true}sofia/gateway/%s/%s", outgoingcall.Auuid, outgoingcall.E164, outgoingcall.E164, outgoingcall.Gateway, outgoingcall.Ani)
				b = fmt.Sprintf("{origination_uuid=%s,origination_caller_id_name=%s,origination_caller_id_number=%s,ignore-early_media=true,codec_string=PCM}sofia/%s/%s", outgoingcall.Buuid, outgoingcall.Ani, outgoingcall.Ani, outgoingcall.Domain, outgoingcall.Id)
			}
			cmd = fmt.Sprintf("originate %s &bridge(%s)", a, b)
			//do bgapi command Async, while job execute success return aleg uuid
			if jobuuid, bgapierr := eslclient.ClientCon.SendBgapiCommandAsync(cmd); bgapierr != nil {
				rtcode = ec.ERROR_SWITCH_EXECUTE_BGAPI_ORIGINATE
				rtmsg = bgapierr.Error()
			} else {
				outgoingcall.Jobuuid = jobuuid
			}
		}
	}

	if rtcode == ec.SUCCESS {
		outgoingcalls = append(outgoingcalls, outgoingcall)
		if rtoutgoingcalls, err := db.InsertOutgoingcalls(outgoingcalls); err != nil {
			rtcode = ec.ERROR_DATABSE_INSERT
			rtmsg = err.Error()
		} else {
			data["len"] = len(rtoutgoingcalls)
			data["lists"] = rtoutgoingcalls
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": gin.H{"rtcode": rtcode, "rtmsg": rtmsg},
		"data": data,
	})
}
