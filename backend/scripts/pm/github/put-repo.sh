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

CONN_ID=${1-1}
SCOPE_ID=${2-384111310}
TR_ID=${3-1}
SCOPE_NAME=${4-"apache/incubator-devlake"}
SCOPE_URL=${5-"https://github.com/apache/incubator-devlake"}

curl -sv -XPUT $LAKE_ENDPOINT/plugins/github/connections/$CONN_ID/scopes \
    -H "Content-Type: application/json" \
    --data @- <<JSON | jq
    {

        "data": [
            {
                "connectionId": $CONN_ID,
                "githubId": $SCOPE_ID,
                "name": "$SCOPE_NAME",
                "htmlUrl": "$SCOPE_URL",
                "transformationRuleId": $TR_ID
            }
        ]
    }
JSON
