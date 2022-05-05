const defaultConfig = [
  [
    {
      Plugin: 'gitlab',
      Options: {
        projectId: 8967944
      }
    },
    {
      Plugin: 'jira',
      Options: {
        boardId: 8,
        connectionId: 1
      }
    },
    {
      Plugin: 'jenkins',
      Options: {}
    },
    {
      Plugin: 'github',
      Options: {
        repo: 'lake',
        owner: 'merico-dev'
      }
    }
  ]
]

export {
  defaultConfig
}
