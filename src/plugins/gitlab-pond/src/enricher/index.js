const projects = require('./projects')
const commits = require('./commits')
const mergeRequests = require('./merge-requests')
const notes = require('./notes')
const projectsCollector = require('../collector/projects')

async function enrich (rawDb, enrichedDb, { projectId }) {
  // verify collected data existence
  const projectsCollection = await projectsCollector.getCollection(rawDb)
  const project = await projectsCollection.findOne({ id: projectId })
  if (!project) {
    throw new Error(`gitlabEnricher unable to find collected data for project ${projectId}`)
  }

  const args = { rawDb, enrichedDb, projectId: Number(projectId) }
  await projects.enrich(args)
  await commits.enrich(args)
  await mergeRequests.enrich(args)
  await notes.enrich(args)
}

module.exports = { enrich }
