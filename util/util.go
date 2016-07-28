package util

import (
	"encoding/base32"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var minute = int64(60)
var hour = minute * 60
var day = hour * 24

func ParseAppAlias(alis string) (string, error) {
	str := strings.Replace(strings.ToUpper(alis), "0", "=", -1)
	data, err := base32.StdEncoding.DecodeString(str)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func ParseTimeLen(timelen int64) string {
	if timelen >= day {
		return fmt.Sprintf("%d天 %d小时 %d分钟 %d秒", int64(timelen/day), int64(timelen%day/hour), int64(timelen%hour/minute), timelen%minute)
	} else if timelen >= hour {
		return fmt.Sprintf("%d小时 %d分钟 %d秒", int64(timelen/hour), int64(timelen%hour/minute), timelen%minute)
	} else if timelen >= minute {
		return fmt.Sprintf("%d分钟 %d秒", int64(timelen/minute), timelen%minute)
	} else {
		return fmt.Sprintf("%d秒", timelen)
	}
}

func ReturnOK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": data})
}

func ReturnDBError(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, gin.H{"code": 18001, "data": "", "error": err.Error()})
}

func ReturnParamError(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{"code": 18002, "data": "", "error": err})
}

func Header(c *gin.Context, key string) (string, bool) {
	v := c.Request.Header.Get(key)
	if v == "" {
		return "", false
	}
	return v, true
}
