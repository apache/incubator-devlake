const issueUtil = require('./issues')

const dbConnector = require('@mongo/connection')
const fetcher = require('./fetcher')

const collectionName = 'jira_issue_changelogs'

module.exports = {
  async collectChangelogs (projectId) {
    const { client, db } = await dbConnector.connect()

    try {
      const issues = await issueUtil.findIssues({ 'fields.project.id': projectId })

      const changelogCollection = await dbConnector.findOrCreateCollection(db, collectionName)

      for (const issue of issues) {
        const changelog = await module.exports.fetchChangelogForIssue(issue.id)

        for (const change of changelog.values) {
          await changelogCollection.insertOne({
            issueId: issue.id,
            ...change
          })
        }
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
  },

  async findChangelogs (where, limit = 999999, sort = { createdAt: 1 }) {
    const { client, db } = await dbConnector.connect()

    let changelogs = []

    try {
      const changelogCollection = await dbConnector.findOrCreateCollection(db, collectionName)
      const changelogsCursor = await changelogCollection.find(where).limit(limit).sort(sort)

      changelogs = await changelogsCursor.toArray()
    } finally {
      dbConnector.disconnect(client)
    }

    return changelogs
  }
}
