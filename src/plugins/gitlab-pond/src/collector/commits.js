require('module-alias/register')
const { findOrCreateCollection } = require('../../../commondb')
const fetcher = require('./fetcher')

const collectionName = 'gitlab_commits'

module.exports = {
  async collect ({ db, projectId, forceAll }) {
    if (!projectId) {
      throw new Error('Failed to collect gitlab data, projectId is required')
    }

    await module.exports.collectByProjectId(db, projectId, forceAll)
  },

  async collectByProjectId (db, projectId, forceAll) {
    const commitsCollection = await findOrCreateCollection(db, collectionName)
    for await (const commit of fetcher.fetchPaged(`projects/${projectId}/repository/commits?with_stats=true`)) {
      commit.projectId = projectId
      await commitsCollection.findOneAndUpdate(
        { id: commit.id },
        { $set: commit },
        { upsert: true }
      )
    }
  }
}
