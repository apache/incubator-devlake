module.exports = {
  jira: {
    host: 'https://merico.atlassian.net',
    basicAuth: 'anVzdGluLmJyYXplYXVAbWVyaWNvLmRldjpMTVlPVnNtZXhXekZlNW5sakdMY0VERDU=',
    proxy: 'http://localhost:4780',
    timeout: 15000,
    maxPagesForTest: 2,
    skipIssueCollection: false,
    skipChangelogCollection: true
  },
  gitlab: {
    host: 'https://gitlab.com',
    apiPath: 'api/v4',
    token: '6RxxyVNVdcjrMewi9sm_',
    maxPagesForTest: 2
  },
  enrichment: {
    connectionString: 'http://localhost:3000/'
  },
  mongo: {
    connectionString: 'mongodb://localhost:27017/test?retryWrites=true&w=majority'
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
}
