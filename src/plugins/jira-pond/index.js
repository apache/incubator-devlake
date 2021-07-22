module.exports = {
  collector: { 
    name: 'jiraCollector',
    exec: function(raw_db, options) {
        // do stuff
        console.log('collection')
        return ['jira_enricher']
    }
  },

  enricher: { 
    name: 'jiraEnricher',
    exec: function(raw_db, enriched_db, options) {
      // do stuff
      console.log('enrichment')
      return []
    } 
  }
}