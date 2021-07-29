const commits = require('./src/collector/commits')
const mergeRequests = require('./src/collector/merge-requests')
const projects = require('./src/collector/projects')
const enrichment = require('./src/enricher')

module.exports = {
  collector: {
    name: 'gitlabCollector',
    exec: async function (rawDb, options) {
      const args = { db: rawDb, ...options }

      console.info('INFO >>> gitlab collecting', options)
      await commits.collect(args)
      await mergeRequests.collect(args)
      await projects.collect(args)
      console.info('INFO >>> gitlab collecting done!')

      return {
        ...options,
        enricher: 'gitlabEnricher'
      }
    }
  },

  enricher: {
    name: 'gitlabEnricher',
    exec: async function (rawDb, enrichedDb, options) {
      await enrichment.enrich(rawDb, enrichedDb, options)
      return []
    }
  }
}

if (require.main === module) {
  require('module-alias/register')
  const dbConnector = require('@mongo/connection')
  const enrichedDb = require('@db/postgres');

  (async function() {
    const { db, client } = await dbConnector.connect()
    try {
      await module.exports.collector.exec(db, { projectId: process.argv[2] })
      // await module.exports.enricher.exec(db, enrichedDb, { projectId: process.argv[2] })
    } finally {
      dbConnector.disconnect(client)
    }
  })()
}
