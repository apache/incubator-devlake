const projects = require('./projects')
const mergeRequests = require('./merge-requests')
const commits = require('./commits')
const notes = require('./notes')
const { gitlab } = require('@config/resolveConfig')

async function collect (db, { projectId, branch, forceAll }) {
  const args = { db, projectId: Number(projectId), branch, forceAll }
  const skipFlags = gitlab.skip
  if (skipFlags) {
    !skipFlags.projects && await projects.collect(args)
    !skipFlags.commits && await commits.collect(args)
    !skipFlags.mergeRequests && await mergeRequests.collect(args)
    !skipFlags.notes && await notes.collect(args)
  } else {
    await projects.collect(args)
    await commits.collect(args)
    await mergeRequests.collect(args)
    await notes.collect(args)
  }
}

module.exports = { collect }
