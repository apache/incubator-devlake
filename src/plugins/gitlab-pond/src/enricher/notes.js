const mongo = require('../util/mongo')

async function enrich ({ rawDb, enrichedDb, projectId }) {
  if (!projectId) {
    throw new Error('Failed to enrich gitlab project, projectId is required')
  }
  
  console.info('INFO >>> gitlab enriching notes for project', projectId)
  await enrichNotesByProjectId(rawDb, enrichedDb, projectId)
  console.info('INFO >>> gitlab enriching notes for project done!', projectId, upsertPromises.length)
}

function findEarliestNote (notes) {
  if (notes && notes.length > 0) {
    const earliestNote = notes.reduce((a, b) => {
      return new Date(a.created_at) < new Date(b.created_at) ? a : b
    })
    return earliestNote
  }
}

// we need a metric that measures a merge request duration as the time from first comment to MR close
async function updateMergeRequestWithFirstCommentTime (notes, mr, enrichedDb) {
  const earliestNote = findEarliestNote(notes)
  if (earliestNote) {
    await enrichedDb.GitlabMergeRequest.update({
      firstCommentTime: earliestNote.created_at
    }, {
      where: {
        id: mr.id
      }
    })
  }
}

/*
  The purpose of this method is to save all the notes from all the merge requests
  into the Postgres db.
  First, we get all MRs from mongo.
  Second, for each MR, we map values from mongo to new values for Postgres.
  Finally, we store GitlabMergeRequestNotes using our PG model.
*/
async function enrichNotesByProjectId (rawDb, enrichedDb, projectId) {
  const mergeRequests = await mongo.findCollection('gitlab_merge_requests',
    { projectId }
    , rawDb)

  const responseNotes = []
  for (const mr of mergeRequests) {
    const mongoNotes = await mongo.findCollection('gitlab_merge_request_notes',
    // { system: false } is necessary to specifically get comments only vs. system notes
      { noteable_id: mr.id, system: false }
      , rawDb)
    responseNotes.push(mongoNotes)
    await updateMergeRequestWithFirstCommentTime(mongoNotes, mr, enrichedDb)
  }
  const mrNotes = responseNotes.flat(1)
  const upsertPromises = []

  mrNotes.forEach(mrNote => {
    const noteToAdd = mapResponseToSchema(mrNote)
    upsertPromises.push(enrichedDb.GitlabMergeRequestNote.upsert(noteToAdd))
  })

  await Promise.all(upsertPromises)
}

function mapResponseToSchema (mrNote) {
  return {
    id: mrNote.id,
    noteableId: mrNote.noteable_id,
    noteableIid: mrNote.noteable_iid,
    authorUsername: mrNote.author && mrNote.author.username,
    body: mrNote.body,
    gitlabCreatedAt: mrNote.created_at,
    noteableType: mrNote.noteable_type,
    confidential: mrNote.confidential
  }
}

module.exports = { enrich, findEarliestNote }
