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
import React, { Fragment, useEffect, useState, useCallback } from 'react'
import { CSSTransition } from 'react-transition-group'
import { useHistory, useLocation, Link } from 'react-router-dom'
import dayjs from '@/utils/time'
import {
  API_PROXY_ENDPOINT,
  ISSUE_TYPES_ENDPOINT,
  ISSUE_FIELDS_ENDPOINT,
  BOARDS_ENDPOINT,
} from '@/config/jiraApiProxy'
import { integrationsData } from '@/data/integrations'
import { Divider, Elevation, Card } from '@blueprintjs/core'
import {
  Providers,
  ProviderTypes,
  ProviderIcons,
  ConnectionStatus,
  ConnectionStatusLabels,
} from '@/data/Providers'
import Nav from '@/components/Nav'
import Sidebar from '@/components/Sidebar'
// import AppCrumbs from '@/components/Breadcrumbs'
import Content from '@/components/Content'

import { DataEntities, DataEntityTypes } from '@/data/DataEntities'
import { NullBlueprint, BlueprintMode } from '@/data/NullBlueprint'
import { NullBlueprintConnection } from '@/data/NullBlueprintConnection'
import { NullConnection } from '@/data/NullConnection'

import {
  WorkflowSteps,
  WorkflowAdvancedSteps,
  DEFAULT_DATA_ENTITIES,
  DEFAULT_BOARDS,
} from '@/data/BlueprintWorkflow'

import useBlueprintManager from '@/hooks/useBlueprintManager'
import usePipelineManager from '@/hooks/usePipelineManager'
import useConnectionManager from '@/hooks/useConnectionManager'
import useBlueprintValidation from '@/hooks/useBlueprintValidation'
import usePipelineValidation from '@/hooks/usePipelineValidation'
import useConnectionValidation from '@/hooks/useConnectionValidation'
import useJIRA from '@/hooks/useJIRA'

import WorkflowStepsBar from '@/components/blueprints/WorkflowStepsBar'
import WorkflowActions from '@/components/blueprints/WorkflowActions'
import ConnectionDialog from '@/components/blueprints/ConnectionDialog'
import CodeInspector from '@/components/pipelines/CodeInspector'
// import NoData from '@/components/NoData'

import DataConnections from '@/components/blueprints/create-workflow/DataConnections'
import DataScopes from '@/components/blueprints/create-workflow/DataScopes'
import DataTransformations from '@/components/blueprints/create-workflow/DataTransformations'
import DataSync from '@/components/blueprints/create-workflow/DataSync'
import AdvancedJSON from '@/components/blueprints/create-workflow/AdvancedJSON'
import AdvancedJSONValidation from '@/components/blueprints/create-workflow/AdvancedJSONValidation'

import { DEVLAKE_ENDPOINT } from '@/utils/config'
import request from '@/utils/request'

// import ConnectionTabs from '@/components/blueprints/ConnectionTabs'

