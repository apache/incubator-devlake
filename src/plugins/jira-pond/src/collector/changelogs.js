const {
  findOrCreateCollection
} = require('commondb')

const issueUtil = require('./issues')

const fetcher = require('./fetcher')

const collectionName = 'jira_issue_changelogs'

module.exports = {
  async collect (options) {
    try {
      console.log('INFO >>> collecting changelogs...', options.projectId)
      const issues = await issueUtil.findIssues({
        'fields.project.id': `${options.projectId}`
      }, options.db)
      console.info('INFO: issues length', issues.length)

      const changelogCollection = await findOrCreateCollection(options.db, collectionName)

      const promises = []
      for (const issue of issues) {
        // todo we cant have this line. It needs to be a promise.all async
        const changelog = await module.exports.fetchChangelogForIssue(issue.id)
        for (const change of changelog.values) {
          const primaryKey = changelog.id

          promises.push(changelogCollection.findOneAndUpdate({
            primaryKey
          }, {
            $set: {
              issueId: issue.id,
              ...change
            }
          }, {
            upsert: true
          }))
        }
      }
      await Promise.all(promises)
    } catch (error) {
      console.log(error)
    }
  },

  async fetchChangelogForIssue (issueId) {
    const requestUri = `issue/${issueId}/changelog`

    return fetcher.fetch(requestUri)
  },

  async findChangelogs (where, db, limit = 999999, sort = {
    createdAt: 1
  }) {
    const changelogCollection = await findOrCreateCollection(db, collectionName)
    const changelogsCursor = await changelogCollection.find(where).limit(limit).sort(sort)

    return await changelogsCursor.toArray()
  }
}
