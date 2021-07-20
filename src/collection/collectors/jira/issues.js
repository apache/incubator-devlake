const axios = require('axios')

const config = require('@config/resolveConfig').jira
const dbConnector = require('@mongo/connection')

const collectionName = 'jira_issues'

module.exports = {
  async collectIssues (projectId) {
    const { client, db } = await dbConnector.connect()

    try {
      const issues = await module.exports.fetchIssues(projectId)

      const issueCollection = await dbConnector.findOrCreateCollection(db, collectionName)

      // Insert issues into mongodb
      await issueCollection.insertMany(issues)
    } catch (error) {
      console.error(error)
    } finally {
      await client.close()
    }
  },

  async findIssues () {
    const { client, db } = await dbConnector.connect()

    let issues = []

    try {
      const issueCollection = await dbConnector.findOrCreateCollection(db, collectionName)
      const foundIssuesCursor = await issueCollection.find()

      issues = await foundIssuesCursor.toArray()
    } finally {
      client.close()
    }

    return issues
  },

  async fetchIssues (project) {
    try {
      const response = await axios.get(`${config.host}/rest/api/3/search?jql=project="${project}"`, {
        headers: {
          Accept: 'application/json',
          Authorization: `Basic ${config.basicAuth}`
        }
      })

      return response.data.issues
    } catch (error) {
      console.error(error)
    }
  }
}
