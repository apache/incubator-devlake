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
import React, { useMemo } from 'react'
import NoData from '@/components/NoData'
import JiraSettings from '@/pages/configure/settings/jira'
import GitlabSettings from '@/pages/configure/settings/gitlab'
import JenkinsSettings from '@/pages/configure/settings/jenkins'
import TapdSettings from '@/pages/configure/settings/tapd'
import GithubSettings from '@/pages/configure/settings/github'
import AzureSettings from '@/pages/configure/settings/azure'
import BitbucketSettings from '@/pages/configure/settings/bitbucket'
import GiteeSettings from '@/pages/configure/settings/gitee'

// Transformation Higher-Order Component (HOC) Settings Loader
const withTransformationSettings = (
  TransformationComponent,
  TransformationProps
) =>
  TransformationComponent ? (
    <TransformationComponent {...TransformationProps} />
  ) : (
    <NoData
      title='No Transformations'
      icon='disable'
      message='This provider does not have additional transformation settings'
      onClick={null}
      actionText={null}
    />
  )

const ProviderTransformationSettings = (props) => {
  const {
    Providers,
    provider,
    blueprint,
    connection,
    transformation = {},
    dataDomainsGroup = {},
    isSaving = false,
    isSavingConnection = false,
    onSettingsChange = () => {},

    // only jira used
    issueTypes = [],
    fields = [],
    jiraProxyError,
    isFetchingJIRA = false
  } = props

  // Provider Transformation Components (LOCAL)
  const TransformationComponents = useMemo(
    () => ({
      [Providers.GITHUB]: GithubSettings,
      [Providers.GITLAB]: GitlabSettings,
      [Providers.JIRA]: JiraSettings,
      [Providers.JENKINS]: JenkinsSettings,
      [Providers.TAPD]: TapdSettings,
      [Providers.AZURE]: AzureSettings,
      [Providers.BITBUCKET]: BitbucketSettings,
      [Providers.GITEE]: GiteeSettings
    }),
    [Providers]
  )

  // Dynamic Transformation Settings via HOC
  const TransformationWithProviderSettings = withTransformationSettings(
    provider?.id && TransformationComponents[provider?.id]
      ? TransformationComponents[provider?.id]
      : null,
    {
      // pass all props but notice default values in Line#47~63 have no effect
      ...props,
      dataDomains: dataDomainsGroup[connection?.id]
    }
  )

  return (
    <div className='transformation-settings' data-provider={provider?.id}>
      {TransformationWithProviderSettings}
    </div>
  )
}

export default ProviderTransformationSettings
