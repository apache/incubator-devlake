const issueUtil = require('./issues')

const dbConnector = require('@mongo/connection')
const fetcher = require('./fetcher')

const collectionName = 'jira_issue_changelogs'

module.exports = {
  async collectChangelogs (projectId) {
    const { client, db } = await dbConnector.connect()

    try {
      const issues = await issueUtil.findIssues({ 'fields.project.name': projectId }, 3)

      const changelogCollection = await dbConnector.findOrCreateCollection(db, collectionName)

      for (const issue of issues) {
        const changelog = await module.exports.fetchChangelogForIssue(issue.id)
        await changelogCollection.insertOne(changelog)
      }
    } catch (error) {
      console.error(error)
    } finally {
      dbConnector.disconnect(client)
    }
  },

  async fetchChangelogForIssue (issueId) {
    const requestUri = `issue/${issueId}/changelog`

    return fetcher.fetch(requestUri)
  }
}
