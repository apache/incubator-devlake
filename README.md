# Dev Lake

## Requirements

- Docker

## Installation

1. Clone this repository
2. From the newly cloned repo directory, run `docker-compose up --build`
3. Test that it's working by opening another terminal and running the test script, `node test/test-docker-compose.js`. You should see `Connected to MongoDB`, `Connected to postgres`,  `Connected to RabbitMQ`.

## Usage

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