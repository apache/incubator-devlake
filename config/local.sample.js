module.exports = {
  // Configuration of lake's own services
  lake: {
    // Enable basic authentication to the lake API
    // token: 'mytoken'
    // Set how often lake will fetch new data from data sources (default every hour)
  },
  // Configuration of MongoDB
  mongo: {
    connectionString:
      "mongodb://lake:lakeIScoming@localhost:27017/lake?authSource=admin",
  },
  // Configuration of rabbitMQ
  rabbitMQ: {
    connectionString: "amqp://guest:guestWhat@localhost:5672/rabbitmq",
  },
  // Configuration of PostgreSQL
  postgres: {
    username: "postgres",
    password: "postgresWhat",
    host: "localhost",
    database: "lake",
    port: 5432,
    dialect: "postgres",
  },
  cron: {
    job: {
      jira: {
        boardId: "<your-board-id>",
      },
      gitlab: {
        projectId: "<your-gitlab-project-id>",
      },
    },
    loopIntervalInMinutes: 60,
  },
};
