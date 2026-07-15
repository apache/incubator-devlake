#!/bin/bash
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

set -e

DATASOURCE_FILE="/etc/grafana/provisioning/datasources/datasource.yml"

# Detect database type
# Priority: DATABASE_TYPE, then legacy MYSQL_*/POSTGRES_* auto-detection
if [ -n "$DATABASE_TYPE" ]; then
  MODE="$DATABASE_TYPE"

  case "$MODE" in
    mysql)
      export MYSQL_URL="${DATABASE_HOST:-mysql}:${DATABASE_PORT:-3306}"
      export MYSQL_DATABASE="${DATABASE_NAME:-lake}"
      export MYSQL_USER="${DATABASE_USER:-merico}"
      export MYSQL_PASSWORD="${DATABASE_PASSWORD:-merico}"
      ;;
    postgresql)
      export POSTGRES_URL="${DATABASE_HOST:-postgres}:${DATABASE_PORT:-5432}"
      export POSTGRES_DATABASE="${DATABASE_NAME:-lake}"
      export POSTGRES_USER="${DATABASE_USER:-merico}"
      export POSTGRES_PASSWORD="${DATABASE_PASSWORD:-merico}"
      ;;
    *)
      echo "ERROR: DATABASE_TYPE must be 'mysql' or 'postgresql'"
      exit 1
      ;;
  esac
else
  # Legacy: auto-detect from MYSQL_*/POSTGRES_* vars
  if [ -n "$POSTGRES_URL" ]; then
    MODE="postgresql"
  elif [ -n "$MYSQL_URL" ]; then
    MODE="mysql"
  else
    echo "WARNING: No database vars. Defaulting to mysql."
    MODE="mysql"
    export MYSQL_URL="mysql:3306"
    export MYSQL_DATABASE="lake"
    export MYSQL_USER="merico"
    export MYSQL_PASSWORD="merico"
  fi
fi

echo "Database type: $MODE"

# Remove unused dashboard folder to prevent confusion
if [ "$MODE" = "mysql" ]; then
  rm -rf /etc/grafana/dashboards/postgresql
else
  rm -rf /etc/grafana/dashboards/mysql
  SSL_MODE="${DATABASE_SSL_MODE:-disable}"
  echo "SSL Mode: ${SSL_MODE}"
fi

# Locate Homepage.json for the home dashboard. The layout differs by deployment:
#  - baked-in image keeps variant subfolders: /etc/grafana/dashboards/<mode>/Homepage.json
#  - dev compose bind-mounts the variant folder directly: /etc/grafana/dashboards/Homepage.json
# Probe both so the home dashboard resolves in either case (a wrong path makes
# Grafana 12+ return HTTP 500 "Failed to load home dashboard").
for _hp in \
  "/etc/grafana/dashboards/${MODE}/Homepage.json" \
  "/etc/grafana/dashboards/Homepage.json"; do
  if [ -f "$_hp" ]; then
    export GF_DASHBOARDS_DEFAULT_HOME_DASHBOARD_PATH="$_hp"
    break
  fi
done

echo "Homepage: $GF_DASHBOARDS_DEFAULT_HOME_DASHBOARD_PATH"

# --- Legacy persistent-volume hygiene -------------------------------------
# This image runs as an arbitrary uid with gid 0 (OpenShift-compatible). Volumes
# created by older dashboard images (uid 472, file mode 0640) can be read-only for
# the current runtime user, which makes Grafana's DB migrations fail with
# "attempt to write a readonly database". Best-effort self-heal here; if it can't
# be fixed (non-root, non-owner) print a clear, actionable error instead of a
# cryptic SQLite crash.
GRAFANA_DATA_DIR="${GF_PATHS_DATA:-/var/lib/grafana}"
GRAFANA_DB="${GRAFANA_DATA_DIR}/grafana.db"

if [ -e "$GRAFANA_DB" ] && [ ! -w "$GRAFANA_DB" ]; then
  echo "WARNING: ${GRAFANA_DB} not writable by $(id -u):$(id -g); attempting permission fix..."
  chgrp -R 0 "$GRAFANA_DATA_DIR" 2>/dev/null || true
  chmod -R g+rwX "$GRAFANA_DATA_DIR" 2>/dev/null || true
fi

if [ -e "$GRAFANA_DB" ] && [ ! -w "$GRAFANA_DB" ]; then
  echo "ERROR: ${GRAFANA_DB} is still read-only for gid 0."
  echo "       The persistent volume was likely created by an older image (uid 472, mode 0640)."
  echo "       Fix it once from the host, then recreate this container:"
  echo "         docker run --rm -v <project>_grafana-storage:/data alpine \\"
  echo "           sh -c 'chgrp -R 0 /data && chmod -R g+rwX /data'"
