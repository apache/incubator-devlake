const issues = require('./issues')
const changelogs = require('./changelogs')

module.exports = {
  async collect ({ projectId }) {
    console.log('Jira Collection, projectId:', projectId)

    await issues.collectIssues(projectId)
    await changelogs.collectChangelogs(projectId)
  },

  issues
}
