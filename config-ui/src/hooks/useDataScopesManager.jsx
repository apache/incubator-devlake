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
import { useCallback, useEffect, useState, useMemo } from 'react'
import { ToastNotification } from '@/components/Toast'
import { DEVLAKE_ENDPOINT } from '@/utils/config'
import request from '@/utils/request'
import { NullBlueprint, BlueprintMode } from '@/data/NullBlueprint'
import { Providers } from '@/data/Providers'

function useDataScopesManager ({ provider, blueprint, /* connection, */ settings = {}, setSettings = () => {} }) {
  const [connections, setConnections] = useState([])

  const [scopeConnection, setScopeConnection] = useState()
  const connection = useMemo(() => scopeConnection, [scopeConnection])

  // const [blueprint, setBlueprint] = useState(NullBlueprint)
  const [boards, setBoards] = useState({})
  const [projects, setProjects] = useState({})
  const [entities, setEntities] = useState({})
  const [transformations, setTransformations] = useState({})
  const [enabledProviders, setEnabledProviders] = useState([])

  const [configuredProject, setConfiguredProject] = useState(null)
  const [configuredBoard, setConfiguredBoard] = useState(null)

  const activeProject = useMemo(() => configuredProject, [configuredProject])
  const activeBoard = useMemo(() => configuredBoard, [configuredBoard])

  // @todo: fix check why these are empty
  const selectedProjects = useMemo(() => projects[connection?.id] || [], [projects, connection?.id])
  const selectedBoards = useMemo(() => boards[connection?.id]?.map(
    (b) => b && b?.id
  ), [boards, connection?.id])

  const storedProjectTransformation = useMemo(() => connection?.transformations[connection?.projects?.findIndex(p => p === configuredProject)], [connection, configuredProject])
  const storedBoardTransformation = useMemo(() => connection?.transformations[connection?.boardIds?.findIndex(b => b === configuredBoard?.id)], [connection, configuredBoard?.id])

  const activeProjectTransformation = useMemo(() => transformations[activeProject], [transformations, activeProject])
  const activeBoardTransformation = useMemo(() => transformations[activeBoard?.id], [transformations, activeBoard?.id])
  const activeTransformation = useMemo(() => transformations[connection?.providerId === Providers.JIRA ? configuredBoard?.id : configuredProject], [transformations, configuredProject, configuredBoard?.id, connection?.providerId])

  const getDefaultTransformations = useCallback((providerId) => {
    let transforms = {}
    switch (providerId) {
      case Providers.GITHUB:
        transforms = {
          prType: '',
          prComponent: '',
          prBodyClosePattern: '',
          issueSeverity: '',
          issueComponent: '',
          issuePriority: '',
          issueTypeRequirement: '',
          issueTypeBug: '',
          issueTypeIncident: '',
          refdiff: null,
        }
        break
      case Providers.JIRA:
        transforms = {
          epicKeyField: '',
          typeMappings: {},
          storyPointField: '',
          remotelinkCommitShaPattern: '',
          bugTags: [],
          incidentTags: [],
          requirementTags: [],
        }
        break
      case Providers.JENKINS:
        // No Transform Settings...
        break
      case Providers.GITLAB:
        // No Transform Settings...
        break
      case Providers.TAPD:
        // No Transform Settings...
        break
    }
    console.log('>>>>> DATA SCOPES MANAGER: Getting Default Transformation Values for PROVIDER Type ', providerId, transforms)
    return transforms
  }, [])

  const initializeTransformations = useCallback((pV, cV, iDx) => ({
    ...pV,
    [cV]: getDefaultTransformations(connection?.providerId, iDx),
  }), [connection, getDefaultTransformations])

  const createProviderScopes = useCallback(
    (
      providerId,
      connection,
      connectionIdx,
      entities = {},
      boards = {},
      projects = {},
      transformations = {},
      defaultScope = { transformation: {}, options: {}, entities: [] }
    ) => {
      console.log(
        '>>> DATA SCOPES MANAGER: CREATING PROVIDER SCOPE FOR CONNECTION...',
        connectionIdx,
        connection
      )
      let newScope = {
        ...defaultScope,
        entities: entities[connection.id]?.map((entity) => entity.value) || [],
      }
      switch (providerId) {
        case Providers.JIRA:
          newScope = boards[connection.id]?.map((b) => ({
            ...newScope,
            options: {
              boardId: Number(b?.id),
              // @todo: verify initial value of since date for jira provider
              // since: new Date(),
            },
            transformation: { ...transformations[b?.id] },
          }))
          break
        case Providers.GITLAB:
          newScope = projects[connection.id]?.map((p) => ({
            ...newScope,
            options: {
              projectId: Number(p),
            },
            transformation: {},
          }))
          break
        case Providers.JENKINS:
          newScope = {
            ...newScope,
            // options: {
            // },
            // transformation: {},
          }
          break
        case Providers.GITHUB:
          newScope = projects[connection.id]?.map((p) => ({
            ...newScope,
            options: {
              owner: p.split('/')[0],
              repo: p.split('/')[1],
            },
            transformation: { ...transformations[p] },
          }))
          break
        case Providers.TAPD:
          newScope = {
            ...newScope,
            // options: {
            // },
            // transformation: {},
          }
          break
      }
      return Array.isArray(newScope) ? newScope.flat() : [newScope]
    },
    []
  )

  const createProviderConnections = useCallback((blueprintConnections) => {
    console.log('>>>>> DATA SCOPES MANAGER: Creating Provider Connection Scopes...', blueprintConnections)
    return blueprintConnections.map((c, cIdx) => ({
      connectionId: c.value || c.connectionId,
      plugin: c.plugin || c.provider,
      scope: createProviderScopes(
        typeof c.provider === 'object' ? c.provider?.id : c.provider,
        c,
        cIdx,
        entities,
        boards,
        projects,
        transformations
      ),
    }))
  }, [boards, projects, entities, transformations, createProviderScopes])

  const modifyConnectionSettings = useCallback(() => {
    const newConnections = createProviderConnections(connections.filter(c => c.providerId === connection?.providerId))
    const existingConnections = blueprint?.settings?.connections?.filter(storedConnection => storedConnection.plugin !== connection?.plugin)
    console.log('>>>>> DATA SCOPES MANAGER: Modifying Connection Scopes...', newConnections)
    console.log('>>>>> DATA SCOPES MANAGER: Filtered Existing connection Scopes...', existingConnections)
    setSettings((currentSettings) => ({
      ...currentSettings,
      connections: [
        ...newConnections,
        ...existingConnections
      ],
    }))
  }, [
    connection,
    connections,
    boards,
    projects,
    entities,
    transformations,
    blueprint?.settings?.connections,
    setSettings,
    createProviderConnections
  ])

  const setTransformationSettings = useCallback(
    (settings, configuredEntity) => {
      console.log(
        '>>>>> DATA SCOPES MANAGER: SETTING TRANSFORMATION SETTINGS PROJECT/BOARD...',
        configuredEntity,
        settings
      )
      setTransformations((existingTransformations) => ({
        ...existingTransformations,
        [configuredEntity]: {
          ...existingTransformations[configuredEntity],
          ...settings,
        },
      }))
    },
    [setTransformations]
  )

  useEffect(() => {
    console.log('>>>>> DATA SCOPES MANAGER: INITIALIZING TRANSFORMATION RULES...', selectedProjects)
    console.log('>>>>> DATA SCOPES MANAGER: INITIALIZING TRANSFORMATION RULES...', selectedBoards)
  }, [selectedProjects, selectedBoards, initializeTransformations])

  useEffect(() => {
    console.log('>>>>> DATA SCOPES MANAGER: CONFIGURED CONNECTION', connection)
    switch (connection?.provider?.id) {
      case Providers.GITHUB:
      case Providers.GITLAB:
        setProjects(p => ({ ...p, [connection?.id]: connection?.projects || [] }))
        setEntities(e => ({ ...e, [connection?.id]: connection?.entityList || [] }))
        connection?.projects.forEach((p, pIdx) => setTransformationSettings(connection.transformations[pIdx], p))
        break
      case Providers.JIRA:
        // fetchBoards()
        // fetchIssueTypes()
        // fetchFields()
        setBoards(b => ({ ...b, [connection?.id]: connection?.boardsList || [] }))
        setEntities(e => ({ ...e, [connection?.id]: connection?.entityList || [] }))
        connection?.boardIds.forEach((bId, bIdx) => setTransformationSettings(connection.transformations[bIdx], bId))
        break
      case Providers.JENKINS:
        setEntities(e => ({ ...e, [connection?.id]: connection?.entityList || [] }))
        break
    }
  }, [
    connection,
    setTransformationSettings
  ])

  useEffect(() => {
    console.log('>>>>> DATA SCOPES MANAGER: Connection List...', connections)
    modifyConnectionSettings()
  }, [
    connections,
    entities,
    projects,
    boards,
    transformations,
    modifyConnectionSettings
  ])

  useEffect(() => {
    console.log('>>>>> DATA SCOPES MANAGER: PROVIDER...', provider)
    switch (provider?.id) {
      case Providers.GITHUB:
        break
      case Providers.GITLAB:
        break
      case Providers.JIRA:
        break
      case Providers.JENKINS:
        break
      case Providers.TAPD:
        break
    }
  }, [provider])

  useEffect(() => {
    console.log('>>>>> DATA SCOPES MANAGER: BOARDS...', boards)
    const boardTransformations = boards[connection?.id]
    if (Array.isArray(boardTransformations) && boardTransformations?.length > 0) {
      // @todo: check & re-enable
      setTransformations((cT) => ({
        ...boardTransformations.reduce(initializeTransformations, {}),
        // Spread Current/Existing Transformations Settings
        ...cT,
      }))
    }
  }, [boards, connection?.id, initializeTransformations])

  useEffect(() => {
    console.log('>>>>> DATA SCOPES MANAGER: PROJECTS...', projects)
    const projectTransformations = projects[connection?.id]
    if (Array.isArray(projectTransformations)) {
      setTransformations((cT) => ({
        ...projectTransformations.reduce(initializeTransformations, {}),
        // Spread Current/Existing Transformations Settings
        ...cT,
      }))
    }
  }, [projects, connection?.id, initializeTransformations])

  useEffect(() => {
    console.log('>>>>> DATA SCOPES MANAGER: DATA ENTITIES...', entities)
  }, [entities])

  useEffect(() => {
    console.log('>>>>> DATA SCOPES MANAGER: TRANSFORMATIONS...', transformations)
  }, [transformations])

  useEffect(() => {
    console.log('>>>>> DATA SCOPES MANAGER: CURRENT BLUEPRINT SETTINGS...', settings)
  }, [settings])

  useEffect(() => {
    console.log('>>>>> DATA SCOPES MANAGER: ACTIVE TRANSFORMATION RULES...', activeTransformation)
  }, [activeTransformation])

  useEffect(() => {
    console.log('>>>>> DATA SCOPES MANAGER: ACTIVE PROJECT TRANSFORMATION RULES...', activeProjectTransformation)
  }, [activeProjectTransformation])

  useEffect(() => {
    console.log('>>>>> DATA SCOPES MANAGER: ACTIVE BOARD TRANSFORMATION RULES...', activeBoardTransformation)
  }, [activeBoardTransformation])

  return {
    connections,
    // blueprint,
    boards,
    projects,
    entities,
    transformations,
    configuredBoard,
    configuredProject,
    storedProjectTransformation,
    storedBoardTransformation,
    activeBoardTransformation,
    activeProjectTransformation,
    activeTransformation,
    scopeConnection,
    enabledProviders,
    // setActiveTransformation,
    setConnections,
    setScopeConnection,
    setConfiguredBoard,
    setConfiguredProject,
    // setBlueprint,
    setBoards,
    setProjects,
    setEntities,
    setTransformations,
    setTransformationSettings,
    initializeTransformations,
    getDefaultTransformations,
    createProviderConnections,
    createProviderScopes,
    modifyConnectionSettings,
    setEnabledProviders
  }
}

export default useDataScopesManager
