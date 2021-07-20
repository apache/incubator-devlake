require('module-alias/register')
const amqp = require('amqplib')

const connectionString = require('@config/resolveConfig').rabbitMQ.connectionString

module.exports = {
  async produce (task, queue) {
    console.log(`Publishing task to ${queue}`, task)

    const conn = await amqp.connect(connectionString, 'heartbeat=60')
    const ch = await conn.createChannel()
    const exch = `${queue}_exchange`
    const rkey = `${queue}_route`

    await ch.assertExchange(exch, 'direct', { durable: true }).catch(console.error)
    await ch.assertQueue(queue, { durable: true })
    await ch.bindQueue(queue, exch, rkey)
    await ch.publish(exch, rkey, Buffer.from(task))

    setTimeout(function () {
      ch.close()
      conn.close()
    }, 500)
  }
}
