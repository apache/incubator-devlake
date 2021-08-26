#!/bin/sh
source ./scripts/compile-plugins.sh
source ./scripts/export-env.sh
go test -v ./test/...
