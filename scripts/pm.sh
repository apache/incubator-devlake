#!/bin/sh

set -e

SCRIPT_DIR="$( cd "$( dirname "$0" )" && pwd )"

LAKE_ENDPOINT=${LAKE_ENDPOINT-'http://localhost:8080'}
LAKE_PIPELINE_URL=$LAKE_ENDPOINT/pipelines

debug() {
    $SCRIPT_DIR/compile-plugins.sh -gcflags=all="-N -l"
    dlv debug
}

run() {
    $SCRIPT_DIR/compile-plugins.sh
    go run $SCRIPT_DIR/../main.go
}

jira_source_post() {
    curl -v -XPOST "$LAKE_ENDPOINT/plugins/jira/sources" --data '
    {
        "name": "test-jira-source",
        "endpoint": "'"$JIRA_ENDPOINT"'",
        "basicAuthEncoded": "'"$JIRA_BASIC_AUTH_ENCODED"'",
        "epicKeyField": "'"$JIRA_ENDPOINT"'",
        "storyPointField": "'"$JIRA_ISSUE_STORYPOINT_FIELD"'",
    }
    ' | jq
}

jira_source_post_full() {
    curl -v -XPOST "$LAKE_ENDPOINT/plugins/jira/sources" --data '
    {
        "name": "test-jira-source",
        "endpoint": "'"$JIRA_ENDPOINT"'",
        "basicAuthEncoded": "'"$JIRA_BASIC_AUTH_ENCODED"'",
        "epicKeyField": "'"$JIRA_ENDPOINT"'",
        "storyPointField": "'"$JIRA_ISSUE_STORYPOINT_FIELD"'",
        "typeMappings": {
            "Story": {
                "standardType": "Requirement",
                "statusMappings": {
                    "已完成": {
                        "standardStatus": "Resolved"
                    },
                    "已解决": {
                        "standardStatus": "Resolved"
                    }
                }
            },
            "Incident": {
                "standardType": "Incident",
                "statusMappings": {
                    "已完成": {
                        "standardStatus": "Resolved"
                    }
                }
            },
            "Bug": {
                "standardType": "Bug",
                "statusMappings": {
                    "已完成": {
                        "standardStatus": "Resolved"
                    }
                }
            }
        }
    }' | jq
}

jira_source_post_fail() {
    curl -v -XPOST "$LAKE_ENDPOINT/plugins/jira/sources" --data @- <<'    JSON' | jq
    {
        "name": "test-jira-source-fail",
        "endpoint": "https://merico.atlassian.net/rest",
        "basicAuthEncoded": "basicAuth",
        "epicKeyField": "epicKeyField",
        "storyPointField": "storyPointField",
        "typeMappings": "ehhlow"
    }
    JSON
}

jira_source_put() {
    curl -v -XPUT "$LAKE_ENDPOINT/plugins/jira/sources/$1" --data @- <<'    JSON' | jq
    {
        "name": "test-jira-source-updated",
        "endpoint": "https://merico.atlassian.net/rest",
        "basicAuthEncoded": "basicAuth",
        "epicKeyField": "epicKeyField",
        "storyPointField": "storyPointField",
    }
    JSON
}

jira_source_put_full() {
    curl -v -XPUT "$LAKE_ENDPOINT/plugins/jira/sources/$1" --data '
    {
        "name": "test-jira-source-updated",
        "endpoint": "'"$JIRA_ENDPOINT"'",
        "basicAuthEncoded": "'"$JIRA_BASIC_AUTH_ENCODED"'",
        "epicKeyField": "'"$JIRA_ENDPOINT"'",
        "storyPointField": "'"$JIRA_ISSUE_STORYPOINT_FIELD"'",
        "typeMappings": {
            "Story": {
                "standardType": "Requirement",
                "statusMappings": {
                    "已完成": {
                        "standardStatus": "Resolved"
                    },
                    "已解决": {
                        "standardStatus": "Resolved"
                    }
                }
            },
            "Incident": {
                "standardType": "Incident",
                "statusMappings": {
                    "已完成": {
                        "standardStatus": "Resolved"
                    }
                }
            },
            "Bug": {
                "standardType": "Bug",
                "statusMappings": {
                    "已完成": {
                        "standardStatus": "Resolved"
                    }
                }
            }
        }
    }' | jq
}

jira_source_list() {
    curl -v "$LAKE_ENDPOINT/plugins/jira/sources" | jq
}

jira_source_get() {
    curl -v "$LAKE_ENDPOINT/plugins/jira/sources/$1" | jq
}

jira_source_delete() {
    curl -v -XDELETE "$LAKE_ENDPOINT/plugins/jira/sources/$1"
}

jira_typemapping_post() {
    curl -v -XPOST "$LAKE_ENDPOINT/plugins/jira/sources/$1/type-mappings" --data @- <<'    JSON' | jq
    {
        "userType": "userType",
        "standardType": "standardType"
    }
    JSON
}

