const collectionManager = require("../src/collector/collection-manager")
const collections = require('../src/collector/collections')
const { groups: { modelName, uriComponents } } = collections

describe('Collection Manager', () => {
  describe('Groups', () => {
    describe.skip('fetchCollectionData', () => {
      it('gets group data', async () => {
        let testGroup1Id = 12848321
        // let testGroup2Id = 8378087
        let res = await collectionManager.fetchCollectionData(modelName, testGroup1Id, '')
        console.log('res', res);
      })
      it.only('gets project data by group', async () => {
        // let testGroupId1 = 12848321
        let testGroupId2 = 12848458
        let res = await collectionManager.fetchCollectionData(modelName, testGroupId2, uriComponents.projects)
        console.log('res', res);
      })
    })
  })
})