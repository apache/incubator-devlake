const axios = require('axios')

const { findOrCreateCollection } = require('../../util/collectionDB')
const config = require('@config/resolveConfig').jira

module.exports = {
  async collectIssues (db, projectId) {
    const issues = await module.exports.fetchIssues(projectId)

    const issueCollection = await findOrCreateCollection(db, 'jira_issues')

    // Insert issues into mongodb
    await issueCollection.insertMany(issues)
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
