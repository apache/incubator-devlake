require('module-alias/register')

const { findOrCreateCollection } = require('commondb')
const fetcher = require('./fetcher')
const dayjs = require('dayjs')

const collectionName = 'jira_issues'

module.exports = {
  async collect ({ db, boardId, forceAll }) {
    if (!boardId) {
      throw new Error('Failed to collect jira issues, boardId is required')
    }
    await module.exports.collectByBoardId(db, boardId, forceAll)
  },

  async collectByBoardId(db, boardId, forceAll) {
    const issuesCollection = await findOrCreateCollection(db, collectionName)
    const latestUpdated = await issuesCollection.find().sort({ 'fields.updated': -1 }).limit(1).next()
    let jql = ''
    if (!forceAll && latestUpdated) {
      const jiraDate = dayjs(latestUpdated.fields.updated).format('YYYY/MM/DD HH:mm')
      jql = encodeURIComponent(`updated >= '${jiraDate}'`)
    }
    for await (const issue of fetcher.fetchPaged(`agile/1.0/board/${boardId}/issue?jql=${jql}`, 'issues')) {
      await issuesCollection.findOneAndUpdate(
        { id: issue.id },
        { $set: issue, $addToSet: { boardId } },
        { upsert: true }
      )
    }
  },

  async findIssues (where, db, limit = 99999999) {
    console.log('INFO >>> findIssues where', where)
    const issueCollection = await findOrCreateCollection(db, collectionName)
    const foundIssuesCursor = await issueCollection.find(where).limit(limit)
    return await foundIssuesCursor.toArray()
  }
}
