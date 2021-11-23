const TEST_DATA = {
  gitlabTriggersJson: {
    Plugin: "gitlab",
    Options: {
      projectId: 8967944,
    },
  },
  gitlabDomainTriggersJson: {
    Plugin: "gitlab-domain",
    Options: {}
  },
  completeTriggersJson: [
    [
      {
        Plugin: "gitlab",
        Options: {
          projectId: 8967944,
        },
      },
      {
        Plugin: "jira",
        Options: {
          boardId: 8,
          sourceId: 1,
        },
      },
      {
        Plugin: "jenkins",
        Options: {},
      },
      {
        Plugin: "github",
        Options: {
          repositoryName: "lake",
          owner: "merico-dev",
        },
      },
    ],
    [
      {
        Options: {},
        Plugin: "gitlab-domain",
      },
      {
        Options: {},
        Plugin: "jiradomain",
      },
      {
        Options: {},
        Plugin: "jenkinsdomain",
      },
      {
        Options: {},
        Plugin: "github-domain",
      },
    ],
  ]
}

module.exports = TEST_DATA