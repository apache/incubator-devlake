// const collectionManager = require('../src/collector/collection-manager')

// describe.skip('Collection Manager', () => {
//   describe('Groups', () => {
//     describe.skip('fetchCollectionData', () => {
//       it('gets group data', async () => {
//         const testGroup1Id = 12848321
//         // let testGroup2Id = 8378087
//         const res = await collectionManager.fetchCollectionData(groupsModelName, testGroup1Id, '')
//         console.log('res', res)
//       })
//       it('gets project data by group', async () => {
//         // let testGroupId1 = 12848321
//         const testGroupId2 = 12848458
//         const res = await collectionManager.fetchCollectionData(groupsModelName, testGroupId2, groupsUriComponents.projects)
//         console.log('res', res)
//       })
//     })
//   })
//   describe('Commits', () => {
//     describe('fetchCollectionData', () => {
//       const projectId = 20103385
//       it('gets commits data', async () => {
//         const res = await collectionManager.fetchCollectionData(commitsModelName, projectId, mergeRequestsUriComponents.commits)
//         console.log('res', res)
//         console.log('res.length', res.length)
//       })
//     })
//   })
//   describe('Merge Requests', () => {
//     describe('fetchCollectionData', () => {
//       const projectId = 20103385
//       it.only('get merge requests', async () => {
//         const res = await collectionManager.fetchCollectionData(mergeRequestsModelName, projectId, mergeRequestsUriComponents.mergeRequests)
//         console.log('res', res)
//         console.log('res.length', res.length)
//       })
//     })
//   })
// })
