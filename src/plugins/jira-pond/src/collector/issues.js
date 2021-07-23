require('module-alias/register')

const { findOrCreateCollection } = require('commondb')

const fetcher = require('./fetcher')

const collectionName = 'jira_issues'

module.exports = {
  async collect (options) {
    try {
      const issuesResponse = await module.exports.fetchIssues(options.projectId)

      await module.exports.save({ issuesResponse, db: options.db })
    } catch (error) {
      console.log(error)
    }
  },

  async save ({
    issuesResponse,
    db
  }) {
    try {
      const promises = []
      const issuesCollection = await findOrCreateCollection(db, collectionName)

      issuesResponse.forEach(issue => {
        issue.primaryKey = Number(issue.id)

        promises.push(issuesCollection.findOneAndUpdate({
          primaryKey: issue.primaryKey
        }, {
          $set: issue
        }, {
          upsert: true
        }))
      })

      await Promise.all(promises)
    } catch (error) {
      console.error(error)
    }
  },

  async fetchIssues (project) {
    let issues = []
    let retry = 0
    const pageSize = 100
    const searchUri = `search?jql=project=${project}`
    const totalResponse = await fetcher.fetch(`${searchUri}&fields=key`)
    let total = totalResponse.total

    const deleteMeFakeTotalForTesting = 300
    total = deleteMeFakeTotalForTesting

    console.log(`INFO: Fetching ${total} issues from project: ${project}`)

    while (issues.length < total) {
      try {
        const pagination = await fetcher.fetch(`${searchUri}&maxResults=${pageSize}&startAt=${issues.length}`)
        issues = issues.concat(pagination.issues)
      } catch (e) {
        console.error(`Jira Get Issue Keys Error start:[${issues.length}] retry:[${retry}]`, { error: e })
        if (retry > 3) {
          throw e
        }
        retry++
        continue
      }
    }
    console.log('JON >>> fetchIssues issues.length', issues.length)
    return issues
  },

  async findIssues (where, db, limit = 99999999) {
    console.log('INFO >>> findIssues where', where)
    const issueCollection = await findOrCreateCollection(db, collectionName)
    const foundIssuesCursor = await issueCollection.find(where).limit(limit)
    return await foundIssuesCursor.toArray()
  }
}
