const projects = require('./projects')
const mergeRequests = require('./merge-requests')
const commits = require('./commits')
const notes = require('./notes')

async function collect (db, { projectId, forceAll }) {
  const args = { db, projectId: Number(projectId), forceAll }
  await projects.collect(args)
  await commits.collect(args)
  await mergeRequests.collect(args)
  await notes.collect(args)
}

module.exports = { collect }
