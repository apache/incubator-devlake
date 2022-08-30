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
}

const DefaultProviderConfig = {
  label: 'NullProvider',
  limit: null,
  icon: (w, h) => <Icon icon='box' size={w || 24} />,
  columns: {
    name: { label: 'Connection Name', placeholder: 'eg. Enter Instance Name', },
    endpoint: { label: 'Endpoint URL', placeholder: 'eg. https://null-api.localhost', },
    proxy: { label: 'Proxy URL', placeholder: 'eg. http://proxy.localhost:8080', },
    token: { label: 'Basic Auth Token', placeholder: 'eg. 3f5cda2a23ff410792e0', },
    username: { label: 'Username', placeholder: 'Enter Username / E-mail', },
    password: { label: 'Password', placeholder: 'Enter Password', },
    rateLimit: { label: 'Rate Limit', placeholder: '1000', },
  },
}

const ProviderConfigMap = {
  [Providers.NULL]: DefaultProviderConfig,
  [Providers.GITLAB]: {
    label: 'GitLab',
    icon: (w, h) => <GitlabProviderIcon width={w || 24} height={h || 24} />,
    columns: {
      name: { label: 'Connection Name', placeholder: 'eg. GitLab', },
      endpoint: { label: 'Endpoint URL', placeholder: 'eg. https://gitlab.com/api/v4/', },
      proxy: { label: 'Proxy URL', placeholder: 'eg. http://proxy.localhost:8080', },
      token: { label: 'Access Token', placeholder: 'eg. ff9d1ad0e5c04f1f98fa', },
      rateLimit: { label: <>Rate Limit <sup>(per hour)</sup></>, placeholder: '1000', },
    },
  },
  [Providers.JENKINS]: {
    label: 'Jenkins',
    icon: (w, h) => <JenkinsProviderIcon width={w || 24} height={h || 24} />,
    columns: {
      name: { label: 'Connection Name', placeholder: 'eg. Jenkins', },
      endpoint: { label: 'Endpoint URL', placeholder: 'URL eg. https://api.jenkins.io/', },
      proxy: { label: 'Proxy URL', placeholder: 'eg. http://proxy.localhost:8080', },
      username: { label: 'Username', placeholder: 'eg. admin', },
      password: { label: 'Password', placeholder: 'eg. ************', },
      rateLimit: { label: <>Rate Limit <sup>(per hour)</sup></>, placeholder: '1000', },
    },
  },
  [Providers.JIRA]: {
    label: 'JIRA',
    icon: (w, h) => <JiraProviderIcon width={w || 24} height={h || 24} />,
    columns: {
      name: { label: 'Connection Name', placeholder: 'eg. JIRA', },
      endpoint: { label: 'Endpoint URL', placeholder: 'eg. https://your-domain.atlassian.net/rest/', },
      proxy: { label: 'Proxy URL', placeholder: 'eg. http://proxy.localhost:8080', },
      username: { label: 'Username / E-mail', placeholder: 'eg. admin', },
      password: {
        label: (
          <>
            Password
            <Tooltip
              content={(
                <span>If you are using JIRA Cloud or JIRA Server, <br />your API Token should be used as password.</span>)}
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
          </>),
        placeholder: 'eg. ************',
      },
      rateLimit: { label: <>Rate Limit <sup>(per hour)</sup></>, placeholder: '1000', },
    },
  },
  [Providers.GITHUB]: {
    label: 'GitHub',
    icon: (w, h) => <GitHubProviderIcon width={w || 24} height={h || 24} />,
    columns: {
      name: { label: 'Connection Name', placeholder: 'eg. GitHub', },
      endpoint: { label: 'Endpoint URL', placeholder: 'eg. https://api.github.com/', },
      proxy: { label: 'Proxy URL', placeholder: 'eg. http://proxy.localhost:8080', },
      token: {
        label: (
          <>
            Auth Token(s)
            <Tooltip
              content={(
                <span>Due to Github's rate limit, input more tokens, <br />comma separated, to accelerate data collection.</span>)}
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
          </>),
        placeholder: 'eg. 4c5cbdb62c165e2b3d18, 40008ebccff9837bb8d2',
      },
      rateLimit: { label: <>Rate Limit <sup>(per hour)</sup></>, placeholder: '1000', },
    },
  },
  [Providers.TAPD]: {
    label: 'TAPD',
    icon: (w, h) => <TapdProviderIcon width={w || 24} height={h || 24} />,
    columns: {
      name: { label: 'Connection Name', placeholder: 'eg. Tapd', },
      endpoint: { label: 'Endpoint URL', placeholder: 'URL eg. https://api.tapd.cn/', },
      proxy: { label: 'Proxy URL', placeholder: 'eg. http://proxy.localhost:8080', },
      token: { label: 'Basic Auth Token', placeholder: 'eg. 6b057ffe68464c93a057', },
      rateLimit: { label: <>Rate Limit <sup>(per hour)</sup></>, placeholder: '1000', },
    },
  },
  [Providers.REFDIFF]: {
    label: 'RefDiff',
    icon: (w, h) => <Icon icon='box' size={w || 24} />
  },
  [Providers.GITEXTRACTOR]: {
    label: 'GitExtractor',
    icon: (w, h) => <Icon icon='box' size={w || 24} />
  },
  [Providers.FEISHU]: {
    label: 'Feishu',
    icon: (w, h) => <img src={FeishuIcon} width={w || 24} height={h || 24} />
  },
  [Providers.AE]: {
    label: 'Analysis Engine (AE)',
    icon: (w, h) => <Icon icon='box' size={w || 24} />
  },
  [Providers.DBT]: {
    label: 'Data Build Tool (DBT)',
    icon: (w, h) => <Icon icon='box' size={w || 24} />
  },
  [Providers.STARROCKS]: {
    label: 'STARROCKS',
    icon: (w, h) => <Icon icon='box' size={w || 24} />
  },
}

const ProviderTypes = {
  PLUGIN: 'plugin',
  INTEGRATION: 'integration',
  PIPELINE: 'pipeline'
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
  DefaultProviderConfig,
  ProviderConfigMap,
  ProviderTypes,
  ConnectionStatus,
  ConnectionStatusLabels
}
