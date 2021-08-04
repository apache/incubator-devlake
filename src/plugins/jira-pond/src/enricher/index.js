const boards = require('./boards')
const issues = require('./issues')

async function enrich (rawDb, enrichedDb, { boardId, forceAll }) {
  const args = { rawDb, enrichedDb, boardId: Number(boardId), forceAll }
  await boards.enrich(args)
  await issues.enrich(args)
}

module.exports = { enrich }
