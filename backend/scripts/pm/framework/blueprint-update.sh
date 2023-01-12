#!/bin/sh
#
# Licensed to the Apache Software Foundation (ASF) under one or more
# contributor license agreements.  See the NOTICE file distributed with
# this work for additional information regarding copyright ownership.
# The ASF licenses this file to You under the Apache License, Version 2.0
# (the "License"); you may not use this file except in compliance with
# the License.  You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

. "$(dirname $0)/../vars/active-vars.sh"

pipeline_id=${1-8}

curl -sv -XPATCH $LAKE_ENDPOINT/blueprints/$pipeline_id \
    -H "Content-Type: application/json" --data @- <<JSON | jq
{
    "name": "MY BLUEPRINT2",
    "cronConfig": "0 0 * * *",
    "settings": {
        "version": "1.0.0",
        "connections": [
            {
                "plugin": "jira",
                "connectionId": 1,
                "scope": [
                    {
                        "transformation": {
                            "epicKeyField": "customfield_10014",
                            "typeMappings": {
                                "缺陷": {
                                    "standardType": "Bug"
                                },
                                "线上事故": {
                                    "standardType": "Incident"
                                },
                                "故事": {
                                    "standardType": "Requirement"
                                }
                            },
                            "storyPointField": "customfield_10024",
                            "remotelinkCommitShaPattern": "/commit/([0-9a-f]{40})$"
                        },
                        "options": {
                            "boardId": 70
                        },
                        "entities": [
                            "TICKET",
                            "CROSS"
                        ]
                    }
                ]
            }
        ]
    },
    "enable": true,
    "mode": "NORMAL",
    "isManual": false
}
JSON
