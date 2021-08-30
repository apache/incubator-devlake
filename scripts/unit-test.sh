#!/bin/sh
set -e
source ./scripts/export-env.sh
source ./scripts/compile-plugins.sh
go test -v $(go list ./... | grep -v /test/)
