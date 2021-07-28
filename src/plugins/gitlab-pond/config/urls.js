const baseUrl = 'https://gitlab.com/api/v3/'

module.exports = {
  groups_path: {
    url: `https://gitlab.com/api/v3/proj/:groupId/mr/:projectId/notes`,
    collection: 'merge_requests'
  },
  mr_path: {
    url: `https://gitlab.com/api/v3/proj/:projectId/mr/:mrId/notes`,
    collection: 'merge_requests'
  },
  proj_path: {
    url: `https://gitlab.com/api/v3/proj/:projectId`,
    collection: 'projects'
  }
}
