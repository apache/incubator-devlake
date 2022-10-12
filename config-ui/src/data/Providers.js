/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */
import React from 'react'
import { Tooltip, Icon } from '@blueprintjs/core'
import { ReactComponent as GitlabProviderIcon } from '@/images/integrations/gitlab.svg'
import { ReactComponent as JenkinsProviderIcon } from '@/images/integrations/jenkins.svg'
import { ReactComponent as JiraProviderIcon } from '@/images/integrations/jira.svg'
import { ReactComponent as GitHubProviderIcon } from '@/images/integrations/github.svg'
import { ReactComponent as TapdProviderIcon } from '@/images/integrations/tapd.svg'
import { ReactComponent as AzureProviderIcon } from '@/images/integrations/azure.svg'
import { ReactComponent as BitbucketProviderIcon } from '@/images/integrations/bitbucket.svg'
import { ReactComponent as GiteeProviderIcon } from '@/images/integrations/gitee.svg'
// import GitExtractorIcon from '@/images/git.png'
// import RefDiffIcon from '@/images/git-diff.png'
import FeishuIcon from '@/images/feishu.png'
// import DBTIcon from '@/images/dbt.png'
// import AEIcon from '@/images/ae.png'

const Providers = {
  NULL: 'null',
  GITLAB: 'gitlab',
  JENKINS: 'jenkins',
  JIRA: 'jira',
  GITHUB: 'github',
  REFDIFF: 'refdiff',
  GITEXTRACTOR: 'gitextractor',
  FEISHU: 'feishu',
  AE: 'ae',
  DBT: 'dbt',
  STARROCKS: 'starrocks',
  TAPD: 'tapd',
  AZURE: 'azure',
  BITBUCKET: 'bitbucket',
  GITEE: 'gitee'
}

const ProviderTypes = {
  PLUGIN: 'plugin',
  INTEGRATION: 'integration',
  PIPELINE: 'pipeline'
}

const ProviderLabels = {
  NULL: 'NullProvider',
  GITLAB: 'GitLab',
  JENKINS: 'Jenkins',
  JIRA: 'JIRA',
  GITHUB: 'GitHub',
  REFDIFF: 'RefDiff',
  GITEXTRACTOR: 'GitExtractor',
  FEISHU: 'Feishu',
  AE: 'Analysis Engine (AE)',
  DBT: 'Data Build Tool (DBT)',
  STARROCKS: 'StarRocks',
  TAPD: 'TAPD',
  AZURE: 'Azure CI',
  BITBUCKET: 'BitBucket',
  GITEE: 'Gitee'
}

const ProviderConnectionLimits = {
  // (All providers are mult-connection, no source limits defined)
  // jenkins: null,
  // jira: null,
  // github: null
  // gitlab: null
}

