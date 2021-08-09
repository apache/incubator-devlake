const producer = require('../queue/producer')

module.exports = {
  async createJobs (message) {
    await producer.produce(JSON.stringify(message), 'collection')
  }
}
