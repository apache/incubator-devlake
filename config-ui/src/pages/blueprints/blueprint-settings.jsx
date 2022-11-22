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
import React, { useEffect, useState, useCallback, useRef, useMemo } from 'react'
import { useParams, useHistory } from 'react-router-dom'
import { ENVIRONMENT } from '@/config/environment'
import dayjs from '@/utils/time'
import {
  JIRA_API_PROXY_ENDPOINT,
  ISSUE_TYPES_ENDPOINT,
  ISSUE_FIELDS_ENDPOINT,
  BOARDS_ENDPOINT
} from '@/config/jiraApiProxy'
import request from '@/utils/request'
import {
  Button,
  Elevation,
  Intent,
  Switch,
  Tag,
  Tooltip,
  Icon,
  Colors,
  Spinner,
  Classes,
  Popover
} from '@blueprintjs/core'

import { integrationsData } from '@/data/integrations'
import JiraBoard from '@/models/JiraBoard'
import DataScopeConnection from '@/models/DataScopeConnection'
import { NullBlueprint, BlueprintMode } from '@/data/NullBlueprint'
import { NullPipelineRun } from '@/data/NullPipelineRun'
import { Providers, ProviderLabels, ProviderIcons } from '@/data/Providers'
import { TaskStatus } from '@/data/Task'

import Nav from '@/components/Nav'
import Sidebar from '@/components/Sidebar'
import Content from '@/components/Content'
import { ToastNotification } from '@/components/Toast'
import BlueprintNameCard from '@/components/blueprints/BlueprintNameCard'
import DataSync from '@/components/blueprints/create-workflow/DataSync'
import CodeInspector from '@/components/pipelines/CodeInspector'

// import { DataEntities, DataEntityTypes } from '@/data/DataEntities'
import { DEFAULT_DATA_ENTITIES } from '@/data/BlueprintWorkflow'
import { DataScopeModes } from '@/data/DataScopes'

import useBlueprintManager from '@/hooks/useBlueprintManager'
import useConnectionManager from '@/hooks/useConnectionManager'
import usePipelineManager from '@/hooks/usePipelineManager'
import useDataScopesManager from '@/hooks/useDataScopesManager'
import useJIRA from '@/hooks/useJIRA'
import useBlueprintValidation from '@/hooks/useBlueprintValidation'
import usePipelineValidation from '@/hooks/usePipelineValidation'

import BlueprintDialog from '@/components/blueprints/BlueprintDialog'
import BlueprintDataScopesDialog from '@/components/blueprints/BlueprintDataScopesDialog'
import BlueprintNavigationLinks from '@/components/blueprints/BlueprintNavigationLinks'
import DataScopesGrid from '@/components/blueprints/DataScopesGrid'
import AdvancedJSON from '@/components/blueprints/create-workflow/AdvancedJSON'
import useGitlab from '@/hooks/useGitlab'
import {
  GITLAB_API_PROXY_ENDPOINT,
  PROJECTS_ENDPOINT
} from '@/config/gitlabApiProxy'

