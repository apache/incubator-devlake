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

set -e

notes() {
    curl -v "$GITLAB_ENDPOINT/projects/8967944/merge_requests/1349/notes?system=false&per_page=100&page=1" \
        -H "Authorization: Bearer $GITLAB_AUTH"
}

commits() {
    SIZE=${1-100}
    PAGE=${2-1}
    PROJ=${3-8967944}
    curl -v "$GITLAB_ENDPOINT/projects/$PROJ/repository/commits?with_stats=true&per_page=$SIZE&page=$PAGE" \
        -H "Authorization: Bearer $GITLAB_AUTH"
}

"$@"
