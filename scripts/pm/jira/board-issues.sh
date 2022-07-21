#!/bin/sh

. "$(dirname $0)/../vars/active-vars.sh"

board_id=${1-"8"}

curl -sv $JIRA_ENDPOINT/agile/1.0/board/$board_id/issue \
    -H "Authorization: Basic $JIRA_BASIC_AUTH" \
    | jq
