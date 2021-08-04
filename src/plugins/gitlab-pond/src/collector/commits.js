require('module-alias/register')
const { findOrCreateCollection } = require('../../../commondb')
const fetcher = require('./fetcher')

async function collect ({ db, projectId, branch, forceAll }) {
  if (!projectId) {
    throw new Error('Failed to collect gitlab data, projectId is required')
  }

  await collectByProjectId(db, projectId, branch, forceAll)
}

async function collectByProjectId (db, projectId, branch, forceAll) {
  console.info('INFO >>> gitlab collecting commits for project', projectId)
  const commitsCollection = await getCollection(db)

  let queryParams = 'withStats=true'
  // in some cases, the user does not want to pull commits from the default branch.
  if (branch) {
    queryParams += `&ref_name=${branch}`
  }

  for await (const commit of fetcher.fetchPaged(`projects/${projectId}/repository/commits?${queryParams}`)) {
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
