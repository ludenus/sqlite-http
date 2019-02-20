#!/bin/bash

# killall sqlite-http
set -e
set -x

gitBranch=`git rev-parse --abbrev-ref HEAD`
gitCommit=`git rev-parse HEAD | cut -c-8`

GOOS=linux 
GOARCH=amd64 

rm -f sqlite-http
go build -ldflags="-s -w -X main.GitBranch=${gitBranch} -X main.GitCommit=${gitCommit}" -v github.com/ludenus/sqlite-http


if [[ "upx" == "$1" ]]; then
    upx --ultra-brute ./sqlite-http
fi