package v1

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/bob1118/fs/db"
	"github.com/bob1118/fs/ec"
	"github.com/bob1118/fs/utils"
	"github.com/gin-gonic/gin"
)

// GetAccounts function
//
// request: GET /api/v1/accounts?uuid=xxx&id=xxx&name=xxx&auth=xxx&group=xxx&domain=xxx&proxy=xxx&pn=xxx&ps=xxx
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
func GetAccounts(c *gin.Context) {
	rtmsg := ``
	rtcode := ec.SUCCESS
	condition := "true"
	data := make(map[string]interface{})

	if uuid := c.Query("uuid"); len(uuid) > 0 {
		condition += fmt.Sprintf(" and account_uuid='%s'", uuid)
	}
	if id := c.Query("id"); len(id) > 0 {
		condition += fmt.Sprintf(" and account_id ='%s'", id)
	}
	if name := c.Query("name"); len(name) > 0 {
		condition += fmt.Sprintf(" and account_name='%s'", name)
	}
	if auth := c.Query("auth"); len(auth) > 0 {
		condition += fmt.Sprintf(" and account_auth='%s'", auth)
	}
	if group := c.Query("group"); len(group) > 0 {
		condition += fmt.Sprintf(" and account_group='%s'", group)
	}
	if domain := c.Query("domain"); len(domain) > 0 {
		condition += fmt.Sprintf(" and account_domain='%s'", domain)
	}
	if proxy := c.Query("proxy"); len(proxy) > 0 {
		condition += fmt.Sprintf(" and account_proxy='%s'", proxy)
	}
	if pageSize := c.Query("ps"); len(pageSize) > 0 {
		if size, err := strconv.Atoi(pageSize); err == nil {
			if size > 0 {
				if pageNumber := c.Query("pn"); len(pageNumber) > 0 {
					if number, err := strconv.Atoi(pageNumber); err == nil {
						if number > 0 {
							offset := (number - 1) * size
							condition += fmt.Sprintf(" offset %d limit %d", offset, size)
						} else {
							rtcode = ec.ERROR_HTTP_REQUEST_URLQUERYKEY
						}
					}
				}
			} else {
				rtcode = ec.ERROR_HTTP_REQUEST_URLQUERYKEY
			}
		}
	}

	if rtcode == ec.SUCCESS {
		if accounts, err := db.SelectAccountsWithCondition(condition); err != nil {
			rtcode = ec.ERROR_DATABSE_QUERY
			rtmsg = err.Error()
		} else {
			data["len"] = len(accounts)
			data["lists"] = accounts
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": gin.H{"rtcode": rtcode, "rtmsg": rtmsg},
		"data": data,
	})
}

// PostAccount function
//
// request: POST /api/v1/account, a Account{} json.
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
func PostAccount(c *gin.Context) {
	rtmsg := ``
	rtcode := ec.SUCCESS
	ua := db.Account{}
	uas := make([]db.Account, 0)
	data := make(map[string]interface{})

	if c.ContentType() != gin.MIMEJSON {
		rtcode = ec.ERROR_HTTP_REQUEST_CONTENTTYPE
	} else {
		if err := c.BindJSON(&ua); err != nil {
			rtcode = ec.ERROR_HTTP_REQUEST_CONTEXTBINDJSON
			rtmsg = err.Error()
		} else {
			if len(ua.Aid) == 0 || len(ua.Adomain) == 0 {
				rtcode = ec.ERROR_HTTP_REQUEST_JSONITEMNULL
			} else {
				if len(ua.Aname) == 0 {
					ua.Aname = ua.Aid
				}
				if len(ua.Aauth) == 0 {
					ua.Aauth = ua.Aid
				}
				if len(ua.Apassword) == 0 {
					ua.Apassword = fmt.Sprintf("%s@%s", ua.Aid, ua.Adomain)
				}
				if len(ua.Aa1hash) == 0 { // md5(user:domain:password)
					s := fmt.Sprintf("%s:%s:%s", ua.Aid, ua.Adomain, ua.Apassword)
					ua.Aa1hash = utils.MakeA1Hash(s)
				}
				if len(ua.Agroup) == 0 {
					ua.Agroup = `default`
				}
				if len(ua.Aproxy) == 0 {
					ua.Aproxy = ua.Adomain
				}
				//cacheable notice, default ""
				//1,cacheable = "true"; mod_xml_curl requst reduce.
				//2,cacheable = "60000" mod_xml_curl request cache timer 60s.
				//3,xml_flush_cache;xml_flush_cache id 1002 10.10.10.250
				if len(ua.Acacheable) == 0 {
					ua.Acacheable = ""
				}
			}
		}
	}

	if rtcode == ec.SUCCESS {
		uas = append(uas, ua)
		if rtuas, err := db.InsertAccounts(uas); err != nil {
			rtcode = ec.ERROR_DATABSE_INSERT
			rtmsg = err.Error()
		} else {
			data["len"] = len(rtuas)
			data["lists"] = rtuas
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": gin.H{"rtcode": rtcode, "rtmsg": rtmsg},
		"data": data,
	})
}

