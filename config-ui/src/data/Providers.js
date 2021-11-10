const Providers = {
  NULL: 'null',
  GITLAB: 'gitlab',
  JENKINS: 'jenkins',
  JIRA: 'jira',
  GITHUB: 'github',
}

const ProviderLabels = {
  NULL: 'NullProvider',
  GITLAB: 'GitLab',
  JENKINS: 'Jenkins',
  JIRA: 'JIRA',
  GITHUB: 'GitHub',
}

const ProviderSourceLimits = {
  gitlab: 1,
  jenkins: 1,
  // jira: null, // (Multi-source, no-limit)
  github: 1
}

export {
  Providers,
  ProviderLabels,
  ProviderSourceLimits
}
