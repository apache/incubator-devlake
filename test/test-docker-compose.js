const config = require('../config/local')

testOutMongoDB = async () => {
  const {
    MongoClient
  } = require('mongodb');
  // const connection = "mongodb+srv://<username>:<password>@<your-cluster-url>/test?retryWrites=true&w=majority";
  const connection = config.mongo.connectionString
  const client = new MongoClient(connection);

  try {
    await client.connect();

    let dbs = await client.db().admin().listDatabases();
    if (dbs.databases.some(db => db.name.length > 0)) {
      console.log('Connected to MongoDB')
    }
  } catch (e) {
    console.error(e);
  } finally {
    // close connection
    await client.close();
  }
}

testOutRabbitMQ = async () => {
  var amqp = require('amqplib/callback_api');
  var conString = config.rabbitMQ.connectionString;
  amqp.connect(conString, function (error0, connection) {
    if (error0) {
      throw error0;
    }
    connection.createChannel(function (error1, channel) {
      if (error1) {
        throw error1;
      }
      var queue = 'hello';
      var msg = 'Hello world';

      channel.assertQueue(queue, {
        durable: false
      });

      channel.sendToQueue(queue, Buffer.from(msg));
      console.log("Connected to RabbitMQ");
    });

    // close connection
    setTimeout(function () {
      connection.close();
      process.exit(0)
    }, 500);
  });
}

testOutPostgres = async () => {
  var pg = require('pg');
  var conString = config.postgres.connectionString;

  var client = new pg.Client(conString);
  client.connect();

  await client
    .query('SELECT NOW() as now')
    .then(res => {
      if (res.rows[0].now instanceof Date) {
        console.log('Connected to postgres')
      }
    })
    .catch(e => console.error(e.stack))

  await client.end(err => {
    if (err) {
      console.log('error during disconnection', err.stack)
    }
  })
}

let main = async () => {
  await testOutMongoDB(config.mongo.connectionString)

  await testOutRabbitMQ(config.mongo.connectionString)

  await testOutPostgres(config.mongo.connectionString)
}

// run the test
main().catch(console.error);