const mongo = require('../mongo')
const projects = require('./projects')
const projectId = '/28270340'
const collectionName = 'gitlab_projects'

module.exports = {
  async collect () {
    try {
      const project = await projects.collectOne(projectId)
      await mongo.storeRawData(project, collectionName)
    } catch (error) {
      console.error(error)
    }
  }
}
