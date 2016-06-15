package main

import (
	"flag"

	"github.com/Dataman-Cloud/omega-billing/config"
	"github.com/Dataman-Cloud/omega-billing/logger"
	"github.com/Dataman-Cloud/omega-billing/router"
	"github.com/Dataman-Cloud/omega-billing/util/mysql"
	"github.com/Dataman-Cloud/omega-billing/util/rabbitmq"
	"github.com/Dataman-Cloud/omega-billing/util/redis"
	log "github.com/cihub/seelog"
)

var (
	envFile = flag.String("config", "env_file", "")
)

func main() {
	flag.Parse()

	config.InitConfig(*envFile)
	logger.LogInit()
	mysql.MysqlInit()
	rabbitmq.MqInit()
	redis.RedisInit()
	log.Debug("omega-billing starting...")
	router.Run()
}
