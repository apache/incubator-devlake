const producer = require('../queue/producer')

module.exports = {
  async createJob (project) {
    const msg = JSON.stringify(project)
    const queue = 'enrichment'

    await producer.produce(msg, queue)
  }
}
