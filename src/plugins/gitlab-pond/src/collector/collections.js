const collections = {
  commits: {
    collectionName: "gitlab_project_repo_commits",
    modelName: "commits",
  },
  mergeRequests: {
    collectionName: "gitlab_project_merge_requests",
    modelName: "merge_requests",
  },
  groups: {
    collectionName: "gitlab_groups",
    modelName: "groups",
    uriComponents: {
      projects: "projects",
    },
  },
}

module.exports = collections;
