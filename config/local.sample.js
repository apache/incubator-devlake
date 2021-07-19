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
  }
}
