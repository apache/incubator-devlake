const fetcher = require('./fetcher')
const projects = require('./projects')
const mergeRequests = require('./merge-requests')
const commits = require('./commits')
const notes = require('./notes')
const { merge } = require('lodash')

const configuration = {
  verified: false,
  fetcher: null,
  skip: {
    commits: false,
    projects: false,
    mergeRequests: false,
    notes: false
  }
}

function configure (config) {
  fetcher.configure(config.fetcher)
  merge(configuration.skip, config.skip)
  configuration.verified = true
}

async function collect (db, { projectId, branch, forceAll }) {
  const args = { db, projectId: Number(projectId), branch, forceAll }
  const skipFlags = configuration.skip
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

module.exports = { configure, collect }
