package dao

import (
	"errors"
	"github.com/Dataman-Cloud/omega-billing/model"
	"github.com/Dataman-Cloud/omega-billing/util"
	"github.com/Dataman-Cloud/omega-billing/util/mysql"
	log "github.com/cihub/seelog"
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
	sql := `insert into app_event(uid, cid, appname, active, createtime,endtime,cpus, mem, instances) values (:uid, :cid, :appname, :active, :createtime, :endtime, :cpus, :mem, :instances)`
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
	count := 0
	err := db.Get(&count, `select count(*) from app_event where uid=? and cid=? and appname=? and active=true`, event.Uid, event.Cid, event.AppName)
	if err != nil {
		log.Error("can't get app_event by uid and cid and appname and active=true")
		return errors.New("can't get app_event by uid and cid and appname and active=true")
	}
	/*if count > 0 {
		sql := `update app_event set endtime = :endtime, active = false where uid = :uid and cid = :cid and appname = :appname and active = true`
		_, err := db.NamedExec(sql, event)
		if err != nil {
			log.Errorf("update app_event error: %v", err)
			return err
		}
	}*/
	tx := db.MustBegin()
	_, err = tx.Exec(`update app_event set endtime=?, active=? where uid=? and cid=? and appname=? and active=true`, event.EndTime, event.Active, event.Uid, event.Cid, event.AppName)
	if err != nil {
		log.Errorf("update app update table app_event error: %v", err)
		tx.Rollback()
		return err
	}
	event.Active = true
	_, err = tx.Exec(`insert into app_event(uid, cid, appname, active, createtime, endtime, cpus, mem, instances) values (?, ?, ?, ?, ?, ?, ?, ?, ?)`, event.Uid, event.Cid, event.AppName, event.Active, event.CreateTime, event.EndTime, event.Cpus, event.Mem, event.Instances)
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

func GetBillings(uid, pcount, pnum uint64, order, sortby, appname, start, end string) ([]model.Event, error) {
	db := mysql.DB()
	if pcount <= 0 || pcount > 100 {
		pcount = 20
	}
	if pnum <= 0 {
		pnum = 1
	}

	if order == "" {
		order = "asc"
	}
	if sortby == "" {
		sortby = "createtime"
	}
	sql := `select * from app_event where uid = ?`
	if appname != "" {
		sql = sql + ` and appname = "` + appname + `"`
	}
	sql = sql + ` and createtime between '` + start + `' and '` + end + `' order by ` + sortby + ` ` + order + ` limit ?,?`
	billings := []model.Event{}
	err := db.Select(&billings, sql, uid, (pnum-1)*pcount, pcount)
	for _, billing := range billings {
		if billing.Active {
			billing.TimeLen = util.ParseTimeLen(time.Now().Unix() - billing.CreateTime.Unix())
		} else {
			billing.TimeLen = util.ParseTimeLen(billing.EndTime.Unix() - billing.CreateTime.Unix())
		}
	}
	return billings, err
}
