const assert = require('assert')

const { findEarliestNote } = require('./notes')

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
      const earlistNote = findEarliestNote(mongoNotes)
      assert.deepStrictEqual(earlistNote.title, 'b')
    })
  })
  describe('mapResponseToSchema', () => {
    it('maps the right object', () => {

    })
  })
})
