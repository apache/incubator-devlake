const issueUtil = require('./issues')

const fetcher = require('./fetcher')

const collectionName = 'jira_issue_changelogs'

module.exports = {
  async collect (options) {
    const issues = await issueUtil.findIssues({ 'fields.project.id': options.projectId }, options.db)

    const changelogCollection = await module.exports.findOrCreateCollection(options.db, collectionName)
    await changelogCollection.deleteMany({}) // temporary

    const promises = []
    for (const issue of issues) {
      // todo we cant have this line. It needs to be a promise.all async
      const changelog = await module.exports.fetchChangelogForIssue(issue.id)
      console.log('JON >>> changelog', changelog)

      for (const change of changelog.values) {
        // todo we need to add our own primary key
        // todo only update based on primary key
        promises.push(changelogCollection.insertOne({
          issueId: issue.id,
          ...change
        }))
      }
    }
    await Promise.all(promises)
  },

  async fetchChangelogForIssue (issueId) {
    const requestUri = `issue/${issueId}/changelog`

    return fetcher.fetch(requestUri)
  },

  async findChangelogs (where, db, limit = 999999, sort = { createdAt: 1 }) {
    const changelogCollection = await module.exports.findOrCreateCollection(db, collectionName)
    const changelogsCursor = await changelogCollection.find(where).limit(limit).sort(sort)

    return await changelogsCursor.toArray()
  },

  async findOrCreateCollection (db, collectionName, options = {}) {
    try {
      const foundCollectionsCursor = await db.listCollections()
      const foundCollections = await foundCollectionsCursor.toArray()

      // check if Jira collection exists
      const collectionExists = foundCollections
        .some(collection => collection.name === collectionName)

      return collectionExists
        ? await db.collection(collectionName)
        : await db.createCollection(collectionName, options)
    } catch (e) {
      console.log('MONGO.DB createCollection() >> ERROR: ', e)
    }
  }
}
