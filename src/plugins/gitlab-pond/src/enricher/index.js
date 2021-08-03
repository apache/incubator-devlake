const projects = require('./projects')
const commits = require('./commits')
const mergeRequests = require('./merge-requests')
const notes = require('./notes')
const { gitlab } = require('@config/resolveConfig')

async function enrich (rawDb, enrichedDb, { projectId }) {
  const args = { rawDb, enrichedDb, projectId: Number(projectId) }
  !gitlab.skip.projects && await projects.enrich(args)
  !gitlab.skip.commits && await commits.enrich(args)
  !gitlab.skip.mergeRequests && await mergeRequests.enrich(args)
  !gitlab.skip.notes && await notes.enrich(args)
}

module.exports = { enrich }


