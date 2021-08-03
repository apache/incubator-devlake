const mongo = require('../util/mongo')

async function enrich ({ rawDb, enrichedDb, projectId }) {
  if (!projectId) {
    throw new Error('Failed to enrich gitlab project, projectId is required')
  }

  await enrichNotesByProjectId(rawDb, enrichedDb, projectId)
}
/*
  The purpose of this method is to save all the notes from all the merge requests
  into the Postgres db.
  First, we get all MRs from mongo.
  Second, for each MR, we map values from mongo to new values for Postgres.
  Finally, we store GitlabMergeRequestNotes using our PG model.
*/
async function enrichNotesByProjectId (rawDb, enrichedDb, projectId) {
  const {
    GitlabMergeRequestNote
  } = enrichedDb

  const mergeRequests = await mongo.findCollection('gitlab_merge_requests',
    { projectId }
    , rawDb)

  const responseNotes = []
  for (const mr of mergeRequests) {
    const res = await mongo.findCollection('gitlab_merge_request_notes',
    // { system: false } is necessary to specifically get comments only vs. system notes
      { noteable_id: mr.id, system: false }
      , rawDb)
    responseNotes.push(res)
  }
  const mrNotes = responseNotes.flat(1)
  const upsertPromises = []

  mrNotes.forEach(mrNote => {
    const noteToAdd = {
      id: mrNote.id,
      noteableId: mrNote.noteable_id,
      noteableIid: mrNote.noteable_iid,
      authorUsername: mrNote.author && mrNote.author.username,
      body: mrNote.body,
      gitlabCreatedAt: mrNote.created_at,
      noteableType: mrNote.noteable_type,
      confidential: mrNote.confidential
    }
    upsertPromises.push(GitlabMergeRequestNote.upsert(noteToAdd))
  })

  await Promise.all(upsertPromises)
}

module.exports = { enrich }
