const commitsCollector = require('../collector/commits')

async function enrich ({ rawDb, enrichedDb, projectId }) {
  if (!projectId) {
    throw new Error('Failed to enrich gitlab commits, projectId is required')
  }

  console.info('INFO >>> gitlab enriching commits for project', projectId)
  await enrichCommitsByProjectId(rawDb, enrichedDb, projectId)
  console.info('INFO >>> gitlab enriching commits for project done!', projectId, counter)
}

async function enrichCommitsByProjectId (rawDb, enrichedDb, projectId) {
  const commitsCollection = await commitsCollector.getCollection(rawDb)
  const cursor = commitsCollection.find({ projectId })
  let counter = 0
  try {
    while (await cursor.hasNext()) {
      const commit = await cursor.next()
<<<<<<< HEAD
      const enriched = {
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
        // some commits can happen with no additions or deletions so stats is not populated
        additions: commit.stats && commit.stats.additions,
        deletions: commit.stats && commit.stats.deletions,
        total: commit.stats && commit.stats.total
      }
=======
      const enriched = mapResponseToSchema(commit)
>>>>>>> cda9245 (chore: set up the tests for kevin to fill in)
      await enrichedDb.GitlabCommit.upsert(enriched)
      counter++
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
