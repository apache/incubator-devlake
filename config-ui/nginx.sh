#!/bin/sh
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
envsubst '${DEVLAKE_ENDPOINT} ${GRAFANA_ENDPOINT} ${SERVER_CONF} ${DNS} ${DNS_VALID}' \
    < /etc/nginx/conf.d/default.conf.tpl \
    > /etc/nginx/conf.d/default.conf
nginx -g 'daemon off;'
