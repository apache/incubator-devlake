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
import React, { useState, useEffect, useCallback, useMemo } from 'react'

import Plugin from '@/models/Plugin'

// LOCAL PLUGIN REGISTRY
// "integration" Plugins a.k.a "Providers"
import JiraPlugin from '@/registry/plugins/jira.json'
import GitHubPlugin from '@/registry/plugins/github.json'
import GitHubGraphqlPlugin from '@/registry/plugins/github_graphql.json'
import GitLabPlugin from '@/registry/plugins/gitlab.json'
import JenkinsPlugin from '@/registry/plugins/jenkins.json'
import TapdPlugin from '@/registry/plugins/tapd.json'
import AzurePlugin from '@/registry/plugins/azure.json'
import BitbucketPlugin from '@/registry/plugins/bitbucket.json'
import GiteePlugin from '@/registry/plugins/gitee.json'
import AePlugin from '@/registry/plugins/ae.json'
import RefdiffPlugin from '@/registry/plugins/refdiff.json'
import DbtPlugin from '@/registry/plugins/dbt.json'
import StarrocksPlugin from '@/registry/plugins/starrocks.json'
import DoraPlugin from '@/registry/plugins/dora.json'

const ProviderTypes = {
  PLUGIN: 'plugin',
  INTEGRATION: 'integration',
  PIPELINE: 'pipeline',
  WEBHOOK: 'webhook'
}

