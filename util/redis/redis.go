package redis

import (
	"fmt"
	. "github.com/Dataman-Cloud/omega-billing/config"
	log "github.com/cihub/seelog"
	"github.com/garyburd/redigo/redis"
)

var pool *redis.Pool

func RedisInit() {
	pool = initPool()
}

func initPool() *redis.Pool {
	log.Debugf("conneciton redis: %s:%d", GetConfig().Rc.Host, GetConfig().Rc.Port)
	return redis.NewPool(func() (redis.Conn, error) {
		c, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", GetConfig().Rc.Host, GetConfig().Rc.Port))
		return c, err
	}, 20)
}

func Open() redis.Conn {
	if pool != nil {
		return pool.Get()
	}
	pool = initPool()
	return pool.Get()
}

func Close() {
	if pool != nil {
		pool.Close()
	}
}
