const issues = require('./issues')

module.exports = {
  async collect ({ projectId }) {
    console.log('Jira Collection, projectId:', projectId)

    await issues.collectIssues(projectId)
  },

  issues
}
