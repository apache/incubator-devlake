const config = require('../../config/local')

testOutMongoDB = async () => {
  const {
    MongoClient
  } = require('mongodb');
  // const connection = "mongodb+srv://<username>:<password>@<your-cluster-url>/test?retryWrites=true&w=majority";
  const connection = config.mongo.connectionString
  const client = new MongoClient(connection);

  console.log('JON >>> attempting to connect to mongodb', connection)
  try {
    await client.connect();

    let dbs = await client.db().admin().listDatabases();
    console.log('INFO: Mongo dbs')
    dbs.databases.forEach(db => console.log(` - ${db.name}`));
  } catch (e) {
    console.error(e);
  } finally {
    // close connection
    await client.close();
  }
}

testOutRabbitMQ = async () => {
  var amqp = require('amqplib/callback_api');
  amqp.connect('amqp://localhost', function (error0, connection) {
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
      console.log(" [x] Sent %s", msg);
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
  var conString = "postgres://YourUserName:YourPassword@localhost:5432/YourDatabase";

  var client = new pg.Client(conString);
  client.connect();

  var query = client.query("SELECT true;");

  query.on('row', function (row) {
    console.log(row);
  });

  // close connection
  query.on('end', function () {
    client.end();
  });
}

let main = async () => {
  await testOutMongoDB(config.mongo.connectionString)

  await testOutRabbitMQ(config.mongo.connectionString)

  await testOutPostgres(config.mongo.connectionString)
}

// run the test
main().catch(console.error);