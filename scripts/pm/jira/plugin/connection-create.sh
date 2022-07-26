#!/bin/sh

. "$(dirname $0)/../../vars/active-vars.sh"

jira_endpoint=${1-$JIRA_ENDPOINT}
jira_username=${2-$JIRA_USERNAME}
jira_password=${3-$JIRA_PASSWORD}

curl -sv $LAKE_ENDPOINT/plugins/jira/connections --data @- <<JSON | jq
{
    "name": "testjira",
    "endpoint": "$jira_endpoint",
    "username": "$jira_username",
    "password": "$jira_password"
}
JSON
