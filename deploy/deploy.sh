#!/bin/bash

#!/bin/bash

PRIVILEGED=false

curl -v -X POST http://10.3.33.5:8080/v2/apps -H Content-Type:application/json -d \
'{
      "id": "shurenyun-'$TASKENV'-'$SERVICE'",
      "cpus": '$CPUS',
      "mem": '$MEM',
      "instances": '$INSTANCES',
      "constraints": [["hostname", "LIKE", "'$DEPLOYIP'"], ["hostname", "UNIQUE"]],
      "container": {
                     "type": "DOCKER",
                     "docker": {
                                     "image": "'$SERVICE_IMAGE'",
                                     "network": "BRIDGE",
                                     "privileged": '$PRIVILEGED',
                                     "forcePullImage": '$FORCEPULLIMAGE',
                                     "portMappings": [
                                             { "containerPort": '$SERVICE_PORT', "hostPort": 0, "protocol": "tcp"}
                                     ]
                                }
                   },
      "env": {
                    "CONFIG_SERVER": "'$CONFIGSERVER'/config/'$TASKENV'",
                    "SERVICE": "cfgfile_'$TASKENV'_'$SERVICE'",
                    "BAMBOO_TCP_PORT": "'$SERVICE_PORT'",
                    "BAMBOO_PRIVATE": "true",
                    "BAMBOO_PROXY":"true",
                    "BAMBOO_BRIDGE": "true",
                    "BAMBOO_HTTP": "true"
             },
      "uris": [
               "'$CONFIGSERVER'/config/demo/config/registry/docker.tar.gz"
       ]
}'
