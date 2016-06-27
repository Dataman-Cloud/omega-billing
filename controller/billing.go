package controller

import (
	"errors"
	"fmt"
	"github.com/Dataman-Cloud/omega-billing/dao"
	"github.com/Dataman-Cloud/omega-billing/util"
	log "github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
	"strconv"
)

func BillingList(c *gin.Context) {
	userid, ok := c.Get("uid")
	if !ok {
		log.Error("can't get userid")
		util.ReturnParamError(c, errors.New("can't get userid"))
		return
	}
	uid, err := strconv.ParseUint(fmt.Sprintf("%s", userid), 10, 64)
	if err != nil {
		log.Errorf("parse uid to uint error: %v", err)
		util.ReturnParamError(c, err)
		return
	}

	pcount := c.Query("per_page")
	if pcount == "" {
		pcount = "20"
	}
	pagecount, err := strconv.ParseUint(pcount, 10, 64)
	if err != nil {
		log.Errorf("parse pagecount to uint error: %v", err)
		util.ReturnParamError(c, err)
		return
	}
	pnum := c.Query("page")
	if pnum == "" {
		pnum = "1"
	}
	pagenum, err := strconv.ParseUint(pnum, 10, 64)
	if err != nil {
		log.Errorf("parse pagenum to uint error: %v", err)
		util.ReturnParamError(c, err)
		return
	}
	order := c.Query("order")
	sortby := c.Query("sort_by")
	appname := c.Query("appname")
	cid := c.Query("cid")
	starttime := c.Query("starttime")
	/*if starttime == "" {
		util.ReturnParamError(c, errors.New("can't find starttime"))
		return
	}*/
	endtime := c.Query("endtime")
	/*if endtime == "" {
		util.ReturnParamError(c, errors.New("can't find endtime"))
		return
	}*/
	billings, count, err := dao.GetBillings(uid, pagecount, pagenum, order, sortby, appname, starttime, endtime, cid)
	if err != nil {
		log.Errorf("get billings error: %v", err)
		util.ReturnDBError(c, err)
		return
	}
	//util.ReturnOK(c, billings)
	util.ReturnOK(c, map[string]interface{}{
		"billings": billings,
		"count":    count,
	})
	return
}
