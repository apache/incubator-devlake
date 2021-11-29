const TEST_DATA = {
  gitlabTriggersJson: {
    Plugin: "gitlab",
    Options: {
      projectId: 8967944,
    },
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
    ]
  ]
}

module.exports = TEST_DATA