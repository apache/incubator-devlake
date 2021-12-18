import React, { Fragment, useEffect, useCallback, useState, useRef } from 'react'
import { CSSTransition } from 'react-transition-group'
import {
  useHistory,
  Link,
  // useParams,
} from 'react-router-dom'
import { GRAFANA_URL } from '@/utils/config'
// import { ToastNotification } from '@/components/Toast'
// import { DEVLAKE_ENDPOINT } from '@/utils/config'
// import request from '@/utils/request'
import {
  // Classes,
  Button, Icon, Intent, Switch,
  // H2, Card, Elevation, Tag,
  Menu,
  FormGroup,
  ButtonGroup,
  InputGroup,
  Popover,
  Tooltip,
  Position,
  // Spinner,
  Colors,
  // Alignment
} from '@blueprintjs/core'
import {
  Providers,
} from '@/data/Providers'
import { integrationsData } from '@/data/integrations'
import usePipelineManager from '@/hooks/usePipelineManager'
import usePipelineValidation from '@/hooks/usePipelineValidation'
import FormValidationErrors from '@/components/messages/FormValidationErrors'
import PipelineIndicator from '@/components/widgets/PipelineIndicator'
import Nav from '@/components/Nav'
import Sidebar from '@/components/Sidebar'
import AppCrumbs from '@/components/Breadcrumbs'
import Content from '@/components/Content'
import { ReactComponent as LayersIcon } from '@/images/layers.svg'
import { ReactComponent as HelpIcon } from '@/images/help.svg'
// import { ReactComponent as PipelineRunningIcon } from '@/images/synchronize.svg'
// import { ReactComponent as PipelineFailedIcon } from '@/images/no-synchronize.svg'
// import { ReactComponent as PipelineCompleteIcon } from '@/images/check-circle.svg'
import { ReactComponent as BackArrowIcon } from '@/images/undo.svg'

import GitlabHelpNote from '@/images/help/gitlab-help.png'
import JiraHelpNote from '@/images/help/jira-help.png'
import GithubHelpNote from '@/images/help/github-help.png'

import '@/styles/pipelines.scss'

