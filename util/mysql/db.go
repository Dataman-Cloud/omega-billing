package mysql

import (
	"fmt"
	. "github.com/Dataman-Cloud/omega-billing/config"
	log "github.com/cihub/seelog"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattes/migrate/driver/mysql"
	"github.com/mattes/migrate/migrate"
	"sync"
)

func init() {
	DB()
	upgradeDB()
}

var db *sqlx.DB

func DB() *sqlx.DB {
	if db != nil {
		return db
	}
	mutex := sync.Mutex{}
	mutex.Lock()
	db, _ = InitDB()
	defer mutex.Unlock()
	return db
}

func upgradeDB() {
	uri := fmt.Sprintf("mysql://%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local",
		GetConfig().Mc.UserName,
		GetConfig().Mc.PassWord,
		GetConfig().Mc.Host,
		GetConfig().Mc.Port,
		GetConfig().Mc.DataBase)
	log.Info("upgrade db mysql drive: ", uri)
	errors, ok := migrate.UpSync(uri, "./sql")
	if errors != nil && len(errors) > 0 {
		for _, err := range errors {
			log.Error("db err", err)
		}
		log.Error("can't upgrade db", errors)
		log.Flush()
		panic(-1)
	}
	if !ok {
		log.Error("can't upgrade db")
		log.Flush()
		panic(-1)
	}
	log.Info("DB upgraded")
	log.Flush()
}

func InitDB() (*sqlx.DB, error) {
	uri := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local",
		GetConfig().Mc.UserName,
		GetConfig().Mc.PassWord,
		GetConfig().Mc.Host,
		GetConfig().Mc.Port,
		GetConfig().Mc.DataBase)
	db, err := sqlx.Open("mysql", uri)
	if err != nil {
		log.Errorf("cat not connection mysql error: %v, uri:%s", err, uri)
		return db, err
	}
	err = db.Ping()
	if err != nil {
		log.Error("can not ping mysql error: ", err)
		return db, err
	}
	db.SetMaxIdleConns(int(GetConfig().Mc.MaxIdleConns))
	db.SetMaxOpenConns(int(GetConfig().Mc.MaxOpenConns))
	return db, err
}
