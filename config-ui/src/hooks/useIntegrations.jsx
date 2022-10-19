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
import { useState, useEffect, useCallback, useMemo } from 'react'

import Plugin from '@/models/Plugin'

// LOCAL PLUGIN REGISTRY
// "integration" Plugins a.k.a "Providers"
import JiraPlugin from '@/registry/plugins/jira.json'
import GitHubPlugin from '@/registry/plugins/github.json'
import GitLabPlugin from '@/registry/plugins/gitlab.json'
import JenkinsPlugin from '@/registry/plugins/jenkins.json'
import TapdPlugin from '@/registry/plugins/tapd.json'
// "plugin" Plugins (Backend / Advanced Mode Plugins)
import AePlugin from '@/registry/plugins/ae.json'
import AzurePlugin from '@/registry/plugins/azure.json'
import BitbucketPlugin from '@/registry/plugins/bitbucket.json'
import GiteePlugin from '@/registry/plugins/gitee.json'
// @todo: import additional backend plugins
// import RefdiffPlugin from '@/registry/plugins/refdiff.json'

const ProviderTypes = {
  PLUGIN: 'plugin',
  INTEGRATION: 'integration',
  PIPELINE: 'pipeline'
}

function useIntegrations(
  pluginRegistry = [
    JiraPlugin,
    GitHubPlugin,
    GitLabPlugin,
    JenkinsPlugin,
    TapdPlugin,
    AePlugin,
    AzurePlugin,
    BitbucketPlugin,
    GiteePlugin
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
      integrations
        .map((P) => P.id)
        .reduce(
          (pV, cV, iDx) => ({ ...pV, [cV.toString()?.toUpperCase()]: cV }),
          {}
        ),
    [integrations]
  )
  const ProviderLabels = useMemo(
    () =>
      integrations
        .map((P) => P)
        .reduce(
          (pV, cV, iDx) => ({
            ...pV,
            [cV?.id.toString()?.toUpperCase()]: cV.name
          }),
          {}
        ),
    [integrations]
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
  const ProviderIcons = useMemo(
    () =>
      integrations
        .map((p) => p.icon)
        .reduce((pV, cV, iDx) => ({ ...pV, [integrations[iDx]?.id]: cV }), {}),
    [integrations]
  )
  const ProviderConnectionLimits = useMemo(() => {}, [])

  const registerPlugin = useCallback((pluginConfig) => {
    console.log(
      '>>> REGISTERING PLUGIN...',
      `${pluginConfig?.name} [${pluginConfig?.type}]`
    )
    // @todo: Validate Plugin before Registration
    return new Plugin(pluginConfig)
  }, [])

  const validatePlugin = useCallback(() => {
    // @todo: Validate JSON Syntax
    // @todo: Validate Required Plugin Property Keys Exist
    // @todo: Validate Property data types
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
      ...registry.map((p) => registerPlugin(p))
    ])
  }, [registry, setPlugins, registerPlugin])

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
      '>>> INTEGRATIONS HOOK: PROVIDER CONFIGURATION CONNECTION FORM LABELS..',
      ProviderFormLabels
    )
    console.log(
      '>>> INTEGRATIONS HOOK: PROVIDER CONFIGURATION CONNECTION  FORM PLACEHOLDERS..',
      ProviderFormPlaceholders
    )
    console.log(
      '>>> INTEGRATIONS HOOK: PROVIDER CONFIGURATION PROVIDER ICONS..',
      ProviderIcons
    )
    console.log(
      '>>> INTEGRATIONS HOOK: PROVIDER DATA SOURCES LIST...',
      DataSources
    )
  }, [
    activeProvider,
    Providers,
    ProviderFormLabels,
    ProviderFormPlaceholders,
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
    ProviderIcons,
    ProviderConnectionLimits,
    ProviderTypes,
    setRegistry,
    setApiRegistry,
    setActiveProvider,
    getPlugin
  }
}

export default useIntegrations
