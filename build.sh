#!/bin/bash

# killall sqlite-http
set -e
set -x

go get github.com/mattn/go-sqlite3

gitBranch=`git rev-parse --abbrev-ref HEAD`
gitCommit=`git rev-parse HEAD`
gitDescribe=`git describe --dirty --tags`

GOOS=linux 
GOARCH=amd64 

rm -f sqlite-http
go build -ldflags="-s -w -X main.GitBranch=${gitBranch} -X main.GitCommit=${gitCommit} -X main.GitDescribe=${gitDescribe}" -v github.com/ludenus/sqlite-http


if [[ "upx" == "$1" ]]; then
    upx --ultra-brute ./sqlite-http
fi