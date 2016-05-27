package config

import (
	"github.com/Dataman-Cloud/omega-billing/model"
	log "github.com/cihub/seelog"
	"github.com/spf13/viper"
)

var config model.Config

func init() {
	viper.SetConfigName("omega-billing")
	viper.AddConfigPath("./")
	viper.AddConfigPath("/etc/omega/")
	viper.AddConfigPath("$HOME/.omega/")
	viper.AddConfigPath("/")
	if err := viper.ReadInConfig(); err != nil {
		log.Errorf("read config omega-billing.yaml error: %v", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Errorf("unmarshal config file error: %v", err)
	}
}

func GetConfig() model.Config {
	return config
}
