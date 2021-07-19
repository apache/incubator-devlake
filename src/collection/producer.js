const amqp = require('amqplib')

const amqpUrl = 'amqp://guest:guest@localhost:5672/rabbitmq'

module.exports = {
  async produce (task, q) {
    console.log('Publishing', task, q)
    const conn = await amqp.connect(amqpUrl, 'heartbeat=60')
    const ch = await conn.createChannel()
    const exch = 'collection_exchange'
    const rkey = 'collection_route'

    await ch.assertExchange(exch, 'direct', { durable: true }).catch(console.error)
    await ch.assertQueue(q, { durable: true })
    await ch.bindQueue(q, exch, rkey)
    await ch.publish(exch, rkey, Buffer.from(task))
    setTimeout(function () {
      ch.close()
      conn.close()
    }, 500)
  }
}
