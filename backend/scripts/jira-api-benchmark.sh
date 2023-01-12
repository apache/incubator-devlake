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
PAGE_SIZE=100
FILENAME=jira-server-benchmark-$(date +"%Y-%m-%d_%H_%M").csv

    # -w "%{http_code}" \
# check if settings are correct
curl -H "Accept: application/json" -vv \
    -u "$USER:$PASSWORD" \
    $HOST/rest/api/2/serverInfo 2> /tmp/jira-server-stderr

STATUS=$(grep -oP '^< HTTP.* \K([0-9]{3})' /tmp/jira-server-stderr | tail -n 1)
if [ "$STATUS" != 200 ]; then
    awk '{ if ($0 ~ /authorization/) $0 = "==AUTHORIZATION FILTERED=="; print }' /tmp/jira-server-stderr
    echo "Failed to fetch jira server information, please make sure USER/PASSWORD/HOST and BOARD are set correct."
    exit 1
fi

echo
echo "Settings looks good, start benchmarking"


# benchmark api response time and save to file
echo "ORDERBY,SINCE,EXPAND,PAGE_SIZE,PAGE,STATUS_CODE,ELLAPSE" > $FILENAME
for ORDERBY in  " ORDER BY created ASC" ""; do
    for SINCE in "-6 month" "-12 month" ""; do
        for EXPAND in "" "changelog"; do
            for PAGE in $(seq 1 $PAGES); do
                JQL=
                if [ -n "$SINCE" ]; then
                    JQL="updated > '$(date --date="$SINCE" +"%Y/%m/%d %H:%M")'"
                fi
                JQL="$JQL$ORDERBY"
                STARTAT=$((PAGE_SIZE * (PAGE - 1)))
                URL="$HOST/rest/agile/1.0/board/$BOARD/issue"
                printf "%s" "$URL" 
                STARTED_TS=$(date +%s)
                STATUS=$(
                curl -H "Accept: application/json" \
                    -u "$USER:$PASSWORD" \
                    -s -o /dev/null -w "%{http_code}" \
                    -G \
                    --data-urlencode "jql=$JQL" \
                    --data-urlencode "expand=$EXPAND" \
                    --data-urlencode "maxResults=100" \
                    --data-urlencode "startAt=$STARTAT" \
                    $URL
                )
                ENDED_TS=$(date +%s)
                ELLAPSE=$(( ENDED_TS - STARTED_TS ))
                printf " %s %s\n" "$STATUS" "$ELLAPSE"
                echo "$ORDERBY,$SINCE,$EXPAND,$PAGE_SIZE,$PAGE,$STATUS,$ELLAPSE" >> $FILENAME
            done
        done
    done
done
