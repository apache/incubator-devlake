const commitsCollector = require('../collector/commits')

async function enrich ({ rawDb, enrichedDb, projectId }) {
  if (!projectId) {
    throw new Error('Failed to enrich gitlab commits, projectId is required')
  }

  console.info('INFO >>> gitlab enriching commits for project', projectId)
  await enrichCommitsByProjectId(rawDb, enrichedDb, projectId)
  console.info('INFO >>> gitlab enriching commits for project done!', projectId)
}

async function enrichCommitsByProjectId (rawDb, enrichedDb, projectId) {
  const commitsCollection = await commitsCollector.getCollection(rawDb)
  const cursor = commitsCollection.find({ projectId })
  try {
    while (await cursor.hasNext()) {
      const commit = await cursor.next()
      const enriched = mapResponseToSchema(commit)
      await enrichedDb.GitlabCommit.upsert(enriched)
    }
  } finally {
    await cursor.close()
  }
}

function mapResponseToSchema (commit) {
  return {
    projectId: commit.projectId,
    id: commit.id,
    shortId: commit.short_id,
    title: commit.title,
    message: commit.message,
    authorName: commit.author_name,
    authorEmail: commit.author_email,
    authoredDate: commit.authored_date,
    committerName: commit.committer_name,
    committerEmail: commit.committer_email,
    committedDate: commit.committed_date,
    webUrl: commit.web_url,
    additions: commit.stats && commit.stats.additions,
    deletions: commit.stats && commit.stats.deletions,
    total: commit.stats && commit.stats.total
  }
}

module.exports = { enrich, mapResponseToSchema }
