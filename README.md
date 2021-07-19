# Dev Lake

## Requirements

- Docker

## Installation

1. Clone this repository
2. From the newly cloned repo directory, run `docker-compose up --build`
3. Test that it's working by opening another terminal and running the test script, `node test/test-docker-compose.js`. You should see `Connected to MongoDB`, `Connected to postgres`,  `Connected to RabbitMQ`.

## Usage

### Create a collection job

1. From the terminal, execute `npm run collect`
2. In another tab, execute `npm run collection-worker`
3. From Postman (or similar), send a request like...
```
POST http://localhost:3001/

[
    {
        "projectId": 126,
        "jira": {
            "apiKey": "abc123",
            "accountUri": "merico.atlassian.net"
        }
    }
]
```
4. See in the collection api terminal that the job was published
5. See in the collection-worker terminal that the job was received by the jira collector

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
