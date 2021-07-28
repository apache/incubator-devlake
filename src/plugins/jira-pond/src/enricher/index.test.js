const assert = require('assert')
const enricher = require('./index')

describe('Jira-pond enricher', () => {
  describe('mapValue()', () => {
    it('returns the lake value for Bug from the jira issue of the user', () => {
      const config = {
        "Bug": "Bugzilla",
        "Incident": "OurJiraIncidentType"
      }
      
      assert(enricher.mapValue('Bugzilla', config), 'Bug')
      assert(enricher.mapValue('OurJiraIncidentType', config), 'Incident')
    })
    it('handles case insensitivity', () => {
      const config = {
        "Bug": "Bugzilla",
        "Incident": "OurJiraIncidentType"
      }
      
      assert(enricher.mapValue('bugzilla', config), 'Bug')
      assert(enricher.mapValue('ourjiraINCIDENTType', config), 'Incident')
    })
    it('returns the passed in issue type if no mapping exists', () => {
      const config = {
        "Bug": "Bugzilla",
        "Incident": "OurJiraIncidentType"
      }
      
      assert(enricher.mapValue('myOtherJiraIssueType', config), 'myOtherJiraIssueType')
    })

    it('when no issue type is passed in, return empty string', () => {
      assert.deepStrictEqual(enricher.mapValue(), '')
      assert.deepStrictEqual(enricher.mapValue(''), '')
    })

    it('handles an array of values', () => {

      const config = {
        "Done": ["Closed", "Complete"]
      }

      assert.deepStrictEqual(enricher.mapValue('Complete', config), 'Done')
      assert.deepStrictEqual(enricher.mapValue('Closed', config), 'Done')
    })
  })
})