fi

# Remove the deprecated Angular grafana-piechart-panel plugin if it lingers in a
# legacy volume. Grafana 12+ dropped Angular support and dashboards now use the
# core "piechart" panel; leaving it causes noisy plugin-validation errors.
LEGACY_PIECHART="${GRAFANA_DATA_DIR}/plugins/grafana-piechart-panel"
if [ -d "$LEGACY_PIECHART" ]; then
  echo "Removing deprecated grafana-piechart-panel from data volume..."
  rm -rf "$LEGACY_PIECHART" 2>/dev/null || true
fi
# --------------------------------------------------------------------------

# Create empty datasource.yml (datasources created via API)
cat > "$DATASOURCE_FILE" << HEADER
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
HEADER

# Admin credentials used by the datasource provisioning API calls below.
# Defaults to Grafana's built-in admin password; override via GF_SECURITY_ADMIN_PASSWORD.
# Exported so Grafana itself and the curl calls use the same value.
export GF_SECURITY_ADMIN_PASSWORD="${GF_SECURITY_ADMIN_PASSWORD:-admin}"

# Start Grafana in background
/run.sh "$@" &
GRAFANA_PID=$!

# Wait for Grafana API
echo "Waiting for Grafana API..."
for i in $(seq 1 60); do
  if curl -f -s http://localhost:3000/api/health >/dev/null 2>&1; then
    echo "Grafana API ready"
    sleep 5  # Extra wait for migrations
    break
  fi
  sleep 2
done

# Grafana 13+ auto-creates empty datasource instances for built-in plugins and
# legacy volumes may contain read-only provisioned datasources with different UIDs.
# Delete ALL existing datasources (by ID, which bypasses read-only flags) so we can
# create a single correctly-configured datasource with the UID our dashboards expect.
echo "Deleting all existing datasources (clean slate)..."
for _id in $(curl -s "http://admin:${GF_SECURITY_ADMIN_PASSWORD}@localhost:3000/api/datasources" 2>/dev/null \
  | grep -o '"id":[0-9]*' | sed 's/"id"://g'); do
  curl -s -X DELETE "http://admin:${GF_SECURITY_ADMIN_PASSWORD}@localhost:3000/api/datasources/${_id}" 2>&1 || true
done
sleep 2

# Create datasource via API (both MySQL and PostgreSQL)
PAYLOAD_FILE="/tmp/datasource-api.json"

if [ "$MODE" = "mysql" ]; then
  cat > "$PAYLOAD_FILE" <<APIJSON
{
  "uid": "devlake-mysql-api",
  "name": "mysql",
  "type": "mysql",
  "url": "${MYSQL_URL}",
  "database": "${MYSQL_DATABASE}",
  "user": "${MYSQL_USER}",
  "secureJsonData": {
    "password": "${MYSQL_PASSWORD}"
  },
  "access": "proxy",
  "isDefault": true,
  "editable": true
}
APIJSON


else
  SSL_MODE="${DATABASE_SSL_MODE:-disable}"
  cat > "$PAYLOAD_FILE" <<APIJSON
{
  "uid": "devlake-postgres-api",
  "name": "postgresql",
  "type": "grafana-postgresql-datasource",
  "url": "${POSTGRES_URL}",
  "database": "${POSTGRES_DATABASE}",
  "user": "${POSTGRES_USER}",
  "secureJsonData": {
    "password": "${POSTGRES_PASSWORD}"
  },
  "jsonData": {
    "sslmode": "${SSL_MODE}",
    "postgresVersion": 1400,
    "database": "${POSTGRES_DATABASE}"
  },
  "access": "proxy",
  "isDefault": true,
  "editable": true
}
APIJSON


fi

echo "Creating datasource via API..."
for i in $(seq 1 10); do
  RESPONSE=$(curl -s -X POST "http://admin:${GF_SECURITY_ADMIN_PASSWORD}@localhost:3000/api/datasources" \
    -H "Content-Type: application/json" \
    -d @"$PAYLOAD_FILE" 2>&1)

  if echo "$RESPONSE" | grep -q '"id"'; then
    echo "Datasource created successfully"
    break
  elif echo "$RESPONSE" | grep -q "already exists"; then
    # A datasource with this name is already present (e.g. provisioned by an
    # older image into a persistent volume). It cannot be recreated via the API
    # but is functional, so treat this as success instead of retrying.
    echo "Datasource already exists; keeping existing one"
    break
  elif echo "$RESPONSE" | grep -q "database is locked"; then
    echo "DB locked, retry $i/10..."
    sleep 3
  else
    echo "API error: $RESPONSE"
    sleep 2
  fi
done

# Wait for Grafana
wait $GRAFANA_PID
