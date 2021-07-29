const collectionManager = require('./src/collector/collection-manager')
const enrichment = require('gitlab-pond/src/enricher')

module.exports = {
  collector: {
    name: 'gitlabCollector',
    exec: async function (rawDb, options) {
      console.log('INFO >>> gitlab collecting')
      await collectionManager.collectProjectsDetails(options, rawDb)
      await collectionManager.collectProjectCommits(options, rawDb)
      // await collectionManager.collectProjectMergeRequests(options, rawDb)
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
      await enrichment.enrich(rawDb, enrichedDb, options)
      return []
    }
  }
}
