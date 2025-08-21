/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/


# DevLake Development Environment Deployment Guide

## Environment Requirements
- Docker v19.03.10+
- Golang v1.19+
- GNU Make
    - Mac (pre-installed)
    - Windows: [Download](http://gnuwin32.sourceforge.net/packages/make.htm)
    - Ubuntu: `sudo apt-get install build-essential libssl-dev`

## How to Set Up the Development Environment
The following guide will explain how to run DevLake's frontend (config-ui) and backend in development mode.

### Clone the Repository
Navigate to where you want to install this project and clone the repository:

```bash
git clone https://github.com/apache/incubator-devlake.git
cd incubator-devlake
```

### Install Plugin Dependencies

RefDiff plugin:
Install Go packages
```bash
cd backend
go get
cd ..
```

### Configure Environment File
Copy the example configuration file to a new local file:

```bash
cp env.example .env
```

Update the following variables in the `.env` file:

- `DB_URL`: Replace `mysql:3306` with `127.0.0.1:3306`
- `DISABLED_REMOTE_PLUGINS`: Set to `True`

### Q Developer Plugin Configuration
The Q Developer plugin requires AWS credentials with access to both S3 and IAM Identity Center:

**Required AWS Permissions:**
- S3: `s3:GetObject`, `s3:ListBucket` for the Q Developer data bucket
- Identity Center: `identitystore:DescribeUser` for user display name resolution

**Required Configuration Fields:**
- AWS Access Key ID and Secret Access Key
- S3 bucket name and region
- IAM Identity Center Store ID (format: `d-xxxxxxxxxx`)
- IAM Identity Center region

### Start MySQL and Grafana Containers

Make sure the Docker daemon is running before this step.

> Grafana needs to rebuild the image, then change the image in docker-compose.datasources.yml to `image: grafana:latest`

```bash
docker-compose -f docker-compose-dev.yml up -d mysql grafana
```

### Run in Development Mode
Run devlake and config-ui in development mode in two separate terminals:

```bash
# Install poetry, follow the guide: https://python-poetry.org/docs/#installation
# Run devlake, only using the q dev plugin here
DEVLAKE_PLUGINS=q_dev nohup make dev &
# Run config-ui
make configure-dev
```

For common errors, please refer to the troubleshooting documentation.

Config UI runs on localhost:4000