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
      projectToSave.numberOfCommits = commits.length
      // await mongo.storeRawData(projectToSave, collectionName)
    } catch (error) {
      console.error(error)
    }
  }
}
