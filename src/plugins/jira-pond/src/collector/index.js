const issues = require('./issues')

async function collect (db, { boardId, forceAll }) {
  const args = { db, boardId, forceAll }
  await issues.collect(args)
}

module.exports = { collect }