// NOTE: Not all fields may be referenced/displayed for a provider,
// ie. JIRA prefers $token over $username and $password
const ProviderFormLabels = {
  null: {
    name: 'Connection Name',
    endpoint: 'Endpoint URL',
    proxy: 'Proxy URL',
    token: 'Basic Auth Token',
    username: 'Username',
    password: 'Password',
    rateLimitPerHour: 'Rate Limit'
  },
  gitlab: {
    name: 'Connection Name',
    endpoint: 'Endpoint URL',
    proxy: 'Proxy URL',
    token: 'Access Token',
    username: 'Username',
    password: 'Password',
    rateLimitPerHour: (
      <>
        Rate Limit <sup>(per hour)</sup>
      </>
    )
  },
  jenkins: {
    name: 'Connection Name',
    endpoint: 'Endpoint URL',
    proxy: 'Proxy URL',
    token: 'Basic Auth Token',
    username: 'Username',
    password: 'Password',
    rateLimitPerHour: (
      <>
        Rate Limit <sup>(per hour)</sup>
      </>
    )
  },
  tapd: {
    name: 'Connection Name',
    endpoint: 'Endpoint URL',
    proxy: 'Proxy URL',
    token: 'Basic Auth Token',
    username: 'Username',
    password: 'Password',
    rateLimitPerHour: (
      <>
        Rate Limit <sup>(per hour)</sup>
      </>
    )
  },
  jira: {
    name: 'Connection Name',
    endpoint: 'Endpoint URL',
    token: 'Basic Auth Token',
    username: 'Username / E-mail',
    proxy: 'Proxy URL',
    password: (
      <>
        Password
        <Tooltip
          content={
            <span>
              If you are using JIRA Cloud or JIRA Server, <br />
              your API Token should be used as password.
            </span>
          }
          intent='primary'
        >
          <Icon
            icon='info-sign'
            size={12}
            style={{
              float: 'left',
              display: 'inline-block',
              alignContent: 'center',
              marginBottom: '4px',
              marginLeft: '8px',
              color: '#999'
            }}
          />
        </Tooltip>
      </>
    ),
    rateLimitPerHour: (
      <>
        Rate Limit <sup>(per hour)</sup>
      </>
    )
  },
  github: {
    name: 'Connection Name',
    endpoint: 'Endpoint URL',
    proxy: 'Proxy URL',
    token: (
      <>
        Auth Token(s)
        <Tooltip
          content={
            <span>
              Due to Github's rate limit, input more tokens, <br />
              comma separated, to accelerate data collection.
            </span>
          }
          intent='primary'
        >
          <Icon
            icon='info-sign'
            size={12}
            style={{
              float: 'left',
              display: 'inline-block',
              alignContent: 'center',
              marginBottom: '4px',
              marginLeft: '8px',
              color: '#999'
            }}
          />
        </Tooltip>
      </>
    ),
    username: 'Username',
    password: 'Password',
    rateLimitPerHour: (
      <>
        Rate Limit <sup>(per hour)</sup>
      </>
    )
  },
  azure: {
    name: 'Connection Name',
    endpoint: 'Endpoint URL',
    proxy: 'Proxy URL',
    token: 'Basic Auth Token',
    username: 'Username',
    password: 'Password',
    rateLimitPerHour: (
      <>
        Rate Limit <sup>(per hour)</sup>
      </>
    )
  },
  bitbucket: {
    name: 'Connection Name',
    endpoint: 'Endpoint URL',
    proxy: 'Proxy URL',
    token: 'Basic Auth Token',
    username: 'Username',
    password: 'Password',
    rateLimitPerHour: (
      <>
        Rate Limit <sup>(per hour)</sup>
      </>
    )
  },
  gitee: {
    name: 'Connection Name',
    endpoint: 'Endpoint URL',
    proxy: 'Proxy URL',
    token: 'Basic Auth Token',
    username: 'Username',
    password: 'Password',
    rateLimitPerHour: (
      <>
        Rate Limit <sup>(per hour)</sup>
      </>
    )
  }
}

