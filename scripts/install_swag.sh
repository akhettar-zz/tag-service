#!/usr/bin/env bash

set -e

echo "Installing swaggo command line tool"
go get -u github.com/swaggo/swag
cd $GOPATH/src/github.com/swaggo/swag/cmd/swag
go build && go install
echo "swag install complete."

