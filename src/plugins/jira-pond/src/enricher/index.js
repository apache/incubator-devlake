const issues = require('./issues')

async function enrich (rawDb, enrichedDb, { forceAll }) {
  const args = { rawDb, enrichedDb, forceAll }
  await issues.enrich(args)
}

module.exports = { enrich }
