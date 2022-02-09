import React, { Fragment, useEffect, useCallback, useState } from 'react'
import { CSSTransition } from 'react-transition-group'
import {
  useHistory,
  useLocation,
  Link,
} from 'react-router-dom'
import { GRAFANA_URL } from '@/utils/config'
import {
  Button, Icon, Intent, Switch,
  FormGroup,
  ButtonGroup,
  InputGroup,
  Popover,
  Tooltip,
  Position,
  Colors,
} from '@blueprintjs/core'
import {
  Providers,
} from '@/data/Providers'
import { integrationsData } from '@/data/integrations'
import usePipelineManager from '@/hooks/usePipelineManager'
import usePipelineValidation from '@/hooks/usePipelineValidation'
import useConnectionManager from '@/hooks/useConnectionManager'
import FormValidationErrors from '@/components/messages/FormValidationErrors'
import PipelineIndicator from '@/components/widgets/PipelineIndicator'
import PipelinePresetsMenu from '@/components/menus/PipelinePresetsMenu'
import ProviderSettings from '@/components/pipelines/ProviderSettings'
import Nav from '@/components/Nav'
import Sidebar from '@/components/Sidebar'
import AppCrumbs from '@/components/Breadcrumbs'
import Content from '@/components/Content'
import { ReactComponent as HelpIcon } from '@/images/help.svg'
import { ReactComponent as BackArrowIcon } from '@/images/undo.svg'
import RunPipelineIcon from '@/images/duplicate.png'

import GitlabHelpNote from '@/images/help/gitlab-help.png'
import JiraHelpNote from '@/images/help/jira-help.png'
import GithubHelpNote from '@/images/help/github-help.png'

import '@/styles/pipelines.scss'

