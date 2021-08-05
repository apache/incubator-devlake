const {mapResponseToSchema} = require("./commits")
const mockRawCommits = require('./mockData/mockRawCommits')
const assert = require('assert')
const stringUtil = require('../../../../util/string')
const db = require('@db/postgres')

const tableName = 'gitlab_commits'

describe('Commits', () => {
  describe('mapResponseToSchema', () => {
    let responseCommit
    let enrichedCommit
    let commitTable
    beforeEach(async () => {
      let queryInterface = db.sequelize.getQueryInterface()
      commitTable = await queryInterface.describeTable(tableName)
      responseCommit = mockRawCommits[0]
      // We always append the project id before storing the MR
      responseCommit.projectId = 1
      enrichedCommit = mapResponseToSchema(responseCommit)
    })
    it('returns an object with the same number of keys as the Postgres table' + 
      'not including created_at and updated_at, which are not relevant here', () => {
      // These two fields are meta fields
      delete commitTable.created_at
      delete commitTable.updated_at
      assert.strictEqual(Object.keys(enrichedCommit).length, Object.keys(commitTable).length)
    })
    it('returns a response object that contains the same keys as the Postgres table (converted to snake case)', () => {
      for(let key in enrichedCommit){
        let snakeCaseKey = stringUtil.convertCamelToSnakeCase(key)
        assert.strictEqual(commitTable.hasOwnProperty(snakeCaseKey), true)
      }
    })
    it('returns an object which does not contain undefined properties', () => {
      for(let key in enrichedCommit){
        assert.notStrictEqual(enrichedCommit[key], undefined)
      }
    })
  })
})
