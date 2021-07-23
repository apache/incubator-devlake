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

      const insertDataPromises = []
      const changelogPromises = []
      for (const issue of issues) {
        changelogPromises.push(module.exports.fetchChangelogForIssue(issue.id))
      }
      let changelogs = await Promise.all(changelogPromises)
      changelogs.forEach((changelog, index) => {
        for (const change of changelog.values) {
          const primaryKey = changelog.id
  
          insertDataPromises.push(changelogCollection.findOneAndUpdate({
            primaryKey
          }, {
            $set: {
              issueId: issues[index].id,
              ...change
            }
          }, {
            upsert: true
          }))
        }
      })
      await Promise.all(insertDataPromises)
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
