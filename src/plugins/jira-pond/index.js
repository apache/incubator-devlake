const issues = require('./src/collector/issues')
const enrichment = require('jira-pond/src/enricher')

module.exports = {
  collector: {
    name: 'jiraCollector',
    exec: async function (rawDb, options) {
      const args = { db: rawDb, ...options }

      console.info('start collecting jira data', options)
      await issues.collect(args)
      console.info('end collecting jira data')

      return {
        ...options,
        enricher: 'jiraEnricher'
      }
    }
  },

  enricher: {
    name: 'jiraEnricher',
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
      await module.exports.collector.exec(db, { boardId: process.argv[2], forceAll: process.argv[3] })
      // await module.exports.enricher.exec(db, enrichedDb, { boardId: process.argv[2], forceAll: process.argv[3] })
    } finally {
      dbConnector.disconnect(client)
    }
  })()
}
