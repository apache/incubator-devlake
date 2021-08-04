module.exports = {
  lake: {
    token: 'mytoken'
  },
  mongo: {
    connectionString: 'mongodb://lake:lakeIScoming@localhost:27017/lake?authSource=admin'
  },
  rabbitMQ: {
    connectionString: 'amqp://guest:guestWhat@localhost:5672/rabbitmq'
  },
  postgres: {
    username: 'postgres',
    password: 'postgresWhat',
    host: 'localhost',
    database: 'lake',
    port: 5432,
    dialect: 'postgres'
  },
  jira: {
    host: 'https://your-domain.atlassian.net',
    basicAuth: '***',
    proxy: 'http://localhost:4780',
    timeout: 10000,
    maxRetry: 10
  },
  gitlab: {
    host: 'https://gitlab.com',
    apiPath: 'api/v4',
    token: '***',
    proxy: 'http://localhost:4780',
    timeout: 10000,
    maxRetry: 10,
    skip: {
      commits: false,
      projects: false,
      mergeRequests: false,
      notes: false
    }
  },
  jiraBoardGitlabProject: {
    8: 8967944
  }
}