const CreatePipeline = (props) => {
  const history = useHistory()
  const location = useLocation()
  // const { providerId } = useParams()
  // const [activeProvider, setActiveProvider] = useState(integrationsData[0])
  const [integrations, setIntegrations] = useState(integrationsData)
  const [jiraIntegration, setJiraIntegration] = useState(integrationsData.find(p => p.id === Providers.JIRA))

  const [today, setToday] = useState(new Date())
  const pipelinePrefixes = ['COLLECT', 'SYNC']
  const pipelineSuffixes = [
    today.getTime(), // 1639630123107
    today.toString(), // Wed Dec 15 2021 23:48:43 GMT-0500 (EST)
    today.toISOString(), // 2021-12-16T04:48:43.107Z
    // eslint-disable-next-line max-len
    `${today.getFullYear()}${today.getMonth() + 1}${today.getDate()}${today.getHours()}${today.getMinutes()}${today.getSeconds()}`, // 202112154936
    today.toUTCString(), // Thu, 16 Dec 2021 04:49:52 GMT
  ]

  const [readyProviders, setReadyProviders] = useState([])
  // const [isRunning, setIsRunning] = useState(false)
  // const [autoRun, setAutoRun] = useState(false)
  // const [enableThrottling, setEnableThrottling] = useState(false)

  const [enabledProviders, setEnabledProviders] = useState([])
  const [runTasks, setRunTasks] = useState([])
  const [existingTasks, setExistingTasks] = useState([])

  const [namePrefix, setNamePrefix] = useState(pipelinePrefixes[0])
  const [nameSuffix, setNameSuffix] = useState(pipelineSuffixes[0])
  const [pipelineName, setPipelineName] = useState(`${namePrefix} ${nameSuffix}`)
  const [projectId, setProjectId] = useState([])
  const [boardId, setBoardId] = useState([])
  const [sourceId, setSourceId] = useState('')
  const [sources, setSources] = useState([])
  const [selectedSource, setSelectedSource] = useState()
  const [repositoryName, setRepositoryName] = useState('')
  const [owner, setOwner] = useState('')

  const [autoRedirect, setAutoRedirect] = useState(true)
  const [restartDetected, setRestartDetected] = useState(false)

  const {
    runPipeline,
    cancelPipeline,
    fetchPipeline,
    pipelineRun,
    isRunning,
    errors: pipelineErrors,
    setSettings: setPipelineSettings,
    lastRunId
  } = usePipelineManager(pipelineName, runTasks)

  const {
    validate,
    errors: validationErrors,
    isValid: isValidPipelineForm
  } = usePipelineValidation({
    enabledProviders,
    pipelineName,
    projectId,
    boardId,
    owner,
    repositoryName,
    sourceId,
    runTasks
  })

  const {
    allConnections,
    fetchAllConnections
  } = useConnectionManager({
    activeProvider: jiraIntegration
  })

  useEffect(() => {
    integrationsData.forEach((i, idx) => {
      setTimeout(() => {
        setReadyProviders(r => [...r, i.id])
      }, idx * 50)
    })
  }, [])

  const isProviderEnabled = (providerId) => {
    return enabledProviders.includes(providerId)
  }

  const isValidPipeline = () => {
    return enabledProviders.length >= 1 &&
      pipelineName !== '' &&
      pipelineName.length > 2 &&
      validationErrors.length === 0
  }

  const getManyProviderOptions = useCallback((providerId, optionProperty, ids, options = {}) => {
    return ids.map(id => {
      return {
        Plugin: providerId,
        Options: {
          [optionProperty]: parseInt(id, 10),
          ...options
        }
      }
    })
  }, [])

  const getProviderOptions = useCallback((providerId) => {
    let options = {}
    switch (providerId) {
      case Providers.JENKINS:
        // NO OPTIONS for Jenkins!
        break
      case Providers.JIRA:
        options = {
          boardId: parseInt(boardId, 10),
          sourceId: parseInt(sourceId, 10)
        }
        break
      case Providers.GITHUB:
        options = {
          repo: repositoryName,
          owner
        }
        break
      case Providers.GITLAB:
        options = {
          projectId: parseInt(projectId, 10)
        }
        break
      default:
        break
    }
    return options
  }, [boardId, owner, projectId, repositoryName, sourceId])

  const configureProvider = useCallback((providerId) => {
    let providerConfig = {}
    switch (providerId) {
      case Providers.GITLAB:
        providerConfig = getManyProviderOptions(providerId, 'projectId', [...projectId])
        break
      case Providers.JIRA:
        providerConfig = getManyProviderOptions(
          providerId,
          'boardId',
          [...boardId],
          {
            sourceId: parseInt(sourceId, 10)
          }
        )
        break
      default:
        providerConfig = {
          Plugin: providerId,
          Options: {
            ...getProviderOptions(providerId)
          }
        }
        break
    }
    return providerConfig
  }, [getProviderOptions, getManyProviderOptions, projectId, boardId, sourceId])

  const resetPipelineName = () => {
    setToday(new Date())
    setPipelineName(`${namePrefix} ${nameSuffix}`)
  }

  const resetConfiguration = () => {
    window.history.replaceState(null, '')
    resetPipelineName()
    setExistingTasks([])
    setEnabledProviders([])
    setProjectId([])
    setBoardId([])
    setSelectedSource(null)
    setRepositoryName('')
    setOwner('')
  }

  const getConnectionName = useCallback((connectionId) => {
    const source = sources.find(s => s.id === connectionId)
    return source ? source.title : '(Instance)'
  }, [sources])

  useEffect(() => {

  }, [pipelineName])

  useEffect(() => {
    console.log('>> PIPELINE RUN TASK SETTINGS FOR PIPELINE MANAGER ....', runTasks)
    setPipelineSettings({
      name: pipelineName,
      tasks: [
        [...runTasks]
      ]
    })
    validate()
  }, [runTasks, pipelineName, setPipelineSettings, validate])

  useEffect(() => {
    console.log('>> ENBALED PROVIDERS = ', enabledProviders)
    const PipelineTasks = enabledProviders.map(p => configureProvider(p))
    setRunTasks(PipelineTasks.flat())
    console.log('>> CONFIGURED PIPELINE TASKS = ', PipelineTasks)
    validate()
    if (enabledProviders.includes(Providers.JIRA)) {
      fetchAllConnections(false)
    } else {
      setSources([])
      setSelectedSource(null)
    }
  }, [enabledProviders, projectId, boardId, sourceId, owner, repositoryName, configureProvider, validate, fetchAllConnections])

  useEffect(() => {
    console.log('>> PIPELINE LAST RUN OBJECT CHANGED!!...', pipelineRun)
    if (pipelineRun.ID && autoRedirect) {
      history.push(`/pipelines/activity/${pipelineRun.ID}`)
    }
  }, [pipelineRun, autoRedirect, history])

  useEffect(() => {
    console.log(namePrefix, nameSuffix)
    setPipelineName(`${namePrefix} ${nameSuffix}`)
    setToday(new Date())
  }, [namePrefix, nameSuffix])

  useEffect(() => {
    console.log('>> JIRA SOURCE ID SELECTED, CONNECTION INSTANCE = ', selectedSource)
    setSourceId(sId => selectedSource ? selectedSource.value : null)
    validate()
  }, [selectedSource, validate])

  useEffect(() => {
    console.log('>> FETCHED ALL JIRA CONNECTIONS... ', allConnections)
    setSources(allConnections.map(c => { return { id: c.ID, title: c.name || 'Instance', value: c.ID } }))
  }, [allConnections])

  useEffect(() => {
    console.log('>> BUILT JIRA INSTANCE SELECT MENU... ', sources)
  }, [sources])

  useEffect(() => {
    if (location.state?.existingTasks) {
      console.log('>> RESTART ATTEMPT: DETECTED EXISTING PIPELINE CONFIGURATION... ', location.state.existingTasks)
      const tasks = location.state.existingTasks
      setRestartDetected(true)
      setExistingTasks(tasks)
      window.history.replaceState(null, '')
      // !WARNING! This logic will only handle ONE STAGE (Stage 1)
      // @todo: refactor later for multi-stage
      const GitLabTask = tasks.filter(t => t.plugin === Providers.GITLAB)
      const GitHubTask = tasks.find(t => t.plugin === Providers.GITHUB)
      const JiraTask = tasks.filter(t => t.plugin === Providers.JIRA)
      const JenkinsTask = tasks.find(t => t.plugin === Providers.JENKINS)
      if (GitLabTask && GitLabTask.length > 0) {
        setEnabledProviders(eP => [...eP, Providers.GITLAB])
        setProjectId(Array.isArray(GitLabTask) ? GitLabTask.map(gT => gT.options?.projectId) : GitLabTask.options?.projectId)
      }
      if (GitHubTask) {
        setEnabledProviders(eP => [...eP, Providers.GITHUB])
        setRepositoryName(GitHubTask.options?.repositoryName || GitHubTask.options?.repo)
        setOwner(GitHubTask.options?.owner)
      }
      if (JiraTask && JiraTask.length > 0) {
        fetchAllConnections(false)
        setEnabledProviders(eP => [...eP, Providers.JIRA])
        setBoardId(Array.isArray(JiraTask) ? JiraTask.map(jT => jT.options?.boardId) : JiraTask.options?.boardId)
        const connSrcId = JiraTask[0].options?.sourceId
        setSelectedSource({
          id: parseInt(connSrcId, 10),
          title: getConnectionName(connSrcId),
          value: parseInt(connSrcId, 10)
        })
      }
      if (JenkinsTask) {
        setEnabledProviders(eP => [...eP, Providers.JENKINS])
      }
    } else {
      setRestartDetected(false)
      setExistingTasks([])
    }

    return () => {
      setRestartDetected(false)
      setExistingTasks([])
      // setProjectId([])
    }
  }, [location])

  return (
    <>
      <div className='container container-create-pipeline'>
        <Nav />
        <Sidebar />
        <Content>
          <main className='main'>
            <AppCrumbs
              items={[
                { href: '/', icon: false, text: 'Dashboard' },
                { href: '/pipelines', icon: false, text: 'Pipelines' },
                { href: '/pipelines/create', icon: false, text: 'RUN Pipeline', current: true },
              ]}
            />

            <div className='headlineContainer'>
              <Link style={{ display: 'flex', fontSize: '14px', float: 'right', marginLeft: '10px', color: '#777777' }} to='/pipelines'>
                <Icon
                  icon={
                    <BackArrowIcon
                      width={16} height={16}
                      fill='rgba(0,0,0,0.25)'
                      style={{
                        marginRight: '6px'
                      }}
                    />
                  } size={16}
                /> Go Back
              </Link>
              <div style={{ display: 'flex' }}>
                <div>
                  <span style={{ marginRight: '10px' }}>
                    <Icon icon={<img src={RunPipelineIcon} width='38' height='38' />} size={38} />
                  </span>
                </div>
                <div>
                  <h1 style={{ margin: 0 }}>
                    Run New Pipeline
                    <Popover
                      key='popover-help-key-create-pipeline'
                      className='trigger-delete-connection'
                      popoverClassName='popover-help-create-pipeline'
                      position={Position.RIGHT}
                      autoFocus={false}
                      enforceFocus={false}
                      usePortal={false}
                    >
                      <a href='#' rel='noreferrer'><HelpIcon width={19} height={19} style={{ marginLeft: '10px' }} /></a>
                      <>
                        <div style={{ textShadow: 'none', fontSize: '12px', padding: '12px', maxWidth: '300px' }}>
                          <div style={{ marginBottom: '10px', fontWeight: 700, fontSize: '14px', fontFamily: '"Montserrat", sans-serif' }}>
                            <Icon icon='help' size={16} /> Run Pipeline
                          </div>
                          <p>Need Help? &mdash; Configure the <strong>Data Providers</strong> you want and click
                            <Icon icon='play' size={12} /> <strong>RUN</strong> to trigger a new Pipeline run.
                          </p>
                        </div>
                      </>
                    </Popover>
                  </h1>

                  <p className='page-description mb-0'>Trigger data collection for one or more Data Providers.</p>
                  <p style={{ margin: 0, padding: 0 }}>
                    In a future release you’ll be able to define Blueprints, and schedule recurring plans.
                  </p>
                </div>
              </div>
            </div>

            <div className='' style={{ width: '100%', marginTop: '10px', alignSelf: 'flex-start', alignContent: 'flex-start' }}>
              <h2 className='headline'>
                <Icon icon='git-pull' height={16} size={16} color='rgba(0,0,0,0.5)' /> Pipeline Name<span className='requiredStar'>*</span>
              </h2>
              <p className='group-caption'>Create a user-friendly name for this Run, or select and use a default auto-generated one.</p>
              <div className='form-group' style={{ maxWidth: '480px', paddingLeft: '22px' }}>
                {isValidPipeline() && (
                  <Icon
                    icon='tick' color={Colors.GREEN5} size={12}
                    style={{ float: 'right', marginTop: '7px', marginLeft: '5px' }}
                  />)}
                {!isValidPipeline() && (
                  <>
                    <Icon
                      icon='exclude-row' color={Colors.RED5}
                      size={12}
                      style={{ float: 'right', marginTop: '7px', marginLeft: '5px' }}
                    />
                  </>
                )}
                <FormGroup
                  disabled={isRunning}
                  label=''
                  labelFor='pipeline-name'
                  className=''
                  contentClassName=''
                  helperText={`RUN DATE = ${today.toLocaleString()}`}
                  fill
                  required
                >

                  <InputGroup
                    id='pipeline-name'
                    disabled={isRunning}
                    placeholder='eg. COLLECTION YYYYMMDDHHMMSS'
                    value={pipelineName}
                    onChange={(e) => setPipelineName(e.target.value)}
                    className={!isValidPipelineForm ? 'input-pipeline-name is-invalid' : 'input-pipeline-name is-valid'}
                    rightElement={
                      <>
                        <Button
                          icon='reset' text='' small
                          minimal
                          onClick={() => resetPipelineName()}
                        />
                        <Popover
                          className='popover-pipeline-menu-trigger'
                          popoverClassName='popover-pipeline-menu'
                          position={Position.RIGHT_BOTTOM}
                          usePortal={true}
                        >
                          <Button text={`${today.toLocaleTimeString()}`} />
                          <>
                            <PipelinePresetsMenu
                              namePrefix={namePrefix}
                              pipelineSuffixes={pipelineSuffixes}
                              setNamePrefix={setNamePrefix}
                              setNameSuffix={setNameSuffix}
                            />
                          </>
                        </Popover>
                        {validationErrors.length > 0 && (
                          <>
                            <div style={{ display: 'block', float: 'right' }}>
                              <Popover
                                key='popover-help-key-validation-errors'
                                className='trigger-validation-errors'
                                popoverClassName='popover-help-validation-errors'
                                position={Position.RIGHT}
                                autoFocus={false}
                                enforceFocus={false}
                                usePortal={false}
                              >
                                <Button
                                  intent={Intent.PRIMARY}
                                  icon={<Icon icon='warning-sign' size={14} color={Colors.ORANGE5} />}
                                  small
                                  style={{ margin: '3px 4px 0 0' }}
                                />
                                <div style={{ padding: '5px', minWidth: '300px', maxWidth: '300px', justifyContent: 'flex-start' }}>
                                  <FormValidationErrors errors={validationErrors} textAlign='left' styles={{ display: 'flex' }} />
                                </div>
                              </Popover>
                            </div>
                          </>
                        )}
                      </>
                    }
                    required
                    // large
                    fill
                  />
                </FormGroup>
              </div>

              <h2 className='headline'>
                <Icon icon='database' height={16} size={16} color='rgba(0,0,0,0.5)' /> Data Providers<span className='requiredStar'>*</span>
              </h2>
              <p className='group-caption'>
                Configure available plugins to enable for this <strong>Pipeline Run</strong>.<br />
                Turn the switch to the ON position to activate.
              </p>
              <div className='data-providers' style={{ marginTop: '8px', width: '100%' }}>
                {integrations.map((provider) => (
                  <CSSTransition
                    key={`fx-key-provider-${provider.id}`}
                    in={readyProviders.includes(provider.id)}
                    timeout={350}
                    classNames='provider-datarow'
                    unmountOnExit
                  >
                    {/* <div key={`provider-${provider.id}`}> */}
                    <div
                      className={`data-provider-row ${enabledProviders.includes(provider.id) ? 'on' : 'off'}`}
                    >
                      <div className='provider-info'>
                        <div className='provider-icon'>{provider.iconDashboard}</div>
                        <span className='provider-name'>{provider.name}</span>
                        <Tooltip
                          intent={Intent.PRIMARY}
                          content={`Enable ${provider.name}`} position={Position.RIGHT} popoverClassName='pipeline-tooltip'
                        >
                          <Switch
                          // alignIndicator={Alignment.CENTER}
                            disabled={isRunning}
                            className='provider-toggle-switch'
                            innerLabel={!enabledProviders.includes(provider.id) ? 'OFF' : null}
                            innerLabelChecked='ON'
                            checked={enabledProviders.includes(provider.id)}
                            onChange={() => setEnabledProviders(p =>
                              enabledProviders.includes(provider.id) ? p.filter(p => p !== provider.id) : [...p, provider.id]
                            )}
                          />
                        </Tooltip>
                      </div>
                      <div className='provider-settings'>
                        {/* showProviderSettings(provider.id) */}
                        <ProviderSettings
                          providerId={provider.id}
                          projectId={projectId}
                          owner={owner}
                          repositoryName={repositoryName}
                          sourceId={sourceId}
                          sources={sources}
                          selectedSource={selectedSource}
                          setSelectedSource={setSelectedSource}
                          boardId={boardId}
                          setProjectId={setProjectId}
                          setOwner={setOwner}
                          setRepositoryName={setRepositoryName}
                          setSourceId={setSourceId}
                          setBoardId={setBoardId}
                          isEnabled={isProviderEnabled}
                          isRunning={isRunning}
                        />
                      </div>
                      <div className='provider-actions'>
                        <ButtonGroup minimal rounded='true'>
                          <Button className='pipeline-action-btn' minimal onClick={() => history.push(`/integrations/${provider.id}`)}>
                            <Icon icon='cog' color={Colors.GRAY4} size={16} />
                          </Button>
                          <Popover
                            key={`popover-help-key-provider-${provider.id}`}
                            className='trigger-provider-help'
                            popoverClassName='popover-provider-help'
                            position={Position.RIGHT}
                            autoFocus={false}
                            enforceFocus={false}
                            usePortal={false}
                          >
                            <Button className='pipeline-action-btn' minimal><Icon icon='help' color={Colors.GRAY4} size={16} /></Button>
                            <>
                              <div style={{ textShadow: 'none', fontSize: '12px', padding: '12px', maxWidth: '300px' }}>
                                <div style={{
                                  marginBottom: '10px',
                                  fontWeight: 700,
                                  fontSize: '14px',
                                  fontFamily: '"Montserrat", sans-serif'
                                }}
                                >
                                  <Icon icon='help' size={16} /> {provider.name} Settings
                                </div>
                                <p>Need Help? &mdash; Please enter the required <strong>Run Settings</strong> for this data provider.</p>
                                {/* specific provider field help notes */}
                                {(() => {
                                  let helpContext = null
                                  switch (provider.id) {
                                    case Providers.GITLAB:
                                      helpContext = (
                                        <img
                                          src={GitlabHelpNote}
                                          alt={provider.name} style={{ maxHeight: '64px', maxWidth: '100%' }}
                                        />
                                      )
                                      break
                                    case Providers.JENKINS:
                                      helpContext = <strong>(Options not required)</strong>
                                      break
                                    case Providers.JIRA:
                                      helpContext = (
                                        <img
                                          src={JiraHelpNote}
                                          alt={provider.name} style={{ maxHeight: '64px', maxWidth: '100%' }}
                                        />
                                      )
                                      break
                                    case Providers.GITHUB:
                                      helpContext = (
                                        <img
                                          src={GithubHelpNote}
                                          alt={provider.name} style={{ maxHeight: '64px', maxWidth: '100%' }}
                                        />
                                      )
                                      break
                                  }
                                  return helpContext
                                })()}
                              </div>
                            </>
                          </Popover>
                        </ButtonGroup>
                      </div>
                    </div>
                    {/* </div> */}
                  </CSSTransition>
                ))}
              </div>

            </div>

            <div style={{ display: 'flex', marginTop: '50px', width: '100%', justifyContent: 'flex-start' }}>
              {validationErrors.length > 0 && (
                <FormValidationErrors errors={validationErrors} />
              )}
            </div>
            <div style={{ display: 'flex', width: '100%', justifyContent: 'flex-start' }}>
              <Button
                className='btn-pipeline btn-run-pipeline' icon='play' intent='primary'
                disabled={!isValidPipeline()}
                onClick={runPipeline}
                loading={isRunning}
              ><strong>Run</strong> Pipeline
              </Button>
              <Tooltip content='Manage Pipelines' position={Position.TOP}>
                <Button
                  onClick={() => history.push('/pipelines')}
                  className='btn-pipeline btn-view-jobs'
                  icon='pulse' minimal style={{ marginLeft: '5px' }}
                >View All Pipelines
                </Button>
              </Tooltip>
              <Button
                className='btn-pipeline btn-reset-pipeline'
                icon='eraser' minimal style={{ marginLeft: '5px' }}
                onClick={resetConfiguration}
              >Reset
              </Button>
              {/* <div style={{ padding: '7px 5px 0 5px' }}>
                <Tooltip content='Manage API Rate Limits' position={Position.TOP}>
                  <Switch
                    intent={Intent.DANGER}
                    checked={enableThrottling}
                    onChange={() => setEnableThrottling(t => !t)}
                    labelElement={<strong style={{ color: !enableThrottling ? Colors.GRAY3 : '' }}>Enable Throttling</strong>}
                  />
                </Tooltip>
              </div> */}
            </div>
            <p style={{ margin: '5px 3px', alignSelf: 'flex-start', fontSize: '10px' }}>
              Visit the <a href='#'><strong>All Jobs</strong></a> section to monitor complete pipeline activity.<br />
              Once you run this pipeline, you’ll be redirected to collection status.
            </p>
          </main>
        </Content>
      </div>
      <PipelineIndicator
        pipeline={pipelineRun}
        graphsUrl={GRAFANA_URL}
        onFetch={fetchPipeline}
        onCancel={cancelPipeline}
        onView={() => history.push(`/pipelines/activity/${pipelineRun.ID}`)}
      />
    </>
  )
}

export default CreatePipeline
