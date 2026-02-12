---
name: start-devlake
description: Start the full DevLake stack locally (MySQL, backend, Grafana, Config UI)
tools:
  - run_in_terminal
---

# Start DevLake Local Development Stack

Run the full Apache DevLake stack locally for development. This prompt starts Docker containers and the Config UI dev server.

## Prerequisites

- Docker Desktop must be running
- Node.js/npm installed for Config UI
- Working directory should be the incubator-devlake repository root

## Workflow

### Step 1: Check container status

Run `docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"` to see which containers are already running.

### Step 2: Start backend containers

If containers are not running, start them with:

```powershell
docker compose -f docker-compose-dev.yml up -d
```

This starts:
- `mysql` on port 3306
- `devlake` backend on port 8080
- `grafana` on port 3002

### Step 3: Verify health endpoints

Check that services are healthy:

**DevLake health:**
```powershell
$ProgressPreference='SilentlyContinue'; Invoke-RestMethod -TimeoutSec 8 'http://localhost:8080/health' | ConvertTo-Json -Compress
```

Expected: `{"code":0,"success":true,"message":"good",...}`

**Grafana health:**
```powershell
$ProgressPreference='SilentlyContinue'; Invoke-RestMethod -TimeoutSec 8 'http://localhost:3002/api/health' | ConvertTo-Json -Compress
```

Expected: `{"database":"ok",...}`

### Step 4: Start Config UI

Check if Config UI is running at http://localhost:4000/. If not accessible, start it:

```powershell
cd config-ui
if (!(Test-Path node_modules)) { npm install }
npm start
```

This starts the Vite dev server in background on port 4000.

### Step 5: Verify Config UI

```powershell
$ProgressPreference='SilentlyContinue'; try { (Invoke-WebRequest -TimeoutSec 5 'http://localhost:4000/' -UseBasicParsing).StatusCode } catch { $_.Exception.Message }
```

Expected: `200`

## Success Output

After all services are running, report the access URLs:

| Service | URL | Status |
|---------|-----|--------|
| **Config UI** | http://localhost:4000/ | ✅ Running |
| **DevLake API** | http://localhost:8080/ | ✅ Healthy |
| **Grafana** | http://localhost:3002/ | ✅ Healthy |

## Troubleshooting

| Symptom | Likely cause | Fix |
|---------|--------------|-----|
| UI not reachable at :4000 | Vite dev server not running | Start config UI in config-ui directory |
| /health fails on :8080 | Backend container stopped | Re-run `docker compose -f docker-compose-dev.yml up -d` |
| Migration banner in UI | New DB migrations needed | Proceed with migration in the UI before configuring connections |
| Connection list empty | Plugin id mismatch | Use `/plugins` endpoint to confirm plugin id |

## Notes

- If you see "New Migration Scripts Detected" banner in the UI, proceed with the migration before configuring connections.
- The Config UI Vite server should be started in a background terminal so it keeps running.
