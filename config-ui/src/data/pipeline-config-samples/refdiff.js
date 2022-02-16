const refdiffConfig = [
  [
    {
      Plugin: 'refdiff',
      Options: {
        repoId: 'github:GithubRepo:384111310',
        pairs: [
          { newRef: 'refs/tags/v0.6.0', oldRef: 'refs/tags/0.5.0' },
          { newRef: 'refs/tags/0.5.0', oldRef: 'refs/tags/0.4.0' }
        ]
      }
    }
  ]
]

export {
  refdiffConfig
}
