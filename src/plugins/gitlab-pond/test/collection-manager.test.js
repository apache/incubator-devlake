const collectionManager = require("../src/collector/collection-manager")
const collections = require('../src/collector/collections')
const { 
  groups: { modelNameForUri: groupsModelName, uriComponents: groupsUriComponents },
  mergeRequests: { modelNameForUri: mergeRequestsModelName, uriComponents: mergeRequestsUriComponents },
  commits: { modelNameForUri: commitsModelName, uriComponents: commitsUriComponents } 
} = collections

describe('Collection Manager', () => {
  describe('Groups', () => {
    describe.skip('fetchCollectionData', () => {
      it('gets group data', async () => {
        let testGroup1Id = 12848321
        // let testGroup2Id = 8378087
        let res = await collectionManager.fetchCollectionData(groupsModelName, testGroup1Id, '')
        console.log('res', res);
      })
      it('gets project data by group', async () => {
        // let testGroupId1 = 12848321
        let testGroupId2 = 12848458
        let res = await collectionManager.fetchCollectionData(groupsModelName, testGroupId2, groupsUriComponents.projects)
        console.log('res', res);
      })
    })
  })
  describe('Commits', () => {
    describe('fetchCollectionData', () => {
      let projectId = 20103385
      it('gets commits data', async () => {
        let res = await collectionManager.fetchCollectionData(commitsModelName, projectId, mergeRequestsUriComponents.commits)
        console.log('res', res);
        console.log('res.length', res.length);
      })
    })
  })
  describe('Merge Requests', () => {
    describe('fetchCollectionData', () => {
      let projectId = 20103385
      it.only('get merge requests', async () => {
        let res = await collectionManager.fetchCollectionData(mergeRequestsModelName, projectId, mergeRequestsUriComponents.mergeRequests)
        console.log('res', res);
        console.log('res.length', res.length);
      })
    })
  })
})