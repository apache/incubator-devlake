#!/bin/sh
set -e
./scripts/export-env.sh
./scripts/compile-plugins.sh
go test -v $(go list ./... | grep -v /test/)
