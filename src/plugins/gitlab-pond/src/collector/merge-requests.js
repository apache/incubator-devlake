require('module-alias/register')

const { findOrCreateCollection } = require('../../../commondb')

const fetcher = require('./fetcher')

const collectionName = 'gitlab_project_merge_requests'

module.exports = {
  async collect (options) {
    try {
      const mergeRequestsResponse = await module.exports.fetchProjectRepoMergeRequests(options.projectId)

      await module.exports.save({ mergeRequestsResponse, db: options.db })
    } catch (error) {
      console.log(error)
    }
  },
  async save ( {response, db} ){
    try {
      const promises = []
      const mergeRequestsCollection = await findOrCreateCollection(db, collectionName)
      response.forEach(mergeRequest => {
        mergeRequest.primaryKey = mergeRequest.id

        promises.push(mergeRequestsCollection.findOneAndUpdate({
          primaryKey: mergeRequest.primaryKey
        }, {
          $set: mergeRequest
        }, {
          upsert: true
        }))
      })

      await Promise.all(promises)
    } catch (error) {
      console.error(error)
    }
  },
  async fetchProjectMergeRequests (projectId) {
    const requestUri = `projects/${projectId}/merge_requests`

   return fetcher.fetch(requestUri)
 },
  async findMergeRequests (where, db, limit = 99999999) {
    console.log('INFO >>> findMergeRequests where', where)
    const mergeRequestsCollection = await findOrCreateCollection(db, collectionName)
    const foundMergeRequestsCursor = await mergeRequestsCollection.find(where).limit(limit)
    return await foundMergeRequestsCursor.toArray()
  }
}