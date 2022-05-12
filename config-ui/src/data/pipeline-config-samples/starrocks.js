const starRocksConfig = [
  [
    {
      Plugin: 'starrocks',
      Options: {
        host: '127.0.0.1',
        port: 9030,
        user: 'root',
        password: '',
        database: 'lake',
        be_port: 8040,
        tables: ['_tool_github_commits']
      }
    }
  ]
]

export {
  starRocksConfig
}
