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
import { useCallback, useEffect, useMemo, useState } from 'react'
import { BlueprintMode } from '@/data/NullBlueprint'
import { integrationsData } from '@/data/integrations'
import JiraBoard from '@/models/JiraBoard'
import GitHubProject from '@/models/GithubProject'
import GitlabProject from '@/models/GitlabProject'
import { DataScopeModes } from '@/data/DataScopes'
import JenkinsJob from '@/models/JenkinsJob'
import useIntegrations from '@/hooks/useIntegrations'
import useTransformationsManager from '@/hooks/data-scope/useTransformationsManager'

function useDataScopesManager({
  mode = DataScopeModes.CREATE,
  provider,
  blueprint,
  /* connection, */ settings = {},
  setSettings = () => {}
}) {
  const {
    integrations: Integrations,
    Providers,
    ProviderLabels,
    ProviderIcons
  } = useIntegrations()

  const [connections, setConnections] = useState([])
  const [newConnections, setNewConnections] = useState([])

  const [scopeConnection, setScopeConnection] = useState()
  const [configuredConnection, setConfiguredConnection] = useState()
  const connection = useMemo(
    () =>
      mode === DataScopeModes.EDIT ? scopeConnection : configuredConnection,
    [scopeConnection, configuredConnection, mode]
  )
  const [scopeEntitiesGroup, setScopeEntitiesGroup] = useState([])

  // this is son set of ALL_DATA_DOMAINS
  // https://devlake.apache.org/docs/DataModels/DevLakeDomainLayerSchema/
  const [dataDomainsGroup, setDataDomainsGroup] = useState({})
  const {
    getTransformation,
    getTransformationScopeOptions,
    changeTransformationSettings,
    initializeDefaultTransformation,
    clearTransformationSettings,
    hasTransformationChanged
  } = useTransformationsManager()
  const [enabledProviders, setEnabledProviders] = useState([])

  const [configuredScopeEntity, setConfiguredScopeEntity] = useState(null)

  const activeTransformation = useMemo(
    () =>
      getTransformation(
        connection?.providerId,
        connection?.id,
        configuredScopeEntity
      ),
    [
      connection?.providerId,
      connection?.id,
      configuredScopeEntity,
      getTransformation
    ]
  )

  const createProviderScopes = useCallback(
    (providerId, connection, connectionIdx, dataDomainsGroup = {}) => {
      console.log(
        '>>> DATA SCOPES MANAGER: CREATING PROVIDER SCOPE FOR CONNECTION...',
        connectionIdx,
        scopeEntitiesGroup,
        connection
      )
      let newScope = {
        // FIXME: entities kept here because it have saved in db like Line#389
        entities: dataDomainsGroup[connection.id]?.map((d) => d.value) || []
      }
      // Generate scopes Dynamically for all Project/Board/Job/... Entities
      newScope =
        scopeEntitiesGroup[connection.id]?.map((e) => ({
          ...newScope,
          options: {
            ...getTransformationScopeOptions(connection?.providerId, e)
          },
          transformation: {
            ...getTransformation(connection?.providerId, connection?.id, e)
          }
        })) || []
      // switch (providerId) {
      //   case Providers.JIRA:
      //     newScope = boards[connection.id]?.map((b) => ({
      //       ...newScope,
      //       options: {
      //         boardId: Number(b?.value),
      //         title: b.title
      //         // @todo: verify initial value of since date for jira provider
      //         // since: new Date(),
      //       },
      //       transformation: {
      //         ...getTransformation(connection?.providerId, connection?.id, b)
      //       }
      //     }))
      //     break
      //   case Providers.GITLAB:
      //     newScope = projects[connection.id]?.map((p) => ({
      //       ...newScope,
      //       options: {
      //         projectId: Number(p.value),
      //         title: p.title
      //       },
      //       transformation: {
      //         ...getTransformation(connection?.providerId, connection?.id, p)
      //       }
      //     }))
      //     break
      //   case Providers.JENKINS:
      //     newScope = projects[connection.id]?.map((p) => ({
      //       ...newScope,
      //       options: {
      //         jobName: p.value
      //       },
      //       transformation: {
      //         ...getTransformation(connection?.providerId, connection?.id, p)
      //       }
      //     }))
      //     break
      //   case Providers.GITHUB:
      //     newScope = projects[connection.id]?.map((p) => ({
      //       ...newScope,
      //       options: {
      //         owner: p.value.split('/')[0],
      //         repo: p.value.split('/')[1]
      //       },
      //       transformation: {
      //         ...getTransformation(connection?.providerId, connection?.id, p)
      //       }
      //     }))
      //     break
      //   case Providers.TAPD:
      //     newScope = {
      //       ...newScope
      //       // options: {
      //       // },
      //       // transformation: {},
      //     }
      //     break
      // }
      return Array.isArray(newScope) ? newScope.flat() : [newScope]
    },
    [
      scopeEntitiesGroup,
      getTransformation,
      getTransformationScopeOptions
      // Providers
    ]
  )

  const createProviderConnections = useCallback(
    (blueprintConnections) => {
      console.log(
        '>>>>> DATA SCOPES MANAGER: Creating Provider Connection Scopes...',
        blueprintConnections
      )
      return blueprintConnections.map((c, cIdx) => ({
        connectionId: c.value || c.connectionId,
        plugin: c.plugin || c.provider,
        scope: createProviderScopes(
          typeof c.provider === 'object' ? c.provider?.id : c.provider,
          c,
          cIdx,
          dataDomainsGroup
        )
      }))
    },
    [dataDomainsGroup, createProviderScopes]
  )

  const modifyConnectionSettings = useCallback(() => {
    const newConnections = createProviderConnections(
      connections.filter((c) => c.providerId === connection?.providerId)
    )
    const existingConnections = blueprint?.settings?.connections?.filter(
      (storedConnection) => storedConnection.plugin !== connection?.plugin
    )
    console.log(
      '>>>>> DATA SCOPES MANAGER: Modifying Connection Scopes...',
      newConnections
    )
    console.log(
      '>>>>> DATA SCOPES MANAGER: Filtered Existing connection Scopes...',
      existingConnections
    )
    setSettings((currentSettings) => ({
      ...currentSettings,
      connections: [...newConnections, ...existingConnections]
    }))
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [
    connection,
    connections,
    blueprint?.settings?.connections,
    setSettings,
    createProviderConnections
  ])

  const getGithubProjects = useCallback(
    (c) =>
      c.scope.map(
        (s) =>
          new GitHubProject({
            id: `${s.options?.owner}/${s.options?.repo}`,
            key: `${s.options?.owner}/${s.options?.repo}`,
            owner: s.options?.owner,
            repo: s.options?.repo,
            value: `${s.options?.owner}/${s.options?.repo}`,
            title: `${s.options?.owner}/${s.options?.repo}`
          })
      ),
    []
  )

  const getGitlabProjects = useCallback(
    (c) =>
      c.scope.map(
        (s) =>
          new GitlabProject({
            id: s.options?.projectId,
            key: s.options?.projectId,
            value: s.options?.projectId,
            title: s.options?.title || `Project ${s.options?.projectId}`
          })
      ),
    []
  )

  const getJenkinsProjects = useCallback(
    (c) =>
      // when s.options?.jobName is empty, it's old jenkins config which collect all job data
      c.scope
        .filter((s) => s.options?.jobName)
        .map(
          (s) =>
            new JenkinsJob({
              id: s.options.jobName,
              key: s.options.jobName,
              value: s.options.jobName,
              title: s.options.jobName
            })
        ),
    []
  )

  const getJiraBoard = useCallback(
    (c) =>
      c.scope.map(
        (s) =>
          new JiraBoard({
            id: s.options?.boardId,
            key: s.options?.boardId,
            value: s.options?.boardId,
            title: s.options?.title || `Board ${s.options?.boardId}`
          })
      ),
    []
  )

  const getProjects = useCallback(
    (c) => {
      switch (c.plugin) {
        case Providers.GITHUB:
          return getGithubProjects(c)
        case Providers.GITLAB:
          return getGitlabProjects(c)
        case Providers.JENKINS:
          return getJenkinsProjects(c)
        case Providers.JIRA:
          return getJiraBoard(c)
        default:
          return []
      }
    },
    [
      getGithubProjects,
      getGitlabProjects,
      getJenkinsProjects,
      getJiraBoard,
      Providers
    ]
  )

  const getAdvancedGithubProjects = useCallback(
    (t) => [
      new GitHubProject({
        id: `${t.options?.owner}/${t.options?.repo}`,
        key: `${t.options?.owner}/${t.options?.repo}`,
        owner: s.options?.owner,
        repo: s.options?.repo,
        value: `${t.options?.owner}/${t.options?.repo}`,
        title: `${t.options?.owner}/${t.options?.repo}`
      })
    ],
    []
  )

  const getAdvancedGitlabProjects = useCallback(
    (t) => [
      new GitlabProject({
        id: t.options?.projectId,
        key: t.options?.projectId,
        value: t.options?.projectId,
        title: `Project ${t.options?.projectId}`
      })
    ],
    []
  )

  const getAdvancedJenkinsProjects = useCallback(
    (t) => [
      new JenkinsJob({
        id: t.options?.jobName,
        key: t.options?.jobName,
        value: t.options?.jobName,
        title: t.options?.jobName
      })
    ],
    []
  )

  const getAdvancedJiraBoards = useCallback(
    (t) => [
      new JiraBoard({
        id: t.options?.boardId,
        key: t.options?.boardId,
        value: t.options?.boardId,
        title: `Board ${t.options?.boardId}`
      })
    ],
    []
  )

  const getDefaultDataDomains = useCallback(
    (providerId) => {
      console.log('GET ENTITIES FOR PROVIDER =', providerId)
      const plugin = Integrations.find((p) => p.id === providerId)
      return plugin ? plugin.getAvailableDataDomains() : []
    },
    [Integrations]
  )

  const getAdvancedScopeEntity = useCallback(
    (c) => {
      switch (c.plugin) {
        case Providers.GITHUB:
          return getAdvancedGithubProjects(c)
        case Providers.GITLAB:
          return getAdvancedGitlabProjects(c)
        case Providers.JENKINS:
          return getAdvancedJenkinsProjects(c)
        case Providers.JIRA:
          return getAdvancedJiraBoards(c)
        default:
          return []
      }
    },
    [
      getAdvancedGithubProjects,
      getAdvancedGitlabProjects,
      getAdvancedJiraBoards,
      getAdvancedJenkinsProjects,
      Providers.GITHUB,
      Providers.GITLAB,
      Providers.JENKINS,
      Providers.JIRA
    ]
  )

  const createNormalConnection = useCallback(
    (
      blueprint,
      c,
      cIdx,
      ALL_DATA_DOMAINS,
      // FIXME: what different between connections and connectionsList ?
      connections = [],
      connectionsList = []
    ) => ({
      ...c,
      mode: BlueprintMode.NORMAL,
      // @IMPORTANT: Preserve Original LIST INDEX ID!
      id: connectionsList.find(
        (lC) => lC.value === c.connectionId && lC.provider === c.plugin
      )?.id,
      connectionId: c.connectionId,
      value: c.connectionId,
      provider: integrationsData.find((i) => i.id === c.plugin),
      providerLabel: ProviderLabels[c.plugin?.toUpperCase()],
      providerId: c.plugin,
      plugin: c.plugin,
      icon: ProviderIcons[c.plugin] ? ProviderIcons[c.plugin](18, 18) : null,
      name:
        connections.find(
          (pC) => pC.connectionId === c.connectionId && pC.plugin === c.plugin
        )?.name ||
        `${ProviderLabels[c.plugin?.toUpperCase()]} #${c.connectionId || cIdx}`,
      // FIXME: entities in `c.scope[0]?.entities` means one of ALL_DATA_DOMAINS and is saved in db,
      // So it kept here.
      dataDomains: c.scope[0]?.entities?.map((e) =>
        ALL_DATA_DOMAINS.find((de) => de.value === e)
      ),
      scopeEntities: getProjects(c),
      transformations: c.scope.map((s) => ({ ...s.transformation })),
      transformationStates: c.scope.map((s) =>
        Object.values(s.transformation ?? {}).some((v) =>
          Array.isArray(v)
            ? v.length > 0
            : v && typeof v === 'object'
            ? Object.keys(v)?.length > 0
            : v?.toString().length > 0
        )
          ? 'Added'
          : '-'
      ),
      scope: c.scope,
      // editable: ![Providers.JENKINS].includes(c.plugin),
      editable: true,
      advancedEditable: false,
      isMultiStage: false,
      isSingleStage: true,
      stage: 1,
      totalStages: 1
    }),
    [getProjects, ProviderLabels, ProviderIcons]
  )

  const createAdvancedConnection = useCallback(
    (
      blueprint,
      c,
      cIdx,
      ALL_DATA_DOMAINS,
      connections = [],
      connectionsList = []
    ) => ({
      ...c,
      mode: BlueprintMode.ADVANCED,
      // @IMPORTANT: Preserve Original LIST INDEX ID!
      id: connectionsList.find(
        (lC) => lC.value === c.options?.connectionId && lC.provider === c.plugin
      )?.id,
      connectionId: c.options?.connectionId,
      value: c.options?.connectionId,
      provider: integrationsData.find((i) => i.id === c.plugin),
      providerLabel: ProviderLabels[c.plugin?.toUpperCase()],
      plugin: c.plugin,
      providerId: c.plugin,
      icon: ProviderIcons[c.plugin] ? ProviderIcons[c.plugin](18, 18) : null,
      name:
        connections.find(
          (pC) =>
            pC.connectionId === c.options?.connectionId &&
            pC.provider === c.plugin
        )?.name || `Connection ID #${c.options?.connectionId || cIdx}`,
      scopeEntities: getAdvancedScopeEntity(c),
      dataDomains: getDefaultDataDomains(c.plugin),
      transformations: [],
      transformationStates:
        typeof c.options?.transformationRules === 'object' &&
        Object.values(c.options?.transformationRules || {}).some(
          (v) => (Array.isArray(v) && v.length > 0) || v.toString().length > 0
        )
          ? ['Added']
          : ['-'],
      scope: c,
      task: c,
      editable: false,
      advancedEditable: true,
      plan: blueprint?.plan,
      isMultiStage:
        Array.isArray(blueprint?.plan) && blueprint?.plan.length > 1,
      isSingleStage:
        Array.isArray(blueprint?.plan) && blueprint?.plan.length === 1,
      stage:
        blueprint?.plan.findIndex((s, sId) =>
          s.find((t) => JSON.stringify(t) === JSON.stringify(c))
        ) + 1,
      totalStages: blueprint?.plan?.length
    }),
    [
      getAdvancedScopeEntity,
      getDefaultDataDomains,
      ProviderIcons,
      ProviderLabels
    ]
  )

  const hasConfiguredEntityTransformationChanged = useCallback(
    (item) => {
      return hasTransformationChanged(
        configuredConnection?.provider,
        configuredConnection?.id,
        item
      )
    },
    [
      configuredConnection?.provider,
      configuredConnection?.id,
      hasTransformationChanged
    ]
  )

  const changeConfiguredEntityTransformation = useCallback(
    (settings) => {
      return changeTransformationSettings(
        configuredConnection?.provider,
        configuredConnection?.id,
        configuredScopeEntity,
        settings
      )
    },
    [
      configuredConnection?.provider,
      configuredConnection?.id,
      configuredScopeEntity,
      changeTransformationSettings
    ]
  )

  useEffect(() => {
    console.log('>>>>> DATA SCOPES MANAGER: CONFIGURED CONNECTION', connection)
    switch (connection?.provider?.id) {
      case Providers.GITHUB:
      case Providers.GITLAB:
      case Providers.JENKINS:
      case Providers.JIRA:
        setScopeEntitiesGroup((g) => ({
          ...g,
          [connection?.id]: connection?.scopeEntities || []
        }))
        setDataDomainsGroup((e) => ({
          ...e,
          [connection?.id]: connection?.dataDomains || []
        }))
        connection?.scopeEntities.forEach((p, pIdx) =>
          changeTransformationSettings(
            connection?.provider?.id,
            connection?.id,
            p,
            connection.transformations[pIdx]
          )
        )
        break
    }
  }, [connection, changeTransformationSettings, Providers])

  useEffect(() => {
    console.log('>>>>> DATA SCOPES MANAGER: Connection List...', connections)
    modifyConnectionSettings()
  }, [
    connections,
    dataDomainsGroup,
    scopeEntitiesGroup,
    modifyConnectionSettings
  ])

  useEffect(() => {
    console.log(
      '>>>>> DATA SCOPES MANAGER: INITIALIZE SCOPE ENTITIES...',
      scopeEntitiesGroup
    )
    const scopeEntities = scopeEntitiesGroup[connection?.id]
    if (Array.isArray(scopeEntities)) {
      for (const scopeEntity of scopeEntities) {
        initializeDefaultTransformation(
          connection?.providerId,
          connection?.id,
          scopeEntity
        )
      }
    }
  }, [
    connection?.providerId,
    scopeEntitiesGroup,
    connection?.id,
    initializeDefaultTransformation
  ])

  useEffect(() => {
    console.log('>>>>> DATA SCOPES MANAGER: DATA ENTITIES...', dataDomainsGroup)
  }, [dataDomainsGroup])

  useEffect(() => {
    console.log(
      '>>>>> DATA SCOPES MANAGER: CURRENT BLUEPRINT SETTINGS...',
      settings
    )
  }, [settings])

  useEffect(() => {
    console.log(
      '>>>>> DATA SCOPES MANAGER: ACTIVE TRANSFORMATION RULES...',
      activeTransformation
    )
  }, [activeTransformation])

  useEffect(() => {
    console.log(
      '>>>>> DATA SCOPES MANAGER: MEMOIZED ACTIVE CONNECTION...',
      connection
    )
  }, [connection])

  useEffect(() => {
    console.log(
      '>>>>> DATA SCOPES MANAGER: SELECTED NEW CONNECTIONS...',
      newConnections
    )
  }, [newConnections])

  return {
    connections,
    newConnections,
    // blueprint,
    scopeEntitiesGroup,
    dataDomainsGroup,
    configuredConnection,
    configuredScopeEntity,
    activeTransformation,
    scopeConnection,
    enabledProviders,
    setNewConnections,
    setConnections,
    setScopeConnection,
    setConfiguredConnection,
    setConfiguredScopeEntity,
    setScopeEntitiesGroup,
    setDataDomainsGroup,
    getTransformation,
    changeTransformationSettings,
    initializeDefaultTransformation,
    clearTransformationSettings,
    hasTransformationChanged,
    hasConfiguredEntityTransformationChanged,
    changeConfiguredEntityTransformation,
    createProviderConnections,
    createProviderScopes,
    getDefaultDataDomains,
    createNormalConnection,
    createAdvancedConnection,
    setEnabledProviders
  }
}

export default useDataScopesManager
