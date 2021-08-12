const collection = require('./src/collector')
const enrichment = require('./src/enricher')

module.exports = {
  configuration: {
    // default configuration which could be overrided by `config/plugins.js`
    collection: null,
    enrichment: null
  },

  configure (configuration) {
    module.exports.configuration = configuration
    collection.configure(configuration.collection)
    enrichment.configure(configuration.enrichment)
  },

  collector: {
    name: 'jiraCollector',
    exec: async function (rawDb, options) {
      console.info('INFO >>> jira collecting', options)
      await collection.collect(rawDb, options)
      console.info('INFO >>> jira collecting done!')

      return {
        ...options,
        enricher: 'jiraEnricher'
      }
    }
  },

  enricher: {
    name: 'jiraEnricher',
    exec: async function (rawDb, enrichedDb, options) {
      console.info('INFO >>> jira enriching', options)
      await enrichment.enrich(rawDb, enrichedDb, options)
      console.info('INFO >>> jira enriching done!', options)
      return []
    }
  }
}

// for debugging only, skip if module being required
if (require.main === module) {
  async function main () {
    require('module-alias/register')
    const dbConnector = require('@mongo/connection')
    const enrichedDb = require('@db/postgres')
    const configuration = require('@config/plugins-conf.js').find(p => p.name === 'jira').configuration

    const boardId = Number(process.argv[2]) || 8
    const forceCollectAll = Number(process.argv[3])
    const forceEnrichAll = Number(process.argv[4])
    const { db, client } = await dbConnector.connect()
    await module.exports.configure(configuration)
    try {
      await module.exports.collector.exec(db, { boardId, forceAll: forceCollectAll })
      await module.exports.enricher.exec(db, enrichedDb, { boardId, forceAll: forceEnrichAll })
    } finally {
      dbConnector.disconnect(client)
    }
  }
  main()
}
