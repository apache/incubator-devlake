import React, { Fragment, useEffect, useCallback, useState, useRef } from 'react'
import { CSSTransition } from 'react-transition-group'
import { useHistory, useParams } from 'react-router-dom'
// import { ToastNotification } from '@/components/Toast'
// import { DEVLAKE_ENDPOINT } from '@/utils/config'
// import request from '@/utils/request'
import {
  Classes,
  Button, Icon, Intent, Switch,
  // H2, Card, Elevation, Tag,
  FormGroup,
  ButtonGroup,
  InputGroup,
  Popover,
  Tooltip,
  Position,
  Spinner,
  Colors,
  // Link,
  // Alignment
} from '@blueprintjs/core'
import {
  Providers,
} from '@/data/Providers'
import { integrationsData } from '@/data/integrations'
import usePipelineManager from '@/hooks/usePipelineManager'
import FormValidationErrors from '@/components/messages/FormValidationErrors'
import Nav from '@/components/Nav'
import Sidebar from '@/components/Sidebar'
import AppCrumbs from '@/components/Breadcrumbs'
import Content from '@/components/Content'
import { ReactComponent as LayersIcon } from '@/images/layers.svg'
import { ReactComponent as HelpIcon } from '@/images/help.svg'
import { ReactComponent as PipelineRunningIcon } from '@/images/synchronize.svg'
import { ReactComponent as PipelineFailedIcon } from '@/images/no-synchronize.svg'
import { ReactComponent as PipelineCompleteIcon } from '@/images/check-circle.svg'

import '@/styles/pipelines.scss'

