const issues = require('./src/issues')
const changelogs = require('./src/changelogs')

module.exports = {
  collector: {
    name: 'jiraCollector',
    exec: async function (rawDb, options) {
      const { projectId } = options

      console.log('paul >>> projectId', projectId)
      await issues.collect({ db: rawDb, projectId })
      await changelogs.collect({ db: rawDb, projectId })

      console.log('INFO >>> done collecting')

      return ['jiraEnricher']
    }
  },

  enricher: {
    name: 'jiraEnricher',
    exec: function (rawDb, enrichedDbs, options) {
      // do stuff
      console.log('enrichment')
      return []
    }
  }
}
