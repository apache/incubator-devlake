const projects = require('./projects')
const commits = require('./commits')
const mergeRequests = require('./merge-requests')

async function enrich (rawDb, enrichedDb, { projectId }) {
  const args = { rawDb, enrichedDb, projectId }
  await projects.enrich(args)
  await commits.enrich(args)
  await mergeRequests.enrich(args)
}

module.exports = { enrich }
