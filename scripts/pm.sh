#!/bin/sh

set -e

LAKE_ENDPOINT=${LAKE_ENDPOINT-'http://localhost:8080'}
LAKE_TASK_URL=$LAKE_ENDPOINT/task

debug() {
    scripts/compile-plugins.sh -gcflags=all="-N -l"
    dlv debug
}

run() {
    scripts/compile-plugins.sh
    go run main.go
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

"$@"
