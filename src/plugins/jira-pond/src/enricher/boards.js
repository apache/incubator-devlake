const boardsCollector = require('../collector/boards')

async function enrich ({ rawDb, enrichedDb, boardId }) {
  if (!boardId) {
    throw new Error('Failed to enrich jira board, boardId is required')
  }

  console.info('INFO >>> jira enriching board', boardId)
  await enrichBoardById(rawDb, enrichedDb, boardId)
  console.info('INFO >>> jira enriching board done!', boardId)
}

async function enrichBoardById (rawDb, enrichedDb, boardId) {
  const boardsCollection = await boardsCollector.getCollection(rawDb)
  const board = await boardsCollection.findOne({ id: boardId })
  const enriched = mapResponseToSchema(board)
  await enrichedDb.JiraBoard.upsert(enriched)
}

function mapResponseToSchema (board) {
  return {
    id: board.id,
    projectId: board.location.projectId,
    name: board.name,
    type: board.type,
    webUrl: board.self
  }
}

module.exports = { enrich }
