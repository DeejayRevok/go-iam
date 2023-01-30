#!/bin/sh
nohup filebeat -e -c /etc/filebeat/filebeat.yml &
nohup metricbeat -e -c /etc/metricbeat/metricbeat.yml &
if [ "${DEVELOPMENT}" = "true" ]; then
    go run github.com/cosmtrek/air@latest
else
    go run main.go
fi
