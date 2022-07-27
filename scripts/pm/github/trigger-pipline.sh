#!/bin/sh

. "$(dirname $0)/../vars/active-vars.sh"

curl -sv $LAKE_ENDPOINT/pipelines --data @- <<JSON | jq
{
    "name": "test-github",
    "plan": [
        [
            {
                "plugin": "github",
                "options": {
                    "connectionId": 1,
                    "owner": "apache",
                    "repo": "incubator-devlake"
                }
            }
        ]
    ]
}
JSON
