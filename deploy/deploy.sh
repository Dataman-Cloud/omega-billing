#!/bin/bash

PRIVILEGED=false

curl -v -X PUT $MARATHON_API_URL/v2/apps/shurenyun-$TASKENV-$SERVICE -H Content-Type:application/json -d \
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
                                             { "containerPort": '$BILLING_NET_PORT', "hostPort": 0, "protocol": "tcp"}
                                     ]
                                }
                   },
      "healthChecks": [{
               "path": "/api/v3/health/billing",
               "protocol": "HTTP",
               "gracePeriodSeconds": 300,
               "intervalSeconds": 60,
               "portIndex": 0,
               "timeoutSeconds": 20,
               "maxConsecutiveFailures": 3
           }],
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
