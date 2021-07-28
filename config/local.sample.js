module.exports = {
  mongo: {
    connectionString: 'mongodb://localhost:27017/test?retryWrites=true&w=majority'
  },
  rabbitMQ: {
    connectionString: 'amqp://guest:guest@localhost:5672/rabbitmq'
  },
  postgres: {
    connectionString: 'postgres://postgres:postgres@localhost:5432/lake',
    username: 'postgres',
    password: 'postgres',
    host: 'localhost',
    database: 'lake',
    port: 5432,
    dialect: 'postgres'
  },
  jira: {
    host: 'https://your-domain.atlassian.net',
    basicAuth: '***',
    proxy: 'http://localhost:4780',
    timeout: 15000
  },
  gitlab: {
    host: 'https://gitlab.com',
    apiPath: 'api/v4',
    token: '***'
  },
  enrichment: {
    connectionString: 'http://localhost:3000/'
  }
}
