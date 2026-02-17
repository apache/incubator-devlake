# Developer Telemetry Plugin for Apache DevLake

## Overview

The Developer Telemetry plugin enables DevLake to ingest, store, and analyze developer productivity metrics collected from local development environments. This plugin operates on a **push model** via HTTP webhooks, where telemetry collectors running on developer machines send aggregated data to DevLake.

## Features

- **Webhook-based ingestion**: Accepts telemetry data via REST API
- **Developer metrics tracking**: Active hours, tool usage, context switching
- **Idempotent updates**: Same-day data can be resent and will be updated
- **Secure authentication**: Optional token-based authentication
- **Connection management**: Support for multiple telemetry sources

## Architecture

```
┌──────────────────┐          HTTP POST           ┌─────────────────┐
│  Telemetry       │  ────────────────────────>   │   DevLake       │
│  Collector       │   JSON Payload               │   Plugin API    │
│  (Local Machine) │                              │                 │
└──────────────────┘                              └─────────────────┘
                                                            │
                                                            ▼
                                                   ┌─────────────────┐
                                                   │   PostgreSQL/   │
                                                   │   MySQL DB      │
                                                   └─────────────────┘
```

## Data Models

### Connections Table: `_tool_developer_telemetry_connections`

Stores connection configurations for telemetry sources.

| Column | Type | Description |
|--------|------|-------------|
| id | BIGINT | Primary key |
| name | VARCHAR(100) | Connection name |
| endpoint | VARCHAR | Not used for push webhooks |
| secret_token | VARCHAR | Optional authentication token |

### Metrics Table: `_tool_developer_metrics`

Stores daily developer telemetry data.

| Column | Type | Description |
|--------|------|-------------|
| connection_id | BIGINT | Foreign key to connections |
| developer_id | VARCHAR(255) | System username |
| date | DATE | Metrics date (YYYY-MM-DD) |
| email | VARCHAR(255) | Developer email |
| name | VARCHAR(255) | Developer name |
| hostname | VARCHAR(255) | Machine hostname |
| active_hours | INT | Active coding hours |
| tools_used | TEXT | JSON array of tools |
| project_context | TEXT | JSON array of projects |
| command_counts | TEXT | JSON map of command usage |
| os_info | VARCHAR(255) | OS information |

Primary key: `(connection_id, developer_id, date)`

## API Endpoints

### Connection Management

#### Create Connection
```
POST /plugins/developer_telemetry/connections
Content-Type: application/json

{
  "name": "Team Dev Fleet",
  "secretToken": "optional-auth-token"
}
```

#### List Connections
```
GET /plugins/developer_telemetry/connections
```

#### Get Connection
```
GET /plugins/developer_telemetry/connections/:connectionId
```

#### Update Connection
```
PATCH /plugins/developer_telemetry/connections/:connectionId
Content-Type: application/json

{
  "name": "Updated Name",
  "secretToken": "new-token"
}
```

#### Delete Connection
```
DELETE /plugins/developer_telemetry/connections/:connectionId
```

### Telemetry Data Ingestion

#### Post Telemetry Report
```
POST /plugins/developer_telemetry/connections/:connectionId/report
Authorization: Bearer <secret_token>  # If secretToken is configured
Content-Type: application/json

{
  "date": "2026-02-11",
  "developer": "irfan.ahmad",
  "email": "irfan@company.com",
  "name": "Irfan Ahmad",
  "hostname": "irfan-macbook-pro",
  "metrics": {
    "active_hours": 9,
    "tools_used": ["go", "vscode", "docker"],
    "commands": {
      "git": 45,
      "docker": 12,
      "make": 8
    },
    "projects": ["incubator-devlake", "developer-telemetry"]
  }
}
```

**Response:**
```json
{
  "success": true,
  "message": "telemetry data received successfully"
}
```

## Building the Plugin

From the `incubator-devlake/backend` directory:

```bash
# Build only the developer_telemetry plugin
DEVLAKE_PLUGINS=developer_telemetry make build-plugin

# Build all plugins including developer_telemetry
make build-plugin
```

The compiled plugin will be at:
```
bin/plugins/developer_telemetry/developer_telemetry.so
```

## Running DevLake with the Plugin

1. Build the plugin (see above)
2. Build and start DevLake:
   ```bash
   make build-server
   make dev
   ```
3. The plugin will be automatically loaded from the `bin/plugins` directory

## Integration with Telemetry Collector

Configure your telemetry collector to send data to the webhook URL:

```bash
# In your collector configuration
DEVLAKE_WEBHOOK_URL="http://localhost:8080/plugins/developer_telemetry/connections/1/report"
DEVLAKE_AUTH_TOKEN="your-secret-token"  # If using authentication
```

Example collector script snippet:
```bash
curl -X POST \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $DEVLAKE_AUTH_TOKEN" \
  -d @telemetry-data.json \
  "$DEVLAKE_WEBHOOK_URL"
```

## Security Considerations

1. **Authentication**: Always configure `secretToken` in production environments
2. **HTTPS**: Use HTTPS in production to encrypt data in transit
3. **Token Rotation**: Regularly rotate secret tokens
4. **Network Security**: Restrict access to the webhook endpoint using firewalls/VPNs

## Idempotency

The plugin supports idempotent updates. If you send telemetry data for the same `developer` + `date` combination multiple times, the latest data will overwrite the previous record. This is useful for:
- Retry mechanisms in collectors
- Incremental updates throughout the day
- Correcting previously sent data

## Testing

### Manual Testing

1. Create a connection:
   ```bash
   curl -X POST http://localhost:8080/plugins/developer_telemetry/connections \
     -H "Content-Type: application/json" \
     -d '{"name": "Test Connection", "secretToken": "test123"}'
   ```

2. Send test telemetry data:
   ```bash
   curl -X POST http://localhost:8080/plugins/developer_telemetry/connections/1/report \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer test123" \
     -d '{
       "date": "2026-02-11",
       "developer": "testuser",
       "email": "test@example.com",
       "name": "Test User",
       "hostname": "test-machine",
       "metrics": {
         "active_hours": 8,
         "tools_used": ["vscode", "git"],
         "commands": {"git": 20, "docker": 5},
         "projects": ["test-project"]
       }
     }'
   ```

3. Verify data in database:
   ```sql
   SELECT * FROM _tool_developer_metrics WHERE developer_id = 'testuser';
   ```

## Troubleshooting

### Plugin not loading
- Check that `developer_telemetry.so` exists in `bin/plugins/developer_telemetry/`
- Check DevLake logs for plugin loading errors
- Verify PLUGIN_DIR environment variable is set correctly

### Authentication errors
- Verify `secretToken` matches in connection and request
- Check `Authorization` header format: `Bearer <token>`
- Ensure connection exists with the specified ID

### Data not appearing
- Check DevLake logs for ingestion errors
- Verify date format is `YYYY-MM-DD`
- Ensure required fields (`date`, `developer`, `metrics`) are present
- Check database connectivity

## License

Licensed under the Apache License, Version 2.0. See LICENSE file for details.
