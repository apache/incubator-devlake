require('module-alias/register')

const issueCollector = require('../../collection/collectors/jira/issues')
const changelogCollector = require('../../collection/collectors/jira/changelogs')

const {
  JiraIssue
} = require('@db/postgres')
const dbConnector = require('@db/mongo/connection')
const closedStatuses = ['Done', 'Closed']

module.exports = {
  async enrich ({
    projectId
  }) {
    const {
      db,
      client
    } = await dbConnector.connect()

    try {
      console.log('Jira Enrichment', projectId)
      await module.exports.enrichLeadTimeOnIssues({
        db,
        projectId
      })
      console.log('Done enriching issues')
    } catch (error) {
      console.log('>>> error', error)
    } finally {
      dbConnector.disconnect(client)
    }
  },

  async enrichLeadTimeOnIssues (options) {
    const issues = await issueCollector.findIssues({
      'fields.project.id': `${options.projectId}`
    }, options.db)

    const creationPromises = []
    const leadTimePromises = []
    const issuesToCreate = []
    issues.forEach(async issue => {
      leadTimePromises.push(module.exports.calculateLeadTime(issue, options.db))
      issuesToCreate.push({
        id: issue.id,
        url: issue.self,
        title: issue.fields.summary,
        projectId: issue.fields.project.id
        // description: issue.fields.description
      })
    })

    const leadTimes = await Promise.all(leadTimePromises)

    leadTimes.forEach((leadTime, index) => {
      console.log('JON >>> leadTime', leadTime)
      let issue = issuesToCreate[index]
      issue = {
        leadTime,
        ...issue
      }
      creationPromises.push(JiraIssue.create(issue))
    })

    await Promise.all(creationPromises)
  },

  async calculateLeadTime (issue, db) {
    console.log('JON >>> db', db)
    const changelogs = await changelogCollector.findChangelogs({
      issueId: `${issue.id}`
    }, db)

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
