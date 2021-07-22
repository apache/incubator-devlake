require('module-alias/register')

const dbConnector = require('@mongo/connection')
const fetcher = require('./fetcher')

const collectionName = 'jira_issues'

module.exports = {
  async collect(options) {
    try {
      const issuesResponse = await module.exports.fetchIssues(options.projectId)
      await module.exports.save({issuesResponse, db: options.db})
    } catch (error) {
      console.log(error)      
    }
  },
  
  async save (options) {
    try {
      let promises = []
      options.issuesResponse.issues.forEach(issue => {
        const id = Number(issue.id)
        promises.push(options.db.collection(collectionName).findOneAndUpdate({
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


  async fetchIssues(project) {
    const requestUri = `search?jql=project="${project}"`

    return fetcher.fetch(requestUri)
  },

  async findIssues(where, db, limit = 99999999) {
    console.log('INFO >>> findIssues where', where)
    const issueCollection = await dbConnector.findOrCreateCollection(db, collectionName)
    const foundIssuesCursor = await issueCollection.find(where).limit(limit)
    return await foundIssuesCursor.toArray()
  }
}