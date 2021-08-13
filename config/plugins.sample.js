module.exports = [
  {
    package: 'jira-pond',
    name: 'jira',
    configuration: {
      collection: {
        fetcher: {
          // ➤➤➤ Replace example host with your own host
          host: 'https://your-domain.atlassian.net',
          // ➤➤➤ Replace *** with your jira API token, please see Jira plugin readme for details
          basicAuth: '***',
          // Set timeout for sending requests to Jira API
          timeout: 10000,
          // Set max retry times for sending requests to Jira API
          maxRetry: 3

          // Enable proxy for interacting with Jira API
          // proxy: 'http://localhost:4780',
        }
      },
      enrichment: {
        issue: {
          //  This maps issue types in your Jira system to the standard issue type in dev lake
          typeMappings: [
            // This maps issue types in your Jira system to the standard issue type in dev lake
            // In lake, we define bugs as issues found in the development process whereas
            // incidents are issues found in production environment

            // This mapping powers the metrics like Bug Count, Bug Age, and etc
            // Replace 'Bug' with your own issue types for bugs.
            {
              originTypes: ['Bug'],
              standardType: 'Bug',
              statusMappings: [
                // This mapping powers the metrics like Bug Age
                { originStatuses: ['Rejected', 'Abandoned', 'Cancelled', 'ByDesign', 'Irreproducible'], standardStatus: 'Rejected' },
                { originStatuses: ['Resolved', 'Approved', 'Verified', 'Done', 'Closed'], standardStatus: 'Resolved' }
              ]
            },
            // This mapping powers the metrics like Incident Count, Incident Age, and etc
            // Replace 'Incident' with your own issue types for incidents
            {
              originTypes: ['Incident'],
              standardType: 'Incident',
              statusMappings: [
                // This mapping powers the metrics like Incident Age
                { originStatuses: ['Rejected', 'Abandoned', 'Cancelled', 'Irreproducible'], standardStatus: 'Rejected' },
                { originStatuses: ['Resolved', 'Approved', 'Verified', 'Done', 'Closed'], standardStatus: 'Resolved' }
              ]
            },
            // This mapping powers the metrics like Requirement Count, Requirement Age, and etc
            // Replace 'Story' with your own issue types for requirements
            {
              originTypes: ['Story'],
              standardType: 'Requirement',
              statusMappings: [
                // This mapping powers the metrics like Requirement Lead Time, Requirement Count
                { originStatuses: ['Rejected', 'Abandoned', 'Cancelled'], standardStatus: 'Rejected' },
                { originStatuses: ['Resolved', 'Done', 'Closed'], standardStatus: 'Resolved' }
              ]
            }
          ],
          // Enables lake to track which epic an issue belongs to
          // ➤➤➤ Replace 'customfiled_10014' with your own field ID for the epic key
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
          // ➤➤➤ Replace example host with your host if your Gitlab is self-hosted
          // Leave this unchanged if you use Gitlab's cloud service
          host: 'https://gitlab.com',
          apiPath: 'api/v4',
          // ➤➤➤ Replace *** with your Gitlab API token, see Gitlab plugin readme for details
          token: '***',
          // Set timeout for sending requests to Jira API
          timeout: 10000,
          // Set max retry times for sending requests to Jira API
          maxRetry: 3

          // Enable proxy for interacting with Gitlab API
          // proxy: 'http://localhost:4780',
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
        // ➤➤➤ Replace example mapping your own mapping
        // Format: <Jira boardID>: <Gitlab projectId>
        jiraBoardId2GitlabProjectId: {
          8: 8967944
        }
      }
    }
  }
]
