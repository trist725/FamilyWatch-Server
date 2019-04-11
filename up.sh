#!/bin/sh
WORK_DIR=$(dirname $(readlink -f $0))

cd $WORK_DIR
go run -mod=vendor main.go -p familywatch
