const producer = require('../queue/producer')

module.exports = {
  async createJob (input) {
    const msg = JSON.stringify(input)
    const queue = 'enrichment'

    await producer.produce(msg, queue)
  }
}
