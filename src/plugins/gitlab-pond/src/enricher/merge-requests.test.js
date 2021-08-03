// const mergeRequests = require("../src/collector/deprecated/merge-requests")
// const mockData = require('./data/merge-requests')
// const dbConnector = require('@mongo/connection')
// const assert = require('assert')

// // SKIP API calls
// describe.skip('MergeRequests', () => {
//   describe('fetchProjectMergeRequests', () => {
//     it('gets mergeRequests for a project', async () => {
//       let projectId = 28270340
//       let foundMergeRequests = await mergeRequests.fetchProjectMergeRequests(projectId)
//       console.log('foundMergeRequests', foundMergeRequests);
//     })
//   })
//   describe('save', () => {
//     it('mergeRequests found in the db have the same length as the mock data provided', async () => {
//       const {
//         db, client
//       } = await dbConnector.connect()
//       const collectionName = 'gitlab_project_merge_requests'
//       try {
//         await dbConnector.clearCollectionData(db, collectionName)
//         await mergeRequests.save({ response: mockData, db})
//         let foundMergeRequests = await mergeRequests.findMergeRequests('', db)
//         assert.equal(foundMergeRequests.length, mockData.length)
//       } catch (error) {
//         console.log('Failed to collect', error)
//       } finally {
//         dbConnector.disconnect(client)
//       }
//     })
//   })
// })
