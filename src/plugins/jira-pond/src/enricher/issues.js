
const issuesCollecotr = require('../collector/issues')
const dayjs = require('dayjs')
const duration = require('dayjs/plugin/duration')
const { merge, isEmpty, isArray } = require('lodash')
dayjs.extend(duration)

const configuration = {
  verified: false,
  typeMappings: [
    {
      originTypes: ['Bug'],
      standardType: 'Bug',
      statusMappings: [
        { originStatuses: ['Rejected', 'Abandoned', 'Cancelled', 'ByDesign', 'Irreproducible'], standardStatus: 'Rejected' },
        { originStatuses: ['Resolved', 'Approved', 'Verified', 'Done', 'Closed'], standardStatus: 'Resolved' }
      ]
    },
    {
      originTypes: ['Incident'],
      standardType: 'Incident',
      statusMappings: [
        { originStatuses: ['Rejected', 'Abandoned', 'Cancelled', 'Irreproducible'], standardStatus: 'Rejected' },
        { originStatuses: ['Resolved', 'Approved', 'Verified', 'Done', 'Closed'], standardStatus: 'Resolved' }
      ]
    },
    {
      originTypes: ['Story'],
      standardType: 'Requirement',
      statusMappings: [
        { originStatuses: ['Rejected', 'Abandoned', 'Cancelled', 'OnHold'], standardStatus: 'Rejected' },
        { originStatuses: ['Resolved', 'Done', 'Closed'], standardStatus: 'Resolved' }
      ]
    }
  ],
  epicKeyField: null
}

function configure (config) {
  merge(configuration, config)
  configuration.verified = false

  const { epicKeyField, typeMappings } = configuration
  if (!epicKeyField || !epicKeyField.startsWith('customfield')) {
    throw new Error('jira enrichment configuration error: issue.epicKeyField is invalid')
  }

  const isValidArray = (a) => isArray(a) && !isEmpty(a)

  typeMappings.forEach(mapping => {
    if (!isValidArray(mapping.originTypes)) {
      throw new Error('jira configuration error: issue.typeMappings.originTypes is invalid: ', mapping.originTypes)
    }
  })
  configuration.verified = true
}

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

  const { epicKeyField, typeMappings } = configuration

  try {
    let counter = 0
    while (await curosr.hasNext()) {
      const issue = await curosr.next()
      const enriched = {
        id: issue.id,
        url: issue.self,
        title: issue.fields.summary,
        projectId: issue.fields.project.id,
        issueType: issue.fields.issuetype.name,
        epicKey: issue.fields[epicKeyField],
        status: issue.fields.status.name,
        issueCreatedAt: issue.fields.created,
        issueUpdatedAt: issue.fields.updated,
        issueResolvedAt: issue.fields.resolutiondate,
        leadTime: null
      }
      const typeMapping = typeMappings.find(tm => tm.originTypes.includes(enriched.issueType))
      if (typeMapping) {
        enriched.issueType = typeMapping.standardType
        const statusMapping = typeMapping.statusMappings.find(sm => sm.originStatuses.includes(enriched.status))
        if (statusMapping) {
          enriched.status = statusMapping.standardStatus
        }
      }
      // by standard, leadtime = days of (resolutiondate - creationdate)
      if (enriched.status === 'Resolved' && issue.fields.resolutiondate) {
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

module.exports = { configure, enrich }