// PostAccounts function
//
// request: POST /api/v1/accounts?domain=mydomain&idstart=8000&idend=8019
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
func PostAccounts(c *gin.Context) {
	rtmsg := ``
	rtcode := ec.SUCCESS
	ua := db.Account{}
	uas := make([]db.Account, 0)
	data := make(map[string]interface{})

	if c.ContentType() != gin.MIMEJSON {
		rtcode = ec.ERROR_HTTP_REQUEST_CONTENTTYPE
	} else {
		if domain := c.Query("domain"); len(domain) == 0 {
			rtcode = ec.ERROR_HTTP_REQUEST_URLQUERYKEY
		} else {
			ids := c.Query("idstart")
			ide := c.Query("idend")
			if len(ids) == 0 || len(ide) == 0 {
				rtcode = ec.ERROR_HTTP_REQUEST_URLQUERYKEY
			} else {
				if start, err := strconv.Atoi(ids); err != nil {
					rtcode = ec.ERROR_HTTP_REQUEST_URLQUERYATOI
					rtmsg = err.Error()
				} else {
					if end, err := strconv.Atoi(ide); err != nil {
						rtcode = ec.ERROR_HTTP_REQUEST_URLQUERYATOI
						rtmsg = err.Error()
					} else {
						for index := start; index <= end; index++ {
							ua.Aid = fmt.Sprintf("%d", index)
							ua.Aname = ua.Aid
							ua.Aauth = ua.Aid
							ua.Apassword = fmt.Sprintf("%s@%s", ua.Aid, domain)
							ua.Aa1hash = utils.MakeA1Hash(fmt.Sprintf("%s:%s:%s", ua.Aid, domain, ua.Apassword))
							ua.Agroup = "default"
							ua.Adomain = domain
							ua.Aproxy = ua.Adomain
							ua.Acacheable = ""
							uas = append(uas, ua)
						}
					}
				}
			}
		}
	}

	if rtcode == ec.SUCCESS {
		if rtuas, err := db.InsertAccounts(uas); err != nil {
			rtcode = ec.ERROR_DATABSE_INSERT
			rtmsg = err.Error()
		} else {
			data["len"] = len(rtuas)
			data["lists"] = rtuas
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": gin.H{"rtcode": rtcode, "rtmsg": rtmsg},
		"data": data,
	})
}

// PutAccount function
//
// request: PUT /api/v1/account/:uuid, a Account{} json.
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
func PutAccount(c *gin.Context) {
	rtmsg := ``
	rtcode := ec.SUCCESS
	ua := db.Account{}
	uas := make([]db.Account, 0)
	data := make(map[string]interface{})

	if c.ContentType() != gin.MIMEJSON {
		rtcode = ec.ERROR_HTTP_REQUEST_CONTENTTYPE
	} else {
		if uuid := c.Param("uuid"); len(uuid) == 0 {
			rtcode = ec.ERROR_HTTP_REQUEST_URLPARAM
		} else {
			if err := c.BindJSON(&ua); err != nil {
				rtcode = ec.ERROR_HTTP_REQUEST_CONTEXTBINDJSON
				rtmsg = err.Error()
			} else {
				if rtua, err := db.UpdateAccountsAccount(uuid, ua); err != nil {
					rtcode = ec.ERROR_DATABSE_UPDATE
					rtmsg = err.Error()
				} else {
					uas = append(uas, rtua)
					data["len"] = len(uas)
					data["lists"] = uas
				}
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": gin.H{"rtcode": rtcode, "rtmsg": rtmsg},
		"data": data,
	})
}

// DeleteAccount function
//
// request: DELETE /api/v1/account/:uuid
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
func DeleteAccount(c *gin.Context) {
	rtmsg := ``
	rtcode := ec.SUCCESS
	uas := make([]db.Account, 0)
	data := make(map[string]interface{})

	if c.ContentType() != gin.MIMEJSON {
		rtcode = ec.ERROR_HTTP_REQUEST_CONTENTTYPE
	} else {
		if uuid := c.Param("uuid"); len(uuid) == 0 {
			rtcode = ec.ERROR_HTTP_REQUEST_URLPARAM
		} else {
			{
				if rtua, err := db.DeleteAccountsAccount(uuid); err != nil {
					rtcode = ec.ERROR_DATABSE_DELETE
					rtmsg = err.Error()
				} else {
					uas = append(uas, rtua)
					data["len"] = len(uas)
					data["lists"] = uas
				}
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": gin.H{"rtcode": rtcode, "rtmsg": rtmsg},
		"data": data,
	})
}
