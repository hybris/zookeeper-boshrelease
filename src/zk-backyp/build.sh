#!/bin/bash
set -e
set -x

go get github.com/mitchellh/goamz/aws
go get github.com/mitchellh/goamz/s3

go build -o bin/zk-backyp src/zk-backyp.go
