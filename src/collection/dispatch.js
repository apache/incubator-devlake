const _has = require('lodash/has')

const producer = require('./producer')

module.exports = {
  async createJobs (project) {
    await producer.produce(JSON.stringify(project), 'collection')
  }
}