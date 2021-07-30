require('module-alias/register')
const { findOrCreateCollection } = require('../../../commondb')
const fetcher = require('./fetcher')

async function collect ({ db, projectId, forceAll }) {
  if (!projectId) {
    throw new Error('Failed to collect gitlab data, projectId is required')
  }

  await collectByProjectId(db, projectId, forceAll)
}

async function collectByProjectId (db, projectId, forceAll) {
  console.info('INFO >>> gitlab collecting commits for project', projectId)
  const commitsCollection = await getCollection(db)
  for await (const commit of fetcher.fetchPaged(`projects/${projectId}/repository/commits?with_stats=true`)) {
    commit.projectId = projectId
    await commitsCollection.findOneAndUpdate(
      { id: commit.id },
      { $set: commit },
      { upsert: true }
    )
  }
  console.info('INFO >>> gitlab collecting commits for project done!', projectId)
}

async function getCollection (db) {
  return await findOrCreateCollection(db, 'gitlab_commits')
}

module.exports = { collect, getCollection }
