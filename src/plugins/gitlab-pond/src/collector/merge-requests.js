require('module-alias/register')

const { findOrCreateCollection } = require('../../../commondb')

const fetcher = require('./fetcher')

const collectionName = 'gitlab_project_merge_requests'

module.exports = {
  async collect ({ db, projectId }) {
    if (!projectId) {
      throw new Error('Failed to collect gitlab data, projectId is required')
    }
    await module.exports.collectMergeRequestsByProjectId(db, projectId)
  },
  async collectByProjectId (db, projectId, forceAll) {
    const mrsCollection = await findOrCreateCollection(db, collectionName)
    for await (const mr of fetcher.fetchPaged(`projects/${projectId}/merge_requests`)) {
      mr.projectId = projectId
      await mrsCollection.findOneAndUpdate(
        { id: mr.id },
        { $set: mr },
        { upsert: true }
      )
    }
  }
}
