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
import React, { Fragment, useEffect, useState, useCallback, useMemo } from 'react'
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
import { Intent } from '@blueprintjs/core'
import {
  Providers,
} from '@/data/Providers'
import Nav from '@/components/Nav'
import Sidebar from '@/components/Sidebar'
// import AppCrumbs from '@/components/Breadcrumbs'
import Content from '@/components/Content'
import { ToastNotification } from '@/components/Toast'

import { BlueprintMode } from '@/data/NullBlueprint'
import { NullBlueprintConnection } from '@/data/NullBlueprintConnection'
import { NullConnection } from '@/data/NullConnection'

import {
  WorkflowSteps,
  WorkflowAdvancedSteps,
  DEFAULT_DATA_ENTITIES
} from '@/data/BlueprintWorkflow'

import useBlueprintManager from '@/hooks/useBlueprintManager'
import usePipelineManager from '@/hooks/usePipelineManager'
import useConnectionManager from '@/hooks/useConnectionManager'
import useDataScopesManager from '@/hooks/useDataScopesManager'
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

import { DEVLAKE_ENDPOINT } from '@/utils/config'
import request from '@/utils/request'

// import ConnectionTabs from '@/components/blueprints/ConnectionTabs'

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

  const [isValidConfiguration, setIsValidConfiguration] = useState(false)
  const [validationAdvancedError, setValidationAdvancedError] = useState()

  const [connectionDialogIsOpen, setConnectionDialogIsOpen] = useState(false)
  const [managedConnection, setManagedConnection] = useState(
    NullBlueprintConnection
  )

  const [dataEntitiesList, setDataEntitiesList] = useState([
    ...DEFAULT_DATA_ENTITIES,
  ])
  const [boardsList, setBoardsList] = useState([])

  const [blueprintConnections, setBlueprintConnections] = useState([])
  const [configuredConnection, setConfiguredConnection] = useState()

  const [activeConnectionTab, setActiveConnectionTab] = useState()

  const [onlineStatus, setOnlineStatus] = useState([])
  const [showBlueprintInspector, setShowBlueprintInspector] = useState(false)
  // const [dataScopes, setDataScopes] = useState([])
  const [dataConnections, setDataConnections] = useState([])
  const [connectionId, setConnectionId] = useState('')

  const [canAdvanceNext, setCanAdvanceNext] = useState(true)
  // eslint-disable-next-line no-unused-vars
  const [canAdvancePrev, setCanAdvancePrev] = useState(true)

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
    rawConfiguration,
    setName: setBlueprintName,
    setCronConfig,
    setCustomCronConfig,
    setTasks: setBlueprintTasks,
    setSettings: setBlueprintSettings,
    // eslint-disable-next-line no-unused-vars
    setDetectedProviderTasks,
    setEnable: setEnableBlueprint,
    setMode: setBlueprintMode,
    setIsManual: setIsManualBlueprint,
    setRawConfiguration,
    // eslint-disable-next-line no-unused-vars
    isFetching: isFetchingBlueprints,
    isSaving,
    createCronExpression: createCron,
    // eslint-disable-next-line no-unused-vars
    getCronSchedule: getSchedule,
    // eslint-disable-next-line no-unused-vars
    getNextRunDate,
    getCronPreset,
    getCronPresetByConfig,
    saveBlueprint,
    // eslint-disable-next-line no-unused-vars
    deleteBlueprint,
    isDeleting: isDeletingBlueprint,
    isManual: isManualBlueprint,
    saveComplete: saveBlueprintComplete
  } = useBlueprintManager()

  const {
    boards,
    projects,
    entities: dataEntities,
    transformations,
    setBoards,
    setProjects,
    setEntities: setDataEntities,
    setTransformations,
    createProviderScopes,
    createProviderConnections,
    getDefaultTransformations,
    initializeTransformations
  } = useDataScopesManager({ connection: configuredConnection, settings: blueprintSettings })

  const {
    pipelineName,
    // pipelines,
    runPipeline,
    // cancelPipeline,
    // fetchPipeline,
    // fetchAllPipelines,
    // pipelineRun,
    // buildPipelineStages,
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
    isValid: isValidPipeline,
    detectedProviders,
    parseJSON
  } = usePipelineValidation({
    enabledProviders,
    pipelineName,
    projects,
    boards,
    connectionId,
    tasks: runTasks,
    tasksAdvanced: runTasksAdvanced,
    advancedMode,
    mode,
    connection: configuredConnection,
    entities: dataEntities,
    rawConfiguration
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
    // eslint-disable-next-line no-unused-vars
    testAllConnections,
    saveConnection,
    // eslint-disable-next-line no-unused-vars
    fetchConnection,
    // eslint-disable-next-line no-unused-vars
    allProviderConnections,
    connectionsList,
    errors: connectionErrors,
    isSaving: isSavingConnection,
    isTesting: isTestingConnection,
    isFetching: isFetchingConnection,
    showError,
    testStatus,
    name: connectionName,
    endpointUrl,
    proxy,
    rateLimit,
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
    setRateLimit,
    setUsername,
    setPassword,
    setToken,
    setInitialTokenStore,
    setTestStatus,
    setTestResponse,
    setAllTestResponses,
    setConnectionsList,
    setSaveComplete: setSaveConnectionComplete,
    fetchAllConnections,
    clearConnection: clearActiveConnection,
    testResponse,
    // eslint-disable-next-line no-unused-vars
    testedConnections,
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
    validate: validateBlueprint,
    errors: blueprintValidationErrors,
    isValid: isValidBlueprint,
    fieldHasError,
    getFieldError,
  } = useBlueprintValidation({
    name,
    boards,
    projects,
    cronConfig,
    customCronConfig,
    enable,
    tasks: blueprintTasks,
    mode,
    connections: blueprintConnections,
    entities: dataEntities,
    activeStep,
    activeProvider: provider,
    activeConnection: configuredConnection
  })

  const {
    validate: validateConnection,
    errors: connectionValidationErrors,
    isValid: isValidConnection,
  } = useConnectionValidation({
    activeProvider,
    name: connectionName,
    endpointUrl,
    proxy,
    rateLimit,
    token,
    username,
    password,
  })

  const [configuredProject, setConfiguredProject] = useState(
    // projects.length > 0 ? projects[0] : null
    null
  )
  const [configuredBoard, setConfiguredBoard] = useState(
    // boards.length > 0 ? boards[0] : null
    null
  )

  const activeTransformation = useMemo(() => transformations[configuredProject || configuredBoard?.id], [transformations, configuredProject, configuredBoard?.id])

  // eslint-disable-next-line no-unused-vars
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
    setConfiguredProject(null)
    setConfiguredBoard(null)
  }, [blueprintSteps])

  const testSelectedConnections = useCallback((connections, savedConnection = {}, callback = () => {}) => {
    const runTest = async () => {
      const results = await Promise.all(connections.map(
        c => {
          const testPayload = c.connectionId === savedConnection?.id && c.name === savedConnection?.name ? {
            endpoint: savedConnection?.endpoint,
            username: savedConnection?.username,
            password: savedConnection?.password,
            token: savedConnection?.token,
            proxy: savedConnection?.proxy
          } : {
            endpoint: c.endpoint,
            username: c.username,
            password: c.password,
            token: c.token,
            proxy: c.proxy
          }
          return request.post(`${DEVLAKE_ENDPOINT}/plugins/${c.plugin}/test`, testPayload)
        })
      )
      setOnlineStatus(results.map(r => r))
    }
    if (mode === BlueprintMode.NORMAL && connections.length > 0) {
      runTest()
    }
    callback()
  }, [mode])

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

  const handleConnectionDialogOpen = useCallback(() => {
    console.log('>>> MANAGING CONNECTION', managedConnection)
  }, [managedConnection])

  const handleConnectionDialogClose = useCallback((savedConnection = {}) => {
    fetchAllConnections(false, true)
    testSelectedConnections(blueprintConnections, savedConnection)
    setConnectionDialogIsOpen(false)
    setManagedConnection(NullBlueprintConnection)
    setTestStatus(0)
    setTestResponse(null)
    setAllTestResponses({})
    setInitialTokenStore({})
    clearActiveConnection()
    setActiveConnection(NullConnection)
    // setSaveConnectionComplete(null)
  }, [
    blueprintConnections,
    testSelectedConnections,
    fetchAllConnections,
    clearActiveConnection,
    setActiveConnection,
    setAllTestResponses,
    setInitialTokenStore,
    // setSaveConnectionComplete,
    setTestResponse,
    setTestStatus
  ])

  const handleTransformationCancel = useCallback(() => {
    setConfiguredProject(null)
    setConfiguredBoard(null)
    console.log('>> Cancel Modify - Transformation Settings')
  }, [setConfiguredProject, setConfiguredBoard])

  const handleTransformationClear = useCallback(() => {
    console.log(
      '>>> CLEARING TRANSFORMATION RULES!',
      '==> PROJECT =',
      configuredProject,
      '==> BOARD =',
      configuredBoard
    )
    setTransformations((existingTransformations) => ({
      ...existingTransformations,
      [configuredProject]: {},
      [configuredBoard?.id]: {},
    }))
    setConfiguredProject(null)
    setConfiguredBoard(null)
  }, [setTransformations, configuredProject, configuredBoard])

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

  const manageConnection = useCallback(
    (connection) => {
      console.log('>> MANAGE CONNECTION...', connection)
      if (connection?.id !== null) {
        setActiveProvider(
          integrationsData.find((p) => p.id === connection.provider)
        )
        setProvider(integrationsData.find((p) => p.id === connection.provider))
        setManagedConnection(connection)
        setConnectionDialogIsOpen(true)
      }
    },
    [setProvider]
  )

  const addProjectTransformation = useCallback((project) => {
    setConfiguredProject(project)
    ToastNotification.clear()
  }, [setConfiguredProject])

  const addBoardTransformation = useCallback((board) => {
    setConfiguredBoard(board)
    ToastNotification.clear()
  }, [setConfiguredBoard])

  const addConnection = useCallback(() => {
    setManagedConnection(NullBlueprintConnection)
    setConnectionDialogIsOpen(true)
  }, [])

  const setTransformationSettings = useCallback(
    (settings, configuredEntity) => {
      console.log(
        '>> SETTING TRANSFORMATION SETTINGS PROJECT/BOARD...',
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

  const handleTransformationSave = useCallback((settings, entity) => {
    console.log('>> SAVING / CLOSING Transformation Settings')
    // manual @save disabled, reactive auto-saving writes settings to transform object...
    // setTransformationSettings(settings, entity)
    setConfiguredProject(null)
    setConfiguredBoard(null)
    ToastNotification.clear()
    ToastNotification.show({ message: 'Transformation Rules Added.', intent: Intent.SUCCESS, icon: 'small-tick' })
  }, [])

  const handleAdvancedMode = (enableAdvanced = true) => {
    setAdvancedMode(enableAdvanced)
  }

  const isValidCode = useCallback(() => {
    let isValid = false
    try {
      const parsedCode = parseJSON(rawConfiguration)
      isValid = true
      setValidationAdvancedError(null)
    } catch (e) {
      console.log('>> FORMAT CODE: Invalid Code Format!', e)
      isValid = false
      setValidationAdvancedError(e.message)
    }
    setIsValidConfiguration(isValid)
    return isValid
  }, [rawConfiguration, parseJSON])

  useEffect(() => {
    console.log('>> ACTIVE STEP CHANGED: ', activeStep)
    if (mode === BlueprintMode.NORMAL && activeStep?.id === 1) {
      const enableNotifications = false
      const getAllSources = true
      fetchAllConnections(enableNotifications, getAllSources)
    }
    if (mode === BlueprintMode.NORMAL &&
        ([2, 3].includes(activeStep?.id)) &&
        enabledProviders.includes(Providers.JIRA)
    ) {
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
    mode
  ])

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
    connectionsList,
    enable,
    validateBlueprint,
  ])

  useEffect(() => {
    setIsManualBlueprint(cronConfig === 'manual')
  }, [cronConfig, setIsManualBlueprint])

  useEffect(() => { }, [activeConnectionTab])

  useEffect(() => {
    console.log('>>>> MY SELECTED BLUEPRINT CONNECTIONS...', blueprintConnections)
    const someConnection = blueprintConnections.find(c => c)
    if (someConnection) {
      setConfiguredConnection(someConnection)
      setActiveConnectionTab(`connection-${someConnection?.id}`)
      setActiveProvider(
        integrationsData.find((p) => p.id === someConnection.provider)
      )
      setProvider(
        integrationsData.find((p) => p.id === someConnection.provider)
      )
    }
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
        case Providers.TAPD:
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

    testSelectedConnections(blueprintConnections)
  }, [
    blueprintConnections,
    setProvider,
    setBoards,
    setDataEntities,
    setProjects,
    testSelectedConnections
  ])

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
  }, [configuredConnection, setActiveConnectionTab])

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
          transformations
        ),
      })),
    }))
    // validatePipeline()
  }, [
    blueprintConnections,
    dataEntities,
    boards,
    projects,
    transformations,
    validatePipeline,
    createProviderScopes,
    setBlueprintSettings,
  ])

  useEffect(() => {
    console.log('>> PROJECTS LIST', projects)
    console.log('>> BOARDS LIST', boards)

    const projectTransformation = projects[configuredConnection?.id]
    const boardTransformation = boards[configuredConnection?.id]?.map(
      (b) => b.id
    )
    if (projectTransformation) {
      setTransformations((cT) => ({
        ...projectTransformation.reduce(initializeTransformations, {}),
        // Spread Current/Existing Transformations Settings
        ...cT,
      }))
    }
    if (boardTransformation) {
      setTransformations((cT) => ({
        ...boardTransformation.reduce(initializeTransformations, {}),
        // Spread Current/Existing Transformations Settings
        ...cT,
      }))
    }
  }, [
    projects,
    boards,
    configuredConnection,
    initializeTransformations,
    setTransformations
  ])

  useEffect(() => {
    console.log(
      '>>> SELECTED PROJECT TO CONFIGURE...',
      configuredProject,
    )
    // setActiveTransformation((aT) =>
    //   configuredProject !== null ? transformations[configuredProject] : {}
    // )
    setCanAdvanceNext(!configuredProject)
  }, [configuredProject, setCanAdvanceNext])

  useEffect(() => {
    console.log(
      '>>> SELECTED BOARD TO CONFIGURE...',
      configuredBoard?.id,
    )
    // setActiveTransformation((aT) =>
    //   configuredBoard ? transformations[configuredBoard?.id] : aT
    // )
    setCanAdvanceNext(!configuredBoard)
  }, [configuredBoard, setCanAdvanceNext])

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
      handleConnectionDialogClose(saveConnectionComplete)
      // fetchAllConnections(false, true)
      // testSelectedConnections(blueprintConnections, saveConnectionComplete)
    }
    return () => setSaveConnectionComplete(null)
  }, [
    connectionDialogIsOpen,
    saveConnectionComplete,
    fetchAllConnections,
    blueprintConnections,
    testSelectedConnections,
    handleConnectionDialogClose,
    setSaveConnectionComplete
  ])

  useEffect(() => {
    console.log('>>> ONLINE STATUS UPDATED...', onlineStatus)
    setDataConnections(blueprintConnections.map((c, cIdx) => ({
      ...c,
      statusResponse: onlineStatus[cIdx],
      status: onlineStatus[cIdx]?.status
    })))
  }, [onlineStatus, blueprintConnections])

  useEffect(() => {
    setConnectionsList(cList => cList.map((c, cIdx) => ({
      ...c,
      statusResponse: dataConnections.find(dC => dC.id === c.id && dC.provider === c.provider),
      status: dataConnections.find(dC => dC.id === c.id && dC.provider === c.provider)?.status
    })))
    setCanAdvanceNext(dataConnections.every(dC => dC.status === 200))
  }, [dataConnections, setConnectionsList])

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
                      blueprintConnections={dataConnections}
                      // blueprintConnections={blueprintConnections}
                      // blueprintConnections={[...blueprintConnections.map((c, cIdx) => ({...c, statusResponse: onlineStatus[cIdx], status: onlineStatus[cIdx]?.status }))]}
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
                      isTesting={isTestingConnection}
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
                      validationErrors={[...validationErrors, ...blueprintValidationErrors]}
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
                      transformations={transformations}
                      activeTransformation={activeTransformation}
                      setTransformations={setTransformations}
                      setTransformationSettings={setTransformationSettings}
                      isSaving={isSaving}
                      isSavingConnection={isSavingConnection}
                      isRunning={isRunning}
                      onSave={handleTransformationSave}
                      onCancel={handleTransformationCancel}
                      onClear={handleTransformationClear}
                      fieldHasError={fieldHasError}
                      getFieldError={getFieldError}
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
              validationErrors={[...validationErrors, ...blueprintValidationErrors]}
              onNext={nextStep}
              onPrev={prevStep}
              onSave={handleBlueprintSave}
              onSaveAndRun={handleBlueprintSaveAndRun}
              isLoading={isSaving || isFetchingJIRA || isFetchingConnection || isTestingConnection}
              isValid={advancedMode ? isValidBlueprint && isValidPipeline : isValidBlueprint}
              canGoNext={canAdvanceNext}
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
        rateLimit={rateLimit}
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
        onRateLimitChange={setRateLimit}
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
                cronConfig: isManualBlueprint ? '0 0 * * *' : (cronConfig === 'custom' ? customCronConfig : cronConfig),
                enable,
                mode,
                isManual: isManualBlueprint
              }
            : {
                name,
                plan: blueprintTasks,
                cronConfig: isManualBlueprint ? '0 0 * * *' : (cronConfig === 'custom' ? customCronConfig : cronConfig),
                enable,
                mode,
                isManual: isManualBlueprint
              }
        }
        onClose={setShowBlueprintInspector}
        hasBackdrop={false}
      />
    </>
  )
}

export default CreateBlueprint
