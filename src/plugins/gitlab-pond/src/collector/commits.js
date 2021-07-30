require('module-alias/register')

const { findOrCreateCollection } = require('../../../commondb')

const fetcher = require('./fetcher')

const collectionName = 'gitlab_project_repo_commits'

module.exports = {
  async collect ({ db, projectId }) {
    if (!projectId) {
      throw new Error('Failed to collect gitlab data, projectId is required')
    }
    await module.exports.collectCommitsByProjectId(db, projectId)
  },
  async collectCommitsByProjectId (db, projectId) {
    const commitsCollection = await findOrCreateCollection(db, collectionName)
    const requestUri = `projects/${projectId}/repository/commits?all=true&with_stats=true`
    for await (const commit of fetcher.fetchPaged(requestUri)) {
      commit.projectId = projectId
      await commitsCollection.findOneAndUpdate(
        { id: commit.id },
        { $set: commit },
        { upsert: true }
      )
    }
  }
}
