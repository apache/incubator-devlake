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
import React, { useEffect, useCallback, useMemo } from 'react'
import // Providers
// ProviderTypes,
// ProviderIcons,
// ConnectionStatus,
// ConnectionStatusLabels,
'@/data/Providers'
// import { DataEntities, DataEntityTypes } from '@/data/DataEntities'
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
  ) : null

const ProviderTransformationSettings = (props) => {
  const {
    Providers = {},
    ProviderLabels = {},
    ProviderIcons = {},
    provider,
    blueprint,
    connection,
    transformation = {},
    boards = {},
    entities = {},
    issueTypes = [],
    fields = [],
    onSettingsChange = () => {},
    isSaving = false,
    isSavingConnection = false,
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
      : // @todo Create <NoTransformations /> Message Component
        null,
    { ...props, entities: props.entities[props?.connection?.id] }
  )

  return (
    <div className='transformation-settings' data-provider={provider?.id}>
      {TransformationWithProviderSettings}
      {/* {provider?.id === Providers.GITHUB && (
        <GithubSettings
          provider={provider}
          connection={connection}
          transformation={transformation}
          onSettingsChange={onSettingsChange}
          entities={entities[connection?.id]}
          isSaving={isSaving}
          isSavingConnection={isSavingConnection}
        />
      )}

      {provider?.id === Providers.GITLAB && (
        <GitlabSettings
          provider={provider}
          connection={connection}
          transformation={transformation}
          onSettingsChange={onSettingsChange}
          entities={entities[connection?.id]}
          isSaving={isSaving}
          isSavingConnection={isSavingConnection}
        />
      )}

      {provider?.id === Providers.JIRA && (
        <JiraSettings
          provider={provider}
          connection={connection}
          issueTypes={issueTypes}
          fields={fields}
          transformation={transformation}
          onSettingsChange={onSettingsChange}
          entities={entities[connection?.id]}
          isSaving={isSaving}
          isSavingConnection={isSavingConnection}
          jiraProxyError={jiraProxyError}
          isFetchingJIRA={isFetchingJIRA}
        />
      )}

      {provider?.id === Providers.JENKINS && (
        <JenkinsSettings
          provider={provider}
          connection={connection}
          transformation={transformation}
          onSettingsChange={onSettingsChange}
          entities={entities[connection?.id]}
          isSaving={isSaving}
          isSavingConnection={isSavingConnection}
        />
      )}
      {provider?.id === Providers.TAPD && (
        <TapdSettings
          provider={provider}
          connection={connection}
          onSettingsChange={onSettingsChange}
          entities={entities[connection?.id]}
          isSaving={isSaving}
          isSavingConnection={isSavingConnection}
        />
      )}
      {provider?.id === Providers.AZURE && (
        <AzureSettings
          provider={provider}
          connection={connection}
          transformation={transformation}
          onSettingsChange={onSettingsChange}
          entities={entities[connection?.id]}
          isSaving={isSaving}
          isSavingConnection={isSavingConnection}
        />
      )}
      {provider?.id === Providers.BITBUCKET && (
        <BitbucketSettings
          provider={provider}
          connection={connection}
          transformation={transformation}
          onSettingsChange={onSettingsChange}
          entities={entities[connection?.id]}
          isSaving={isSaving}
          isSavingConnection={isSavingConnection}
        />
      )}
      {provider?.id === Providers.GITEE && (
        <GiteeSettings
          provider={provider}
          connection={connection}
          transformation={transformation}
          onSettingsChange={onSettingsChange}
          entities={entities[connection?.id]}
          isSaving={isSaving}
          isSavingConnection={isSavingConnection}
        />
      )} */}
    </div>
  )
}

export default ProviderTransformationSettings
