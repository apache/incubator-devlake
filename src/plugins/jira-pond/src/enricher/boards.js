const boardsCollector = require('../collector/boards')

async function enrich ({ rawDb, enrichedDb, boardId }) {
  if (!boardId) {
    throw new Error('Failed to enrich jira board, boardId is required')
  }

  await enrichBoardById(rawDb, enrichedDb, boardId)
}

async function enrichBoardById (rawDb, enrichedDb, boardId) {
  console.info('INFO >>> jira enriching board', boardId)
  const boardsCollection = await boardsCollector.getCollection(rawDb)
  const board = await boardsCollection.findOne({ id: boardId })
  const enriched = {
    id: board.id,
    projectId: board.location.projectId,
    name: board.name,
    type: board.type,
    webUrl: board.self
  }
  await enrichedDb.JiraBoard.upsert(enriched)
  console.info('INFO >>> jira enriching board done!', boardId)
}

module.exports = { enrich }