const CreatePipeline = (props) => {
  const history = useHistory()
  // const { providerId } = useParams()
  const [activeProvider, setActiveProvider] = useState(integrationsData[0])
  const [integrations, setIntegrations] = useState(integrationsData)

  const today = new Date()
  const [autoRun, setAutoRun] = useState(false)
  const [enableThrottling, setEnableThrottling] = useState(false)
  const [ready, setReady] = useState([])
  // const [isRunning, setIsRunning] = useState(false)

  const [enabledProviders, setEnabledProviders] = useState([])
  const [runTasks, setRunTasks] = useState([])

  const [pipelineName, setPipelineName] = useState(`COLLECTION ${Date.now()}`)
  const [projectId, setProjectId] = useState('')
  const [boardId, setBoardId] = useState('')
  const [sourceId, setSourceId] = useState('')
  const [repositoryName, setRepositoryName] = useState('')
  const [owner, setOwner] = useState('')

  const [validationErrors, setValidationErrors] = useState([])

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

  // useEffect(() => {
  //   setActiveProvider(providerId ? integrationsData.find(p => p.id === providerId) : integrationsData[0])
  // }, [providerId])

  useEffect(() => {
    integrationsData.forEach((i, idx) => {
      setTimeout(() => {
        setReady(r => [...r, i.id])
      }, idx * 50)
    })
  }, [])

  useEffect(() => {
    console.log('>> READY LIST = ', ready)
  }, [ready])

  const isProviderEnabled = (providerId) => {
    return enabledProviders.includes(providerId)
  }

  const isValidPipeline = () => {
    return enabledProviders.length >= 1 &&
      pipelineName.length >= 2 &&
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
          boardId,
          sourceId
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
          projectId
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
              label={<strong>Board IDs <span className='requiredStar'>*</span></strong>}
              labelInfo={<span style={{ display: 'block' }}>Enter JIRA Board(s)</span>}
              inline={false}
              labelFor='board-ids'
              className=''
              contentClassName=''
              fill
            >
              <InputGroup
                id='pipeline-name'
                disabled={isRunning || !isProviderEnabled(providerId)}
                placeholder='eg. 8'
                value={boardId}
                onChange={(e) => setBoardId(e.target.value)}
                className='input-board-ids'
                fill={false}
              />
            </FormGroup>
            <FormGroup
              disabled={isRunning || !isProviderEnabled(providerId)}
              label={<strong>Source ID<span className='requiredStar'>*</span></strong>}
              labelInfo={<span style={{ display: 'block' }}>Enter Connection Instance ID</span>}
              inline={false}
              labelFor='source-id'
              className=''
              contentClassName=''
              style={{ marginLeft: '12px' }}
              fill
            >
              <InputGroup
                id='pipeline-name'
                disabled={isRunning || !isProviderEnabled(providerId)}
                placeholder='eg. 54'
                value={sourceId}
                onChange={(e) => setSourceId(e.target.value)}
                className='input-source-id'
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
              label={<strong>Repository Name<span className='requiredStar'>*</span></strong>}
              labelInfo={<span style={{ display: 'block' }}>Enter Git repository</span>}
              inline={false}
              labelFor='repository-name'
              className=''
              contentClassName=''
              fill
            >
              <InputGroup
                id='repository-name'
                disabled={isRunning || !isProviderEnabled(providerId)}
                placeholder='eg. lake'
                value={repositoryName}
                onChange={(e) => setRepositoryName(e.target.value)}
                className='input-board-ids'
                fill={false}
              />
            </FormGroup>
            <FormGroup
              disabled={isRunning || !isProviderEnabled(providerId)}
              label={<strong>Owner<span className='requiredStar'>*</span></strong>}
              labelInfo={<span style={{ display: 'block' }}>Enter Project Owner</span>}
              inline={false}
              labelFor='owner'
              className=''
              contentClassName=''
              style={{ marginLeft: '12px' }}
              fill
            >
              <InputGroup
                id='owner'
                disabled={isRunning || !isProviderEnabled(providerId)}
                placeholder='eg. merio-dev'
                value={owner}
                onChange={(e) => setOwner(e.target.value)}
                className='input-owner'
                // fill={false}
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
              label={<strong>Project IDs <span className='requiredStar'>*</span></strong>}
              labelInfo={<span style={{ display: 'block' }}>Enter the GitLab Projects to map</span>}
              inline={false}
              labelFor='project-ids'
              className=''
              contentClassName=''
              // fill
            >
              <InputGroup
                id='project-ids'
                disabled={isRunning || !isProviderEnabled(providerId)}
                placeholder='eg. 937810831'
                value={projectId}
                onChange={(e) => setProjectId(e.target.value)}
                className='input-project-ids'
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
    setValidationErrors(errors => pipelineName.length <= 2
      ? [...errors, 'Name: Enter a valid Pipeline Name']
      : errors.filter(e => !e.startsWith('Name:')))
  }, [pipelineName])

  useEffect(() => {
    console.log('>> PIPELINE RUN TASK SETTINGS FOR PIPELINE MANAGER ....', runTasks)
    setPipelineSettings({
      name: pipelineName,
      tasks: [
        [...runTasks]
      ]
    })
  }, [runTasks, pipelineName, setPipelineSettings])

  useEffect(() => {
    console.log('>> ENBALED PROVIDERS = ', enabledProviders)
    const PipelineTasks = enabledProviders.map(p => configureProvider(p))
    setRunTasks(PipelineTasks)
    console.log('>> CONFIGURED PIPELINE TASKS = ', PipelineTasks)
    setValidationErrors(errors => enabledProviders.includes(Providers.GITLAB) && !projectId ? ['GitLab: Specify Run Settings'] : errors.filter(e => !e.startsWith('GitLab:')))
    setValidationErrors(errors => enabledProviders.includes(Providers.JIRA) && (!boardId || !sourceId) ? ['JIRA: Specify Run Settings'] : errors.filter(e => !e.startsWith('JIRA:')))
    setValidationErrors(errors => enabledProviders.includes(Providers.GITHUB) && (!repositoryName || !owner) ? ['GitHub: Specify Run Settings'] : errors.filter(e => !e.startsWith('GitHub:')))
    setValidationErrors(errors => enabledProviders.length === 0 ? ['Pipeline: Invalid/Empty Configuration'] : errors.filter(e => !e.startsWith('Pipeline:')))
  }, [enabledProviders, projectId, boardId, sourceId, owner, repositoryName, configureProvider])

  useEffect(() => {
    console.log('>> PIPELINE VALIDATION ERRORS...', validationErrors)
  }, [validationErrors])

  useEffect(() => {
    console.log('>> PIPELINE LAST RUN OBJECT CHANGED!!...', pipelineRun)
  }, [pipelineRun])

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
                { href: '/pipelines', icon: false, text: 'Pipelines' },
                { href: '/pipelines/create', icon: false, text: 'RUN Pipeline', current: true },
              ]}
            />
            <div className='headlineContainer'>
              {/* <Link style={{ float: 'right', marginLeft: '10px', color: '#777777' }} to='/integrations'>
                <Icon icon='fast-backward' size={16} /> Go Back
              </Link> */}
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

                  <p className='page-description'>Trigger data collection for one or more Data Providers.</p>
                  <p style={{ margin: 0, padding: 0 }}>In a future release you’ll be able to define Blueprints, and schedule recurring plans.</p>
                </div>
              </div>
            </div>

            <div className='' style={{ width: '100%', marginTop: '20px', alignSelf: 'flex-start', alignContent: 'flex-start' }}>
              <h3 className='group-header'>
                <Icon icon='git-pull' height={16} size={16} color='rgba(0,0,0,0.5)' /> Pipeline Name <span className='requiredStar'>*</span>
              </h3>
              <p className='group-caption'>Create a user-friendly name for this Run, or use the default auto-generated one.</p>

              <div className='form-group' style={{ maxWidth: '420px', paddingLeft: '22px' }}>
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
                    className='input-pipeline-name'
                    rightElement={
                      <>
                        <Button
                          icon='reset' text='' small
                          minimal
                          onClick={() => setPipelineName(`COLLECTION ${Date.now()}`)}
                        />
                        <Button text={`${today.toLocaleTimeString()}`} />
                      </>
                    }
                    required
                    // large
                    fill
                  />
                </FormGroup>
              </div>

              <h3 className='group-header'>
                <Icon icon='database' height={16} size={16} color='rgba(0,0,0,0.5)' /> Data Providers <span className='requiredStar'>*</span>
              </h3>
              <p className='group-caption'>Configure available plugins to enable for this <strong>Pipeline Run</strong>.<br />Turn the switch to the ON position to activate.</p>

              <div className='data-providers' style={{ marginTop: '24px', width: '100%' }}>
                {integrations.map((provider) => (
                  <CSSTransition
                    key={`fx-key-provider-${provider.id}`}
                    in={ready.includes(provider.id)}
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
                          <Button className='pipeline-action-btn' minimal><Icon icon='help' color={Colors.GRAY4} size={16} /></Button>
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
              <div style={{ padding: '7px 5px 0 5px' }}>
                <Tooltip content='Manage API Rate Limits' position={Position.TOP}>
                  <Switch
                    intent={Intent.DANGER}
                    checked={enableThrottling}
                    onChange={() => setEnableThrottling(t => !t)}
                    labelElement={<strong style={{ color: !enableThrottling ? Colors.GRAY3 : '' }}>Enable Throttling</strong>}
                  />
                </Tooltip>
              </div>
            </div>
            <p style={{ margin: '5px 3px', alignSelf: 'flex-start', fontSize: '10px' }}>
              Visit the <a href='#'><strong>All Jobs</strong></a> section to monitor complete pipeline activity.<br />
              Once you run this pipeline, you’ll be redirected to collection status.
            </p>
          </main>
        </Content>
      </div>
      {/* {pipelineRun && pipelineRun.ID !== null && ( */}
      <CSSTransition
        in={pipelineRun && pipelineRun.ID !== null}
        timeout={300} classNames='lastrun-module'
        unmountOnExit
      >
        <div
          className='trigger-module-lastrun'
          style={{
            position: 'fixed',
            borderRadius: '40px',
            backgroundColor: '#ffffff',
            width: '40px',
            height: '40px',
            right: '30px',
            bottom: '20px',
            zIndex: 500,
            boxShadow: '0px 0px 6px rgba(0, 0, 0, 0.25)',
            display: 'flex',
            alignItems: 'center',
            alignContent: 'center',
            justifyContent: 'center'
          }}
        >
          <div style={{
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            width: '40px',
            height: '40px',
            position: 'relative'
          }}
          >
            <Popover
              key='popover-lastrun-info'
              className='popover-trigger-lastrun'
              popoverClassName='popover-lastrun'
              position={Position.LEFT_BOTTOM}
              autoFocus={false}
              enforceFocus={false}
              usePortal={false}
              onOpening={() => fetchPipeline(pipelineRun.ID, true)}
              // onOpened=
            >
              <>
                <Spinner
                  value={pipelineRun.status === 'TASK_COMPLETED' ? 100 : null}
                  className='lastrun-spinner'
                  intent={pipelineRun.status === 'TASK_COMPLETED' ? Intent.WARNING : Intent.PRIMARY}
                  size={40} style={{

                  }}
                />
                <Icon
                  icon={pipelineRun.status === 'TASK_COMPLETED'
                    ? <PipelineCompleteIcon width={40} height={40} style={{ marginTop: '3px', display: 'flex', alignSelf: 'center' }} />
                    : <PipelineRunningIcon width={40} height={40} style={{ marginTop: '3px', display: 'flex', alignSelf: 'center' }} />}
                  size={40}
                />
              </>
              <>
                <div style={{ fontSize: '12px', padding: '12px', minWidth: '420px', maxWidth: '420px' }}>
                  <h3 className='group-header' style={{ marginTop: '0', marginBottom: '6px' }}>
                    <Icon icon='help' size={16} /> {pipelineRun.name || 'Last Pipeline Run'}
                  </h3>
                  <p style={{ fontSize: '11px' }}>{pipelineRun.message}</p>
                  <div style={{ display: 'flex', width: '100%', justifyContent: 'space-between' }}>
                    <div>
                      <label><strong>Pipeline ID</strong></label>
                      <div style={{ fontSize: '13px', fontWeight: 800 }}>{pipelineRun.ID}</div>
                    </div>
                    <div style={{ padding: '0 12px' }}>
                      <label><strong>Tasks</strong></label>
                      <div style={{ fontSize: '13px' }}>{pipelineRun.finishedTasks}/{pipelineRun.totalTasks}</div>
                    </div>
                    <div>
                      <label><strong>Status</strong></label>
                      <div style={{ fontSize: '13px' }}>{pipelineRun.status}</div>
                    </div>
                    <div style={{ paddingLeft: '10px', justifyContent: 'flex-end', alignSelf: 'flex-end' }}>
                      {pipelineRun.status === 'TASK_COMPLETED' && (
                        <Button
                          intent='primary'
                          icon='doughnut-chart' text='Graphs'
                          style={{ backgroundColor: '#3bd477', color: '#ffffff' }}
                          small
                        />
                      )}
                      {pipelineRun.status === 'TASK_RUNNING' && (
                        <Button
                          className='btn-cancel-pipeline'
                          small icon='stop' text='CANCEL' intent='primary'
                          onClick={() => cancelPipeline(pipelineRun.ID)}
                        />
                      )}
                      <Button
                        minimal
                        className={`btn-ok ${Classes.POPOVER_DISMISS}`}
                        small text='OK'
                        style={{ marginLeft: '3px' }}
                      />
                    </div>
                  </div>
                </div>
              </>
            </Popover>
          </div>
        </div>
      </CSSTransition>
      {/* )} */}
    </>
  )
}

export default CreatePipeline
