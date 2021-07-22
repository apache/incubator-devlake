const dbConnector = require('@mongo/connection')
const issues = require('./issues')
// const changelogs = require('./changelogs')

module.exports = {
  async collect ({ projectId }) {
    console.log('Jira Collection, projectId:', projectId)

    const { db } = await dbConnector.connect()
    const jiraCollectors = []

    // ^ Issues
    const jiraIssues = await issues.collectIssues(projectId)

    const setJiraIssues = new Promise(() => {
      jiraIssues.forEach(issue => {
        const id = Number(issue.id)
        const collectionName = 'jira_issues'
        const doc = db.collection(collectionName).findOneAndUpdate({ id }, { $set: issue }, { upsert: true })
        jiraCollectors.push(doc)
      })
    })

    // ^ Changelogs
    // await changelogs.collectChangelogs(projectId)

    // ^ All Jira
    Promise.all(setJiraIssues)
  }
}
