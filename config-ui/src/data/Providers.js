import React from 'react'
import { Tooltip } from '@blueprintjs/core'
import { ReactComponent as GitlabProviderIcon } from '@/images/integrations/gitlab.svg'
import { ReactComponent as JenkinsProviderIcon } from '@/images/integrations/jenkins.svg'
import { ReactComponent as JiraProviderIcon } from '@/images/integrations/jira.svg'
import { ReactComponent as GitHubProviderIcon } from '@/images/integrations/github.svg'

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
    // token: 'Auth Token(s)',
    token: (
      <Tooltip
        content={(<span>Due to Github’s rate limit, input more tokens, <br />comma separated, to accelerate data collection.</span>)}
        intent='primary'
      >
        Auth Token(s)
      </Tooltip>),
    username: 'Username',
    password: 'Password'
  },
}

const ProviderFormPlaceholders = {
  null: {
    name: 'Enter Instance Name',
    endpoint: 'Enter Endpoint URL eg. https://null-api.localhost',
    token: 'Enter Auth Token eg. 3f5cda2a23ff410792e0',
    username: 'Enter Username / E-mail',
    password: 'Enter Password'
  },
  gitlab: {
    name: 'Enter Instance Name',
    endpoint: 'Enter Endpoint URL eg. https://gitlab.com/api/v4',
    token: 'Enter Auth Token eg. ff9d1ad0e5c04f1f98fa',
    username: 'Enter Username / E-mail',
    password: 'Enter Password'
  },
  jenkins: {
    name: 'Enter Instance Name',
    endpoint: 'Enter Endpoint URL eg. https://api.jenkins.io',
    token: 'Enter Auth Token eg. 6b057ffe68464c93a057',
    username: 'Enter Username / E-mail',
    password: 'Enter Password'
  },
  jira: {
    name: 'Enter Instance Name',
    endpoint: 'Enter Endpoint URL eg. https://your-domain.atlassian.net/rest/',
    token: 'Enter Auth Token eg. 8c06a7cc50b746bfab30',
    username: 'Enter Username / E-mail',
    password: 'Enter Password'
  },
  github: {
    name: 'Enter Instance Name',
    endpoint: 'Enter Endpoint URL eg. https://api.github.com',
    token: 'Enter Auth Token(s) eg. 4c5cbdb62c165e2b3d18, 40008ebccff9837bb8d2',
    username: 'Enter Username / E-mail',
    password: 'Enter Password'
  }
}

const ProviderIcons = {
  [Providers.GITLAB]: (w, h) => <GitlabProviderIcon width={w || 24} height={h || 24} />,
  [Providers.JENKINS]: (w, h) => <JenkinsProviderIcon width={w || 24} height={h || 24} />,
  [Providers.JIRA]: (w, h) => <JiraProviderIcon width={w || 24} height={h || 24} />,
  [Providers.GITHUB]: (w, h) => <GitHubProviderIcon width={w || 24} height={h || 24} />,
}

export {
  Providers,
  ProviderIcons,
  ProviderLabels,
  ProviderSourceLimits,
  ProviderFormLabels,
  ProviderFormPlaceholders
}
