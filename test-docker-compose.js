const config = require('../config/local')

testOutMongoDB = (mongoConnectionString) => {
  // connect to mongodb
  // save a record
  // get a record
  // delete a record
}

testOutRabbitMQ = (mongoConnectionString) => {
  // connect
  // save a message
  // receive a message
  // disconnect
}

testOutPostgres = (mongoConnectionString) => {
  // connect
  // save a record
  // get a record
  // delete a record
}

testOutMongoDB(config.mongo.connectionString)

testOutRabbitMQ(config.mongo.connectionString)

testOutPostgres(config.mongo.connectionString)
