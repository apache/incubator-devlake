require('module-alias/register')
const config = require('@config/resolveConfig').jiraBoardGitlabProject

module.exports = {
  configuration: {
    // default configuration which could be overrided by `config/plugins.js`
  },

  async initialize (rawDb, enrichedDb, plugins) {
    if (!config) {
      return
    }
    const { JiraBoardGitlabProject } = enrichedDb
    // remove all board-to-repo mapping
    await JiraBoardGitlabProject.destroy({ where: {}, truncate: true })
    // sync from configuration to database for JOIN query
    for (const [boardId, projectId] of Object.entries(config)) {
      await JiraBoardGitlabProject.create({ boardId, projectId })
    }
  },
  enricher: {
    name: 'compoundFiguresEnricher',
    exec: async function (rawDb, enrichedDb, options) {
      // await enrichment.enrich(rawDb, enrichedDb, options)
      return []
    }
  }
}

if (require.main === module) {
  const dbConnector = require('@mongo/connection');
  const enrichedDb = require('@db/postgres');

  (async function () {
    const { db, client } = await dbConnector.connect()
    try {
      await module.exports.initialize(db, enrichedDb, {})
    } finally {
      dbConnector.disconnect(client)
    }
  })()
}
