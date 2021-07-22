module.exports = {
  collector: {
    name: 'jiraCollector',
    exec: function (rawDb, options) {
      // do stuff
      console.log('collection')
      return ['jira_enricher']
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
