// // The purpose of this file is to make it easy for the collection-manager
// // to construct the proper API requests to gitlab api.

// // For example:

// // ${modelNameForUri}/${projectId}/${uriComponents.commits}

// // Developer will still need some knowledge of the api to use it properly.

// const collections = {
//   commits: {
//     collectionName: "gitlab_project_repo_commits",
//     modelNameForUri: "projects",
//     uriComponents: {
//       commits: 'repository/commits'
//     }
//   },
//   mergeRequests: {
//     collectionName: "gitlab_project_merge_requests",
//     modelNameForUri: "projects",
//     uriComponents: {
//       mergeRequests: 'merge_requests'
//     }
//   },
//   groups: {
//     collectionName: "gitlab_groups",
//     modelNameForUri: "groups",
//     uriComponents: {
//       projects: "projects",
//     },
//   },
//   notes: {
//     collectionName: 'gitlab_notes',
//     modelNameForUri: 'projects',
//     uriComponents: {
//       onMergeRequests: 'merge_requests/:merge_request_iid/notes'
//     }
//   }
// }

// module.exports = collections;
