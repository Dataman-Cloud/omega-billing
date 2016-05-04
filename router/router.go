package router

import (
	"fmt"
	. "github.com/Dataman-Cloud/omega-billing/config"
	"github.com/Dataman-Cloud/omega-billing/controller"
	"github.com/Dataman-Cloud/omega-billing/util"
	mq "github.com/Dataman-Cloud/omega-billing/util/rabbitmq"
	rd "github.com/Dataman-Cloud/omega-billing/util/redis"
	log "github.com/cihub/seelog"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func Run() {
	go mq.Run()
	defer func() {
		mq.Close()
		rd.Close()
	}()
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// Global middleware
	r.Use(gin.Logger(), gin.Recovery())

	log.Debugf("http listening %s:%d", GetConfig().Host, GetConfig().Port)

	v3 := r.Group("/api/v3/billing", Authenticate)
	{
		v3.GET("/list", controller.BillingList)
	}

	s := &http.Server{
		Addr:           fmt.Sprintf("%s:%d", GetConfig().Host, GetConfig().Port),
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}

func Authenticate(c *gin.Context) {
	authenticated := false
	if val, ok := util.Header(c, "Authorization"); ok {
		conn := rd.Open()
		defer conn.Close()
		uid, err := redis.String(conn.Do("HGET", "s:"+val, "user_id"))
		if err == nil {
			authenticated = true
			c.Set("uid", uid)
		} else if err != redis.ErrNil {
			log.Errorf("billing get error:", err)
		}
	}
	if authenticated {
		c.Next()
	} else {
		c.String(http.StatusUnauthorized, "Invalid authentication")
		c.Abort()
	}
}
