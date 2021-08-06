require('module-alias/register')

const assert = require('assert')
const db = require('@db/postgres')
const tableName = 'gitlab_merge_request_notes'
const stringUtil = require('../../../../util/string')
const { findEarliestNote, mapResponseToSchema } = require('./notes')
const mockRawMergeRequestNotes = require('./mockData/mockRawMergeRequestNotes')

describe('Notes', () => {
  describe('findEarliestNote', () => {
    it('finds the earliest note in a collection of notes', () => {
      const mongoNotes = [
        {
          title: 'a',
          created_at: '2021-02-17T06:21:37.665Z'
        },
        {
          title: 'b',
          created_at: '2010-02-17T06:21:37.665Z'
        },
        {
          title: 'c',
          created_at: '2019-02-17T06:21:37.665Z'
        },
        {
          title: 'd',
          created_at: '2018-02-17T06:21:37.665Z'
        }
      ]
      const earliestNote = findEarliestNote(mongoNotes)
      assert.deepStrictEqual(earliestNote.title, 'b')
    })
  })
  describe('mapResponseToSchema', () => {
    let mrNotesTable
    let responseMrNote
    let enrichedMrNote
    beforeEach(async () => {
      let queryInterface = db.sequelize.getQueryInterface()
      mrNotesTable = await queryInterface.describeTable(tableName)
      responseMrNote = mockRawMergeRequestNotes[0]
      // We always append the project id before storing the MR
      responseMrNote.projectId = 1
      enrichedMrNote = mapResponseToSchema(responseMrNote)
    })
    it('returns an object with the same number of keys as the Postgres table' + 
      'not including created_at and updated_at, which are not relevant here', () => {
      // These two fields are meta fields
      delete mrNotesTable.created_at
      delete mrNotesTable.updated_at
      assert.strictEqual(Object.keys(enrichedMrNote).length, Object.keys(mrNotesTable).length)
    })
    it('returns a response object that contains the same keys as the Postgres table (converted to snake case)', () => {
      for(let key in enrichedMrNote){
        let snakeCaseKey = stringUtil.convertCamelToSnakeCase(key)
        assert.strictEqual(mrNotesTable.hasOwnProperty(snakeCaseKey), true)
      }
    })
    it('returns an object which does not contain undefined properties', () => {
      for(let key in enrichedMrNote){
        assert.notStrictEqual(enrichedMrNote[key], undefined)
      }
    })
  })
})
