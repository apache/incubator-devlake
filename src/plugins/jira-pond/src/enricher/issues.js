
const issuesCollecotr = require('../collector/issues')
const constants = require('@config/constants.json').jira
const { mapValue } = require('@src/util/mapping')
const dayjs = require('dayjs')
const duration = require('dayjs/plugin/duration')
dayjs.extend(duration)

async function enrich ({ rawDb, enrichedDb, boardId, forceAll }) {
  // TODO: parameter checking

  await enrichIssues(rawDb, enrichedDb, boardId, forceAll)
}

async function enrichIssues (rawDb, enrichedDb, boardId, forceAll) {
  console.info(`INFO >>> Jira enriching issues for board #${boardId}, forceAll: ${forceAll}`)
  const issueCollection = await issuesCollecotr.getCollection(rawDb)
  const { JiraIssue, JiraBoardIssue } = enrichedDb
  // filtering out portion of records that need to be enriched
  const curosr = (
    forceAll
      ? issueCollection.find({ boardIds: boardId })
      : issueCollection.find({ $where: 'this.enriched < this.fields.updated || !this.enriched', boardIds: boardId })
  )

  try {
    let counter = 0
    while (await curosr.hasNext()) {
      const issue = await curosr.next()
      const enriched = {
        id: issue.id,
        url: issue.self,
        title: issue.fields.summary,
        projectId: issue.fields.project.id,
        issueType: mapValue(issue.fields.issuetype.name, constants.mappings),
        epicKey: issue.fields[constants.epicKeyField],
        status: issue.fields.status.name,
        issueCreatedAt: issue.fields.created,
        issueUpdatedAt: issue.fields.updated,
        issueResolvedAt: issue.fields.resolutiondate,
        leadTime: null
      }
      // by standard, leadtime = days of (resolutiondate - creationdate)
      if (issue.fields.resolutiondate) {
        enriched.leadTime = dayjs.duration(dayjs(issue.fields.resolutiondate) - dayjs(issue.fields.created)).days()
      }
      await JiraIssue.upsert(enriched)
      // update board-issue ManyToMany relationship
      await JiraBoardIssue.destroy({ where: { issueId: issue.id } })
      if (issue.boardIds) {
        for (const boardId of issue.boardIds) {
          await JiraBoardIssue.create({ boardId, issueId: issue.id })
        }
      }
      // update enrichment timestamp
      await issueCollection.updateOne(
        { id: issue.id },
        { $set: { enriched: issue.fields.updated } }
      )
      counter++
    }
    console.log('INFO >>> Jira total enriched issues: ', counter)
  } finally {
    await curosr.close()
  }
  console.info('INFO >>> Jira enriching issues done!')
}

module.exports = { enrich }
