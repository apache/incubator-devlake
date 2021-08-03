const projects = require('./projects')
const commits = require('./commits')
const mergeRequests = require('./merge-requests')
const notes = require('./notes')
const { maybeSkip } = require('../util/async')

async function enrich (rawDb, enrichedDb, { projectId }) {
  const args = { rawDb, enrichedDb, projectId: Number(projectId) }
  await maybeSkip(projects.enrich(args), 'projects')
  await maybeSkip(commits.enrich(args), 'commits')
  await maybeSkip(mergeRequests.enrich(args), 'mergeRequests')
  await maybeSkip(notes.enrich(args), 'notes')
}

module.exports = { enrich }
