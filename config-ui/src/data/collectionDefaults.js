const collectionDefaults = {
  github: {
    options: {
      repositoryName: 'lake',
      owner: 'merico-dev'
    },
    domainPlugin: {
      name: 'github-domain',
      options: {}
    }
  },
  gitlab: {
    options: {
      projectId: 8967944
    },
    domainPlugin: {
      name: 'gitlab-domain',
      options: {}
    }
  },
  jira: {
    options: {
      boardId: 8,
      sourceId: 1,
    },
    domainPlugin: {
      name: 'jiradomain',
      options: {}
    }
  },
  jenkins: {
    options: {},
    domainPlugin: {
      name: 'jenkinsdomain',
      options: {}
    }
  },
}

module.exports = collectionDefaults
