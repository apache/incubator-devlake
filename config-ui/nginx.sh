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
if [ -n "$ADMIN_USER" ] && [ -n "$ADMIN_PASS" ]; then
    htpasswd -c -b /etc/nginx/.htpasswd "$ADMIN_USER" "$ADMIN_PASS"
    export SERVER_CONF='
    auth_basic           "DevLake Config UI";
    auth_basic_user_file /etc/nginx/.htpasswd;
    '
fi
export DNS=$(grep nameserver /etc/resolv.conf | awk '{print $2}')
export DNS_VALID=${DNS_VALID:-300s}
export DEVLAKE_ENDPOINT_PROTO=${DEVLAKE_ENDPOINT_PROTO:-http}
export GRAFANA_ENDPOINT_PROTO=${GRAFANA_ENDPOINT_PROTO:-http}
envsubst '${DEVLAKE_ENDPOINT} ${DEVLAKE_ENDPOINT_PROTO} ${GRAFANA_ENDPOINT} ${GRAFANA_ENDPOINT_PROTO} ${USE_EXTERNAL_GRAFANA} ${SERVER_CONF} ${DNS} ${DNS_VALID}' \
    < /etc/nginx/conf.d/default.conf.tpl \
    > /etc/nginx/conf.d/default.conf
nginx -g 'daemon off;'
