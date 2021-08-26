#!/bin/sh
# for PLUG in $(find plugins/* -type d -not -name core); do
#   NAME=$(basename $PLUG)
#   # go build -buildmode=plugin -o $PLUG/$NAME.so $PLUG/$NAME.go
# done
  go build -buildmode=plugin -o plugins/jira/jira.so plugins/jira/jira.go
  # go build -buildmode=plugin -o plugins/jira/init.so plugins/jira/init.go
