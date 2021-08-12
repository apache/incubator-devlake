const collection = require('./src/collector')
const enrichment = require('./src/enricher')

module.exports = {
  configuration: {
    collection: null
  },

  configure (configuration) {
    module.exports.configuration = configuration
    collection.configure(configuration.collection)
  },

  collector: {
    name: 'gitlabCollector',
    exec: async function (rawDb, options) {
      console.info('INFO >>> gitlab collecting', options)
      await collection.collect(rawDb, options)
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
      console.info('INFO >>> gitlab enriching', options)
      await enrichment.enrich(rawDb, enrichedDb, options)
      console.info('INFO >>> gitlab enriching done!')
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
    const configuration = require('@config/plugins-conf.js').find(p => p.name === 'gitlab').configuration

    const projectId = process.argv[2] || 24547305
    const { db, client } = await dbConnector.connect()
    await module.exports.configure(configuration)
    try {
      await module.exports.collector.exec(db, { projectId })
      await module.exports.enricher.exec(db, enrichedDb, { projectId })
    } finally {
      dbConnector.disconnect(client)
    }
  }
  main()
}
