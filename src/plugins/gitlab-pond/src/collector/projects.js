require('module-alias/register')
const { findOrCreateCollection } = require('../../../commondb')
const fetcher = require('./fetcher')

const collectionName = 'gitlab_projects'

module.exports = {
  async collect ({ db, projectId, forceAll }) {
    if (!projectId) {
      throw new Error('Failed to collect gitlab data, projectId is required')
    }

    await module.exports.collectByProjectId(db, projectId, forceAll)
  },

  async collectByProjectId (db, projectId, forceAll) {
    const projectsCollection = await findOrCreateCollection(db, collectionName)
    const response = await fetcher.fetch(`projects/${projectId}`)
    const project = response.data
    await projectsCollection.findOneAndUpdate(
      { id: project.id },
      { $set: project },
      { upsert: true }
    )
  }
}
