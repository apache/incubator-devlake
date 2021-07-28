const collectionManager = require('./src/collector/collection-manager')

module.exports = {
  collector: {
    name: 'gitlabCollector',
    exec: async function (rawDb, options) {
      const { projectId } = options
      console.log('rawDb', rawDb)
      console.log('projectId', projectId)
      await collectionManager.collectAll()
      console.log('INFO >>> done collecting')

      return {
        ...options,
        enricher: 'gitlabEnricher'
      }
    }
  },

  enricher: {
    name: 'gitlabEnricher',
    exec: async function (rawDb, enrichedDb, options) {
      const { projectId } = options
      console.log('rawDb, enrichedDb, projectId', rawDb, enrichedDb, projectId)
      // await enrichment.enrich(rawDb, enrichedDb, projectId)

      return []
    }
  }
}
