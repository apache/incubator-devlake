const boards = require('./boards')
const issues = require('./issues')
const boardsCollector = require('../collector/boards')

async function enrich (rawDb, enrichedDb, { boardId, forceAll }) {
  // verify collected data existence
  const boardsCollection = await boardsCollector.getCollection(rawDb)
  const board = await boardsCollection.findOne({ id: boardId })
  if (!board) {
    throw new Error(`jiraEnricher: unable to find collected data for board ${boardId}`)
  }

  const args = { rawDb, enrichedDb, boardId: Number(boardId), forceAll }
  await boards.enrich(args)
  await issues.enrich(args)
}

module.exports = { enrich }
