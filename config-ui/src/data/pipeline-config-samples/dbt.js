const dbtConfig = [
  [
    {
      Plugin: 'dbt',
      Options: {
        projectPath: '/var/www/html/my-project',
        projectName: 'myproject',
        projectTarget: 'dev',
        selectedModels: ['model_one', 'model_two'],
        projectVars: {
          demokey1: 'demovalue1',
          demokey2: 'demovalue2'
        }
      }
    }
  ]
]

export {
  dbtConfig
}
