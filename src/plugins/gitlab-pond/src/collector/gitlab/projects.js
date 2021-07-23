require('module-alias/register')

const fetcher = require('../util/fetcher')
const modelName = 'projects'
const { gitlab: { host, apiPath, token } } = require('../../../../../../config/resolveConfig')
const privateTokenHeader = {"PRIVATE-TOKEN": token}

module.exports = {
  async fetchProject (projectId) {
    return fetcher.fetchOne(host, apiPath, modelName, projectId)
  },
  async fetchAllProjects () {
    return fetcher.fetchAll(host, apiPath, modelName)
  },
  async fetchProjectRepoCommits (projectId) {
    const routeName = 'repository/commits?all=true&with_stats=true'
    return fetcher.fetch(`${host}/${apiPath}/${modelName}/${projectId}/${routeName}`, privateTokenHeader)
  },
  async fetchProjectRepoTree (projectId) {
    const routeName = 'repository/tree'
    return fetcher.fetch(`${host}/${apiPath}/${modelName}/${projectId}/${routeName}`)
  },
  async fetchProjectFiles (projectId, tree) {
    const routeName = 'repository/files'
    const files = []
    for (const treeNode of tree) {
      const path = treeNode.path.replace('.', '%2E').replace('/', '%2F')
      const url = `${host}/${apiPath}/${modelName}/${projectId}/${routeName}/${path}/raw`
      const file = await fetcher.fetch(url)
      files.push(file)
    }
    return files
  },
  async fetchMergeRequests (projectId) {
    const routeName = 'merge_requests'
    const url = `${host}/${apiPath}/${modelName}/${projectId}/${routeName}`
    return fetcher.fetch(url, privateTokenHeader) 
  }
}
