require('module-alias/register')

const {mapResponseToSchema} = require("./merge-requests")
const mockRawMergeRequests = require('./mockData/mockRawMergeRequests')
const assert = require('assert')
const db = require('@db/postgres')
const tableName = 'gitlab_merge_requests'
const stringUtil = require('../../../../util/string')

describe('Merge Requests', () => {
  describe('mapResponseToSchema', () => {
    let enrichedMergeRequest
    let responseMergeRequest
    let mrTable
    beforeEach(async () => {
      let queryInterface = db.sequelize.getQueryInterface()
      mrTable = await queryInterface.describeTable(tableName)
      responseMergeRequest = mockRawMergeRequests[0]
      // We always append the project id before storing the MR
      responseMergeRequest.projectId = 1
      enrichedMergeRequest = mapResponseToSchema(responseMergeRequest)
    })
    it('returns an object with the same number of keys as the Postgres table' + 
      'not including created_at, updated_at, or first_comment_time, which are not relevant here', () => {
      // These two fields are meta fields
      delete mrTable.created_at
      delete mrTable.updated_at
      // This field is enriched when notes are collected, so they are not relevant here
      delete mrTable.first_comment_time
      assert.strictEqual(Object.keys(enrichedMergeRequest).length, Object.keys(mrTable).length)
    })
    it('returns a response object that contains the same keys as the Postgres table (converted to snake case)', () => {
      for(let key in enrichedMergeRequest){
        let snakeCaseKey = stringUtil.convertCamelToSnakeCase(key)
        assert.strictEqual(mrTable.hasOwnProperty(snakeCaseKey), true)
      }
    })
    it('returns an object which does not contain undefined properties', () => {
      for(let key in enrichedMergeRequest){
        assert.notStrictEqual(enrichedMergeRequest[key], undefined)
      }
    })
  })
})