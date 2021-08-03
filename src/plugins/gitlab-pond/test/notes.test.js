const assert = require('assert')

const { findEarliestNote } = require('../src/enricher/notes')

describe('Notes', () => {
  describe('', () => {
    it('findEarliestNote', () => {
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
      const earlistNote = findEarliestNote(mongoNotes)
      assert.deepStrictEqual(earlistNote.title, 'b')
    })
  })
})
