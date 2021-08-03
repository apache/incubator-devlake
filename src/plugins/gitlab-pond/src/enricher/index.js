const projects = require('./projects')
const commits = require('./commits')
const mergeRequests = require('./merge-requests')
const notes = require('./notes')

async function enrich (rawDb, enrichedDb, { projectId }) {
  const args = { rawDb, enrichedDb, projectId: Number(projectId) }
  await projects.enrich(args)
  await commits.enrich(args)
  await mergeRequests.enrich(args)
  await notes.enrich(args)
}

module.exports = { enrich }
