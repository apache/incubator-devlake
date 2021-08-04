const assert = require('assert')
const { mapValue } = require('./mapping')

describe('Mapping', () => {
  describe('mapValue()', () => {
    it('returns the lake value for Bug from the jira issue of the user', () => {
      const config = {
        Bug: 'Bugzilla',
        Incident: 'OurJiraIncidentType'
      }

      assert(mapValue('Bugzilla', config), 'Bug')
      assert(mapValue('OurJiraIncidentType', config), 'Incident')
    })
    it('handles case insensitivity', () => {
      const config = {
        Bug: 'Bugzilla',
        Incident: 'OurJiraIncidentType'
      }

      assert(mapValue('bugzilla', config), 'Bug')
      assert(mapValue('ourjiraINCIDENTType', config), 'Incident')
    })
    it('returns the passed in issue type if no mapping exists', () => {
      const config = {
        Bug: 'Bugzilla',
        Incident: 'OurJiraIncidentType'
      }

      assert(mapValue('myOtherJiraIssueType', config), 'myOtherJiraIssueType')
    })

    it('when no issue type is passed in, return empty string', () => {
      assert.deepStrictEqual(mapValue(), '')
      assert.deepStrictEqual(mapValue(''), '')
    })

    it('handles an array of values', () => {
      const config = {
        Done: ['Closed', 'Complete']
      }

      assert.deepStrictEqual(mapValue('Complete', config), 'Done')
      assert.deepStrictEqual(mapValue('Closed', config), 'Done')
    })
  })
})
