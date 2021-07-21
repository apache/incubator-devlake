require('module-alias/register')

const dbConnector = require('@mongo/connection')
const fetcher = require('./fetcher')

const collectionName = 'jira_issues'

module.exports = {
  async collectIssues (projectId) {
    const { client, db } = await dbConnector.connect()

    try {
      const { issues } = await module.exports.fetchIssues(projectId)

      await db.collection(collectionName).remove()

      const issueCollection = await dbConnector.findOrCreateCollection(db, collectionName)

      // Insert issues into mongodb
      await issueCollection.insertMany(issues)
    } catch (error) {
      console.error(error)
    } finally {
      dbConnector.disconnect(client)
    }
  },

  async fetchIssues (project) {
    const requestUri = `search?jql=project="${project}"`

    return fetcher.fetch(requestUri)
  },

  async findIssues (where, limit = 99999999) {
    const { client, db } = await dbConnector.connect()

    let issues = []

    try {
      const issueCollection = await dbConnector.findOrCreateCollection(db, collectionName)
      const foundIssuesCursor = await issueCollection.find(where).limit(limit)

      issues = await foundIssuesCursor.toArray()
    } finally {
      dbConnector.disconnect(client)
    }

    return issues
  }
}
