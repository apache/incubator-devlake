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
import { Icon } from '@blueprintjs/core'
import { Providers, ProviderLabels, ProviderTypes } from '@/data/Providers'

// import JiraSettings from '@/pages/configure/settings/jira'
// import GitlabSettings from '@/pages/configure/settings/gitlab'
// import JenkinsSettings from '@/pages/configure/settings/jenkins'
// import TapdSettings from '@/pages/configure/settings/tapd'
// import GithubSettings from '@/pages/configure/settings/github'

import { ReactComponent as GitlabProvider } from '@/images/integrations/gitlab.svg'
import { ReactComponent as JenkinsProvider } from '@/images/integrations/jenkins.svg'
import { ReactComponent as JiraProvider } from '@/images/integrations/jira.svg'
import { ReactComponent as GitHubProvider } from '@/images/integrations/github.svg'
import { ReactComponent as TapdProvider } from '@/images/integrations/tapd.svg'
import { ReactComponent as AzureProvider } from '@/images/integrations/azure.svg'
import { ReactComponent as BitbucketProvider } from '@/images/integrations/bitbucket.svg'
import { ReactComponent as GiteeProvider } from '@/images/integrations/gitee.svg'
// import GitExtractorProvider from '@/images/git.png'
// import RefDiffProvider from '@/images/git-diff.png'
// import { ReactComponent as NullProvider } from '@/images/integrations/null.svg'

// @todo: TO-BE replaced with Integrations Hook
const integrationsData = [
  {
    id: Providers.GITLAB,
    type: ProviderTypes.INTEGRATION,
    enabled: true,
    multiConnection: true,
    isBeta: false,
    name: ProviderLabels.GITLAB,
    icon: (
      <GitlabProvider
        className='providerIconSvg'
        width='30'
        height='30'
        style={{ float: 'left', marginTop: '5px' }}
      />
    ),
    iconDashboard: (
      <GitlabProvider className='providerIconSvg' width='40' height='40' />
    ),
    // @todo: relocated to ProviderTransformationSettings since v0.12.0
    settings: {}
  },
  {
    id: Providers.JENKINS,
    type: ProviderTypes.INTEGRATION,
    enabled: true,
    multiConnection: true,
    isBeta: false,
    name: ProviderLabels.JENKINS,
    icon: (
      <JenkinsProvider
        className='providerIconSvg'
        width='30'
        height='30'
        style={{ float: 'left', marginTop: '5px' }}
      />
    ),
    iconDashboard: (
      <JenkinsProvider className='providerIconSvg' width='40' height='40' />
    ),
    // @todo: relocated to ProviderTransformationSettings since v0.12.0
    settings: {}
  },
  {
    id: Providers.TAPD,
    type: ProviderTypes.INTEGRATION,
    enabled: true,
    multiConnection: true,
    isBeta: true,
    name: ProviderLabels.TAPD,
    icon: (
      <TapdProvider
        className='providerIconSvg'
        width='30'
        height='30'
        style={{ float: 'left', marginTop: '5px' }}
      />
    ),
    iconDashboard: (
      <TapdProvider className='providerIconSvg' width='40' height='40' />
    ),
    // relocated to ProviderTransformationSettings since v0.12.0
    settings: {}
  },
  {
    id: Providers.JIRA,
    type: ProviderTypes.INTEGRATION,
    enabled: true,
    multiConnection: true,
    isBeta: false,
    name: ProviderLabels.JIRA,
    icon: (
      <JiraProvider
        className='providerIconSvg'
        width='30'
        height='30'
        style={{ float: 'left', marginTop: '5px' }}
      />
    ),
    iconDashboard: (
      <JiraProvider className='providerIconSvg' width='40' height='40' />
    ),
    // relocated to ProviderTransformationSettings since v0.12.0
    settings: {}
  },
  {
    id: Providers.GITHUB,
    type: ProviderTypes.INTEGRATION,
    enabled: true,
    multiConnection: true,
    isBeta: false,
    name: ProviderLabels.GITHUB,
    icon: (
      <GitHubProvider
        className='providerIconSvg'
        width='30'
        height='30'
        style={{ float: 'left', marginTop: '5px' }}
      />
    ),
    iconDashboard: (
      <GitHubProvider className='providerIconSvg' width='40' height='40' />
    ),
    // relocated to ProviderTransformationSettings since v0.12.0
    settings: {}
  },
  {
    id: Providers.AZURE,
    type: ProviderTypes.INTEGRATION,
    enabled: true,
    multiConnection: true,
    isBeta: false,
    name: ProviderLabels.AZURE,
    icon: (
      <AzureProvider
        className='providerIconSvg'
        width='30'
        height='30'
        style={{ float: 'left', marginTop: '5px' }}
      />
    ),
    iconDashboard: (
      <AzureProvider className='providerIconSvg' width='40' height='40' />
    ),
    // relocated to ProviderTransformationSettings since v0.12.0
    settings: {}
  },
  {
    id: Providers.BITBUCKET,
    type: ProviderTypes.INTEGRATION,
    enabled: true,
    multiConnection: true,
    isBeta: false,
    name: ProviderLabels.BITBUCKET,
    icon: (
      <BitbucketProvider
        className='providerIconSvg'
        width='30'
        height='30'
        style={{ float: 'left', marginTop: '5px' }}
      />
    ),
    iconDashboard: (
      <BitbucketProvider className='providerIconSvg' width='40' height='40' />
    ),
    // relocated to ProviderTransformationSettings since v0.12.0
    settings: {}
  },
  {
    id: Providers.GITEE,
    type: ProviderTypes.INTEGRATION,
    enabled: true,
    multiConnection: true,
    isBeta: false,
    name: ProviderLabels.GITEE,
    icon: (
      <GiteeProvider
        className='providerIconSvg'
        width='30'
        height='30'
        style={{ float: 'left', marginTop: '5px' }}
      />
    ),
    iconDashboard: (
      <GiteeProvider className='providerIconSvg' width='40' height='40' />
    ),
    // relocated to ProviderTransformationSettings since v0.12.0
    settings: {}
  }
]

// @todo: deprecate this var, used for legacy V11 Pipeline
const pluginsData = [
  {
    id: Providers.GITEXTRACTOR,
    type: ProviderTypes.PIPELINE,
    enabled: true,
    multiConnection: false,
    name: ProviderLabels.GITEXTRACTOR,
    icon: <Icon icon='box' size={30} />,
    iconDashboard: <Icon icon='box' size={32} />,
    // relocated to ProviderTransformationSettings since v0.12.0
    settings: {}
  },
  {
    id: Providers.REFDIFF,
    type: ProviderTypes.PIPELINE,
    enabled: true,
    multiConnection: false,
    name: ProviderLabels.REFDIFF,
    icon: <Icon icon='box' size={30} />,
    iconDashboard: <Icon icon='box' size={32} />,
    // relocated to ProviderTransformationSettings since v0.12.0
    settings: {}
  }
]

export { integrationsData, pluginsData }
