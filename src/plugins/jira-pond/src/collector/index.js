const boards = require('./boards')
const issues = require('./issues')

async function collect (db, { boardId, forceAll }) {
  const args = { db, boardId: Number(boardId), forceAll }
  await boards.collect(args)
  await issues.collect(args)
}

module.exports = { collect }
