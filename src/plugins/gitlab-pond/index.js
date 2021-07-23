// const issues = require('./src/collector/issues')
// const changelogs = require('./src/collector/changelogs')
// const enrichment = require('jira-pond/src/enricher')


module.exports = {
  collector: {
    name: 'gitlabCollector',
    exec: async function (rawDb, options) {
      const { projectId } = options

      await issues.collect({ db: rawDb, projectId })
      await changelogs.collect({ db: rawDb, projectId })

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

      await enrichment.enrich(rawDb, enrichedDb, projectId)

      return []
    }
  }
}
