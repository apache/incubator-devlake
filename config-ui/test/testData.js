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
          connectionId: 1,
        },
      },
      {
        Plugin: "jenkins",
        Options: {},
      },
      {
        Plugin: "github",
        Options: {
          repo: "lake",
          owner: "merico-dev",
        },
      },
    ]
  ]
}

module.exports = TEST_DATA
