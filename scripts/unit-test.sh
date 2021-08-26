#!/bin/sh
source ./scripts/compile-plugins.sh
go test -v $(go list ./... | grep -v /test/)