// manage transformations in one place
const useTransformationsManager = () => {
  const [transformations, setTransformations] = useState({})

  const generateKey = (connection, projectNameOrBoard) => {
    console.log(
      '>> generateKey generateKeygenerateKeygenerateKeygenerateKeygenerateKey',
      connection,
      projectNameOrBoard
    )
    return `${connection?.provider}/${connection?.id}/${projectNameOrBoard.id || projectNameOrBoard}`
  }

  const changeTransformationSettings = useCallback((settings, connection, projectNameOrBoard) => {
    const key = generateKey(connection, projectNameOrBoard)
    console.log(
      '>> SETTING TRANSFORMATION SETTINGS PROJECT/BOARD...',
      key,
      settings
    )
    setTransformations((existingTransformations) => ({
      ...existingTransformations,
      [key]: {
        ...existingTransformations[key],
        ...settings,
      },
    }))
  }, [setTransformations])

  const getDefaultTransformations = (provider) => {
    let transforms = {}
    switch (provider) {
      case Providers.GITHUB:
        transforms = {
          prType: '',
          prComponent: '',
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
    }
    return transforms
  }

  const initDefaultTransformationSettingsIfNotExist = useCallback((connection, projectNameOrBoard) => {
    const key = generateKey(connection, projectNameOrBoard)
    console.log(
      '>> INIT DEFAULT TRANSFORMATION SETTINGS PROJECT/BOARD...',
      key,
    )
    if (!transformations[key]) {
      setTransformations(old => ({
        ...old,
        [key]: getDefaultTransformations(connection?.provider),
      }))
    }
  }, [setTransformations, transformations])

  const getTransformation = useCallback((connection, projectNameOrBoard) => {
    const key = generateKey(connection, projectNameOrBoard)
    return transformations[key]
  }, [transformations])

  const clearTransformationSettings = useCallback((connection, projectNameOrBoard) => {
    const key = generateKey(connection, projectNameOrBoard)
    console.log(
      '>> CLEAR TRANSFORMATION SETTINGS PROJECT/BOARD...',
      key,
    )
    setTransformations((existingTransformations) => ({
      ...existingTransformations,
      [key]: null,
    }))
  }, [setTransformations])

  const checkTransformationIsExist = useCallback((connection, projectNameOrBoard) => {
    const key = generateKey(connection, projectNameOrBoard)
    const storedTransform = transformations[key]
    return Object.values(storedTransform).some(v => v && v.length > 0)
  }, [transformations])

  return {
    getTransformation,
    changeTransformationSettings,
    initDefaultTransformationSettingsIfNotExist,
    clearTransformationSettings,
    checkTransformationIsExist,
  }
}

const CreateBlueprint = (props) => {
  const history = useHistory()
  // const dispatch = useDispatch()

  const [blueprintAdvancedSteps, setBlueprintAdvancedSteps] = useState(
    WorkflowAdvancedSteps
  )
  const [blueprintNormalSteps, setBlueprintNormalSteps] =
    useState(WorkflowSteps)
  const [blueprintSteps, setBlueprintSteps] = useState(blueprintNormalSteps)
  const [advancedMode, setAdvancedMode] = useState(false)
  const [activeStep, setActiveStep] = useState(
    blueprintSteps.find((s) => s.id === 1)
  )
  const [activeProvider, setActiveProvider] = useState(integrationsData[0])

  const [enabledProviders, setEnabledProviders] = useState([])
  const [runTasks, setRunTasks] = useState([])
  const [runTasksAdvanced, setRunTasksAdvanced] = useState([])
  const [runNow, setRunNow] = useState(false)
  const [newBlueprintId, setNewBlueprintId] = useState()
  const [existingTasks, setExistingTasks] = useState([])
  const [rawConfiguration, setRawConfiguration] = useState(
    JSON.stringify([runTasks], null, '  ')
  )
  const [isValidConfiguration, setIsValidConfiguration] = useState(false)
  const [validationAdvancedError, setValidationAdvancedError] = useState()

  const [connectionDialogIsOpen, setConnectionDialogIsOpen] = useState(false)
  const [managedConnection, setManagedConnection] = useState(
    NullBlueprintConnection
  )

  const [connectionsList, setConnectionsList] = useState([])

  const [dataEntitiesList, setDataEntitiesList] = useState([
    ...DEFAULT_DATA_ENTITIES,
  ])
  const [boardsList, setBoardsList] = useState([])

  const [blueprintConnections, setBlueprintConnections] = useState([])
  const [configuredConnection, setConfiguredConnection] = useState()
  const [dataEntities, setDataEntities] = useState({})
  const [activeConnectionTab, setActiveConnectionTab] = useState()

  const [onlineStatus, setOnlineStatus] = useState({})
  useEffect(async () => {
    const results = await Promise.all(blueprintConnections.map(
      c => request.post(`${DEVLAKE_ENDPOINT}/plugins/${c.plugin}/test`, c))
    )
    setOnlineStatus(results.map(r => r.status === 200 ? "Online" : "Offline"))
  }, [blueprintConnections])

  const [showBlueprintInspector, setShowBlueprintInspector] = useState(false)

  const {
    getTransformation,
    changeTransformationSettings,
    initDefaultTransformationSettingsIfNotExist,
    clearTransformationSettings,
    checkTransformationIsExist,
  } = useTransformationsManager()
  const [activeTransformation, setActiveTransformation] = useState()

  // @todo: replace with $projects
  const [projectId, setProjectId] = useState([])
  const [projects, setProjects] = useState({})
  const [boards, setBoards] = useState({})
  // @todo: replace with $boards
  const [boardId, setBoardId] = useState([])
  const [connectionId, setConnectionId] = useState('')
  const [connections, setConnections] = useState([])
  const [repositories, setRepositories] = useState([])
  const [selectedConnection, setSelectedConnection] = useState()
  const [repositoryName, setRepositoryName] = useState('')
  const [owner, setOwner] = useState('')
  const [gitExtractorUrl, setGitExtractorUrl] = useState('')
  const [gitExtractorRepoId, setGitExtractorRepoId] = useState('')
  const [selectedGithubRepo, setSelectedGithubRepo] = useState()
  const [refDiffRepoId, setRefDiffRepoId] = useState('')
  const [refDiffPairs, setRefDiffPairs] = useState([])
  const [refDiffTasks, setRefDiffTasks] = useState([
    'calculateCommitsDiff',
    'calculateIssuesDiff',
  ])

  const [configuredProject, setConfiguredProject] = useState(
    projects.length > 0 ? projects[0] : null
  )
  const [configuredBoard, setConfiguredBoard] = useState(
    boards.length > 0 ? boards[0] : null
  )

  const {
    activeConnection,
    // eslint-disable-next-line no-unused-vars
    blueprint,
    // eslint-disable-next-line no-unused-vars
    blueprints,
    name,
    cronConfig,
    customCronConfig,
    // eslint-disable-next-line no-unused-vars
    cronPresets,
    tasks: blueprintTasks,
    settings: blueprintSettings,
    detectedProviderTasks,
    enable,
    mode,
    setName: setBlueprintName,
    setCronConfig,
    setCustomCronConfig,
    setTasks: setBlueprintTasks,
    setSettings: setBlueprintSettings,
    setDetectedProviderTasks,
    setEnable: setEnableBlueprint,
    setMode: setBlueprintMode,
    // eslint-disable-next-line no-unused-vars
    isFetching: isFetchingBlueprints,
    isSaving,
    createCronExpression: createCron,
    // eslint-disable-next-line no-unused-vars
    getCronSchedule: getSchedule,
    getNextRunDate,
    getCronPreset,
    getCronPresetByConfig,
    saveBlueprint,
    deleteBlueprint,
    isDeleting: isDeletingBlueprint,
    saveComplete: saveBlueprintComplete,
  } = useBlueprintManager()

  const {
    // eslint-disable-next-line no-unused-vars
    validate: validateBlueprint,
    // eslint-disable-next-line no-unused-vars
    errors: blueprintValidationErrors,
    // setErrors: setBlueprintErrors,
    isValid: isValidBlueprint,
    fieldHasError,
    getFieldError,
  } = useBlueprintValidation({
    name,
    cronConfig,
    customCronConfig,
    enable,
    tasks: blueprintTasks,
    mode,
  })

  const {
    pipelineName,
    pipelines,
    runPipeline,
    cancelPipeline,
    fetchPipeline,
    fetchAllPipelines,
    pipelineRun,
    buildPipelineStages,
    isRunning,
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
    detectPipelineProviders,
  } = usePipelineManager(null, runTasks)

  const {
    validate: validatePipeline,
    validateAdvanced: validateAdvancedPipeline,
    errors: validationErrors,
    setErrors: setPipelineErrors,
    isValid: isValidPipelineForm,
    detectedProviders,
  } = usePipelineValidation({
    enabledProviders,
    pipelineName,
    projectId,
    projects,
    boardId,
    boards,
    owner,
    repositoryName,
    connectionId,
    gitExtractorUrl,
    gitExtractorRepoId,
    refDiffRepoId,
    refDiffTasks,
    refDiffPairs,
    tasks: runTasks,
    tasksAdvanced: runTasksAdvanced,
    advancedMode,
    mode,
    connection: configuredConnection,
    entities: dataEntities,
  })

  const {
    fetchIssueTypes,
    fetchFields,
    fetchBoards,
    boards: jiraApiBoards,
    issueTypes: jiraApiIssueTypes,
    fields: jiraApiFields,
    isFetching: isFetchingJIRA,
    error: jiraProxyError,
  } = useJIRA(
    {
      apiProxyPath: API_PROXY_ENDPOINT,
      issuesEndpoint: ISSUE_TYPES_ENDPOINT,
      fieldsEndpoint: ISSUE_FIELDS_ENDPOINT,
      boardsEndpoint: BOARDS_ENDPOINT,
    },
    configuredConnection
  )

  const {
    testConnection,
    saveConnection,
    fetchConnection,
    allProviderConnections,
    errors: connectionErrors,
    isSaving: isSavingConnection,
    isTesting: isTestingConnection,
    showError,
    testStatus,
    name: connectionName,
    endpointUrl,
    proxy,
    token,
    initialTokenStore,
    username,
    password,
    provider,
    setActiveConnection,
    setProvider,
    setName,
    setEndpointUrl,
    setProxy,
    setUsername,
    setPassword,
    setToken,
    setInitialTokenStore,
    setTestStatus,
    setTestResponse,
    setAllTestResponses,
    setSaveComplete: setSaveConnectionComplete,
    fetchAllConnections,
    connectionLimitReached,
    clearConnection: clearActiveConnection,
    testResponse,
    allTestResponses,
    saveComplete: saveConnectionComplete
  } = useConnectionManager(
    {
      activeProvider,
      connectionId: managedConnection?.connectionId,
    },
    managedConnection && managedConnection?.id !== null
  )

  const {
    validate: validateConnection,
    errors: connectionValidationErrors,
    isValid: isValidConnection,
  } = useConnectionValidation({
    activeProvider,
    name: connectionName,
    endpointUrl,
    proxy,
    token,
    username,
    password,
  })

  const isValidStep = useCallback((stepId) => { }, [])

  const nextStep = useCallback(() => {
    setActiveStep((aS) =>
      blueprintSteps.find(
        (s) => s.id === Math.min(aS.id + 1, blueprintSteps.length)
      )
    )
  }, [blueprintSteps])

  const prevStep = useCallback(() => {
    setActiveStep((aS) =>
      blueprintSteps.find((s) => s.id === Math.max(aS.id - 1, 1))
    )
  }, [blueprintSteps])

  const handleConnectionTabChange = useCallback(
    (tab) => {
      console.log('>> CONNECTION TAB CHANGED', tab)
      const selectedConnection = blueprintConnections.find(
        (c) => c.id === Number(tab.split('-')[1])
      )
      setActiveConnectionTab(tab)
      setActiveProvider(
        integrationsData.find((p) => p.id === selectedConnection.provider)
      )
      setProvider(
        integrationsData.find((p) => p.id === selectedConnection.provider)
      )
      setConfiguredConnection(selectedConnection)
    },
    [blueprintConnections, setProvider]
  )

  const handleConnectionDialogOpen = () => {
    console.log('>>> MANAGING CONNECTION', managedConnection)
  }

  const handleConnectionDialogClose = useCallback(() => {
    setConnectionDialogIsOpen(false)
    setManagedConnection(NullBlueprintConnection)
    setTestStatus(0)
    setTestResponse(null)
    setAllTestResponses({})
    setInitialTokenStore({})
    clearActiveConnection()
    setActiveConnection(NullConnection)
    setSaveConnectionComplete(null)
  }, [
    clearActiveConnection,
    setActiveConnection,
    setAllTestResponses,
    setInitialTokenStore,
    setSaveConnectionComplete,
    setTestResponse,
    setTestStatus
  ])

  const handleTransformationCancel = useCallback(() => {
    setConfiguredProject(null)
    setConfiguredBoard(null)
    console.log('>> Cancel Modify - Transformation Settings')
  }, [])

  const handleTransformationClear = useCallback(() => {
    console.log(
      '>>> CLEARING TRANSFORMATION RULES!',
      '==> PROJECT =',
      configuredProject,
      '==> BOARD =',
      configuredBoard
    )
    if (configuredProject) {
      clearTransformationSettings(configuredConnection, configuredProject)
    }
    if (configuredBoard) {
      clearTransformationSettings(configuredConnection, configuredBoard)
    }
    setConfiguredProject(null)
    setConfiguredBoard(null)
  }, [clearTransformationSettings, configuredConnection, configuredProject, configuredBoard])

  const handleBlueprintSave = useCallback(() => {
    console.log('>>> SAVING BLUEPRINT!!')
    setRunNow(false)
    saveBlueprint()
  }, [saveBlueprint])

  const handleBlueprintSaveAndRun = useCallback(() => {
    console.log('>>> SAVING BLUEPRINT & RUNNING NOW!!')
    setRunNow(true)
    saveBlueprint()
  }, [saveBlueprint])

  const getRestrictedDataEntities = useCallback(() => {
    let items = []
    switch (configuredConnection.provider) {
      case Providers.GITLAB:
      case Providers.JIRA:
      case Providers.GITHUB:
        items = dataEntitiesList.filter((d) => d.name !== 'ci-cd')
        break
      case Providers.JENKINS:
        items = dataEntitiesList.filter((d) => d.name === 'ci-cd')
        break
    }
    return items
  }, [dataEntitiesList, configuredConnection])

  const createProviderScopes = useCallback(
    (
      providerId,
      connection,
      connectionIdx,
      entities = [],
      boards = [],
      projects = [],
      getTransformation = (connection, projectOrBoard) => {},
      defaultScope = { transformation: {}, options: {}, entities: [] }
    ) => {
      console.log(
        '>>> CREATING PROVIDER SCOPE FOR CONNECTION...',
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
              boardId: Number(b.id),
              // @todo: verify initial value of since date for jira provider
              // since: new Date(),
            },
            transformation: { ...getTransformation(connection, b) },
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
          }
          break
        case Providers.GITHUB:
          newScope = projects[connection.id]?.map((p) => ({
            ...newScope,
            options: {
              owner: p.split('/')[0],
              repo: p.split('/')[1],
            },
            transformation: { ...getTransformation(connection, p) },
          }))
          break
      }
      return Array.isArray(newScope) ? newScope.flat() : [newScope]
    },
    []
  )

  const manageConnection = useCallback(
    (connection) => {
      console.log('>> MANAGE CONNECTION...', connection)
      if (connection?.id !== null) {
        setActiveProvider(
          integrationsData.find((p) => p.id === connection.provider)
        )
        setProvider(integrationsData.find((p) => p.id === connection.provider))
        // fetchConnection(true, false, connection.id)
        setManagedConnection(connection)
        setConnectionDialogIsOpen(true)
      }
    },
    [setProvider]
  )

  const addProjectTransformation = useCallback((project) => {
    setConfiguredProject(project)
  }, [])

  const addBoardTransformation = useCallback((board) => {
    setConfiguredBoard(board)
  }, [])

  const addConnection = useCallback(() => {
    setManagedConnection(NullBlueprintConnection)
    setConnectionDialogIsOpen(true)
  }, [])

  const handleTransformationSave = useCallback((settings, connection, projectNameOrBoard) => {
    console.log('>> SAVING / CLOSING Transformation Settings', settings, connection, projectNameOrBoard)
    changeTransformationSettings(settings, connection, projectNameOrBoard)
    setConfiguredProject(null)
    setConfiguredBoard(null)
  }, [changeTransformationSettings])

  const handleAdvancedMode = (enableAdvanced = true) => {
    setAdvancedMode(enableAdvanced)
  }

  const parseJSON = (jsonString = '') => {
    try {
      return JSON.parse(jsonString)
    } catch (e) {
      console.log('>> PARSE JSON ERROR!', e)
      throw e
    }
  }

  const isValidCode = useCallback(() => {
    let isValid = false
    try {
      const parsedCode = parseJSON(rawConfiguration)
      isValid = true
    } catch (e) {
      console.log('>> FORMAT CODE: Invalid Code Format!', e)
      isValid = false
      setValidationAdvancedError(e.message)
    }
    setIsValidConfiguration(isValid)
    return isValid
  }, [rawConfiguration])

  useEffect(() => {
    console.log('>> ACTIVE STEP CHANGED: ', activeStep)
    if (activeStep?.id === 1) {
      const enableNotifications = false
      const getAllSources = true
      fetchAllConnections(enableNotifications, getAllSources)
    }
    if (activeStep?.id === 2 || activeStep?.id === 3) {
      fetchBoards()
      fetchIssueTypes()
      fetchFields()
    }
    setBlueprintNormalSteps((bS) => [
      ...bS.map((s) =>
        s.id < activeStep?.id
          ? { ...s, complete: true }
          : { ...s, complete: false }
      ),
    ])
    setBlueprintAdvancedSteps((bS) => [
      ...bS.map((s) =>
        s.id < activeStep?.id
          ? { ...s, complete: true }
          : { ...s, complete: false }
      ),
    ])
  }, [
    activeStep,
    fetchAllConnections,
    fetchBoards,
    fetchFields,
    fetchIssueTypes,
  ])

  useEffect(() => {
    console.log('>>> ALL DATA PROVIDER CONNECTIONS...', allProviderConnections)
    setConnectionsList(
      allProviderConnections?.map((c, cIdx) => ({
        ...c,
        id: cIdx,
        name: c.name,
        title: c.name,
        value: c.id,
        status:
          ConnectionStatusLabels[c.status] ||
          ConnectionStatusLabels[ConnectionStatus.OFFLINE],
        provider: c.provider,
        plugin: c.provider,
      }))
    )
  }, [allProviderConnections])

  useEffect(() => {
    console.log(
      '>> PIPELINE RUN TASK SETTINGS FOR PIPELINE MANAGER ....',
      runTasks
    )
    setPipelineSettings({
      name: pipelineName,
      // blueprintId: saveBlueprintComplete?.id || 0,
      plan: advancedMode ? runTasksAdvanced : [[...runTasks]],
    })
    // setRawConfiguration(JSON.stringify(buildPipelineStages(runTasks, true), null, '  '))
    if (advancedMode) {
      validateAdvancedPipeline()
      setBlueprintTasks(runTasksAdvanced)
    } else {
      validatePipeline()
      setBlueprintTasks([[...runTasks]])
    }
  }, [
    pipelineName,
    advancedMode,
    runTasks,
    runTasksAdvanced,
    setPipelineSettings,
    validatePipeline,
    validateAdvancedPipeline,
    setBlueprintTasks,
    // saveBlueprintComplete?.id
  ])

  useEffect(() => {
    console.log(
      '>> BLUEPRINT SETTINGS FOR PIPELINE MANAGER ....',
      blueprintSettings
    )
  }, [blueprintSettings])

  useEffect(() => {
    validateBlueprint()
  }, [
    name,
    cronConfig,
    customCronConfig,
    blueprintTasks,
    enable,
    validateBlueprint,
  ])

  useEffect(() => {
    setConfiguredConnection(
      blueprintConnections.length > 0 ? blueprintConnections[0] : null
    )
    const getDefaultEntities = (providerId) => {
      let entities = []
      switch (providerId) {
        case Providers.GITHUB:
        case Providers.GITLAB:
          entities = DEFAULT_DATA_ENTITIES.filter((d) => d.name !== 'ci-cd')
          break
        case Providers.JIRA:
          entities = DEFAULT_DATA_ENTITIES.filter((d) => d.name === 'issue-tracking' || d.name === 'cross-domain')
          break
        case Providers.JENKINS:
          entities = DEFAULT_DATA_ENTITIES.filter((d) => d.name === 'ci-cd')
          break
      }
      return entities
    }
    const initializeEntities = (pV, cV) => ({
      ...pV,
      [cV.id]: !pV[cV.id] ? getDefaultEntities(cV?.provider) : [],
    })
    const initializeProjects = (pV, cV) => ({ ...pV, [cV.id]: [] })
    const initializeBoards = (pV, cV) => ({ ...pV, [cV.id]: [] })
    setDataEntities((dE) => ({
      ...blueprintConnections.reduce(initializeEntities, {}),
    }))
    setProjects((p) => ({
      ...blueprintConnections.reduce(initializeProjects, {}),
    }))
    setBoards((b) => ({
      ...blueprintConnections.reduce(initializeBoards, {}),
    }))
    setEnabledProviders([
      ...new Set(blueprintConnections.map((c) => c.provider)),
    ])
  }, [blueprintConnections])

  useEffect(() => {
    console.log('>> CONFIGURING CONNECTION', configuredConnection)
    if (configuredConnection) {
      setConfiguredProject(null)
      setConfiguredBoard(null)
      switch (configuredConnection.provider) {
        case Providers.GITLAB:
        case Providers.GITHUB:
          setDataEntitiesList(
            DEFAULT_DATA_ENTITIES.filter((d) => d.name !== 'ci-cd')
          )
          // setConfiguredProject(projects.length > 0 ? projects[0] : null)
          break
        case Providers.JIRA:
          setDataEntitiesList(
            DEFAULT_DATA_ENTITIES.filter((d) => d.name === 'issue-tracking' || d.name === 'cross-domain')
          )
          break
        case Providers.JENKINS:
          setDataEntitiesList(
            DEFAULT_DATA_ENTITIES.filter((d) => d.name === 'ci-cd')
          )
          break
        default:
          setDataEntitiesList(DEFAULT_DATA_ENTITIES)
          break
      }
    }
  }, [configuredConnection])

  useEffect(() => {
    console.log('>> DATA ENTITIES', dataEntities)
  }, [dataEntities])

  useEffect(() => {
    console.log('>> BOARDS', boards)
  }, [boards])

  useEffect(() => {
    setBlueprintSettings((currentSettings) => ({
      ...currentSettings,
      connections: blueprintConnections.map((c, cIdx) => ({
        ...NullBlueprintConnection,
        connectionId: c.value,
        plugin: c.plugin || c.provider,
        scope: createProviderScopes(
          c.provider,
          c,
          cIdx,
          dataEntities,
          boards,
          projects,
          getTransformation
        ),
      })),
    }))
    // validatePipeline()
  }, [
    blueprintConnections,
    dataEntities,
    boards,
    projects,
    getTransformation,
    validatePipeline,
    createProviderScopes,
    setBlueprintSettings,
  ])

  useEffect(() => {
    console.log('>> PROJECTS LIST', projects)
    console.log('>> BOARDS LIST', boards)
    projects[configuredConnection?.id]?.map(
      (p) => initDefaultTransformationSettingsIfNotExist(configuredConnection, p)
    )
    boards[configuredConnection?.id]?.map(
      (b) => initDefaultTransformationSettingsIfNotExist(configuredConnection, b)
    )
  }, [initDefaultTransformationSettingsIfNotExist, projects, boards, configuredConnection])

  useEffect(() => {
    console.log(
      '>>> SELECTED PROJECT TO CONFIGURE...',
      configuredProject,
    )
    setActiveTransformation((aT) =>
      configuredProject ? getTransformation(configuredConnection, configuredProject) : aT
    )
  }, [configuredProject, getTransformation])

  useEffect(() => {
    console.log(
      '>>> SELECTED BOARD TO CONFIGURE...',
      configuredBoard?.id,
    )
    setActiveTransformation((aT) =>
      configuredBoard ? getTransformation(configuredConnection, configuredBoard) : aT
    )
  }, [configuredBoard, getTransformation])

  // useEffect(() => {
  //   console.log(
  //     '>>> ACTIVE/MODIFYING TRANSFORMATION RULES...',
  //     activeTransformation
  //   )
  // }, [activeTransformation])

  useEffect(() => {
    console.log('>>> BLUEPRINT WORKFLOW STEPS...', blueprintSteps)
  }, [blueprintSteps])

  useEffect(() => {
    setBlueprintSteps(
      advancedMode ? blueprintAdvancedSteps : blueprintNormalSteps
    )
    setBlueprintMode(
      advancedMode ? BlueprintMode.ADVANCED : BlueprintMode.NORMAL
    )
  }, [
    advancedMode,
    blueprintNormalSteps,
    blueprintAdvancedSteps,
    setBlueprintMode,
  ])

  useEffect(() => {
    if (isValidCode()) {
      setRunTasksAdvanced(JSON.parse(rawConfiguration))
    }
  }, [rawConfiguration, isValidCode])

  useEffect(() => {
    if (saveBlueprintComplete?.id) {
      setNewBlueprintId(saveBlueprintComplete?.id)
    }
  }, [
    saveBlueprintComplete,
    setPipelineSettings,
    setPipelineSettings,
    setNewBlueprintId,
    runNow,
    runPipeline,
    history
  ])

  useEffect(() => {
    if (runNow && newBlueprintId) {
      const newPipelineConfiguration = {
        name: `${saveBlueprintComplete?.name} ${Date.now()}`,
        blueprintId: saveBlueprintComplete?.id,
        plan: saveBlueprintComplete?.plan,
      }
      runPipeline(newPipelineConfiguration)
      setRunNow(false)
      history.push(`/blueprints/detail/${saveBlueprintComplete?.id}`)
    } else if (newBlueprintId) {
      history.push(`/blueprints/detail/${saveBlueprintComplete?.id}`)
    }
  }, [
    runNow,
    saveBlueprintComplete,
    newBlueprintId,
    runPipeline,
    history
  ])

  useEffect(() => {
    console.log('>>> FETCHED JIRA API BOARDS FROM PROXY...', jiraApiBoards)
    setBoardsList(jiraApiBoards)
  }, [jiraApiBoards])

  useEffect(() => {
    if (saveConnectionComplete?.id && connectionDialogIsOpen) {
      handleConnectionDialogClose()
      fetchAllConnections(false, true)
    }
  }, [connectionDialogIsOpen, saveConnectionComplete, fetchAllConnections, handleConnectionDialogClose])

  useEffect(() => {
    console.log('>>> CONNECTIONS SELECTOR LIST UPDATED...', connectionsList)
  }, [connectionsList])

  return (
    <>
      <div className='container'>
        <Nav />
        <Sidebar />
        <Content>
          <main className='main'>
            <WorkflowStepsBar activeStep={activeStep} steps={blueprintSteps} />

            <div
              className={`workflow-content workflow-step-id-${activeStep?.id}`}
            >
              {advancedMode ? (
                <>
                  {activeStep?.id === 1 && (
                    <AdvancedJSON
                      activeStep={activeStep}
                      advancedMode={advancedMode}
                      runTasksAdvanced={runTasksAdvanced}
                      blueprintConnections={blueprintConnections}
                      connectionsList={connectionsList}
                      name={name}
                      setBlueprintName={setBlueprintName}
                      // setBlueprintConnections={setBlueprintConnections}
                      fieldHasError={fieldHasError}
                      getFieldError={getFieldError}
                      // addConnection={addConnection}
                      // manageConnection={manageConnection}
                      onAdvancedMode={handleAdvancedMode}
                      // @todo add multistage checker method
                      isMultiStagePipeline={() => { }}
                      rawConfiguration={rawConfiguration}
                      setRawConfiguration={setRawConfiguration}
                      isSaving={isSaving}
                      isValidConfiguration={isValidConfiguration}
                      validationAdvancedError={validationAdvancedError}
                      validationErrors={validationErrors}
                    />
                  )}

                  {activeStep?.id === 2 && (
                    <DataSync
                      activeStep={activeStep}
                      advancedMode={advancedMode}
                      cronConfig={cronConfig}
                      customCronConfig={customCronConfig}
                      createCron={createCron}
                      setCronConfig={setCronConfig}
                      getCronPreset={getCronPreset}
                      fieldHasError={fieldHasError}
                      getFieldError={getFieldError}
                      setCustomCronConfig={setCustomCronConfig}
                      getCronPresetByConfig={getCronPresetByConfig}
                    />
                  )}
                </>
              ) : (
                <>
                  {activeStep?.id === 1 && (
                    <DataConnections
                      activeStep={activeStep}
                      advancedMode={advancedMode}
                      blueprintConnections={blueprintConnections}
                      onlineStatus={onlineStatus}
                      connectionsList={connectionsList}
                      name={name}
                      setBlueprintName={setBlueprintName}
                      setBlueprintConnections={setBlueprintConnections}
                      fieldHasError={fieldHasError}
                      getFieldError={getFieldError}
                      addConnection={addConnection}
                      manageConnection={manageConnection}
                      onAdvancedMode={handleAdvancedMode}
                    />
                  )}

                  {activeStep?.id === 2 && (
                    <DataScopes
                      provider={provider}
                      activeStep={activeStep}
                      advancedMode={advancedMode}
                      activeConnectionTab={activeConnectionTab}
                      blueprintConnections={blueprintConnections}
                      dataEntitiesList={dataEntitiesList}
                      boardsList={boardsList}
                      boards={boards}
                      dataEntities={dataEntities}
                      projects={projects}
                      configuredConnection={configuredConnection}
                      handleConnectionTabChange={handleConnectionTabChange}
                      setDataEntities={setDataEntities}
                      setProjects={setProjects}
                      setBoards={setBoards}
                      prevStep={prevStep}
                      isSaving={isSaving}
                      isRunning={isRunning}
                    />
                  )}

                  {activeStep?.id === 3 && (
                    <DataTransformations
                      provider={provider}
                      activeStep={activeStep}
                      advancedMode={advancedMode}
                      activeConnectionTab={activeConnectionTab}
                      blueprintConnections={blueprintConnections}
                      dataEntities={dataEntities}
                      projects={projects}
                      boardsList={boardsList}
                      boards={boards}
                      issueTypes={jiraApiIssueTypes}
                      fields={jiraApiFields}
                      configuredConnection={configuredConnection}
                      configuredProject={configuredProject}
                      configuredBoard={configuredBoard}
                      handleConnectionTabChange={handleConnectionTabChange}
                      prevStep={prevStep}
                      addBoardTransformation={addBoardTransformation}
                      addProjectTransformation={addProjectTransformation}
                      checkTransformationIsExist={checkTransformationIsExist}
                      activeTransformation={activeTransformation}
                      isSaving={isSaving}
                      isSavingConnection={isSavingConnection}
                      isRunning={isRunning}
                      onSave={handleTransformationSave}
                      onCancel={handleTransformationCancel}
                      onClear={handleTransformationClear}
                      jiraProxyError={jiraProxyError}
                      isFetchingJIRA={isFetchingJIRA}
                    />
                  )}

                  {activeStep?.id === 4 && (
                    <DataSync
                      activeStep={activeStep}
                      advancedMode={advancedMode}
                      cronConfig={cronConfig}
                      customCronConfig={customCronConfig}
                      createCron={createCron}
                      setCronConfig={setCronConfig}
                      getCronPreset={getCronPreset}
                      fieldHasError={fieldHasError}
                      getFieldError={getFieldError}
                      setCustomCronConfig={setCustomCronConfig}
                      getCronPresetByConfig={getCronPresetByConfig}
                    />
                  )}
                </>
              )}
            </div>

            <WorkflowActions
              activeStep={activeStep}
              blueprintSteps={blueprintSteps}
              advancedMode={advancedMode}
              setShowBlueprintInspector={setShowBlueprintInspector}
              validationErrors={validationErrors}
              onNext={nextStep}
              onPrev={prevStep}
              onSave={handleBlueprintSave}
              onSaveAndRun={handleBlueprintSaveAndRun}
              isLoading={isSaving}
            />
          </main>
        </Content>
      </div>

      <ConnectionDialog
        integrations={integrationsData}
        activeProvider={activeProvider}
        setProvider={setActiveProvider}
        setTestStatus={setTestStatus}
        setTestResponse={setTestResponse}
        connection={managedConnection}
        errors={connectionErrors}
        validationErrors={connectionValidationErrors}
        endpointUrl={endpointUrl}
        name={connectionName}
        proxy={proxy}
        token={token}
        initialTokenStore={initialTokenStore}
        username={username}
        password={password}
        isOpen={connectionDialogIsOpen}
        isTesting={isTestingConnection}
        isSaving={isSavingConnection}
        isValid={isValidConnection}
        onClose={handleConnectionDialogClose}
        onOpen={handleConnectionDialogOpen}
        onTest={testConnection}
        onSave={saveConnection}
        onValidate={validateConnection}
        onNameChange={setName}
        onEndpointChange={setEndpointUrl}
        onProxyChange={setProxy}
        onTokenChange={setToken}
        onUsernameChange={setUsername}
        onPasswordChange={setPassword}
        testStatus={testStatus}
        testResponse={testResponse}
        allTestResponses={allTestResponses}
      />

      <CodeInspector
        title={name}
        titleIcon='add'
        subtitle='JSON CONFIGURATION'
        isOpen={showBlueprintInspector}
        activePipeline={
          !advancedMode
            ? {
              // ID: 0,
              name,
              // tasks: blueprintTasks,
              settings: blueprintSettings,
              cronConfig,
              enable,
              mode,
            }
            : {
              name,
              plan: blueprintTasks,
              cronConfig,
              enable,
              mode,
            }
        }
        onClose={setShowBlueprintInspector}
        hasBackdrop={false}
      />
    </>
  )
}

export default CreateBlueprint
