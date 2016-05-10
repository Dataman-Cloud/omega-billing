package controller

import (
	"errors"
	"github.com/Dataman-Cloud/omega-billing/dao"
	"github.com/Dataman-Cloud/omega-billing/util"
	"github.com/Dataman-Cloud/omega-billing/util/mysql"
	rd "github.com/Dataman-Cloud/omega-billing/util/redis"
	"github.com/garyburd/redigo/redis"
	//log "github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

func Health(c *gin.Context) {
	hstart := time.Now().Unix()
	health := map[string]map[string]interface{}{
		"omegaBilling": map[string]interface{}{
			"status": 0,
		},
	}
	conn := rd.Open()
	defer conn.Close()
	_, err := redis.String(conn.Do("PING"))
	if err == nil {
		health["redis"] = map[string]interface{}{"status": 0}
	} else {
		health["redis"] = map[string]interface{}{"status": 1}
		health["omegaBilling"]["status"] = 1
	}
	health["redis"]["time"] = time.Now().Unix() - hstart
	health["omegaBilling"]["time"] = time.Now().Unix() - hstart

	err = mysql.DB().Ping()
	if err == nil {
		health["mysql"] = map[string]interface{}{"status": 0}
	} else {
		health["mysql"] = map[string]interface{}{"status": 1}
		health["omegaBilling"]["status"] = 1
	}
	health["mysql"]["time"] = time.Now().Unix() - hstart

	util.ReturnOK(c, health)
	return
}

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
	appname := c.Query("appname")
	starttime := c.Query("starttime")
	if starttime == "" {
		util.ReturnParamError(c, errors.New("can't find starttime"))
		return
	}
	endtime := c.Query("endtime")
	if endtime == "" {
		util.ReturnParamError(c, errors.New("can't find endtime"))
		return
	}
	billings, err := dao.GetBillings(userid.(uint64), uint64(pagecount), uint64(pagenum), order, sortby, appname, starttime, endtime)
	if err != nil {
		util.ReturnDBError(c, err)
		return
	}
	util.ReturnOK(c, billings)
	return
}
