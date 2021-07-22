const issues = require('./src/collector/issues')
const changelogs = require('./src/collector/changelogs')
const enrichment = require('jira-pond/src/enricher')

module.exports = {
  collector: {
    name: 'jiraCollector',
    exec: async function (rawDb, options) {
      const { projectId } = options

      console.log('collection >>> projectId', projectId)
      await issues.collect({ db: rawDb, projectId })
      await changelogs.collect({ db: rawDb, projectId })

      console.log('INFO >>> done collecting')

      return ['jiraEnricher']
    }
  },

  enricher: {
    name: 'jiraEnricher',
    exec: async function (rawDb, enrichedDb, options) {
      const { projectId } = options

      console.log('enricher >>> projectId', projectId)
      await enrichment.enrich(rawDb, enrichedDb, projectId)
      console.log('enrichment')
      return []
    }
  }
}