jira_typemapping_put() {
    curl -v -XPUT "$LAKE_ENDPOINT/plugins/jira/sources/$1/type-mappings/$2" --data @- <<'    JSON' | jq
    {
        "standardType": "standardTypeUpdated"
    }
    JSON
}

jira_typemapping_delete() {
    curl -v -XDELETE "$LAKE_ENDPOINT/plugins/jira/sources/$1/type-mappings/$2"
}

jira_typemapping_list() {
    curl -v "$LAKE_ENDPOINT/plugins/jira/sources/$1/type-mappings" | jq
}

jira_statusmapping_post() {
    curl -v -XPOST "$LAKE_ENDPOINT/plugins/jira/sources/$1/type-mappings/$2/status-mappings" --data @- <<'    JSON' | jq
    {
        "userStatus": "userStatus",
        "standardStatus": "standardStatus"
    }
    JSON
}

jira_statusmapping_put() {
    curl -v -XPUT "$LAKE_ENDPOINT/plugins/jira/sources/$1/type-mappings/$2/status-mappings/$3" --data @- <<'    JSON' | jq
    {
        "standardStatus": "standardStatusUpdated"
    }
    JSON
}

jira_statusmapping_delete() {
    curl -v -XDELETE "$LAKE_ENDPOINT/plugins/jira/sources/$1/type-mappings/$2/status-mappings/$3"
}

jira_statusmapping_list() {
    curl -v "$LAKE_ENDPOINT/plugins/jira/sources/$1/type-mappings/$2/status-mappings" | jq
}

jira_echo() {
    curl -v -XPOST "$LAKE_ENDPOINT/plugins/jira/echo" --data @- <<'    JSON' | jq
    {
        "plugin": "jira",
        "options": {
            "boardId": 8
        }
    }
    JSON
}

pipeline_new() {
    curl -v -XPOST $LAKE_PIPELINE_URL --data @- <<'    JSON' | jq
    {
        "name": "test-all",
        "tasks": [
            [
                {
                    "plugin": "jira",
                    "options": {
                        "sourceId": 1,
                        "boardId": 8
                    }
                },
                {
                    "plugin": "jenkins",
                    "options": {}
                }
            ]
        ]
    }
    JSON
}

pipelines() {
    curl -v $LAKE_PIPELINE_URL'?'$1 | jq
}

pipeline() {
    curl -v $LAKE_PIPELINE_URL/$1 | jq
}

pipeline_cancel() {
    curl -v -XDELETE $LAKE_PIPELINE_URL/$1
}

pipeline_tasks() {
    curl -v $LAKE_PIPELINE_URL/$1/tasks'?'$2 | jq
}

jira() {
    curl -v -XPOST $LAKE_PIPELINE_URL --data '
    {
        "name": "test-jira",
        "tasks": [
            [
                {
                    "plugin": "jira",
                    "options": {
                        "sourceId": '$1',
                        "boardId": '$2',
                        "tasks": ['"$3"']
                    }
                }
            ]
        ]
    }
    ' | jq
}

gitlab() {
    curl -v -XPOST $LAKE_PIPELINE_URL --data @- <<'    JSON'
    {
        "name": "test-gitlab",
        "tasks": [
            [
                {
                    "plugin": "gitlab",
                    "options": {
                        "projectId": 8967944,
                        "tasks": ["collectMrs"]
                    }
                }
            ]
        ]
    }
    JSON
}

github() {
    curl -v -XPOST $LAKE_PIPELINE_URL --data @- <<'    JSON'
    {
        "name": "test-github",
        "tasks": [
            [
                {
                    "plugin": "github",
                    "options": {
                        "repositoryName": "lake",
                        "owner": "merico-dev",
                        "tasks": ["collectIssues"]
                    }
                }
            ]
        ]
    }
    JSON
}

jenkins() {
    curl -v -XPOST $LAKE_PIPELINE_URL --data @- <<'    JSON'
    {
        "name": "test-jenkins",
        "tasks": [
            [
                {
                    "plugin": "jenkins",
                    "options": {}
                }
            ]
        ]
    }
    JSON
}

ae() {
    curl -v -XPOST $LAKE_PIPELINE_URL --data @- <<'    JSON'
    {
        "name": "test-ae",
        "tasks": [
            [
                {
                    "plugin": "ae",
                    "options": {
                        "projectId": 13
                    }
                }
            ]
        ]
    }
    JSON
}

truncate() {
    SQL=$()
    echo "SET FOREIGN_KEY_CHECKS=0;"
    echo 'show tables' | mycli local-lake | tail -n +2 | xargs -I{} -n 1 echo "truncate table {};"
    echo "SET FOREIGN_KEY_CHECKS=1;"
}

lint() {
    golangci-lint run -v
}

"$@"
