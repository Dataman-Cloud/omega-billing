package controller

import (
	"errors"
	"github.com/Dataman-Cloud/omega-billing/dao"
	"github.com/Dataman-Cloud/omega-billing/util"
	//log "github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
	"strconv"
)

func BillingList(c *gin.Context) {
	userid, ok := c.Get("uid")
	if !ok {
		util.ReturnParamError(c, errors.New("can't get userid"))
		return
	}

	pcount := c.Query("per_page")
	pagecount, err := strconv.ParseInt(pcount, 10, 64)
	if err != nil {
		util.ReturnParamError(c, err)
		return
	}
	pnum := c.Query("page")
	pagenum, err := strconv.ParseInt(pnum, 10, 64)
	if err != nil {
		util.ReturnParamError(c, err)
		return
	}
	order := c.Query("order")
	sortby := c.Query("sort_by")
	billings, err := dao.GetBilling(userid.(uint64), uint64(pagecount), uint64(pagenum), order, sortby)
	if err != nil {
		util.ReturnDBError(c, err)
		return
	}
	util.ReturnOK(c, billings)
	return
}
