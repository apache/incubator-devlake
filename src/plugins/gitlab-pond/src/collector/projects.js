require('module-alias/register')
const { findOrCreateCollection } = require('../../../commondb')
const fetcher = require('./fetcher')

async function collect ({ db, projectId, forceAll }) {
  if (!projectId) {
    throw new Error('Failed to collect gitlab data, projectId is required')
  }

  console.info('INFO >>> gitlab collecting project', projectId)
  await collectByProjectId(db, projectId, forceAll)
  console.info('INFO >>> gitlab collecting project done!', projectId)
}

async function collectByProjectId (db, projectId, forceAll) {
  const projectsCollection = await getCollection(db)
  const response = await fetcher.fetch(`projects/${projectId}`)
  const project = response.data
  await projectsCollection.findOneAndUpdate(
    { id: project.id },
    { $set: project },
    { upsert: true }
  )
}

async function getCollection (db) {
  return await findOrCreateCollection(db, 'gitlab_projects')
}

module.exports = { collect, getCollection }
