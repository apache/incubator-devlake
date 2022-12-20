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

request() {
  response="$(curl -s "$LAKE_ENDPOINT/$1" | jq --color-output)"
  echo "
  INDEX: $2 URL: $LAKE_ENDPOINT/$1
  $response
  " | less --RAW-CONTROL-CHARS
  read PAUSE
}

start=$1
length=$2
counter=0
for url in \
"pipelines" \
"pipelines?pageSize=1&page=1" \
"pipelines?pageSize=1&page=2" \
"pipelines?blueprint_id=1" \
"pipelines?status=TASK_FAILED" \
"pipelines?pending=1" \
"pipelines?label=foobar" \
"blueprints" \
"blueprints?pageSize=1&page=1" \
"blueprints?pageSize=1&page=2" \
"blueprints?enable=1" \
"blueprints?is_manual=0" \
"blueprints?label=foobar" \
; do
  counter=$(($counter+1))
  if [ "$start" = "-l" ]; then
    printf "#%-3s %s\n" "$counter" "$url"
    continue
  fi
  if [ -n "$start" ] && [ "$counter" -lt "$start" ]; then
    continue
  fi
  if [ -n "$length" ] && [ "$counter" -ge $(($start+$length)) ]; then
    continue
  fi
  request "$url" $counter
done