const collectionManager = require('./src/collector/collection-manager')
const enrichment = require('gitlab-pond/src/enricher')
const commits = require('./src/collector/commits')
const mergeRequests = require('./src/collector/merge-requests')

module.exports = {
  collector: {
    name: 'gitlabCollector',
    exec: async function (rawDb, options) {
      console.log('INFO >>> gitlab collecting')
      console.log(options)
      await collectionManager.collectProjectsDetails(options, rawDb)
      for (const projectId of options.projectIds) {
        await commits.collect({ rawDb, projectId })
      }
      for (const projectId of options.projectIds) {
        await mergeRequests.collect({ rawDb, projectId })
      }
      // await collectionManager.collectProjectCommits(options, rawDb)
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
