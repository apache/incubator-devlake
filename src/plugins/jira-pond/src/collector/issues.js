require('module-alias/register')

const { findOrCreateCollection } = require('commondb')
const fetcher = require('./fetcher')
const dayjs = require('dayjs')

async function collect ({ db, boardId, forceAll }) {
  if (!boardId) {
    throw new Error('Failed to collect jira issues, boardId is required')
  }
  console.info('INFO >>> jira collecting issues for board', boardId)
  await collectByBoardId(db, boardId, forceAll)
  console.info('INFO >>> jira collecting issues for board done!', boardId, counter)
}

async function collectByBoardId (db, boardId, forceAll) {
  const issuesCollection = await getCollection(db)
  const latestUpdated = await issuesCollection.find().sort({ 'fields.updated': -1 }).limit(1).next()
  const $addToSet = { boardIds: boardId }
  let jql = ''
  if (!forceAll && latestUpdated) {
    const jiraDate = dayjs(latestUpdated.fields.updated).format('YYYY/MM/DD HH:mm')
    jql = encodeURIComponent(`updated >= '${jiraDate}' ORDER BY updated ASC`)
  }
  let counter = 0
  for await (const issue of fetcher.fetchPaged(`agile/1.0/board/${boardId}/issue?jql=${jql}`, 'issues')) {
    await issuesCollection.findOneAndUpdate(
      { id: issue.id },
      { $set: issue, $addToSet },
      { upsert: true }
    )
    counter++
  }
}

async function getCollection (db) {
  return await findOrCreateCollection(db, 'jira_issues')
}

module.exports = { collect, getCollection }
