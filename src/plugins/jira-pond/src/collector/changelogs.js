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

      const fetchPromises = []
      const savePromises = []

      // fetch changelogs for all issues
      for (const issue of issues) {
        fetchPromises.push(module.exports.fetchChangelogForIssue(issue.id))
      }
      const changelogs = await Promise.all(fetchPromises)

      // save changelogs into mongodb
      changelogs.forEach((changelog, index) => {
        changelog && changelog.values && changelog.values.forEach(change => {
          // set changelog id as its primaryKey
          const primaryKey = change.id
          change.primaryKey = primaryKey

          savePromises.push(changelogCollection.findOneAndUpdate({
            primaryKey
          }, {
            $set: {
              issueId: issues[index].id,
              ...change
            }
          }, {
            upsert: true
          }))
        })
      })

      await Promise.all(savePromises)
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
