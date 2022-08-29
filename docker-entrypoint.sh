#!/bin/sh
if [ "${DEVELOPMENT}" = "true" ]; then
    go run github.com/cosmtrek/air@latest
else 
    go run main.go 
fi