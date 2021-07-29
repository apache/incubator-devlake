# Dev Lake

![badge](https://github.com/merico-dev/lake/actions/workflows/main.yml/badge.svg)

## Requirements

- Node.js
- Docker

## Installation

1. Clone this repository
1-1. Create the rabbitMQ log directory: `mkdir -p /tmp/rabbitmq/logs/`
2. From the newly cloned repo directory, run `docker-compose up --build`

    > NOTE: If you get an error like this:

    >"Error response from daemon: invalid mount config for type "bind": bind source path does not exist: /tmp/rabbitmq/etc/"

    >You can fix it by adding the directory in the terminal:
```
    mkdir /tmp/rabbitmq/etc
```
3. Run `docker-compose ps` to see containers runnning.
4. Install dependencies with `npm i`
5. Run migration with `npx sequelize-cli db:migrate`

## Configuration

1. Make a copy of `config/local.sample.js` under the name of `config/local.js`
2. We can use default values for most fields except the Jira section.

For how to set up basic authorization with Jira, and many more things, please see [CONFIGURATION.md](CONFIGURATION.md)

### Grafana Connection For Data Visualization (https://localhost:3002)

Connect to the Grafana database:
Inside `docker-compose.yml` edit the environment variables as needed to connect to your local postgres instance, specifically:
- `GF_DATABASE_NAME`
- `GF_DATABASE_USER`
- `GF_DATABASE_PASSWORD`

Connect the Grafana data source:
Additionally to use the postgres database as data source inside grafana, ensure postgres config options are correct in `./grafana/datasources/datasource.yml`, specifically:
- `database`
- `user`
- `secureJsonData/password`

## Usage

### Create a Collection Job

1. From the terminal, execute `npm run all`
2. From Postman (or similar), send a request like...

```json

POST http://localhost:3001/

{
    "jira": {
        "projectId": "10003",
        "accountUri": "merico.atlassian.net"
    },
    "gitlab": {
        "projectIds": [19688130, 8967944, 20103385]
    }
}

```

3. See that the collection job was published, jira collection ran, the enrichment job was published, and jira enrichment ran

To run only the enrichment job on existing collections: `POST http://localhost:3000/`

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
