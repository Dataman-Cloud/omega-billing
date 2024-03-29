package dao

import (
	"errors"
	"fmt"
	"github.com/Dataman-Cloud/omega-billing/model"
	"github.com/Dataman-Cloud/omega-billing/util"
	"github.com/Dataman-Cloud/omega-billing/util/mysql"
	log "github.com/cihub/seelog"
	"strconv"
	"time"
)

func AddEvent(event *model.Event) (uint64, error) {
	db := mysql.DB()
	count := 0
	err := db.Get(&count, `select count(*) from app_event where uid=? and cid=? and appname=? and active=true`, event.Uid, event.Cid, event.AppName)
	if err != nil {
		log.Errorf("add event error or illegal data: %v", err)
		return 0, err
	}
	if count > 0 {
		sql := `update app_event set endtime = :endtime, active = false where uid = :uid and cid = :cid and appname = :appname and active = true`
		_, err := db.NamedExec(sql, event)
		if err != nil {
			log.Errorf("update app_event error: %v", err)
			return 0, err
		}
	}
	sql := `insert into app_event(uid, cid, appname, active, starttime,endtime,cpus, mem, instances) values (:uid, :cid, :appname, :active, :starttime, :endtime, :cpus, :mem, :instances)`
	stmt, err := db.PrepareNamed(sql)
	if err != nil {
		log.Errorf("add event stmt sql error: %v", err)
		return 0, err
	}
	defer func() {
		if stmt != nil {
			err = stmt.Close()
			if err != nil {
				log.Errorf("insert into event stmt close error: %v", err)
			}
		}
	}()
	result, err := stmt.Exec(event)
	if err != nil {
		log.Errorf("insert into event error: %v", err)
		return 0, err
	}
	id, err := result.LastInsertId()
	return uint64(id), err
}

func DeleteApp(event *model.Event) error {
	db := mysql.DB()
	sql := `update app_event set endtime = :endtime, active = :active where uid = :uid and cid = :cid and appname = :appname and active = true`
	_, err := db.NamedExec(sql, event)
	if err != nil {
		log.Errorf("update app_event error: %v", err)
		return err
	}
	return nil
}

func GetBilling(event *model.Event) (model.Event, error) {
	db := mysql.DB()
	billing := model.Event{}
	sql := `select * from app_event where uid=? and cid=? and appname=? and active=true`
	err := db.Get(&billing, sql, event.Uid, event.Cid, event.AppName)
	return billing, err
}

func UpdateApp(event *model.Event) error {
	db := mysql.DB()
	tx := db.MustBegin()
	_, err := tx.NamedExec(`update app_event set endtime=:endtime, active=false where uid=:uid and cid=:cid and appname=:appname and active=true`, event)
	if err != nil {
		log.Errorf("update app update table app_event error: %v", err)
		tx.Rollback()
		return err
	}
	_, err = tx.NamedExec(`insert into app_event(uid, cid, appname, active, starttime, endtime, cpus, mem, instances) values (:uid, :cid, :appname, :active, :starttime, :endtime, :cpus, :mem, :instances)`, event)
	if err != nil {
		log.Errorf("update app insert into table app_event error: %v", err)
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		log.Errorf("update app commit error: %v", err)
		tx.Rollback()
		return err
	}
	return nil
}

func GetBillings(uid, pcount, pnum uint64, order, sortby, appname, start, end, cid string) ([]model.Event, int, error) {
	db := mysql.DB()
	if pcount <= 0 || pcount > 100 {
		pcount = 20
	}
	if pnum <= 0 {
		pnum = 1
	}

	if order == "" {
		order = "desc"
	}
	if sortby == "" {
		sortby = "starttime"
	}
	sql := `select * from app_event where uid = ?`
	sql1 := `select count(*) from app_event where uid = ?`
	if appname != "" {
		if cid == "" {
			return nil, 0, errors.New("not found clusterid")
		}
		sql = sql + ` and appname = "` + appname + `" and cid = ` + cid
		sql1 = sql1 + ` and appname = "` + appname + `" and cid = ` + cid
	}
	if start != "" && end != "" {
		starttime, err := strconv.ParseInt(start, 10, 64)
		if err != nil {
			log.Errorf("parse start to int64 error: %v", err)
			return nil, 0, err
		}
		endtime, err := strconv.ParseInt(end, 10, 64)
		if err != nil {
			log.Errorf("parse end to int64 error: %v", err)
			return nil, 0, err
		}
		sql = sql + ` and (starttime between '` + time.Unix(starttime, 0).Format(time.RFC3339) + `' and '` + time.Unix(endtime, 0).Format(time.RFC3339) + `' or endtime between '` + time.Unix(starttime, 0).Format(time.RFC3339) + `' and '` + time.Unix(endtime, 0).Format(time.RFC3339) + `')`

		sql1 = sql1 + ` and (starttime between '` + time.Unix(starttime, 0).Format(time.RFC3339) + `' and '` + time.Unix(endtime, 0).Format(time.RFC3339) + `' or endtime between '` + time.Unix(starttime, 0).Format(time.RFC3339) + `' and '` + time.Unix(endtime, 0).Format(time.RFC3339) + `')`
	} else {
		if start != "" {
			starttime, err := strconv.ParseInt(start, 10, 64)
			if err != nil {
				log.Errorf("parse start to int64 error: %v", err)
				return nil, 0, err
			}
			sql = sql + ` and starttime >= '` + time.Unix(starttime, 0).Format(time.RFC3339) + `'`
			sql1 = sql1 + ` and starttime >= '` + time.Unix(starttime, 0).Format(time.RFC3339) + `'`
		} else if end != "" {
			endtime, err := strconv.ParseInt(end, 10, 64)
			if err != nil {
				log.Errorf("parse end to int64 error: %v", err)
				return nil, 0, err
			}
			sql = sql + ` and endtime <= '` + time.Unix(endtime, 0).Format(time.RFC3339) + `'`
			sql1 = sql1 + ` and endtime <= '` + time.Unix(endtime, 0).Format(time.RFC3339) + `'`
		}
	}
	count := 0
	err := db.Get(&count, sql1, uid)
	if err != nil {
		log.Errorf("get billing count error: %v", err)
		return nil, 0, err
	}
	if order != "desc" {
		order = "asc"
	}
	sql = sql + ` order by ` + sortby + ` ` + order + ` ,id ` + order
	sql = sql + fmt.Sprintf(" limit %d, %d", (pnum-1)*pcount, pcount)
	log.Debug("---------: ", sql)
	billings := []model.Event{}
	err = db.Select(&billings, sql, uid)
	for v, billing := range billings {
		if billing.Active {
			billings[v].TimeLen = util.ParseTimeLen(time.Now().Unix() - billing.StartTime.Unix())
			//billings[v].EndTime = time.Now()
		} else {
			billings[v].TimeLen = util.ParseTimeLen(billing.EndTime.Unix() - billing.StartTime.Unix())
		}
	}
	//return billings, len(billings), err
	return billings, count, err
}
