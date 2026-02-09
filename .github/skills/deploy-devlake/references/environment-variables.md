# Environment Variables Reference

## Required Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `DB_URL` | **Yes** | Database connection URL. See format below. |
| `ENCRYPTION_SECRET` | **Yes** | 32-character secret. **Backend panics without it.** |
| `PORT` | Yes | API port (default: 8080) |

## Recommended Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `MODE` | `debug` | Set to `release` for production |
| `PLUGIN_DIR` | `bin/plugins` | Path to Go plugins |
| `REMOTE_PLUGIN_DIR` | `python/plugins` | Path to Python plugins |
| `LOGGING_DIR` | `/app/logs` | Log output directory |
| `TZ` | `UTC` | Timezone |

## DB_URL Format

### MySQL (Azure Database for MySQL)

```
mysql://USER:PASSWORD@SERVER.mysql.database.azure.com:3306/DATABASE?charset=utf8mb4&parseTime=True&loc=UTC&tls=true
```

**Required query parameters:**
- `parseTime=True` - Without this, datetime fields fail with scan error
- `loc=UTC` - Timezone for datetime parsing
- `tls=true` - Required for Azure MySQL SSL connection

### PostgreSQL (Azure Database for PostgreSQL)

```
postgres://USER:PASSWORD@SERVER.postgres.database.azure.com:5432/DATABASE?sslmode=require
```

**Required query parameters:**
- `sslmode=require` - Required for Azure PostgreSQL SSL connection

## ENCRYPTION_SECRET

Must be exactly 32 characters. Generate with:

**Bash:**
```bash
openssl rand -base64 24 | tr -dc 'a-zA-Z0-9' | head -c 32
```

**PowerShell:**
```powershell
[guid]::NewGuid().ToString().Replace('-','').Substring(0,32)
```

**Important:** Backend will panic immediately on startup if this is missing or not 32 characters.

## Grafana Variables

| Variable | Description |
|----------|-------------|
| `GF_SERVER_ROOT_URL` | Public URL for Grafana (used for links in emails, etc.) |
| `GF_SECURITY_ADMIN_PASSWORD` | Admin password (optional, defaults to admin) |

## Config UI Variables

| Variable | Description |
|----------|-------------|
| `DEVLAKE_ENDPOINT` | URL to DevLake backend API |
| `GRAFANA_ENDPOINT` | URL to Grafana instance |

## Example: Complete Backend Environment

```bash
DB_URL=mysql://merico:MyPassword123@devlakemysql.mysql.database.azure.com:3306/lake?charset=utf8mb4&parseTime=True&loc=UTC&tls=true
ENCRYPTION_SECRET=a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6
PORT=8080
MODE=release
PLUGIN_DIR=bin/plugins
REMOTE_PLUGIN_DIR=python/plugins
LOGGING_DIR=/app/logs
TZ=UTC
```
