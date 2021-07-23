# Dev Lake

![badge](https://github.com/merico-dev/lake/actions/workflows/main.yml/badge.svg)

## Requirements

- Node.js
- Docker

## Installation

1. Clone this repository
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
2. We can use default values for most fields except the Jira section. For how to set up basic authorization with Jira, please see this [section](#jira) below

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
    }
}

```

3. See that the collection job was published, jira collection ran, the enrichment job was published, and jira enrichment ran

To run only the enrichment job on existing collections: `POST http://localhost:3000/`

## Connection Information

### Postgres Connection

- DB Name: lake
- Hostname: localhost
- Port: 5432
- Username: postgres
- Password: postgres

### MongoDB Connection

- DB Name: test
- Hostname: localhost
- Port: 27017
- Username: (none required)
- Password: (none required)

### RabbitMQ Connection

- Vhost Name: rabbitmq
- Hostname: localhost
- Port: 5672
- Username: guest
- Password: guest

### Grafana Connection

- Hostname: localhost
- Port: 3002
- Username: admin
- Password: admin

## Services

### Jira

__Jira auth setup__

1. Create an API key on Jira
3. Create a __basic auth header__ from your API key - [Jira Docs](https://developer.atlassian.com/cloud/jira/platform/basic-auth-for-rest-apis/#supply-basic-auth-headers)
3. Copy your __basic auth header__ into the `jira.basicAuth` field in `/config/local.js` file
4. Add your jira hostname to the `jira.host` field in the `/config/local.js` file
