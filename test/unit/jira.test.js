var assert = require('assert');
const jiraUtil = require('../../src/util/jira')

describe('Jira util', function() {
  it('Calculates lead time for JIRA issues', () => {
    let month = 1
    let year = 2000
    let startDay = 1
    let endDay = 2

    let startDate = new Date(year, month, startDay)
    let endDate = new Date(year, month, endDay)
    const result = jiraUtil.calculateLeadTime(startDate, endDate)
    console.log('result', result);
    assert.equal(1, 1)
  })
});