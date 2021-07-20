require('module-alias/register')

const issueCollector = require('../../collection/collectors/jira/issues')
const changelogCollector = require('../../collection/collectors/jira/changelogs')

const { JiraIssue } = require('@db/postgres')

const closedStatuses = ['Done', 'Closed']

module.exports = {
  async enrich ({ projectId }) {
    console.log('Jira Enrichment', projectId)

    const issues = await issueCollector.findIssues({ 'fields.project.name': `${projectId}` })

    issues.forEach(async issue => {
      const leadTime = await module.exports.calculateLeadTime(issue)

      await JiraIssue.create({
        id: issue.id,
        url: issue.self,
        title: issue.fields.summary,
        projectId: issue.fields.project.id,
        description: issue.fields.description,
        leadTime
      })
    })

    console.log('Done enriching issues')
  },

  async calculateLeadTime (issue) {
    const changelogs = await changelogCollector.findChangelogs({ issueId: `${issue.id}` })

    let leadTime = 0
    let lastTime = new Date(issue.fields.created).getTime()
    let isDone = false

    for (const change of changelogs) {
      for (const item of change.items) {
        if (item.field === 'status') {
          const changeTime = new Date(change.created).getTime()

          if (!closedStatuses.includes(item.fromString)) {
            const elapsedTime = changeTime - lastTime

            leadTime += elapsedTime
          }

          lastTime = changeTime
          isDone = closedStatuses.includes(item.toString)
        }
      }
    }

    return isDone
      ? Math.round(leadTime / 1000)
      : 0
  }
}
