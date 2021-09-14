#!/bin/sh

set -e

SCRIPT_DIR="$( cd "$( dirname "$0" )" && pwd )"

LAKE_ENDPOINT=${LAKE_ENDPOINT-'http://localhost:8080'}
LAKE_TASK_URL=$LAKE_ENDPOINT/task

debug() {
    $SCRIPT_DIR/compile-plugins.sh -gcflags=all="-N -l"
    dlv debug
}

run() {
    $SCRIPT_DIR/compile-plugins.sh
    go run $SCRIPT_DIR/../main.go
}

jira() {
    curl -v -XPOST $LAKE_TASK_URL --data @- <<'    JSON'
    [
        {
            "plugin": "jira",
            "options": {
                "boardId": 8
            }
        }
    ]
    JSON
}

jira_enrich_issues() {
    curl -v -XPOST $LAKE_TASK_URL --data @- <<'    JSON'
    [
        {
            "plugin": "jira",
            "options": {
                "boardId": 8,
                "tasks": [ "enrichIssues" ]
            }
        }
    ]
    JSON
}

all() {
    curl -v -XPOST $LAKE_TASK_URL --data @- <<'    JSON'
    [
            {
                "plugin": "gitlab",
                "options": {
                    "projectId": 8967944
                }
            },
            {
                "plugin": "jira",
                "options": {
                    "boardId": 8
                }
            },
            {
                "plugin": "jenkins",
                "options": {}
            }
    ]
    JSON
}

gitlab() {
    curl -v -XPOST $LAKE_TASK_URL --data @- <<'    JSON'
    [
            {
                "plugin": "gitlab",
                "options": {
                    "projectId": 8967944
                }
            }
    ]
    JSON
}

jenkins() {
    curl -v -XPOST $LAKE_TASK_URL --data @- <<'    JSON'
    [
            {
                "plugin": "jenkins",
                "options": {}
            }
    ]
    JSON
}

truncate() {
    SQL=$()
    echo "SET FOREIGN_KEY_CHECKS=0;"
    echo 'show tables' | mycli local-lake | tail -n +2 | xargs -I{} -n 1 echo "truncate table {};"
    echo "SET FOREIGN_KEY_CHECKS=1;"
}

tasks() {
    curl -v $LAKE_TASK_URL?status=$1 | jq
}

"$@"
