const _has = require('lodash/has')

const jira = require('./jira')

module.exports = {
  async createJobs (project) {
    if (_has(project, 'jira')) {
      await jira.collect(project.jira)
    }
  }
}