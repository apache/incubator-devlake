const issueCollector = require('../collection/collectors/jira/issues')

module.exports = {
  async enrich (config) {
    console.log('Jira Enrichment', config)

    const issues = await issueCollector.findIssues(1)
    console.log('paul >>> issues', issues)
  }
}
