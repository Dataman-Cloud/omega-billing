package model

type Config struct {
	Host string
	Port uint16
	Log  LogConfig   `mapstructure:"log"`
	Mq   MqConfig    `mapstructure:"mq"`
	Mc   MysqlConfig `mapstructure:"mysql"`
	Rc   RedisConfig `mapstructure:"redis"`
}

type LogConfig struct {
	Console    bool
	AppendFile bool
	File       string
	Level      string
	Formatter  string
	MaxSize    uint64
}

type MysqlConfig struct {
	UserName     string
	PassWord     string
	Host         string
	Port         uint16
	DataBase     string
	MaxIdleConns uint16
	MaxOpenConns uint16
}

type MqConfig struct {
	User        string
	PassWord    string
	Host        string
	Port        uint16
	QueueTTL    uint64
	MessageTTL  uint64
	ConsumeName string
}

type RedisConfig struct {
	Host string
	Port uint16
}
