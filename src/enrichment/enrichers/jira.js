const issueCollector = require('../../collection/collectors/jira/issues')
const { JiraIssue } = require('@db/postgres')

module.exports = {
  async enrich (config) {
    console.log('Jira Enrichment', config)
    const limit = 99999
    const issues = await issueCollector.findIssues(limit)

    issues.forEach(async issue => {
      await JiraIssue.create({
        id: issue.id,
        url: issue.self,
        title: issue.fields.summary,
        projectId: issue.fields.project.id,
        description: issue.fields.description
      })
    })

    console.log('Done enriching issues')
  }
}
