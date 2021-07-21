require('module-alias/register')

const fetcher = require('./fetcher')
const projectUri = '/projects'

module.exports = {
  async fetchProject (projectId) {
    const requestUri = `${projectUri}${projectId}`

    return fetcher.fetch(requestUri)
  },
}