#!/bin/sh
mkdir -p /go/src/github.com/Dataman-Cloud/omega-billing
mkdir -p $HOME/.omega/
mkdir  /etc/omega/
export GOPATH=/go
export GO15VENDOREXPERIMENT=1

cp -r /src/* /go/src/github.com/Dataman-Cloud/omega-billing
#rm /etc/localtime && cd /go/src/github.com/Dataman-Cloud/omega-billing && mv localtime /etc
cd /go/src/github.com/Dataman-Cloud/omega-billing && \
mv start.sh /bin/ && \
mv omega-billing.yaml.sample /etc/omega/omega-billing.yaml && \
mv sql /bin/sql && \
go build && mv omega-billing /bin/omega-billing
#apk del go
rm -rf /go
rm -rf /src
rm -rf /var/cache/apk/*
