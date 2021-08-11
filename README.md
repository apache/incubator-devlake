# Dev Lake

![badge](https://github.com/merico-dev/lake/actions/workflows/test.yml/badge.svg?branch=ts-main)
[![codecov](https://codecov.io/gh/merico-dev/lake/branch/ts-main/graph/badge.svg?token=UN126GAU9D)](https://codecov.io/gh/merico-dev/lake)

## Requirements

- Node.js
- Docker

## Setup

### Local

1. Clone this repository
2. Install dependencies with `npm i`
3. Install `postgres` and `redis` and startup
4. Create config file `cp .env.sample .env`. change `DB_URL` and `REDIS_URL` to your local db
5. Start services `npm run all`

### Docker

1. Clone this repository
2. From the newly cloned repo directory, run `docker-compose -f ./devops/docker-compose.yml --project-directory ./ --build up `
3. Run `docker-compose ps` to see containers runnning.
4. Install dependencies with `npm i`
5. Create config file `cp .env.sample .env`. change `DB_URL` and `REDIS_URL` to your docker container db
6. Start services `npm run all`

## Configuration

1. Make a copy of `.env.sample` to `.env`
2. SET `REDIS_URL` AND `DB_URL`. Queue service require redis for the task queue. Plugins write datas into configed DB.

### Grafana Connection For Data Visualization (https://localhost:3002)

Connect to the Grafana database:
Inside `./devops/docker-compose.yml` edit the environment variables as needed to connect to your local postgres instance, specifically:
- `GF_DATABASE_NAME`
- `GF_DATABASE_USER`
- `GF_DATABASE_PASSWORD`

Connect the Grafana data source:
Additionally to use the postgres database as data source inside grafana, ensure postgres config options are correct in `./grafana/datasources/datasource.yml`, specifically:
- `database`
- `user`
- `secureJsonData/password`

## Usage

### Add Data Source

```
POST  /sources
{
    type: 'Jira',
    options: {
        host: '',
        username: '',
        password: 'password or api token'
    }
}

response:
{ source: 'source id' }
```

### Add Data Source Task

```
POST  /sources/${source id}
{
    collector: ['Issue'],
    enricher: ['LeadTime'],
    options: {
        projects: ['ProjectName'],
        boards: ['Scrum Board Id']
    }
}

response:
{ task: 'task id' }
```

### Waiting Task Finished

```
GET  /tasks/${task id}

response:
{
    task: {
        id: 'task id',
        status: 'finished',
        progress: {
            collector: [
                {name: 'Issue', status: 'finished'}
            ],
            enricher: [
                {name: 'LeadTime', status: 'finished'}
            ]
        }
    } 
}
```

### Using Grafana

**Login Credentials**

- Visit: `http://localhost:3002`
- Username: `admin`
- Password: `admin`

**Provisioning a Grafana Dashboard**

To save a dashboard in the `lake` repo and load it:
1. Create a dashboard in browser (visit `/dashboard/new`, or use sidebar)
2. Save dashboard (in top right of screen)
3. Go to dashboard settings (in top right of screen)
4. Click on _JSON Model_ in sidebar
5. Copy code into a new `.json` file in `/grafana/dashboards`

## Contributing

[CONTRIBUTING.md](CONTRIBUTING.md)