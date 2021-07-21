require('module-alias/register')

const fetcher = require('../util/fetcher')
const modelName = 'projects'
const host = 'https://gitlab.com'
const path = 'api/v4'

module.exports = {
  async fetchProject (projectId) {
    return fetcher.fetchOne(host, path, modelName, projectId)
  },
  async fetchAllProjects () {
    return fetcher.fetchAll(host, path, modelName)
  }
}