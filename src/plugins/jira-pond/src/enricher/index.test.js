const assert = require('assert')
const enricher = require('./index')

describe('Jira-pond enricher', () => {
  describe('mapIssueTypeFromConfiguration()', () => {
    it('returns the lake value for Bug from the jira issue of the user', () => {
      const issueTypes = {
        "Bug": "Bugzilla",
        "Incident": "OurJiraIncidentType"
      }
      
      assert(enricher.mapIssueTypeFromConfiguration('Bugzilla', issueTypes), 'Bug')
      assert(enricher.mapIssueTypeFromConfiguration('OurJiraIncidentType', issueTypes), 'Incident')
    })
    it('handles case insensitivity', () => {
      const issueTypes = {
        "Bug": "Bugzilla",
        "Incident": "OurJiraIncidentType"
      }
      
      assert(enricher.mapIssueTypeFromConfiguration('bugzilla', issueTypes), 'Bug')
      assert(enricher.mapIssueTypeFromConfiguration('ourjiraINCIDENTType', issueTypes), 'Incident')
    })
    it('returns the passed in issue type if no mapping exists', () => {
      const issueTypes = {
        "Bug": "Bugzilla",
        "Incident": "OurJiraIncidentType"
      }
      
      assert(enricher.mapIssueTypeFromConfiguration('myOtherJiraIssueType', issueTypes), 'myOtherJiraIssueType')
    })

    it('when no issue type is passed in, return empty string', () => {
      assert.deepStrictEqual(enricher.mapIssueTypeFromConfiguration(), '')
      assert.deepStrictEqual(enricher.mapIssueTypeFromConfiguration(''), '')
    })
  })
})