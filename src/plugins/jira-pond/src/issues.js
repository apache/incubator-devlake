require('module-alias/register')

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
      const issuesCollection = await module.exports.findOrCreateCollection(db, collectionName)
      await issuesCollection.deleteMany({}) // temporary

      issuesResponse.issues.forEach(issue => {
        const id = Number(issue.id)

        promises.push(issuesCollection.findOneAndUpdate({
          id
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
    const requestUri = `search?jql=project="${project}"`

    return fetcher.fetch(requestUri)
  },

  async findIssues (where, db, limit = 99999999) {
    console.log('INFO >>> findIssues where', where)
    const issueCollection = await module.exports.findOrCreateCollection(db, collectionName)
    const foundIssuesCursor = await issueCollection.find(where).limit(limit)
    return await foundIssuesCursor.toArray()
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