const CreatePipeline = (props) => {
  const history = useHistory()
  // const { providerId } = useParams()
  const [activeProvider, setActiveProvider] = useState(integrationsData[0])
  const [integrations, setIntegrations] = useState(integrationsData)

  const [today, setToday] = useState(new Date())
  const pipelinePrefixes = ['COLLECT', 'SYNC']
  const pipelineSuffixes = [
    today.getTime(), // 1639630123107
    today.toString(), // Wed Dec 15 2021 23:48:43 GMT-0500 (EST)
    today.toISOString(), // 2021-12-16T04:48:43.107Z
    `${today.getFullYear()}${today.getMonth() + 1}${today.getDate()}${today.getMinutes()}${today.getSeconds()}`, // 202112154936
    today.toUTCString(), // Thu, 16 Dec 2021 04:49:52 GMT
  ]
  // const [autoRun, setAutoRun] = useState(false)
  // const [enableThrottling, setEnableThrottling] = useState(false)
  const [readyProviders, setReadyProviders] = useState([])
  // const [isRunning, setIsRunning] = useState(false)

  const [enabledProviders, setEnabledProviders] = useState([])
  const [runTasks, setRunTasks] = useState([])

  const [namePrefix, setNamePrefix] = useState(pipelinePrefixes[0])
  const [nameSuffix, setNameSuffix] = useState(pipelineSuffixes[0])
  const [pipelineName, setPipelineName] = useState(`${namePrefix} ${nameSuffix}`)
  const [projectId, setProjectId] = useState('')
  const [boardId, setBoardId] = useState('')
  const [sourceId, setSourceId] = useState('')
  const [repositoryName, setRepositoryName] = useState('')
  const [owner, setOwner] = useState('')

  // const [validationErrors, setValidationErrors] = useState([])

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

  // useEffect(() => {
  //   setActiveProvider(providerId ? integrationsData.find(p => p.id === providerId) : integrationsData[0])
  // }, [providerId])

  useEffect(() => {
    integrationsData.forEach((i, idx) => {
      setTimeout(() => {
        setReadyProviders(r => [...r, i.id])
      }, idx * 50)
    })
  }, [])

  // useEffect(() => {
  //   console.log('>> READY LIST = ', readyProviders)
  // }, [readyProviders])

  const isProviderEnabled = (providerId) => {
    return enabledProviders.includes(providerId)
  }

  const isValidPipeline = () => {
    return enabledProviders.length >= 1 &&
      pipelineName !== '' &&
      pipelineName.length > 2 &&
      validationErrors.length === 0
  }

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
          repositoryName,
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
    return {
      Plugin: providerId,
      Options: {
        ...getProviderOptions(providerId)
      }
    }
  }, [getProviderOptions])

  const showProviderSettings = (providerId) => {
    let providerSettings = null
    switch (providerId) {
      case Providers.JENKINS:
        providerSettings = <p><strong style={{ fontWeight: 900 }}>AUTO-CONFIGURED</strong><br />No Additional Settings</p>
        break
      case Providers.JIRA:
        providerSettings = (
          <>
            <FormGroup
              disabled={isRunning || !isProviderEnabled(providerId)}
              label={<strong>Source ID<span className='requiredStar'>*</span></strong>}
              labelInfo={<span style={{ display: 'block' }}>Enter Connection Instance ID</span>}
              inline={false}
              labelFor='source-id'
              className=''
              contentClassName=''
              fill
            >
              <InputGroup
                id='source-id'
                disabled={isRunning || !isProviderEnabled(providerId)}
                placeholder='eg. 54'
                value={sourceId}
                onChange={(e) => setSourceId(e.target.value)}
                className='input-source-id'
                autoComplete='off'
                fill={false}
              />
            </FormGroup>
            <FormGroup
              disabled={isRunning || !isProviderEnabled(providerId)}
              label={<strong>Board ID<span className='requiredStar'>*</span></strong>}
              labelInfo={<span style={{ display: 'block' }}>Enter JIRA Board No.</span>}
              inline={false}
              labelFor='board-id'
              className=''
              contentClassName=''
              style={{ marginLeft: '12px' }}
              fill
            >
              <InputGroup
                id='board-id'
                disabled={isRunning || !isProviderEnabled(providerId)}
                placeholder='eg. 8'
                value={boardId}
                onChange={(e) => setBoardId(e.target.value)}
                className='input-board-id'
                autoComplete='off'
                fill={false}
              />
            </FormGroup>
          </>
        )
        break
      case Providers.GITHUB:
        providerSettings = (
          <>
            <FormGroup
              disabled={isRunning || !isProviderEnabled(providerId)}
              label={<strong>Owner<span className='requiredStar'>*</span></strong>}
              labelInfo={<span style={{ display: 'block' }}>Enter Project Owner</span>}
              inline={false}
              labelFor='owner'
              className=''
              contentClassName=''
              fill
            >
              <InputGroup
                id='owner'
                disabled={isRunning || !isProviderEnabled(providerId)}
                placeholder='eg. merio-dev'
                value={owner}
                onChange={(e) => setOwner(e.target.value)}
                className='input-owner'
                autoComplete='off'
                // fill={false}
              />
            </FormGroup>
            <FormGroup
              disabled={isRunning || !isProviderEnabled(providerId)}
              label={<strong>Repository Name<span className='requiredStar'>*</span></strong>}
              labelInfo={<span style={{ display: 'block' }}>Enter Git repository</span>}
              inline={false}
              labelFor='repository-name'
              className=''
              contentClassName=''
              style={{ marginLeft: '12px' }}
              fill
            >
              <InputGroup
                id='repository-name'
                disabled={isRunning || !isProviderEnabled(providerId)}
                placeholder='eg. lake'
                value={repositoryName}
                onChange={(e) => setRepositoryName(e.target.value)}
                className='input-repository-name'
                autoComplete='off'
                fill={false}
              />
            </FormGroup>
          </>
        )
        break
      case Providers.GITLAB:
        providerSettings = (
          <>
            <FormGroup
              disabled={isRunning || !isProviderEnabled(providerId)}
              label={<strong>Project ID<span className='requiredStar'>*</span></strong>}
              labelInfo={<span style={{ display: 'block' }}>Enter the GitLab Project ID No.</span>}
              inline={false}
              labelFor='project-id'
              className=''
              contentClassName=''
              // fill
            >
              <InputGroup
                id='project-id'
                disabled={isRunning || !isProviderEnabled(providerId)}
                placeholder='eg. 937810831'
                value={projectId}
                onChange={(e) => setProjectId(pId => e.target.value)}
                className='input-project-id'
                autoComplete='off'
                // fill={false}
              />
            </FormGroup>
          </>
        )
        break
      default:
        break
    }

    return providerSettings
  }

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
  }, [runTasks, pipelineName, setPipelineSettings])

  useEffect(() => {
    console.log('>> ENBALED PROVIDERS = ', enabledProviders)
    const PipelineTasks = enabledProviders.map(p => configureProvider(p))
    setRunTasks(PipelineTasks)
    console.log('>> CONFIGURED PIPELINE TASKS = ', PipelineTasks)
    validate()
  }, [enabledProviders, projectId, boardId, sourceId, owner, repositoryName, configureProvider])

  useEffect(() => {
    console.log('>> PIPELINE LAST RUN OBJECT CHANGED!!...', pipelineRun)
  }, [pipelineRun])

  useEffect(() => {
    console.log(namePrefix, nameSuffix)
    setPipelineName(`${namePrefix} ${nameSuffix}`)
    setToday(new Date())
  }, [namePrefix, nameSuffix])

  return (
    <>
      <div className='container'>
        <Nav />
        <Sidebar />
        <Content>
          <main className='main'>
            <AppCrumbs
              items={[
                { href: '/', icon: false, text: 'Dashboard' },
                { href: '/pipelines/create', icon: false, text: 'Pipelines', disabled: true },
                { href: '/pipelines/create', icon: false, text: 'RUN Pipeline', current: true },
              ]}
            />

            <div className='headlineContainer'>
              <Link style={{ display: 'flex', fontSize: '14px', float: 'right', marginLeft: '10px', color: '#777777' }} to='/'>
                <Icon icon={<BackArrowIcon width={16} height={16} fill='rgba(0,0,0,0.25)' style={{ marginRight: '6px' }} />} size={16} /> Go Back
              </Link>
              <div style={{ display: 'flex' }}>
                <div>
                  <span style={{ marginRight: '10px' }}>
                    <Icon icon={<LayersIcon width={42} height={42} />} size={42} />
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
                          <p>Need Help? &mdash; Configure the <strong>Data Providers</strong> you want and click <Icon icon='play' size={12} /> <strong>RUN</strong> to trigger a new Pipeline run.</p>
                        </div>
                      </>
                    </Popover>
                  </h1>

                  <p className='page-description mb-0'>Trigger data collection for one or more Data Providers.</p>
                  <p style={{ margin: 0, padding: 0 }}>In a future release you’ll be able to define Blueprints, and schedule recurring plans.</p>
                </div>
              </div>
            </div>

            <div className='' style={{ width: '100%', marginTop: '10px', alignSelf: 'flex-start', alignContent: 'flex-start' }}>
              <h2 className='headline'>
                <Icon icon='git-pull' height={16} size={16} color='rgba(0,0,0,0.5)' /> Pipeline Name<span className='requiredStar'>*</span>
              </h2>
              <p className='group-caption'>Create a user-friendly name for this Run, or select and use a default auto-generated one.</p>

              <div className='form-group' style={{ maxWidth: '480px', paddingLeft: '22px' }}>
                {isValidPipeline() && (<Icon icon='tick' color={Colors.GREEN5} size={12} style={{ float: 'right', marginTop: '7px', marginLeft: '5px' }} />)}
                {!isValidPipeline() && (<Icon icon='exclude-row' color={Colors.RED5} size={12} style={{ float: 'right', marginTop: '7px', marginLeft: '5px' }} />)}
                <FormGroup
                  disabled={isRunning}
                  label=''
                  labelFor='pipeline-name'
                  className=''
                  contentClassName=''
                  // label={<strong>Name</strong>}
                  // labelInfo={<span style={{ display: 'block' }}>{`RUN DATE = ${today.toLocaleString()}`}</span>}
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
                          onClick={() => setPipelineName(`${namePrefix} ${nameSuffix}`)}
                        />
                        <Popover
                          className='popover-pipeline-menu-trigger'
                          popoverClassName='popover-pipeline-menu'
                          position={Position.RIGHT_BOTTOM}
                          usePortal={true}
                        >
                          <Button text={`${today.toLocaleTimeString()}`} />
                          <>
                            <Menu className='pipeline-presets-menu'>
                              <label style={{
                                fontSize: '10px',
                                fontWeight: 800,
                                fontFamily: '"Montserrat", sans-serif',
                                textTransform: 'uppercase',
                                padding: '6px 8px',
                                display: 'block'
                              }}
                              >Preset Naming Options
                              </label>
                              <Menu.Item text='COLLECTION ...' active={namePrefix === 'COLLECT'}>
                                <Menu.Item icon='key-option' text='COLLECT [UNIXTIME]' onClick={() => setNamePrefix('COLLECT') | setNameSuffix(pipelineSuffixes[0])} />
                                <Menu.Item icon='key-option' text='COLLECT [YYYYMMDDHHMMSS]' onClick={() => setNamePrefix('COLLECT') | setNameSuffix(pipelineSuffixes[3])} />
                                <Menu.Item icon='key-option' text='COLLECT [ISO]' onClick={() => setNamePrefix('COLLECT') | setNameSuffix(pipelineSuffixes[2])} />
                                <Menu.Item icon='key-option' text='COLLECT [UTC]' onClick={() => setNamePrefix('COLLECT') | setNameSuffix(pipelineSuffixes[4])} />
                              </Menu.Item>
                              <Menu.Item text='SYNCHRONIZE ...' active={namePrefix === 'SYNC'}>
                                <Menu.Item icon='key-option' text='SYNC [UNIXTIME]' onClick={() => setNamePrefix('SYNC') | setNameSuffix(pipelineSuffixes[0])} />
                                <Menu.Item icon='key-option' text='SYNC [YYYYMMDDHHMMSS]' onClick={() => setNamePrefix('SYNC') | setNameSuffix(pipelineSuffixes[3])} />
                                <Menu.Item icon='key-option' text='SYNC [ISO]' onClick={() => setNamePrefix('SYNC') | setNameSuffix(pipelineSuffixes[2])} />
                                <Menu.Item icon='key-option' text='SYNC [UTC]' onClick={() => setNamePrefix('SYNC') | setNameSuffix(pipelineSuffixes[4])} />
                              </Menu.Item>
                              <Menu.Item text='RUN ...' active={namePrefix === 'RUN'}>
                                <Menu.Item icon='key-option' text='RUN [UNIXTIME]' onClick={() => setNamePrefix('RUN') | setNameSuffix(pipelineSuffixes[0])} />
                                <Menu.Item icon='key-option' text='RUN [YYYYMMDDHHMMSS]' onClick={() => setNamePrefix('RUN') | setNameSuffix(pipelineSuffixes[3])} />
                                <Menu.Item icon='key-option' text='RUN [ISO]' onClick={() => setNamePrefix('RUN') | setNameSuffix(pipelineSuffixes[2])} />
                                <Menu.Item icon='key-option' text='RUN [UTC]' onClick={() => setNamePrefix('RUN') | setNameSuffix(pipelineSuffixes[4])} />
                              </Menu.Item>
                              <Menu.Divider />
                              <Menu.Item text='Advanced Options' icon='cog'>
                                <Menu.Item icon='new-object' text='Save Pipeline Blueprint' disabled />
                              </Menu.Item>
                            </Menu>
                          </>
                        </Popover>
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
                        {showProviderSettings(provider.id)}
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
              <Button className='btn-pipeline btn-view-jobs' icon='eye-open' minimal style={{ marginLeft: '5px' }}>View All Jobs</Button>
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
      <PipelineIndicator pipeline={pipelineRun} graphsUrl={GRAFANA_URL} onFetch={fetchPipeline} />
    </>
  )
}

export default CreatePipeline
