// const mongo = require('../mongo')
// const groups = require('./groups')
// const standalone = require('./standalone')
const projects = require('./projects')
// const collectionName = 'gitlab_projects'

module.exports = {
  async collect (projectId) {
    console.log('Gitlab collection: projectId', projectId)
    try {
      const projectToSave = {}
      const commits = await projects.fetchProjectRepoCommits(projectId)
      const merge_requests = await projects.fetchMergeRequests(projectId)
      projectToSave.commits = commits
      projectToSave.merge_requests = merge_requests
      // await mongo.storeRawData(projectToSave, collectionName)
    } catch (error) {
      console.error(error)
    }
  }
}
