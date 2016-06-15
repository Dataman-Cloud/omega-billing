package config

import (
	"bufio"
	"errors"
	"github.com/Dataman-Cloud/omega-billing/model"
	log "github.com/Sirupsen/logrus"
	"os"
	"reflect"
	"strconv"
	"strings"
)

var config model.Config

func GetConfig() model.Config {
	return config
}

type EnvEntry struct {
	BILLING_HOST                string `required:"true"`
	BILLING_PORT                uint16 `required:"true"`
	BILLING_LOG_CONSOLE         bool   `required:"true"`
	BILLING_LOG_APPEND_FILE     bool   `required:"true"`
	BILLING_LOG_FILE            string `required:"true"`
	BILLING_LOG_LEVEL           string `required:"true"`
	BILLING_LOG_FORMATTER       string `required:"true"`
	BILLING_LOG_MAX_SIZE        uint64 `required:"true"`
	BILLING_MQ_USER             string `required:"true"`
	BILLING_MQ_PASSWD           string `required:"true"`
	BILLING_MQ_HOST             string `required:"true"`
	BILLING_MQ_PORT             uint16 `required:"true"`
	BILLING_MQ_QUEUE_TTL        int64  `required:"true"`
	BILLING_MQ_MSG_TTL          int64  `required:"true"`
	BILLING_MQ_EXCHANGE         string `required:"true"`
	BILLING_MQ_ROUTE_KEY        string `required:"true"`
	BILLING_MQ_CONSUME_NAME     string `required:"true"`
	BILLING_MSQL_USER           string `required:"true"`
	BILLING_MSQL_PASSWD         string `required:"true"`
	BILLING_MSQL_HOST           string `required:"true"`
	BILLING_MSQL_PORT           uint16 `required:"true"`
	BILLING_MSQL_DB             string `required:"true"`
	BILLING_MSQL_MAX_IDLE_CONNS uint16 `required:"true"`
	BILLING_MSQL_MAX_OPEN_CONNS uint16 `required:"true"`
	BILLING_REDIS_HOST          string `required:"true"`
	BILLING_REDIS_PORT          uint16 `required:"true"`
}

func InitConfig(envFile string) *model.Config {
	loadEnvFile(envFile)

	envEntry := NewEnvEntry()

	config.Host = envEntry.BILLING_HOST
	config.Port = envEntry.BILLING_PORT

	config.Log.AppendFile = envEntry.BILLING_LOG_APPEND_FILE
	config.Log.Console = envEntry.BILLING_LOG_CONSOLE
	config.Log.Formatter = envEntry.BILLING_LOG_FORMATTER
	config.Log.Level = envEntry.BILLING_LOG_LEVEL
	config.Log.File = envEntry.BILLING_LOG_FILE
	config.Log.MaxSize = envEntry.BILLING_LOG_MAX_SIZE

	config.Mc.DataBase = envEntry.BILLING_MSQL_DB
	config.Mc.Host = envEntry.BILLING_MSQL_HOST
	config.Mc.Port = envEntry.BILLING_MSQL_PORT
	config.Mc.UserName = envEntry.BILLING_MSQL_USER
	config.Mc.PassWord = envEntry.BILLING_MSQL_PASSWD
	config.Mc.MaxIdleConns = envEntry.BILLING_MSQL_MAX_IDLE_CONNS
	config.Mc.MaxOpenConns = envEntry.BILLING_MSQL_MAX_OPEN_CONNS

	config.Mq.ConsumeName = envEntry.BILLING_MQ_CONSUME_NAME
	config.Mq.Exchange = envEntry.BILLING_MQ_EXCHANGE
	config.Mq.Host = envEntry.BILLING_MQ_HOST
	config.Mq.MessageTTL = envEntry.BILLING_MQ_MSG_TTL
	config.Mq.PassWord = envEntry.BILLING_MQ_PASSWD
	config.Mq.Port = envEntry.BILLING_MQ_PORT
	config.Mq.QueueTTL = envEntry.BILLING_MQ_QUEUE_TTL
	config.Mq.Routingkey = envEntry.BILLING_MQ_ROUTE_KEY
	config.Mq.User = envEntry.BILLING_MQ_USER

	config.Rc.Host = envEntry.BILLING_REDIS_HOST
	config.Rc.Port = envEntry.BILLING_REDIS_PORT

	return &config
}

func NewEnvEntry() *EnvEntry {
	envEntry := &EnvEntry{}

	val := reflect.ValueOf(envEntry).Elem()

	for i := 0; i < val.NumField(); i++ {
		typeField := val.Type().Field(i)
		required := typeField.Tag.Get("required")

		env := os.Getenv(typeField.Name)

		if env == "" && required == "true" {
			exitMissingEnv(typeField.Name)
		}

		var envEntryValue interface{}
		var err error
		valueFiled := val.Field(i).Interface()
		value := val.Field(i)
		switch valueFiled.(type) {
		case int64:
			envEntryValue, err = strconv.ParseInt(env, 10, 64)

		case int16:
			envEntryValue, err = strconv.ParseInt(env, 10, 16)
			_, ok := envEntryValue.(int64)
			if !ok {
				continue
			}
			envEntryValue = int16(envEntryValue.(int64))
		case uint16:
			envEntryValue, err = strconv.ParseUint(env, 10, 16)

			_, ok := envEntryValue.(uint64)
			if !ok {
				continue
			}
			envEntryValue = uint16(envEntryValue.(uint64))
		case uint64:
			envEntryValue, err = strconv.ParseUint(env, 10, 64)
		case bool:
			envEntryValue, err = strconv.ParseBool(env)
		default:
			envEntryValue = env
		}

		if err != nil {
			exitCheckEnv(typeField.Name, err)
		}
		value.Set(reflect.ValueOf(envEntryValue))
	}

	return envEntry
}

func loadEnvFile(envfile string) {
	// load the environment file
	f, err := os.Open(envfile)
	if err == nil {
		defer f.Close()

		r := bufio.NewReader(f)
		for {
			line, _, err := r.ReadLine()
			if err != nil {
				break
			}

			key, val, err := parseln(string(line))
			if err != nil {
				continue
			}

			if len(os.Getenv(strings.ToUpper(key))) == 0 {
				err1 := os.Setenv(strings.ToUpper(key), val)
				if err1 != nil {
					log.Errorln(err1.Error())
				}
			}
		}
	}
}

// helper function to parse a "key=value" environment variable string.
func parseln(line string) (key string, val string, err error) {
	line = removeComments(line)
	if len(line) == 0 {
		return
	}
	splits := strings.SplitN(line, "=", 2)

	if len(splits) < 2 {
		err = errors.New("missing delimiter '='")
		return
	}

	key = strings.Trim(splits[0], " ")
	val = strings.Trim(splits[1], ` "'`)
	return

}

// helper function to trim comments and whitespace from a string.
func removeComments(s string) (_ string) {
	if len(s) == 0 || string(s[0]) == "#" {
		return
	} else {
		index := strings.Index(s, " #")
		if index > -1 {
			s = strings.TrimSpace(s[0:index])
		}
	}
	return s
}

func exitMissingEnv(env string) {
	log.Fatalf("program exit missing config for env %s", env)
	os.Exit(1)
}

func exitCheckEnv(env string, err error) {
	log.Fatalf("Check env %s, %s", env, err.Error())

}