function useIntegrations(
  pluginRegistry = [
    JiraPlugin,
    GitHubPlugin,
    GitHubGraphqlPlugin,
    GitLabPlugin,
    JenkinsPlugin,
    TapdPlugin,
    AzurePlugin,
    BitbucketPlugin,
    GiteePlugin,
    AePlugin,
    RefdiffPlugin,
    DbtPlugin,
    StarrocksPlugin,
    DoraPlugin
  ]
) {
  const [registry, setRegistry] = useState(pluginRegistry || [])
  const [plugins, setPlugins] = useState([])

  // @todo: fetch live/dynamic plugins from API
  const [apiPlugins, setApiPlugins] = useState([])
  const [apiRegistry, setApiRegistry] = useState([])

  const [activeProvider, setActiveProvider] = useState()

  const integrations = useMemo(
    () => plugins.filter((p) => p.type === ProviderTypes.INTEGRATION),
    [plugins]
  )

  const DataSources = useMemo(
    () =>
      integrations.map((P, iDx) => ({
        id: iDx + 1,
        name: P.name,
        title: P.name,
        value: P.id
      })),
    [integrations]
  )

  const Providers = useMemo(
    () =>
      plugins
        .map((P) => P.id)
        .reduce(
          (pV, cV, iDx) => ({ ...pV, [cV.toString()?.toUpperCase()]: cV }),
          {}
        ),
    [plugins]
  )

  const ProviderLabels = useMemo(
    () =>
      plugins
        .map((P) => P)
        .reduce(
          (pV, cV, iDx) => ({
            ...pV,
            [cV?.id.toString()?.toUpperCase()]: cV.name
          }),
          {}
        ),
    [plugins]
  )

  const ProviderFormLabels = useMemo(
    () =>
      integrations
        .map((P) => P.getConnectionFormLabels())
        .reduce((pV, cV, iDx) => ({ ...pV, [integrations[iDx]?.id]: cV }), {}),
    [integrations]
  )

  const ProviderFormPlaceholders = useMemo(
    () =>
      integrations
        .map((P) => P.getConnectionFormPlaceholders())
        .reduce((pV, cV, iDx) => ({ ...pV, [integrations[iDx]?.id]: cV }), {}),
    [integrations]
  )

  const ProviderFormTooltips = useMemo(
    () =>
      integrations
        .map((P) => P.getConnectionFormTooltips())
        .reduce((pV, cV, iDx) => ({ ...pV, [integrations[iDx]?.id]: cV }), {}),
    [integrations]
  )

  const ProviderIcons = useMemo(
    () =>
      integrations
        .map((p) => p.icon)
        .reduce(
          (pV, cV, iDx) => ({
            ...pV,
            [integrations[iDx]?.id]: (w, h) => (
              <img
                src={'/' + cV}
                style={{ width: `${w}px`, height: `${h}px` }}
              />
            )
          }),
          {}
        ),
    [integrations]
  )

  const ProviderConnectionLimits = useMemo(
    () =>
      integrations
        .map((P) => parseInt(P.connectionLimit, 10))
        .reduce((pV, cV, iDx) => ({ ...pV, [integrations[iDx]?.id]: cV }), {}),
    [integrations]
  )

  const ProviderTransformations = useMemo(
    () =>
      integrations
        .map((P) => P.getDefaultTransformations())
        .reduce((pV, cV, iDx) => ({ ...pV, [integrations[iDx]?.id]: cV }), {}),
    [integrations]
  )

  const registerPlugin = useCallback((pluginConfig) => {
    console.log(
      '>>> REGISTERING PLUGIN...',
      `${pluginConfig?.name} [${pluginConfig?.type}]`
    )
    // @todo: Validate Plugin before Registration
    return new Plugin(pluginConfig)
  }, [])

  const validatePlugin = useCallback((pluginConfig) => {
    let isValid = false
    const requiredProperties = [
      'id',
      'name',
      'type',
      'enabled',
      'multiConnection',
      'icon'
    ]
    // todo: enhance plugin validation
    try {
      console.log('>>> INTEGRATIONS HOOK: VALIDATING PLUGIN...', pluginConfig)
      JSON.parse(JSON.stringify(pluginConfig))
      isValid = requiredProperties.every((p) =>
        Object.prototype.hasOwnProperty.call(pluginConfig, p)
      )
      if (!isValid) {
        console.log(
          '>>> INTEGRATIONS HOOK: PLUGIN SCHEMA INCOMPLETE, MISSING REQUIRED PROPERTIES!',
          pluginConfig
        )
      }
    } catch (e) {
      console.log(
        '>>> INTEGRATIONS HOOK: PLUGIN VALIDATION FAILED!',
        e,
        pluginConfig
      )
    }
    return isValid
  }, [])

  const getPlugin = useCallback(
    (pluginId) => {
      return plugins.find((p) => p.id === pluginId)
    },
    [plugins]
  )

  useEffect(() => {
    console.log(
      '>>> INTEGRATIONS HOOK: PLUGIN REGISTRY CONFIGURATION!!!',
      registry
    )
    setPlugins((aP) => [
      // ...aP,
      ...registry
        .filter((p) => p.enabled)
        .filter((p) => validatePlugin(p))
        .map((p) => registerPlugin(p))
    ])
  }, [registry, setPlugins, validatePlugin, registerPlugin])

  useEffect(() => {
    console.log(
      '>>> INTEGRATIONS HOOK: REGISTERED PLUGIN OBJECT CLASSES...',
      plugins
    )
    setActiveProvider(plugins[0])
  }, [plugins])

  useEffect(() => {
    console.log(
      '>>> INTEGRATIONS HOOK: REGISTERED LIVE API PLUGIN OBJECT CLASSES...',
      apiPlugins
    )
  }, [apiPlugins])

  useEffect(() => {
    console.log('>>> INTEGRATIONS HOOK: ACTIVE PROVIDER..', activeProvider)
  }, [activeProvider])

  useEffect(() => {
    console.log(
      '>>> INTEGRATIONS HOOK: PROVIDERS CONFIGURATION LIST ...',
      Providers
    )
    console.log(
      '>>> INTEGRATIONS HOOK: PROVIDER LABELS CONFIGURATION LIST ...',
      ProviderLabels
    )
    console.log(
      '>>> INTEGRATIONS HOOK: PROVIDER CONFIGURATION CONNECTION FORM LABELS..',
      ProviderFormLabels
    )
    console.log(
      '>>> INTEGRATIONS HOOK: PROVIDER CONFIGURATION CONNECTION  FORM PLACEHOLDERS..',
      ProviderFormPlaceholders
    )
    console.log(
      '>>> INTEGRATIONS HOOK: PROVIDER CONFIGURATION CONNECTION  FORM TOOLTIPS..',
      ProviderFormTooltips
    )
    console.log(
      '>>> INTEGRATIONS HOOK: PROVIDER CONFIGURATION PROVIDER ICONS..',
      ProviderIcons
    )
    console.log(
      '>>> INTEGRATIONS HOOK: PROVIDER CONNECTION LIMITS...',
      ProviderConnectionLimits
    )
    console.log(
      '>>> INTEGRATIONS HOOK: PROVIDER DATA SOURCES LIST...',
      DataSources
    )
  }, [
    activeProvider,
    Providers,
    ProviderLabels,
    ProviderFormLabels,
    ProviderFormPlaceholders,
    ProviderConnectionLimits,
    ProviderFormTooltips,
    ProviderIcons,
    DataSources
  ])

  return {
    activeProvider,
    registry,
    plugins,
    apiPlugins,
    apiRegistry,
    integrations,
    DataSources,
    Providers,
    ProviderLabels,
    ProviderFormLabels,
    ProviderFormPlaceholders,
    ProviderFormTooltips,
    ProviderIcons,
    ProviderConnectionLimits,
    ProviderTypes,
    ProviderTransformations,
    setRegistry,
    setApiRegistry,
    setActiveProvider,
    getPlugin
  }
}

export default useIntegrations
