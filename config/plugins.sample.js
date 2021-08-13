module.exports = [
  {
    package: "jira-pond",
    name: "jira",
    configuration: {
      collection: {
        fetcher: {
          // ➤➤➤ Replace example host with your own host
          host: "https://your-domain.atlassian.net",
          // ➤➤➤ Replace *** with your jira API token, please see Jira plugin readme for details
          basicAuth: "<your-jira-token>",
          // Set timeout for sending requests to Jira API
          timeout: 10000,
          // Set max retry times for sending requests to Jira API
          maxRetry: 3,

          // Enable proxy for interacting with Jira API
          // proxy: 'http://localhost:4780',
        },
      },
      enrichment: {
        issue: {
          mapping: {
            // This maps issue types in your Jira system to the standard issue type in dev lake
            // In lake, we define bugs as issues found in the development process whereas
            // incidents are issues found in production environment
            // Format: <Standard Type>: [<Jira Type>]
            type: {
              // This mapping powers the metrics like Bug Count, Bug Age, and etc
              // Replace 'Bug' with your own issue types for bugs.
              Bug: ["Bug"],
              // This mapping powers the metrics like Incident Count, Incident Age, and etc
              // Replace 'Incident' with your own issue types for incidents
              Incident: ["Incident"],
            },
          },
          // Enables lake to track which epic an issue belongs to
          // ➤➤➤ Replace 'customfiled_10014' with your own field ID for the epic key
          epicKeyField: "customfield_10014",
        },
      },
    },
  },
  {
    package: "gitlab-pond",
    name: "gitlab",
    configuration: {
      collection: {
        fetcher: {
          // ➤➤➤ Replace example host with your host if your Gitlab is self-hosted
          // Leave this unchanged if you use Gitlab's cloud service
          host: "https://gitlab.com",
          apiPath: "api/v4",
          // ➤➤➤ Replace *** with your Gitlab API token, see Gitlab plugin readme for details
          token: "<your-gitlab-token>",
          // Set timeout for sending requests to Jira API
          timeout: 10000,
          // Set max retry times for sending requests to Jira API
          maxRetry: 3,

          // Enable proxy for interacting with Gitlab API
          // proxy: 'http://localhost:4780',
        },
        // Customize what data you'd like to collect
        // Leave this unchanged if you wish to collect everything
        skip: {
          commits: false,
          projects: false,
          mergeRequests: false,
          notes: false,
        },
      },
    },
  },
  {
    package: "compound-figures",
    name: "compound-figures",
    configuration: {
      enrichment: {
        // ➤➤➤ Replace example mapping your own mapping
        // Format: <Jira boardID>: <Gitlab projectId>
        jiraBoardId2GitlabProjectId: {
          8: 8967944,
        },
      },
    },
  },
];