const ProviderFormPlaceholders = {
  null: {
    name: 'eg. Enter Instance Name',
    endpoint: 'eg. https://null-api.localhost',
    proxy: 'eg. http://proxy.localhost:8080',
    token: 'eg. 3f5cda2a23ff410792e0',
    username: 'Enter Username / E-mail',
    password: 'Enter Password',
    rateLimitPerHour: '1000'
  },
  gitlab: {
    name: 'eg. GitLab',
    endpoint: 'eg. https://gitlab.com/api/v4/',
    proxy: 'eg. http://proxy.localhost:8080',
    token: 'eg. ff9d1ad0e5c04f1f98fa',
    username: 'Enter Username / E-mail',
    password: 'Enter Password',
    rateLimitPerHour: '1000'
  },
  jenkins: {
    name: 'eg. Jenkins',
    endpoint: 'URL eg. https://api.jenkins.io/',
    proxy: 'eg. http://proxy.localhost:8080',
    token: 'eg. 6b057ffe68464c93a057',
    username: 'eg. admin',
    password: 'eg. ************',
    rateLimitPerHour: '1000'
  },
  tapd: {
    name: 'eg. Tapd',
    endpoint: 'URL eg. https://api.tapd.cn/',
    proxy: 'eg. http://proxy.localhost:8080',
    token: 'eg. 6b057ffe68464c93a057',
    username: 'eg. admin',
    password: 'eg. ************',
    rateLimitPerHour: '1000'
  },
  jira: {
    name: 'eg. JIRA',
    endpoint: 'eg. https://your-domain.atlassian.net/rest/',
    proxy: 'eg. http://proxy.localhost:8080',
    token: 'eg. 8c06a7cc50b746bfab30',
    username: 'eg. admin',
    password: 'eg. ************',
    rateLimitPerHour: '1000'
  },
  github: {
    name: 'eg. GitHub',
    endpoint: 'eg. https://api.github.com/',
    proxy: 'eg. http://proxy.localhost:8080',
    token: 'eg. 4c5cbdb62c165e2b3d18, 40008ebccff9837bb8d2',
    username: 'eg. admin',
    password: 'eg. ************',
    rateLimitPerHour: '1000'
  },
  azure: {
    name: 'eg. Azure',
    endpoint: 'eg. https://api.azure.com/',
    proxy: 'eg. http://proxy.localhost:8080',
    token: 'eg. 4c5cbdb62c165e2b3d18, 40008ebccff9837bb8d2',
    username: 'eg. admin',
    password: 'eg. ************',
    rateLimitPerHour: '1000'
  },
  bitbucket: {
    name: 'eg. Bitbucket',
    endpoint: 'eg. https://api.bitbucket.com/',
    proxy: 'eg. http://proxy.localhost:8080',
    token: 'eg. 4c5cbdb62c165e2b3d18, 40008ebccff9837bb8d2',
    username: 'eg. admin',
    password: 'eg. ************',
    rateLimitPerHour: '1000'
  },
  gitee: {
    name: 'eg. Gitee',
    endpoint: 'eg. https://api.gitee.com/',
    proxy: 'eg. http://proxy.localhost:8080',
    token: 'eg. 4c5cbdb62c165e2b3d18, 40008ebccff9837bb8d2',
    username: 'eg. admin',
    password: 'eg. ************',
    rateLimitPerHour: '1000'
  }
}

const ProviderIcons = {
  [Providers.GITLAB]: (w, h) => (
    <GitlabProviderIcon width={w || 24} height={h || 24} />
  ),
  [Providers.JENKINS]: (w, h) => (
    <JenkinsProviderIcon width={w || 24} height={h || 24} />
  ),
  [Providers.TAPD]: (w, h) => (
    <TapdProviderIcon width={w || 24} height={h || 24} />
  ),
  [Providers.JIRA]: (w, h) => (
    <JiraProviderIcon width={w || 24} height={h || 24} />
  ),
  [Providers.GITHUB]: (w, h) => (
    <GitHubProviderIcon width={w || 24} height={h || 24} />
  ),
  [Providers.REFDIFF]: (w, h) => <Icon icon='box' size={w || 24} />,
  [Providers.GITEXTRACTOR]: (w, h) => <Icon icon='box' size={w || 24} />,
  [Providers.FEISHU]: (w, h) => (
    <img src={FeishuIcon} width={w || 24} height={h || 24} />
  ),
  [Providers.AE]: (w, h) => <Icon icon='box' size={w || 24} />,
  [Providers.DBT]: (w, h) => <Icon icon='box' size={w || 24} />,
  // @todo: update with svg icons
  [Providers.AZURE]: (w, h) => (
    <AzureProviderIcon width={w || 24} height={h || 24} />
  ),
  [Providers.BITBUCKET]: (w, h) => (
    <BitbucketProviderIcon width={w || 24} height={h || 24} />
  ),
  [Providers.GITEE]: (w, h) => (
    <GiteeProviderIcon width={w || 24} height={h || 24} />
  )
}

const ConnectionStatus = {
  OFFLINE: 0,
  ONLINE: 1,
  DISCONNECTED: 2,
  TESTING: 3
}

const ConnectionStatusLabels = {
  [ConnectionStatus.OFFLINE]: 'Offline',
  [ConnectionStatus.ONLINE]: 'Online',
  [ConnectionStatus.DISCONNECTED]: 'Disconnected',
  [ConnectionStatus.TESTING]: 'Testing'
}

export {
  Providers,
  ProviderTypes,
  ProviderIcons,
  ProviderLabels,
  ProviderConnectionLimits,
  ProviderFormLabels,
  ProviderFormPlaceholders,
  ConnectionStatus,
  ConnectionStatusLabels
}
