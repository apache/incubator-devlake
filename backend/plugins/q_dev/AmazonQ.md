# Apache DevLake Q Developer Plugin - Project Analysis

## Project Overview

This is a plugin for **Apache DevLake**, an open-source dev data platform that ingests, analyzes, and visualizes fragmented data from DevOps tools. The Q Developer plugin specifically focuses on retrieving and analyzing AWS Q Developer usage data from AWS S3.

### Apache DevLake Context
- **Purpose**: Extract insights for engineering excellence, developer experience, and community growth
- **Target Users**: Engineering Leads, Open Source Software Maintainers, development teams
- **Core Functionality**: Brings together dev data from multiple silos to provide complete SDLC view
- **Metrics Support**: Implements DORA metrics and other frameworks with prebuilt dashboards
- **Visualization**: Integrated dashboards powered by Grafana

## Tech Stack

### Backend Technologies
- **Language**: Go 1.20+
- **Framework**: Apache DevLake plugin architecture
- **Web Framework**: Gin (gin-gonic/gin v1.9.1)
- **Database**: PostgreSQL (with pgx/v5 driver) and MySQL support
- **AWS Integration**: AWS SDK for Go (aws-sdk-go v1.55.6)
- **Data Processing**: GoCSV for CSV parsing
- **Configuration**: Mapstructure for configuration mapping
- **Validation**: Go Playground Validator v10.19.0
- **UUID Generation**: Google UUID v1.3.0

### Infrastructure & Deployment
- **Containerization**: Docker v19.03.10+
- **Build System**: GNU Make
- **Development Environment**: Docker Compose
- **Frontend**: Config UI (separate component)
- **Visualization**: Grafana dashboards

### AWS Services Used
- **AWS S3**: Primary data source for Q Developer usage metrics
- **AWS SDK**: For S3 operations and authentication
- **IAM**: Access key and secret key authentication

## Project Structure

```
q_dev/
├── README.md                    # Plugin documentation
├── Q_DEV_deploy_guide.md       # Development environment setup
├── q_dev.go                    # Main plugin entry point (standalone mode)
├── img.png                     # Configuration UI screenshot
├── api/                        # REST API endpoints
│   └── connection.go           # Connection CRUD operations
├── impl/                       # Plugin implementation
│   └── impl.go                 # Core plugin logic and interfaces
├── models/                     # Data models
│   ├── connection.go           # AWS S3 connection model
│   ├── s3_file_meta.go        # S3 file metadata model
│   ├── user_data.go           # User usage data model
│   ├── user_metrics.go        # Aggregated user metrics model
│   └── migrationscripts/      # Database migrations
├── tasks/                      # Data processing tasks
│   ├── task_data.go           # Task data structures
│   ├── s3_client.go           # S3 client wrapper
│   ├── s3_file_collector.go   # S3 file metadata collection
│   ├── s3_data_extractor.go   # CSV data extraction and parsing
│   └── user_metrics_converter.go # User metrics aggregation
└── plans/                      # Execution plans
```

## Core Functionality

### Data Pipeline (3-Stage Process)
1. **Collection**: `collectQDevS3Files` - Collects S3 file metadata without downloading content
2. **Extraction**: `extractQDevS3Data` - Downloads CSV files and parses user data into database
3. **Conversion**: `convertQDevUserMetrics` - Aggregates user data into metrics (averages, totals)

### Database Schema
- `_tool_q_dev_connections`: AWS S3 connection configurations
- `_tool_q_dev_s3_file_meta`: S3 file metadata storage
- `_tool_q_dev_user_data`: Raw user data from CSV files
- `_tool_q_dev_user_metrics`: Aggregated user metrics and analytics

### API Endpoints
- `POST /plugins/q_dev/connections` - Create new connection
- `GET /plugins/q_dev/connections` - List all connections
- `GET /plugins/q_dev/connections/:id` - Get specific connection
- `PATCH /plugins/q_dev/connections/:id` - Update connection
- `DELETE /plugins/q_dev/connections/:id` - Delete connection
- `POST /plugins/q_dev/test` - Test connection
- `POST /plugins/q_dev/connections/:id/test` - Test existing connection

## Configuration

### Connection Parameters
- **AWS Access Key ID**: Authentication credential
- **AWS Secret Access Key**: Authentication credential (sanitized in responses)
- **AWS Region**: S3 bucket region
- **S3 Bucket Name**: Target bucket containing Q Developer data
- **Rate Limit**: API requests per hour (default: 20000)

### Blueprint Configuration
```json
[
  [
    {
      "plugin": "q_dev",
      "subtasks": null,
      "options": {
        "connectionId": 5,
        "s3Prefix": ""
      }
    }
  ]
]
```

## Development Environment

### Requirements
- Docker v19.03.10+
- Golang v1.19+
- GNU Make
- MySQL/PostgreSQL database
- Grafana for visualization

### Setup Commands
```bash
# Clone repository
git clone https://github.com/apache/incubator-devlake.git
cd incubator-devlake

# Install dependencies
cd backend && go get && cd ..

# Configure environment
cp env.example .env
# Update DB_URL and DISABLED_REMOTE_PLUGINS

# Start services
docker-compose -f docker-compose-dev.yml up -d mysql grafana

# Run in development mode
DEVLAKE_PLUGINS=q_dev make dev
make configure-dev
```

### Standalone Mode
The plugin supports standalone debugging mode:
```bash
go run q_dev.go --connectionId=1 --s3Prefix="path/to/data"
```

## Plugin Architecture

### Interface Implementation
The plugin implements multiple DevLake interfaces:
- `PluginMeta`: Basic plugin metadata
- `PluginInit`: Initialization logic
- `PluginTask`: Task execution
- `PluginApi`: REST API endpoints
- `PluginModel`: Database models
- `PluginSource`: Data source operations
- `PluginMigration`: Database migrations
- `CloseablePluginTask`: Resource cleanup

### Task Data Flow
1. **Connection Setup**: Establish AWS S3 client with credentials
2. **File Discovery**: Scan S3 bucket for CSV files matching prefix
3. **Metadata Storage**: Store file information without downloading
4. **Data Extraction**: Download and parse CSV files into structured data
5. **Metrics Calculation**: Aggregate user data into meaningful metrics
6. **Storage**: Persist all data in DevLake database for visualization

## Security Considerations
- AWS credentials are sanitized in API responses
- Secret access keys are masked with placeholder values
- Rate limiting prevents API abuse
- Connection testing validates credentials before storage

## Integration Points
- **Grafana**: Visualization of collected metrics
- **DevLake Core**: Plugin framework and database operations
- **AWS S3**: Primary data source
- **Config UI**: Web interface for plugin configuration (localhost:4000)

This plugin serves as a bridge between AWS Q Developer usage analytics and the broader DevLake ecosystem, enabling teams to incorporate AI coding assistant metrics into their overall development analytics and DORA metrics tracking.