const BlueprintSettings = (props) => {
  // eslint-disable-next-line no-unused-vars
  const history = useHistory()
  const { bId } = useParams()

  const [activeProvider, setActiveProvider] = useState(integrationsData[0])
  // @disabled Provided By Data Scopes Manager
  // const [activeTransformation, setActiveTransformation] = useState()

  const [blueprintId, setBlueprintId] = useState()
  const [activeBlueprint, setActiveBlueprint] = useState(NullBlueprint)
  const [currentRun, setCurrentRun] = useState(NullPipelineRun)
  const [dataEntitiesList, setDataEntitiesList] = useState([
    ...DEFAULT_DATA_ENTITIES
  ])

  // @disabled Provided By Data Scopes Manager
  // const [connections, setConnections] = useState([])
  const [blueprintConnections, setBlueprintConnections] = useState([])
  // const [configuredConnection, setConfiguredConnection] = useState()

  // @todo: relocate or discard
  const [newConnectionScopes, setNewConnectionScopes] = useState({})

  const [blueprintDialogIsOpen, setBlueprintDialogIsOpen] = useState(false)
  const [blueprintScopesDialogIsOpen, setBlueprintScopesDialogIsOpen] =
    useState(false)

  const [activeSetting, setActiveSetting] = useState({
    id: null,
    title: '',
    payload: {}
  })

  const [showBlueprintInspector, setShowBlueprintInspector] = useState(false)
  const [runTasks, setRunTasks] = useState([])
  const [runTasksAdvanced, setRunTasksAdvanced] = useState([])

  const [boardSearch, setBoardSearch] = useState()

  const {
    // eslint-disable-next-line no-unused-vars
    activeStep,
    blueprint,
    name: blueprintName,
    cronConfig,
    customCronConfig,
    enable,
    tasks: blueprintTasks,
    settings: blueprintSettings,
    rawConfiguration,
    mode,
    interval,
    isSaving,
    isFetching: isFetchingBlueprint,
    activateBlueprint,
    deactivateBlueprint,
    getNextRunDate,
    // eslint-disable-next-line no-unused-vars
    fetchBlueprint,
    patchBlueprint,
    setName: setBlueprintName,
    setCronConfig,
    setCustomCronConfig,
    setEnable,
    setMode,
    setInterval,
    setIsManual,
    setTasks: setBlueprintTasks,
    setSettings: setBlueprintSettings,
    setRawConfiguration,
    createCron,
    getCronPreset,
    getCronPresetByConfig,
    detectCronInterval,
    fetchAllBlueprints,
    saveBlueprint,
    saveComplete,
    errors: blueprintErrors
  } = useBlueprintManager()

  const {
    connections,
    boards,
    projects,
    entities,
    transformations,
    scopeConnection,
    activeBoardTransformation,
    activeProjectTransformation,
    activeTransformation,
    configuredConnection,
    configuredBoard,
    configuredProject,
    configurationKey,
    enabledProviders,
    setConfiguredConnection,
    setConfiguredBoard,
    setConfiguredProject,
    setBoards,
    setProjects,
    setEntities,
    setTransformations,
    setTransformationSettings,
    // setActiveTransformation,
    setConnections,
    setScopeConnection,
    setEnabledProviders,
    createProviderConnections,
    createProviderScopes,
    createNormalConnection,
    createAdvancedConnection,
    getJiraMappedBoards,
    getDefaultEntities
  } = useDataScopesManager({
    mode: DataScopeModes.EDIT,
    blueprint: activeBlueprint,
    provider: activeProvider,
    // connection: scopeConnection,
    settings: blueprintSettings,
    setSettings: setBlueprintSettings
  })

  const {
    fetchConnection,
    allProviderConnections,
    connectionsList,
    isFetching: isFetchingConnection,
    fetchAllConnections
  } = useConnectionManager(
    {
      activeProvider,
      connectionId: configuredConnection?.connectionId
    },
    configuredConnection && configuredConnection?.id !== null
  )

  const {
    // eslint-disable-next-line no-unused-vars
    pipelineName,
    // pipelines,
    // runPipeline,
    // cancelPipeline,
    // fetchPipeline,
    // fetchAllPipelines,
    // pipelineRun,
    // buildPipelineStages,
    // isRunning,
    isFetchingAll: isFetchingAllPipelines,
    // eslint-disable-next-line no-unused-vars
    errors: pipelineErrors,
    setSettings: setPipelineSettings,
    setPipelineName,
    // eslint-disable-next-line no-unused-vars
    lastRunId,
    // eslint-disable-next-line no-unused-vars
    allowedProviders,
    // eslint-disable-next-line no-unused-vars
    detectPipelineProviders
  } = usePipelineManager(null, runTasks)

  const {
    validate: validateBlueprint,
    errors: blueprintValidationErrors,
    isValid: isValidBlueprint,
    fieldHasError,
    getFieldError,
    isValidCronExpression,
    isValidJSON,
    isValidConfiguration: isValidJSONConfiguration,
    validateAdvancedError,
    validateBlueprintName,
    validateRepositoryName,
    validateNumericSet
  } = useBlueprintValidation({
    name: blueprintName,
    boards,
    projects,
    cronConfig,
    customCronConfig,
    enable,
    tasks: blueprintTasks,
    mode,
    // connections: blueprintConnections,
    // entities: dataEntities,
    activeStep
    // activeProvider: provider,
    // activeConnection: configuredConnection
  })

  const {
    validate: validatePipeline,
    validateAdvanced: validateAdvancedPipeline,
    errors: pipelineValidationErrors,
    setErrors: setPipelineErrors,
    isValid: isValidPipeline,
    detectedProviders,
    parseJSON
  } = usePipelineValidation({
    enabledProviders,
    pipelineName: activeBlueprint?.name,
    projects,
    boards,
    connectionId: scopeConnection?.id,
    tasks: runTasks,
    tasksAdvanced: runTasksAdvanced,
    advancedMode: activeBlueprint?.mode === BlueprintMode.ADVANCED,
    mode,
    connection: configuredConnection,
    entities: entities,
    rawConfiguration
  })

  const {
    fetchIssueTypes,
    fetchFields,
    fetchBoards,
    fetchAllResources,
    allResources: allJiraResources,
    boards: jiraApiBoards,
    issueTypes: jiraApiIssueTypes,
    fields: jiraApiFields,
    isFetching: isFetchingJIRA,
    error: jiraProxyError
  } = useJIRA(
    {
      apiProxyPath: JIRA_API_PROXY_ENDPOINT,
      issuesEndpoint: ISSUE_TYPES_ENDPOINT,
      fieldsEndpoint: ISSUE_FIELDS_ENDPOINT,
      boardsEndpoint: BOARDS_ENDPOINT
    },
    configuredConnection
  )

  const {
    fetchProjects: fetchGitlabProjects,
    projects: gitlabProjects,
    isFetching: isFetchingGitlab,
    error: gitlabProxyError
  } = useGitlab(
    {
      apiProxyPath: GITLAB_API_PROXY_ENDPOINT,
      projectsEndpoint: PROJECTS_ENDPOINT
    },
    configuredConnection
  )

  const handleBlueprintActivation = useCallback(
    (blueprint) => {
      if (blueprint.enable) {
        deactivateBlueprint(blueprint)
      } else {
        activateBlueprint(blueprint)
      }
    },
    [activateBlueprint, deactivateBlueprint]
  )

  const handleBlueprintDialogClose = useCallback(() => {
    console.log('>>> CLOSING BLUEPRINT DIALOG & REVERTING SETTINGS...')
    setBlueprintDialogIsOpen(false)
    setBlueprintName(activeBlueprint?.name)
  }, [activeBlueprint, setBlueprintName])

  const handleBlueprintScopesDialogClose = useCallback(() => {
    console.log(
      '>>> CLOSING BLUEPRINT CONNECTION SCOPES DIALOG & REVERTING SETTINGS...'
    )
    setBlueprintScopesDialogIsOpen(false)
    setScopeConnection(null)
    // restore/revert data scope + settings on close (cancel)
    // setTransformations({})
    // switch (scopeConnection?.provider?.id) {
    //   case Providers.GITHUB:
    //     setActiveTransformation(activeProjectTransformation)
    //     setProjects(p => ({ ...p, [configuredConnection?.id]: scopeConnection?.projects }))
    //     setEntities(e => ({ ...e, [configuredConnection?.id]: scopeConnection?.entityList }))
    //     break
    //   case Providers.JIRA:
    //     setActiveTransformation(activeBoardTransformation)
    //     setBoards(b => ({...b, [configuredConnection?.id]: scopeConnection?.boardsList }))
    //     setEntities(e => ({ ...e, [configuredConnection?.id]: scopeConnection?.entityList }))
    //     break
    // }
  }, [
    // activeProjectTransformation,
    // activeBoardTransformation,
    // configuredConnection,
    setScopeConnection
    // scopeConnection
  ])

  const handleBlueprintScopesDialogOpening = useCallback(
    () => {
      console.log('>>> OPENING BLUEPRINT CONNECTION SCOPES DIALOG...')
    },
    [
      // activeProvider?.id
    ]
  )

  const handleBlueprintSave = useCallback(() => {
    ToastNotification.clear()
    patchBlueprint(activeBlueprint, activeSetting?.payload, (response) => {
      console.log('>>> MY BP RESPONSE!!', response)
      if (response?.status === 200) {
        switch (activeSetting?.id) {
          case 'scopes':
          case 'transformations':
            handleBlueprintScopesDialogClose()
            break
          default:
            handleBlueprintDialogClose()
            break
        }
      } else {
        ToastNotification.show({
          message: response.message || 'Unable to modify Blueprint',
          intent: 'danger',
          icon: 'error'
        })
      }
    })
  }, [
    activeSetting,
    activeBlueprint,
    patchBlueprint,
    handleBlueprintDialogClose,
    handleBlueprintScopesDialogClose
  ])

  const handleConnectionStepChange = useCallback((newStepId, lastStepId, e) => {
    console.log(
      '>>> CONNECTION SETTINGS STEP CHANGED...',
      newStepId,
      lastStepId,
      e
    )
    switch (newStepId) {
      case 'scopes':
        setActiveSetting((aS) => ({
          ...aS,
          id: 'scopes',
          title: 'Change Data Scope'
        }))
        break
      case 'transformations':
        setActiveSetting((aS) => ({
          ...aS,
          id: 'transformations',
          title: 'Change Transformation Rules'
        }))
        break
    }
  }, [])

  const viewBlueprintStatus = useCallback(() => {
    history.push(`/blueprints/detail/${blueprintId}`)
  }, [history, blueprintId])

  const viewBlueprintSettings = useCallback(() => {
    history.push(`/blueprints/settings/${blueprintId}`)
  }, [history, blueprintId])

  const viewBlueprints = useCallback(() => {
    history.push('/blueprints')
  }, [history])

  const modifySetting = useCallback(
    (settingId) => {
      let title = null
      switch (settingId) {
        case 'name':
          title = 'Change Blueprint Name'
          break
        case 'cronConfig':
          title = 'Change Sync Frequency'
          break
        case 'plan':
          title = 'Change Task Configuration'
          break
        default:
          break
      }
      setActiveSetting((aS) => ({ ...aS, id: settingId, title }))
      setBlueprintDialogIsOpen(true)
      fetchBlueprint(blueprintId)
    },
    [blueprintId, fetchBlueprint]
  )

  const modifyConnection = useCallback(
    (connectionIdx, connectionId, provider) => {
      const connection = connectionsList.find(
        (c) => c.connectionId === connectionId && c.provider === provider?.id
      )
      const connectionWithScope = connections.find(
        (c) =>
          c.connectionId === connectionId && c.provider?.id === provider?.id
      )
      console.log('>>> MODIFYING DATA CONNECTION SCOPE...', connectionWithScope)
      setActiveProvider((aP) =>
        connection
          ? integrationsData.find((i) => i.id === connection?.provider)
          : aP
      )
      setActiveSetting((aS) => ({
        ...aS,
        id: 'scopes',
        title: 'Change Data Scope'
      }))
      setConfiguredConnection({
        ...connection,
        transformations: connectionWithScope.transformations
      })
      setScopeConnection({ ...connection, ...connectionWithScope })
    },
    [
      // activeProvider,
      connectionsList,
      connections,
      setScopeConnection,
      setConfiguredConnection
    ]
  )

  const validateActiveSetting = useCallback(() => {
    let isValid = false
    if (activeBlueprint?.mode === BlueprintMode.NORMAL) {
      switch (activeSetting?.id) {
        case 'name':
          isValid = validateBlueprintName(blueprintName)
          break
        case 'cronConfig':
          isValid =
            cronConfig === 'custom'
              ? isValidCronExpression(customCronConfig)
              : ['manual', 'custom'].includes(cronConfig) ||
                isValidCronExpression(cronConfig)
          break
        case 'scopes':
        case 'transformations':
          switch (activeProvider?.id) {
            case Providers.GITHUB:
              isValid =
                Array.isArray(projects[configuredConnection?.id]) &&
                validateRepositoryName(projects[configuredConnection?.id]) &&
                projects[configuredConnection?.id]?.length > 0 &&
                Array.isArray(entities[configuredConnection?.id]) &&
                entities[configuredConnection?.id]?.length > 0
              break
            case Providers.GITLAB:
              isValid =
                Array.isArray(projects[configuredConnection?.id]) &&
                projects[configuredConnection?.id]?.length > 0 &&
                entities[configuredConnection?.id]?.length > 0
              break
            case Providers.JIRA:
              isValid =
                Array.isArray(boards[configuredConnection?.id]) &&
                boards[configuredConnection?.id]?.length > 0 &&
                Array.isArray(entities[configuredConnection?.id]) &&
                entities[configuredConnection?.id]?.length > 0
              break
            case Providers.JENKINS:
              isValid = entities[configuredConnection?.id]?.length > 0
              break
            case Providers.TAPD:
              isValid = entities[configuredConnection?.id]?.length > 0
              break
            case Providers.ZENTAO:
              isValid = entities[configuredConnection?.id]?.length > 0
              break
            default:
              isValid = true
          }
          break
      }
    } else if (activeBlueprint?.mode === BlueprintMode.ADVANCED) {
      isValid = isValidBlueprint && isValidPipeline
    }

    return isValid
  }, [
    activeSetting?.id,
    blueprintName,
    cronConfig,
    customCronConfig,
    validateBlueprintName,
    // validateNumericSet,
    validateRepositoryName,
    isValidCronExpression,
    isValidBlueprint,
    isValidPipeline,
    projects,
    boards,
    entities,
    configuredConnection,
    activeProvider?.id,
    activeBlueprint?.mode
  ])

  // @note: lifted higher to dsm hook
  // const getDefaultEntities = useCallback((providerId) => {
  //   let entities = []
  //   switch (providerId) {
  //     case Providers.GITHUB:
  //     case Providers.GITLAB:
  //       entities = DEFAULT_DATA_ENTITIES.filter((d) => d.name !== 'ci-cd')
  //       break
  //     case Providers.JIRA:
  //       entities = DEFAULT_DATA_ENTITIES.filter((d) => d.name === 'issue-tracking' || d.name === 'cross-domain')
  //       break
  //     case Providers.JENKINS:
  //       entities = DEFAULT_DATA_ENTITIES.filter((d) => d.name === 'ci-cd')
  //       break
  //     case Providers.TAPD:
  //       entities = DEFAULT_DATA_ENTITIES.filter((d) => d.name === 'ci-cd')
  //       break
  //   }
  //   return entities
  // }, [])

  const addProjectTransformation = useCallback(
    (project) => {
      setConfiguredProject(project)
      // ToastNotification.clear()
    },
    [setConfiguredProject]
  )

  const addBoardTransformation = useCallback(
    (board) => {
      setConfiguredBoard(board)
      // ToastNotification.clear()
    },
    [setConfiguredBoard]
  )

  // @todo: lift higher to dsm hook
  // const getJiraMappedBoards = useCallback((options = []) => {
  //   return options.map(({ boardId, title }, sIdx) => {
  //     return {
  //       id: boardId,
  //       key: boardId,
  //       value: boardId,
  //       title: title || `Board ${boardId}`,
  //     }
  //   })
  // }, [scopeConnection?.endpoint])

  useEffect(() => {
    console.log('>>> ACTIVE PROVIDER!', activeProvider)
    setDataEntitiesList((deList) =>
      activeProvider ? getDefaultEntities(activeProvider?.id) : deList
    )
  }, [activeProvider, getDefaultEntities, setDataEntitiesList])

  useEffect(() => {
    setBlueprintId(bId)
    console.log('>>> REQUESTED SETTINGS for BLUEPRINT ID ===', bId)
  }, [bId])

  useEffect(() => {
    if (!isNaN(blueprintId)) {
      console.log('>>>> FETCHING BLUEPRINT ID...', blueprintId)
      fetchBlueprint(blueprintId)
      fetchAllConnections(false, true)
    }
  }, [blueprintId, fetchBlueprint, fetchAllConnections])

  useEffect(() => {
    console.log('>>>> SETTING ACTIVE BLUEPRINT...', blueprint)
    if (blueprint?.id) {
      setActiveBlueprint((b) => ({
        ...b,
        ...blueprint
      }))
    }
  }, [blueprint])

  useEffect(() => {
    console.log('>>> ACTIVE BLUEPRINT ....', activeBlueprint)
    if (activeBlueprint?.id && activeBlueprint?.mode === BlueprintMode.NORMAL) {
      setConnections(
        activeBlueprint?.settings?.connections.map(
          (c, cIdx) =>
            new DataScopeConnection(
              createNormalConnection(
                activeBlueprint,
                c,
                cIdx,
                DEFAULT_DATA_ENTITIES,
                allProviderConnections,
                connectionsList,
                [Providers.JIRA].includes(c.plugin)
                  ? getJiraMappedBoards(
                      c.scope?.map((s) => s.options?.boardId),
                      [
                        ...allJiraResources?.boards,
                        ...c.scope.map((s) => s.options)
                      ]
                    )
                  : []
              )
            )
        )
      )
    } else if (
      activeBlueprint?.id &&
      activeBlueprint?.mode === BlueprintMode.ADVANCED
    ) {
      setConnections(
        activeBlueprint?.plan?.flat().map(
          (c, cIdx) =>
            new DataScopeConnection(
              createAdvancedConnection(
                activeBlueprint,
                c,
                cIdx,
                DEFAULT_DATA_ENTITIES,
                allProviderConnections,
                connectionsList,
                [Providers.JIRA].includes(c.plugin)
                  ? getJiraMappedBoards(
                      c.scope?.map((s) => s.options?.boardId),
                      [
                        ...allJiraResources?.boards,
                        ...c.scope.map((s) => s.options)
                      ]
                    )
                  : []
              )
            )
        )
      )
    }
    setBlueprintName(activeBlueprint?.name)
    setCronConfig(
      [
        getCronPreset('hourly').cronConfig,
        getCronPreset('daily').cronConfig,
        getCronPreset('weekly').cronConfig,
        getCronPreset('monthly').cronConfig
      ].includes(activeBlueprint?.cronConfig)
        ? activeBlueprint?.cronConfig
        : activeBlueprint?.isManual
        ? 'manual'
        : 'custom'
    )
    setCustomCronConfig(
      !['custom', 'manual'].includes(activeBlueprint?.cronConfig)
        ? activeBlueprint?.cronConfig
        : '0 0 * * *'
    )
    setInterval(detectCronInterval(activeBlueprint?.cronConfig))
    setMode(activeBlueprint?.mode)
    setEnable(activeBlueprint?.enable)
    setIsManual(activeBlueprint?.isManual)
    setRawConfiguration(
      JSON.stringify(activeBlueprint?.plan, null, '  ') ||
        JSON.stringify([[]], null, '  ')
    )
    // setBlueprintSettings(activeBlueprint?.settings)
  }, [
    activeBlueprint,
    setBlueprintName,
    setConnections,
    detectCronInterval,
    getCronPreset,
    setCronConfig,
    setCustomCronConfig,
    setEnable,
    setInterval,
    setIsManual,
    setMode,
    setBlueprintSettings,
    // jiraApiBoards,
    allJiraResources?.boards,
    allProviderConnections,
    isFetchingJIRA,
    connectionsList,
    getDefaultEntities,
    getJiraMappedBoards,
    setRawConfiguration,
    createAdvancedConnection,
    createNormalConnection
  ])

  useEffect(() => {
    console.log('>>> SETTING ACTIVE SETTINGS PAYLOAD....')
    const isCustomCron = cronConfig === 'custom'
    const isManualCron = cronConfig === 'manual'

    switch (activeSetting?.id) {
      case 'name':
        setActiveSetting((aS) => ({
          ...aS,
          payload: {
            name: blueprintName
          }
        }))
        break
      case 'cronConfig':
        setActiveSetting((aS) => ({
          ...aS,
          payload: {
            isManual: !!isManualCron,
            cronConfig: isManualCron
              ? getCronPreset('daily').cronConfig
              : isCustomCron
              ? customCronConfig
              : cronConfig
          }
        }))
        break
      case 'scopes':
      case 'transformations':
        setActiveSetting((aS) => ({
          ...aS,
          payload: {
            settings: blueprintSettings
          }
        }))
        break
      case 'plan':
        setActiveSetting((aS) => ({
          ...aS,
          payload: {
            // plan: JSON.parse(rawConfiguration)
            plan: runTasksAdvanced
          }
        }))
        break
    }
  }, [
    blueprintName,
    cronConfig,
    customCronConfig,
    activeSetting?.id,
    getCronPreset,
    blueprintSettings,
    transformations,
    runTasksAdvanced
  ])

  useEffect(() => {
    console.log(
      '>>> RECEIVED ACTIVE SETTINGS PAYLOAD....',
      activeSetting?.payload
    )
  }, [activeSetting?.payload])

  useEffect(() => {
    console.log('>>> ACTIVE UI SETTING OBJECT...', activeSetting)
  }, [activeSetting])

  // useEffect(() => {
  //   validateBlueprint()
  // }, [
  //   blueprintName,
  //   // @todo: fix dependency warning with validateBlueprint
  //   // validateBlueprint
  // ])

  useEffect(() => {
    console.log('>>> DATA SCOPE CONNECTIONS...', connections)
    setBlueprintConnections(
      connections.map((c) =>
        connectionsList.find(
          (cItem) =>
            cItem.connectionId === c.connectionId &&
            cItem.provider === c.provider?.id
        )
      )
    )
  }, [connections, connectionsList])

  useEffect(() => {
    console.log('>>> AVAILABLE DATA ENTITIES...', dataEntitiesList)
  }, [dataEntitiesList])

  useEffect(() => {
    console.log('>>> SELECTED BLUEPRINT CONNECTIONS...', blueprintConnections)
  }, [blueprintConnections])

  useEffect(() => {
    console.log(
      '>>> CONNECTION SCOPE SELECTED, LOADING BLUEPRINT SETTINGS...',
      scopeConnection
    )
    const isJIRAProvider = scopeConnection?.providerId === Providers.JIRA
    if (scopeConnection) {
      if (isJIRAProvider) {
        setBlueprintScopesDialogIsOpen(true)
      } else {
        setBlueprintScopesDialogIsOpen(true)
      }
    }
  }, [
    // loadBlueprint,
    // activeProvider,
    // isFetchingJIRA,
    // jiraApiBoards,
    scopeConnection
    // configuredProject, configuredBoard
  ])

  // useEffect(() => {
  //   if (allJiraResources?.boards?.length > 0) {
  //     // setBlueprintScopesDialogIsOpen(true)
  //   }
  // }, [allJiraResources])

  // useEffect(() => {
  //   console.log('>>> CONFIGURING / MODIFYING CONNECTION', configuredConnection)
  //   if (configuredConnection?.id) {
  //     // setBoards({ [configuredConnection?.id]: [] })
  //     // setProjects({ [configuredConnection?.id]: [] })
  //     // setEntities({ [configuredConnection?.id]: [] })
  //   }
  // }, [configuredConnection])

  useEffect(() => {
    if (
      scopeConnection?.providerId === Providers.JIRA &&
      scopeConnection?.connectionId &&
      activeBlueprint?.mode === BlueprintMode.NORMAL
    ) {
      fetchIssueTypes()
      fetchFields()
      fetchBoards(undefined, (boards) =>
        setConnections((Cs) =>
          Cs.map(
            (c) =>
              new DataScopeConnection({
                ...c,
                boardsList: getJiraMappedBoards(c.boardIds, [
                  ...(boards ?? []),
                  ...c.scope.map((s) => s.options)
                ])
              })
          )
        )
      )
    }
  }, [
    activeBlueprint?.mode,
    fetchIssueTypes,
    fetchFields,
    fetchBoards,
    scopeConnection?.connectionId,
    scopeConnection?.providerId,
    getJiraMappedBoards,
    setConnections
  ])

  useEffect(() => {
    if (
      scopeConnection?.providerId === Providers.JIRA &&
      scopeConnection?.connectionId &&
      activeBlueprint?.mode === BlueprintMode.NORMAL &&
      boardSearch
    ) {
      fetchBoards(boardSearch, (boards) =>
        setConnections((Cs) =>
          Cs.map(
            (c) =>
              new DataScopeConnection({
                ...c,
                boardsList: getJiraMappedBoards(c.boardIds, [
                  ...(boards ?? []),
                  ...c.scope.map((s) => s.options)
                ])
              })
          )
        )
      )
    }
  }, [
    activeBlueprint?.mode,
    fetchBoards,
    scopeConnection?.connectionId,
    scopeConnection?.providerId,
    getJiraMappedBoards,
    setConnections,
    boardSearch
  ])

  useEffect(() => {
    console.log('>>> PROJECTS SELECTED...', projects)
    if (configuredConnection?.id) {
      setNewConnectionScopes((cS) => ({
        ...cS,
        [configuredConnection?.id]: {
          ...cS[[configuredConnection?.id]],
          projects: projects[configuredConnection?.id] || []
        }
      }))
    }
  }, [projects, configuredConnection?.id])

  useEffect(() => {
    console.log('>>> BOARDS SELECTED...', boards)
    if (configuredConnection?.id) {
      setNewConnectionScopes((cS) => ({
        ...cS,
        [configuredConnection?.id]: {
          ...cS[[configuredConnection?.id]],
          boards: boards[configuredConnection?.id] || []
        }
      }))
    }
  }, [boards, configuredConnection?.id])

  useEffect(() => {
    console.log('>>> ENTITIES SELECTED...', entities)
    if (configuredConnection?.id) {
      setNewConnectionScopes((cS) => ({
        ...cS,
        [configuredConnection?.id]: {
          ...cS[[configuredConnection?.id]],
          entities: entities[configuredConnection?.id] || []
        }
      }))
    }
  }, [entities, configuredConnection?.id])

  useEffect(() => {
    console.log('>>> NEW CONNECTION SCOPES', newConnectionScopes)
    // setScopeConnection(sC => ({...sC, projects: newConnectionScopes[configuredConnection?.id]?.projects }))
  }, [newConnectionScopes, configuredConnection?.id])

  useEffect(() => {
    console.log('>>>> SELECTED BOARDS!', boards)
  }, [boards])

  useEffect(() => {
    console.log(
      '>> PIPELINE RUN TASK SETTINGS FOR PIPELINE MANAGER ....',
      runTasks
    )
    setPipelineSettings({
      name: pipelineName,
      plan:
        activeBlueprint?.mode === BlueprintMode.ADVANCED
          ? runTasksAdvanced
          : [[...runTasks]]
    })
    if (activeBlueprint?.mode === BlueprintMode.ADVANCED) {
      validateAdvancedPipeline()
      setBlueprintTasks(runTasksAdvanced)
    } else {
      validatePipeline()
      setBlueprintTasks([[...runTasks]])
    }
  }, [
    pipelineName,
    activeBlueprint?.mode,
    runTasks,
    runTasksAdvanced,
    setPipelineSettings,
    validatePipeline,
    validateAdvancedPipeline,
    setBlueprintTasks
    // saveBlueprintComplete?.id
  ])

  useEffect(() => {
    if (isValidJSON(rawConfiguration)) {
      setRunTasksAdvanced(JSON.parse(rawConfiguration))
    }
  }, [rawConfiguration, isValidJSON])

  return (
    <>
      <div className='container'>
        <Nav />
        <Sidebar />
        <Content>
          <main className='main'>
            {activeBlueprint?.id !== null && blueprintErrors.length === 0 && (
              <div
                className='blueprint-header'
                style={{
                  display: 'flex',
                  width: '100%',
                  justifyContent: 'space-between',
                  marginBottom: '10px',
                  whiteSpace: 'nowrap'
                }}
              >
                <div className='blueprint-name' style={{}}>
                  <h2
                    style={{
                      fontWeight: 'bold',
                      display: 'flex',
                      alignItems: 'center',
                      color: !activeBlueprint?.enable ? Colors.GRAY1 : 'inherit'
                    }}
                  >
                    {activeBlueprint?.name}
                    <Tag
                      minimal
                      intent={
                        activeBlueprint.mode === BlueprintMode.ADVANCED
                          ? Intent.DANGER
                          : Intent.PRIMARY
                      }
                      style={{ marginLeft: '10px' }}
                    >
                      {activeBlueprint?.mode?.toString().toUpperCase()}
                    </Tag>
                  </h2>
                </div>
                <div
                  className='blueprint-info'
                  style={{ display: 'flex', alignItems: 'center' }}
                >
                  <div className='blueprint-schedule'>
                    {activeBlueprint?.isManual ? (
                      <strong>Manual Mode</strong>
                    ) : (
                      <span
                        className='blueprint-schedule-interval'
                        style={{
                          textTransform: 'capitalize',
                          padding: '0 10px'
                        }}
                      >
                        {activeBlueprint?.interval} (at{' '}
                        {dayjs(
                          getNextRunDate(activeBlueprint?.cronConfig)
                        ).format(
                          `hh:mm A ${
                            activeBlueprint?.interval !== 'Hourly'
                              ? ' MM/DD/YYYY'
                              : ''
                          }`
                        )}
                        )
                      </span>
                    )}{' '}
                    <span className='blueprint-schedule-nextrun'>
                      {!activeBlueprint?.isManual && (
                        <>
                          Next Run{' '}
                          {dayjs(
                            getNextRunDate(activeBlueprint?.cronConfig)
                          ).fromNow()}
                        </>
                      )}
                    </span>
                  </div>
                  <div
                    className='blueprint-actions'
                    style={{ padding: '0 10px' }}
                  >
                    {/* <Button
                      intent={Intent.PRIMARY}
                      small
                      text='Run Now'
                      onClick={runBlueprint}
                      disabled={!activeBlueprint?.enable || currentRun?.status === TaskStatus.RUNNING}
                    /> */}
                  </div>
                  <div className='blueprint-enabled'>
                    <Switch
                      id='blueprint-enable'
                      name='blueprint-enable'
                      checked={activeBlueprint?.enable}
                      label={
                        activeBlueprint?.enable
                          ? 'Blueprint Enabled'
                          : 'Blueprint Disabled'
                      }
                      onChange={() =>
                        handleBlueprintActivation(activeBlueprint)
                      }
                      style={{
                        marginBottom: 0,
                        marginTop: 0,
                        color: !activeBlueprint?.enable
                          ? Colors.GRAY3
                          : 'inherit'
                      }}
                      disabled={currentRun?.status === TaskStatus.RUNNING}
                    />
                  </div>
                  <div style={{ padding: '0 10px' }}>
                    <Button
                      intent={Intent.PRIMARY}
                      icon='trash'
                      small
                      minimal
                      disabled
                    />
                  </div>
                </div>
              </div>
            )}

            {blueprintErrors?.length > 0 && (
              <div className='bp3-non-ideal-state blueprint-non-ideal-state'>
                <div className='bp3-non-ideal-state-visual'>
                  <Icon icon='warning-sign' size={32} color={Colors.RED5} />
                </div>
                <h4 className='bp3-heading'>Invalid Blueprint</h4>
                <div>{blueprintErrors[0]}</div>
                <button
                  className='bp3-button bp3-intent-primary'
                  onClick={viewBlueprints}
                >
                  Continue
                </button>
              </div>
            )}

            {activeBlueprint?.id !== null && blueprintErrors.length === 0 && (
              <>
                <BlueprintNavigationLinks blueprint={activeBlueprint} />

                <div
                  className='blueprint-main-settings'
                  style={{
                    display: 'flex',
                    alignSelf: 'flex-start',
                    color: !activeBlueprint?.enable ? Colors.GRAY2 : 'inherit'
                  }}
                >
                  <div className='configure-settings-name'>
                    <h3>Name</h3>
                    <div style={{ display: 'flex', alignItems: 'center' }}>
                      <div className='blueprint-name'>
                        {activeBlueprint?.name}
                      </div>
                      <Button
                        icon='annotation'
                        intent={Intent.PRIMARY}
                        size={12}
                        small
                        minimal
                        onClick={() => modifySetting('name')}
                      />
                    </div>
                  </div>
                  <div
                    className='configure-settings-frequency'
                    style={{ marginLeft: '40px' }}
                  >
                    <h3>Sync Frequency</h3>
                    <div style={{ display: 'flex', alignItems: 'center' }}>
                      <div className='blueprint-frequency'>
                        {activeBlueprint?.isManual ? (
                          'Manual'
                        ) : (
                          <span>
                            {activeBlueprint?.interval} (at{' '}
                            {dayjs(
                              getNextRunDate(activeBlueprint?.cronConfig)
                            ).format('hh:mm A')}
                            )
                          </span>
                        )}
                      </div>
                      <Button
                        icon='annotation'
                        intent={Intent.PRIMARY}
                        size={12}
                        small
                        minimal
                        onClick={() => modifySetting('cronConfig')}
                      />
                    </div>
                  </div>
                </div>

                {activeBlueprint?.id &&
                  activeBlueprint?.mode === BlueprintMode.NORMAL && (
                    <div
                      className='data-scopes-grid'
                      style={{
                        width: '100%',
                        marginTop: '40px',
                        alignSelf: 'flex-start'
                      }}
                    >
                      <h2
                        style={{
                          fontWeight: 'bold',
                          color: !activeBlueprint?.enable
                            ? Colors.GRAY1
                            : 'inherit'
                        }}
                      >
                        Data Scope and Transformation
                      </h2>
                      <DataScopesGrid
                        connections={connections}
                        blueprint={activeBlueprint}
                        onModify={modifyConnection}
                        mode={activeBlueprint?.mode}
                        loading={
                          isFetchingBlueprint ||
                          isFetchingJIRA ||
                          isFetchingGitlab
                        }
                      />
                    </div>
                  )}

                {activeBlueprint?.id && mode === BlueprintMode.ADVANCED && (
                  <div
                    className='data-advanced'
                    style={{
                      width: '100%',
                      maxWidth: '100%',
                      marginTop: '40px',
                      alignSelf: 'flex-start'
                    }}
                  >
                    <div style={{ display: 'flex', alignItems: 'center' }}>
                      <h2 style={{ fontWeight: 'bold' }}>
                        Data Scope and Transformation
                      </h2>
                      <div>
                        <Button
                          icon='annotation'
                          text='Edit JSON'
                          intent={Intent.PRIMARY}
                          small
                          minimal
                          onClick={() => modifySetting('plan')}
                          style={{ fontSize: '12px' }}
                        />
                      </div>
                    </div>
                    <DataScopesGrid
                      connections={connections}
                      blueprint={activeBlueprint}
                      onModify={() => modifySetting('plan')}
                      mode={activeBlueprint?.mode}
                      classNames={['advanced-mode-grid']}
                      loading={
                        isFetchingBlueprint ||
                        isFetchingJIRA ||
                        isFetchingGitlab
                      }
                    />
                  </div>
                )}

                {ENVIRONMENT !== 'production' && (
                  <Button
                    // loading={isLoading}
                    intent={Intent.PRIMARY}
                    icon='code'
                    text='Inspect'
                    onClick={() => setShowBlueprintInspector(true)}
                    style={{ margin: '12px auto' }}
                    minimal
                    small
                  />
                )}
              </>
            )}
          </main>
        </Content>
      </div>

      <BlueprintDialog
        isOpen={blueprintDialogIsOpen}
        title={activeSetting?.title}
        blueprint={activeBlueprint}
        onSave={handleBlueprintSave}
        isSaving={isSaving}
        isValid={validateActiveSetting()}
        onClose={handleBlueprintDialogClose}
        onCancel={handleBlueprintDialogClose}
        errors={[...pipelineValidationErrors, ...blueprintValidationErrors]}
        content={(() => {
          let Settings = null
          switch (activeSetting?.id) {
            case 'name':
              Settings = (
                <BlueprintNameCard
                  name={blueprintName}
                  setBlueprintName={setBlueprintName}
                  fieBldHasError={fieldHasError}
                  getFieldError={getFieldError}
                  elevation={Elevation.ZERO}
                  enableDivider={false}
                  cardStyle={{ padding: 0 }}
                  isSaving={isSaving}
                />
              )
              break
            case 'cronConfig':
              Settings = (
                <DataSync
                  cronConfig={cronConfig}
                  customCronConfig={customCronConfig}
                  createCron={createCron}
                  setCronConfig={setCronConfig}
                  getCronPreset={getCronPreset}
                  fieldHasError={fieldHasError}
                  getFieldError={getFieldError}
                  setCustomCronConfig={setCustomCronConfig}
                  getCronPresetByConfig={getCronPresetByConfig}
                  elevation={Elevation.ZERO}
                  enableHeader={false}
                  cardStyle={{ padding: 0 }}
                />
              )
              break
            case 'plan':
              Settings = (
                <AdvancedJSON
                  // activeStep={activeStep}
                  advancedMode={mode === BlueprintMode.ADVANCED}
                  runTasksAdvanced={runTasksAdvanced}
                  blueprintConnections={blueprintConnections}
                  connectionsList={connectionsList}
                  name={name}
                  setBlueprintName={setBlueprintName}
                  fieldHasError={fieldHasError}
                  getFieldError={getFieldError}
                  // onAdvancedMode={handleAdvancedMode}
                  // @todo add multistage checker method
                  isMultiStagePipeline={() => {}}
                  rawConfiguration={rawConfiguration}
                  setRawConfiguration={setRawConfiguration}
                  isSaving={isSaving}
                  // @todo re-enable validation
                  isValidConfiguration={true}
                  validationAdvancedError={null}
                  validationErrors={pipelineValidationErrors}
                  elevation={Elevation.ZERO}
                  enableHeader={false}
                  useBlueprintName={false}
                  showTemplates={true}
                  showModeNotice={false}
                  cardStyle={{ padding: 0 }}
                  descriptionText='Enter Advanced JSON Tasks'
                />
              )
          }
          return Settings
        })()}
      />

      <BlueprintDataScopesDialog
        isOpen={blueprintScopesDialogIsOpen}
        title={activeSetting?.title}
        dataEntitiesList={dataEntitiesList}
        blueprint={activeBlueprint}
        blueprintConnections={blueprintConnections}
        configuredConnection={configuredConnection}
        configuredProject={configuredProject}
        configuredBoard={configuredBoard}
        configurationKey={configurationKey}
        scopeConnection={scopeConnection}
        activeTransformation={activeTransformation}
        addProjectTransformation={addProjectTransformation}
        addBoardTransformation={addBoardTransformation}
        provider={activeProvider}
        entities={entities}
        boards={boards}
        setBoardSearch={setBoardSearch}
        boardsList={jiraApiBoards}
        projects={projects}
        issueTypesList={jiraApiIssueTypes}
        fieldsList={jiraApiFields}
        isFetching={isFetchingBlueprint}
        isFetchingJIRA={isFetchingJIRA}
        fetchGitlabProjects={fetchGitlabProjects}
        gitlabProjects={gitlabProjects}
        isFetchingGitlab={isFetchingGitlab}
        gitlabProxyError={gitlabProxyError}
        setConfiguredProject={setConfiguredProject}
        setConfiguredBoard={setConfiguredBoard}
        setBoards={setBoards}
        setProjects={setProjects}
        setEntities={setEntities}
        setTransformationSettings={setTransformationSettings}
        onOpening={handleBlueprintScopesDialogOpening}
        onSave={handleBlueprintSave}
        isSaving={isSaving}
        // @todo: validation status
        isValid={validateActiveSetting()}
        onClose={handleBlueprintScopesDialogClose}
        onCancel={handleBlueprintScopesDialogClose}
        onStepChange={handleConnectionStepChange}
        fieldHasError={fieldHasError}
        getFieldError={getFieldError}
        jiraProxyError={jiraProxyError}
        errors={[...pipelineValidationErrors, ...blueprintValidationErrors]}
      />

      <CodeInspector
        title={<>&nbsp; {blueprintName}</>}
        titleIcon={
          activeProvider ? (
            <Icon icon={activeProvider?.icon} size={16} />
          ) : (
            'add'
          )
        }
        subtitle='JSON CONFIGURATION'
        isOpen={showBlueprintInspector}
        activePipeline={
          activeBlueprint?.mode === BlueprintMode.ADVANCED
            ? {
                name: activeBlueprint?.name,
                plan: activeBlueprint?.plan
              }
            : {
                name: activeBlueprint?.name,
                settings: blueprintSettings
              }
        }
        onClose={setShowBlueprintInspector}
        hasBackdrop={false}
      />
    </>
  )
}

export default BlueprintSettings
