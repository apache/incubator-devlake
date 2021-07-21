const mongo = require('./mongo')
const projects = require('./projects')
const projectId = '/28270340'

module.exports = {
  async collect () {
    try {
      const project = await projects.collectOne(projectId)
      await mongo.storeRawData(project)
    } catch (error) {
      console.error(error)
    }
  }
}
