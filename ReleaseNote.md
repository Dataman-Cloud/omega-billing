#Release Note

更新版本v0.1.050900
 * 修改配置文件`omega-billing.yaml`
```
host: "0.0.0.0"
port: 5013
log:
  console: true
  appendfile: false
  file: "./log/omega-audit.log"
  level: "debug"
  formatter: "text"
  maxSize: 1024000
mq:
  user: "guest"
  password: "guest"
  host: "10.3.20.53"
  port: 5672
  queueTTL: 86400000
  messageTTL: 300000
  exchange: "application"
  routingkey: "cluster"
  consumeName: "billing"
mysql:
  username: "root"
  password: "dataman1234"
  host: "127.0.0.1"
  port: 3306
  database: "billing"
  maxIdleConns: 5
  maxOpenConns: 50
redis:
  host: "localhost"
  port: 6379
```
 * 增加配置文件项
```
location ~ /api/v3/billing {
  if ($request_method = OPTIONS ) {
      add_header Access-Control-Allow-Origin "*" ;
      add_header Access-Control-Allow-Methods "GET,PUT,POST,DELETE,OPTIONS,PATCH";
      add_header Access-Control-Allow-Headers "Content-Type, Depth, User-Agent, X-File-Size, X-Requested-With, X-Requested-By, If-Modified-Since, X-File-Name, Cache-Control, X-XSRFToken, Authorization" ;
      add_header Access-Control-Allow-Credentials "true" ;
      add_header Content-Length 0 ;
      add_header Content-Type application/json ;
      return 204;
  }
  if ($request_method != 'OPTIONS') {
      add_header 'Access-Control-Allow-Origin' '*' always;
      add_header 'Access-Control-Allow-Credentials' 'true' always;
      add_header 'Access-Control-Allow-Methods' 'GET,PUT,POST,DELETE,OPTIONS,PATCH' always;
      add_header 'Access-Control-Allow-Headers' 'DNT,X-CustomHeader,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type' always;
  }
  auth_request    /_auth;
  proxy_pass      http://10.3.20.53:5013;
}
```
