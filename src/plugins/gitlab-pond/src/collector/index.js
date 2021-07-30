const commits = require('./commits')
const mergeRequests = require('./merge-requests')
const projects = require('./projects')

async function collect (db, { projectId, forceAll }) {
  const args = { db, projectId, forceAll }
  await projects.collect(args)
  await commits.collect(args)
  await mergeRequests.collect(args)
}

module.exports = { collect }
