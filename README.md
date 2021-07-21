# Dev Lake

![badge](https://github.com/merico-dev/lake/actions/workflows/main.yml/badge.svg)

## Requirements

- Docker

## Installation

1. Clone this repository
2. From the newly cloned repo directory, run `docker-compose up --build`
3. Test that it's working by opening another terminal and running the test script, `node test/test-docker-compose.js`. You should see `Connected to MongoDB`, `Connected to postgres`,  `Connected to RabbitMQ`.

## Usage

### Create a Collection or Enrichment Job

1. From the terminal, execute `npm run all`
2. From Postman (or similar), send a request like...

- Collection: `POST http://localhost:3001/`
- Enrichment: `POST http://localhost:3000/`

```json
{
    "projectId": 555555,
    "jira": {
        "projectId": "test-api",
        "accountUri": "merico.atlassian.net"
    }
}
```

3. See that the collection job was published, jira collection ran, the enrichment job was published, and jira enrichment ran

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

## Services

### Jira

__Jira auth setup__

1. Create an API key on Jira
3. Create a __basic auth header__ from your API key - [Jira Docs](https://developer.atlassian.com/cloud/jira/platform/basic-auth-for-rest-apis/#supply-basic-auth-headers)
3. Copy your __basic auth header__ into the `jira.basicAuth` field in `/config/local.js` file
4. Add your jira hostname to the `jira.host` field in the `/config/local.js` file
