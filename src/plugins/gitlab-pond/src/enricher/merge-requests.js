const mergeRequestsCollector = require('../collector/merge-requests')

async function enrich ({ rawDb, enrichedDb, projectId }) {
  if (!projectId) {
    throw new Error('Failed to enrich gitlab merge-requests, projectId is required')
  }

  await enrichMergeRequestsByProjectId(rawDb, enrichedDb, projectId)
}

async function enrichMergeRequestsByProjectId (rawDb, enrichedDb, projectId) {
  console.info('INFO >>> gitlab enriching merge-requests for project', projectId)
  const mergeRequestsCollection = await mergeRequestsCollector.getCollection(rawDb)
  const cursor = mergeRequestsCollection.find({ projectId })
  let counter = 0
  try {
    while (await cursor.hasNext()) {
      const mergeRequest = await cursor.next()
      const enriched = {
        projectId: mergeRequest.project_id,
        id: mergeRequest.id,
        numberOfReviewers: mergeRequest.reviewers && mergeRequest.reviewers.length,
        state: mergeRequest.state,
        title: mergeRequest.title,
        webUrl: mergeRequest.web_url,
        userNotesCount: mergeRequest.user_notes_count,
        workInProgress: mergeRequest.work_in_progress,
        sourceBranch: mergeRequest.source_branch,
        mergedAt: mergeRequest.merged_at,
        gitlabCreatedAt: mergeRequest.created_at,
        closedAt: mergeRequest.closed_at,
        mergedByUsername: mergeRequest.merged_by && mergeRequest.merged_by.username,
        description: mergeRequest.description,
        reviewers: mergeRequest.reviewers && mergeRequest.reviewers.map(reviewer => reviewer.username),
        authorUsername: mergeRequest.author && mergeRequest.author.username
      }
      await enrichedDb.GitlabMergeRequest.upsert(enriched)
      counter++
    }
  } finally {
    await cursor.close()
  }
  console.info('INFO >>> gitlab enriching merge-requests for project done!', projectId, counter)
}

module.exports = { enrich }
