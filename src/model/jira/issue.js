const Schema = require('../core/schema')

class JiraIssue extends Schema {
  get valiator () {
    return {
      primary: ['id', 'key'],
      properties: {
        id: {
          type: 'int',
          description: 'jira issue id',
          required: true
        },
        key: {
          type: 'string',
          description: 'jira issue key',
          required: true
        },
        fields: {
          type: 'object',
          description: 'jira issue fields',
          properties: {
            created: {
              type: 'date',
              description: 'jira created date'
            },
            updated: {
              type: 'date',
              description: 'jira updated date'
            },
            status: {
              type: 'object'
            },
            issuetype: {
              type: 'object'
            },
            creator: {
              type: 'object'
            },
            assignee: {
              type: 'object'
            }
          }
        }
      }
    }
  }
}

module.exports = JiraIssue
