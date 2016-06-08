package controller

import (
	"errors"
	"fmt"
	"github.com/Dataman-Cloud/omega-billing/dao"
	"github.com/Dataman-Cloud/omega-billing/util"
	"github.com/Dataman-Cloud/omega-billing/util/mysql"
	rd "github.com/Dataman-Cloud/omega-billing/util/redis"
	log "github.com/cihub/seelog"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

func Health(c *gin.Context) {
	hstart := time.Now().UnixNano()
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
	health["redis"]["time"] = (time.Now().UnixNano() - hstart) / 1000000
	health["omegaBilling"]["time"] = (time.Now().UnixNano() - hstart) / 1000000

	err = mysql.DB().Ping()
	if err == nil {
		health["mysql"] = map[string]interface{}{"status": 0}
	} else {
		health["mysql"] = map[string]interface{}{"status": 1}
		health["omegaBilling"]["status"] = 1
	}
	health["mysql"]["time"] = (time.Now().UnixNano() - hstart) / 1000000

	util.ReturnOK(c, health)
	return
}

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
