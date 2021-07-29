require('module-alias/register')

const issueCollector = require('../collector/issues')
const changelogCollector = require('../collector/changelogs')
const constants = require('@config/constants.json').jira
const { mapValue } = require('@src/util/mapping')

module.exports = {
  async enrich (rawDb, enrichedDb, projectId) {
    console.log('INFO: Starting Jira Enrichment for projectId: ', projectId)
    await module.exports.enrichLeadTimeOnIssues(
      rawDb,
      enrichedDb,
      projectId
    )
    console.log('INFO: Done enriching Jira issues')
  },

  async enrichLeadTimeOnIssues (rawDb, enrichedDb, projectId) {
    const { JiraIssue } = enrichedDb

    const issues = await issueCollector.findIssues({
      'fields.project.id': `${projectId}`
    }, rawDb)

    const upsertPromises = []
    const leadTimePromises = []
    const issuesToCreate = []
    issues.forEach(async issue => {
      leadTimePromises.push(module.exports.calculateLeadTime(issue, rawDb))
      issuesToCreate.push({
        id: issue.id,
        url: issue.self,
        title: issue.fields.summary,
        projectId: issue.fields.project.id,
        issueType: mapValue(issue.fields.issuetype.name, constants.mappings)
        // description: issue.fields.description
      })
    })

    const leadTimes = await Promise.all(leadTimePromises)

    leadTimes.forEach((leadTime, index) => {
      let issue = issuesToCreate[index]
      console.log('INFO: issueId & leadTime', issue.id, leadTime)
      issue = {
        leadTime,
        ...issue
      }
      // Create all new records
      upsertPromises.push(JiraIssue.upsert(issue))
    })

    await Promise.all(upsertPromises)
  },

  async calculateLeadTime (issue, db) {
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

          if (!constants.mappings.Closed.includes(item.fromString)) {
            const elapsedTime = changeTime - lastTime

            leadTime += elapsedTime
          }

          lastTime = changeTime
          isDone = constants.mappings.Closed.includes(item.toString)
        }
      }
    }

    return isDone
      ? Math.round(leadTime / 1000)
      : 0
  }
}
