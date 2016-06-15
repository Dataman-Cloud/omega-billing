##omega-billing环境变量说明
注释括号里面代表原来相对应配置文件的字段


```
BILLING_NET_HOST=0.0.0.0 				#billing服务监听地址  (host)
BILLING_NET_PORT=5013					#billing服务监听端口  (port)
BILLING_LOG_CONSOLE=true				#服务日志是否打开标准输出 (log.console)
BILLING_LOG_APPEND_FILE=false			#服务日志是否输出到文件  (log.appendfile)
BILLING_LOG_FILE=./log/omega-audit.log #服务日志输出文件位置  (log.file)
BILLING_LOG_LEVEL=debug					#服务日志打印等级   (log.level)
BILLING_LOG_FORMATTER=text				#服务日志打印格式   (log.formatter)
BILLING_LOG_MAX_SIZE=1024000			#服务日志最大字节数  (log.maxSize)
BILLING_MQ_USER=guest					#rabbitmq用户名    (mq.user)
BILLING_MQ_PASSWD=guest					#rabbitmq密码     (mq.password)
BILLING_MQ_HOST=10.3.20.53				#rabbitmq地址     (mq.host)
BILLING_MQ_PORT=5672					#rabbitmq端口      (mq.port)
BILLING_MQ_QUEUE_TTL=86400000			#rabbitmq 队列ttl  (mq.queueTTL)
BILLING_MQ_MSG_TTL=300000				#rabbitmq 消息     (mq.messageTTL)
BILLING_MQ_EXCHANGE=cluster				#rabbitmq exchange名称
BILLING_MQ_ROUTE_KEY=cluster			#rabbitmq routingkey名称(mq.routingkey)
BILLING_MQ_CONSUME_NAME=billing			#rabbitmq 消费队列名称(mq.consumeName)
BILLING_MYSQL_USER=root					#mysql用户名  (mysql.username)
BILLING_MYSQL_PASSWD=111111				#mysql密码  (mysql.password)
BILLING_MYSQL_HOST=127.0.0.1			#mysql地址   (mysql.host)
BILLING_MYSQL_PORT=3306					#mysql端口   (mysql.port)
BILLING_MYSQL_DB=billing				#mysql数据库名(mysql.db)
BILLING_MYSQL_MAX_IDLE_CONNS=5			#mysql最大闲置链接数(mysql.maxIdleConns)
BILLING_MYSQL_MAX_OPEN_CONNS=50			#mysql最大打开链接数(mysql.maxOpenConns)
BILLING_REDIS_HOST=localhost			#redis地址
BILLING_REDIS_PORT=6379					#redis端口
```
