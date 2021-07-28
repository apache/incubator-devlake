require('module-alias/register')

const { findOrCreateCollection } = require('commondb')

const fetcher = require('./fetcher')

const collectionName = 'jira_issues'

module.exports = {
  async collect (options) {
    try {
      const issues = await module.exports.fetchIssues(options.projectId)

      console.log('THE ISSUES', issues)

      await module.exports.save({ issues, db: options.db })
    } catch (error) {
      console.log(error)
    }
  },

  async save ({ issues, db }) {
    try {
      const promises = []
      const issuesCollection = await findOrCreateCollection(db, collectionName)

      issues.forEach(issue => {
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
    const pageSize = 100
    const searchUri = `search?jql=project=${project}`
    let issues = []
    let startAt = 0
    let retry = 0
    const totalResults = await fetcher.fetch(`${searchUri}&fields=key`)

    while (issues.length < totalResults.total) {
      console.log('INFO >> fetching issues ', issues.length)
      try {
        const pagination = await fetcher.fetch(`${searchUri}&maxResults=${pageSize}&startAt=${startAt}`)
        issues = issues.concat(pagination.issues)
        startAt += pageSize
      } catch (e) {
        console.error(`Jira Get Issue Keys Error start:[${issues.length}] retry:[${retry}]`, { error: e })
        if (retry > 3) {
          throw e
        }
        retry++
        continue
      }
    }

    return issues
  },

  async findIssues (where, db, limit = 99999999) {
    console.log('INFO >>> findIssues where', where)
    const issueCollection = await findOrCreateCollection(db, collectionName)
    const foundIssuesCursor = await issueCollection.find(where).limit(limit)
    return await foundIssuesCursor.toArray()
  }
}
