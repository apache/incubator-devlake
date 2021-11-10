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

// NOTE: Not all fields may be referenced/displayed for a provider,
// ie. JIRA prefers $token over $username and $password
const ProviderFormLabels = {
  null: {
    name: 'Connection Name',
    endpoint: 'Endpoint URL',
    token: 'Basic Auth Token',
    username: 'Username',
    password: 'Password'
  },
  gitlab: {
    name: 'Connection Name',
    endpoint: 'Endpoint URL',
    token: 'Basic Auth Token',
    username: 'Username',
    password: 'Password'
  },
  jenkins: {
    name: 'Connection Name',
    endpoint: 'Endpoint URL',
    token: 'Basic Auth Token',
    username: 'Username',
    password: 'Password'
  },
  jira: {
    name: 'Connection Name',
    endpoint: 'Endpoint URL',
    token: 'Basic Auth Token',
    username: 'Username',
    password: 'Password'
  },
  github: {
    name: 'Connection Name',
    endpoint: 'Endpoint URL',
    token: 'Auth Token(s)',
    username: 'Username',
    password: 'Password'
  },
}

const ProviderFormPlaceholders = {
  null: {
    name: 'Enter Instance Name',
    endpoint: 'Enter Endpoint URL eg. https://null-api.localhost',
    token: 'Enter Auth Token eg. EJrLG8DNeXADQcGOaaaX4B47',
    username: 'Enter Username / E-mail',
    password: 'Enter Password'
  },
  gitlab: {
    name: 'Enter Instance Name',
    endpoint: 'Enter Endpoint URL eg. https://gitlab.com/api/v4',
    token: 'Enter Auth Token eg. EJrLG8DNeXADQcGOaaaX4B47',
    username: 'Enter Username / E-mail',
    password: 'Enter Password'
  },
  jenkins: {
    name: 'Enter Instance Name',
    endpoint: 'Enter Endpoint URL eg. https://api.jenkins.io',
    token: 'Enter Auth Token eg. EJrLG8DNeXADQcGOaaaX4B47',
    username: 'Enter Username / E-mail',
    password: 'Enter Password'
  },
  jira: {
    name: 'Enter Instance Name',
    endpoint: 'Enter Endpoint URL eg. https://your-domain.atlassian.net/rest/api',
    token: 'Enter Auth Token eg. EJrLG8DNeXADQcGOaaaX4B47',
    username: 'Enter Username / E-mail',
    password: 'Enter Password'
  },
  github: {
    name: 'Enter Instance Name',
    endpoint: 'Enter Endpoint URL eg. https://api.github.com',
    token: 'Enter Auth Token eg. EJrLG8DNeXADQcGOaaaX4B47',
    username: 'Enter Username / E-mail',
    password: 'Enter Password'
  }
}

export {
  Providers,
  ProviderLabels,
  ProviderSourceLimits,
  ProviderFormLabels,
  ProviderFormPlaceholders
}
