const projects = require('./projects')
const mergeRequests = require('./merge-requests')
const commits = require('./commits')
const notes = require('./notes')
const { gitlab } = require('@config/resolveConfig')

async function collect (db, { projectId, forceAll }) {
  const args = { db, projectId: Number(projectId), forceAll }
  !gitlab.skip.projects && await projects.collect(args)
  !gitlab.skip.commits && await commits.collect(args)
  !gitlab.skip.mergeRequests && await mergeRequests.collect(args)
  !gitlab.skip.notes && await notes.collect(args)
}

module.exports = { collect }
