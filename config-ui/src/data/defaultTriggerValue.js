const someJson = [
  [
    {
      plugin: 'jira',
      Options: {
        boardId: 8,
        sourceId: 1
      }
    },
    {
      plugin: 'gitlab',
      Options: {
        projectId: 8967944
      }
    },
    {
      plugin: 'jenkins',
      Options: {}
    },
    {
      Plugin: 'github',
      Options: {
        repositoryName: 'lake',
        owner: 'merico-dev'
      }
    }
  ],
  [
    { plugin: 'jiradomain', options: {} },
    { plugin: 'gitlab-domain', options: {} },
    { plugin: 'jenkinsdomain', options: {} },
    { plugin: 'github-domain', options: {} }
  ]
]

export default someJson
