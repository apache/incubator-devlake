const producer = require('../queue/producer')

module.exports = {
  async createJobs (project) {
    await producer.produce(JSON.stringify(project), 'enrichment')
  }
}
