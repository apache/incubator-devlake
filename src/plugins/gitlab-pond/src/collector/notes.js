require('module-alias/register')
const { findOrCreateCollection } = require('../../../commondb')
const mongo = require('../util/mongo')
const fetcher = require('./fetcher')

const collectionName = 'gitlab_merge_request_notes'

module.exports = {
  async collect ({ db, projectId, forceAll }) {
    if (!projectId) {
      throw new Error('Failed to collect GitLab data, projectId is required')
    }

    console.info('INFO >>> GitLab collecting notes for project. projectId: ', projectId)
    const mergeRequests = await mongo.findCollection('gitlab_merge_requests', { projectId }, db)
    console.log(`INFO >>> GitLab collecting notes, found ${mergeRequests.length} MRs for project. projectId: ${projectId}`)
    for (const mr of mergeRequests) {
      await module.exports.collectByMergeRequestIid(db, projectId, mr.iid, forceAll)
    }
    console.info('INFO >>> GitLab collecting notes for project done! projectId: ', projectId)
  },

  async collectByMergeRequestIid (db, projectId, mergeRequestIid, forceAll) {
    const notesCollection = await findOrCreateCollection(db, collectionName)
    const response = await fetcher.fetch(`projects/${projectId}/merge_requests/${mergeRequestIid}/notes`)
    const notes = response.data

    for (const note of notes) {
      await notesCollection.findOneAndUpdate(
        { id: note.id },
        { $set: note },
        { upsert: true }
      )
    }
  }
}
