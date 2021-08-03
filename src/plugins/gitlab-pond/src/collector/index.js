const projects = require('./projects')
const mergeRequests = require('./merge-requests')
const commits = require('./commits')
const notes = require('./notes')
const { maybeSkip } = require('../util/async')

async function collect (db, { projectId, forceAll }) {
  const args = { db, projectId: Number(projectId), forceAll }
  await maybeSkip(projects.collect(args), 'projects')
  await maybeSkip(commits.collect(args), 'commits')
  await maybeSkip(mergeRequests.collect(args), 'mergeRequests')
  await maybeSkip(notes.collect(args), 'notes')
}

module.exports = { collect }
