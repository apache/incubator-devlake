#!/bin/sh

. "$(dirname $0)/../vars/active-vars.sh"

curl -sv $LAKE_ENDPOINT/pipelines --data @- <<JSON | jq
{
    "name": "test-all",
    "plan": [
        [
            {
                "plugin": "jira",
                "options": {
                    "connectionId": 1,
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
