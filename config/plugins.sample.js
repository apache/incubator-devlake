module.exports = [
  {
    package: 'jira-pond',
    name: 'jira',
    configuration: {
      collection: {
        fetcher: {
          // Replace example host with your own host
          host: 'https://your-domain.atlassian.net',
          // Replace *** with your jira API token, please see Jira plugin readme for details
          basicAuth: '***',
          // Enable proxy for interacting with Jira API
          // proxy: 'http://localhost:4780',
          // Set timeout for sending requests to Jira API
          timeout: 10000,
          // Set max retry times for sending requests to Jira API
          maxRetry: 3
        }
      },
      enrichment: {
        issue: {
          mapping: {
            status: {
            // Format: <Standard Status>: <Jira Status>
              Closed: ['Done', 'Closed']
            },
            type: {
            // Format: <Standard Type>: <Jira Type>
              Bug: ['Bug'],
              Incident: ['Incident']
            }
          },
          epicKeyField: 'customfield_10014'
        }
      }
    }
  },
  {
    package: 'gitlab-pond',
    name: 'gitlab',
    configuration: {
      collection: {
        fetcher: {
          // Replace example host with your host if your Gitlab is self-hosted
          // Leave this unchanged if you use Gitlab's cloud service
          host: 'https://gitlab.com',
          apiPath: 'api/v4',
          // Replace *** with your Gitlab API token, please see Gitlab plugin readme for details
          token: '***',
          // Enable proxy for interacting with Gitlab API
          // proxy: 'http://localhost:4780',
          // Set timeout for sending requests to Jira API
          timeout: 10000,
          // Set max retry times for sending requests to Jira API
          maxRetry: 3
        },
        // Customize what data you'd like to collect
        // Leave this unchanged if you wish to collect everything
        skip: {
          commits: false,
          projects: false,
          mergeRequests: false,
          notes: false
        }
      }
    }
  },
  {
    package: 'compound-figures',
    name: 'compound-figures',
    configuration: {
      enrichment: {
        // Replace example mapping your own mapping
        // Format: <Jira boardID>: <Gitlab projectId>
        jiraBoardId2GitlabProjectId: {
          8: 8967944
        }
      }
    }
  }
]
