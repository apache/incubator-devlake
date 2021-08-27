#!/bin/sh
set -e
./scripts/compile-plugins.sh
./scripts/export-env.sh
go test -v $(go list ./... | grep -v /test/)
