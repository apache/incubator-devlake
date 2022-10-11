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


# jira server user name
USER=
# jira server password
PASSWORD=
# FQDN of your jira server, i.e. https://jira.example.com:8080
HOST=
# id of a Scrum Board with more than 3000 issues
BOARD=8
PAGES=30

    # -w "%{http_code}" \
# check if settings are correct
curl -H "Accept: application/json" -vv \
    -u "$USER:$PASSWORD" \
    $HOST/rest/api/2/serverInfo 2> /tmp/jira-server-stderr

STATUS=$(grep -oP '^< HTTP.* \K([0-9]{3})' /tmp/jira-server-stderr)
if [ "$STATUS" != 200 ]; then
    awk '{ if ($0 ~ /authorization/) $0 = "==AUTHORIZATION FILTERED=="; print }' /tmp/jira-server-stderr
    echo "Failed to fetch jira server information, please make sure USER/PASSWORD/HOST and BOARD are set correct."
    exit 1
fi

echo
echo "Settings looks good, start benchmarking"


# benchmark api response time and save to result.csv
echo "EXPAND,PAGE_SIZE,PAGE,STATUS_CODE,ELLAPSE" > result.csv
for PAGE_SIZE in 50 100; do
    for EXPAND in "" "changelog"; do
        for PAGE in $(seq 1 $PAGES); do
            STARTAT=$((PAGE_SIZE * (PAGE - 1)))
            URL="$HOST/rest/agile/1.0/board/$BOARD/issue?expand=$EXPAND&jql=ORDER+BY+created+ASC&maxResults=100&startAt=$STARTAT"
            printf "%s" "$URL" 
            STARTED_TS=$(date +%s)
            STATUS=$(
            curl -H "Accept: application/json" \
                -u "$USER:$PASSWORD" \
                -s -o /dev/null -w "%{http_code}" \
                $URL
            )
            ENDED_TS=$(date +%s)
            ELLAPSE=$(( ENDED_TS - STARTED_TS ))
            printf " %s %s\n" "$STATUS" "$ELLAPSE"
            echo "$EXPAND","$PAGE_SIZE","$PAGE","$STATUS","$ELLAPSE" >> result.csv
        done
    done
done
