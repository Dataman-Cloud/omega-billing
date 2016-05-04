package main

import (
	_ "github.com/Dataman-Cloud/omega-billing/config"
	_ "github.com/Dataman-Cloud/omega-billing/logger"
	"github.com/Dataman-Cloud/omega-billing/router"
	_ "github.com/Dataman-Cloud/omega-billing/util/mysql"
	_ "github.com/Dataman-Cloud/omega-billing/util/rabbitmq"
	_ "github.com/Dataman-Cloud/omega-billing/util/redis"
	log "github.com/cihub/seelog"
)

func main() {
	log.Debug("omega-billing starting...")
	router.Run()
}
