const boards = require('./boards')
const issues = require('./issues')
const fetcher = require('./fetcher')
const { merge } = require('lodash')

const configuration = {
  verified: false,
  fetcher: null
}

async function configure (config) {
  await fetcher.configure(config.fetcher)
  merge(configuration, config)
  configuration.verified = true
}

async function collect (db, { boardId, forceAll }) {
  const args = { db, boardId: Number(boardId), forceAll }
  await boards.collect(args)
  await issues.collect(args)
}

module.exports = { configure, collect }
