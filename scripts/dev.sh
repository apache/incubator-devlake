#!/bin/sh
go build -o lake
source ./scripts/export-env.sh
./lake
