require('module-alias/register')

const fetcher = require('../util/fetcher')
const modelName = 'projects'
const { gitlab: { host, apiPath, token } } = require('../../../../config/resolveConfig')

module.exports = {
  async fetchProject (projectId) {
    return fetcher.fetchOne(host, apiPath, modelName, projectId)
  },
  async fetchAllProjects () {
    return fetcher.fetchAll(host, apiPath, modelName)
  },
  async fetchProjectRepoCommits (projectId) {
    const routeName = 'repository/commits'
    return fetcher.fetch(`${host}/${apiPath}/${modelName}/${projectId}/${routeName}`)
  },
  async fetchProjectRepoTree (projectId) {
    const routeName = 'repository/tree'
    return fetcher.fetch(`${host}/${apiPath}/${modelName}/${projectId}/${routeName}`)
  },
  async fetchProjectFiles (projectId, tree, defaultBranch) {
    const routeName = 'repository/files'
    let files = []
    for(let treeNode of tree) {
      let path = treeNode.path.replace('.', '%2E').replace('/', '%2F')
      console.log('path', path);
      let url = `${host}/${apiPath}/${modelName}/${projectId}/${routeName}/${path}/raw`
      console.log('url', url);
      let file = await fetcher.fetch(url)
      files.push(file)
    }
    return files
  }
}