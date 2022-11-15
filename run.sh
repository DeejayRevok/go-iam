#!/bin/sh
service filebeat start
service metricbeat start
if [ "${DEVELOPMENT}" = "true" ]; then
    go run github.com/cosmtrek/air@latest
else 