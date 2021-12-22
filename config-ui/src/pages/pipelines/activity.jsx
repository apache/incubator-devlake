import React, { Fragment, useEffect, useState, useRef } from 'react'
import { CSSTransition } from 'react-transition-group'
import { useHistory, useParams } from 'react-router-dom'
import { GRAFANA_URL } from '@/utils/config'
import dayjs from '@/utils/time'
// import * as dayjs from 'dayjs'
// import * as relativeTime from 'dayjs/plugin/relativeTime'
// import * as updateLocale from 'dayjs/plugin/updateLocale'
import {
  Button, Icon, Intent,
  Card, Elevation,
  Popover,
  Tooltip,
  Position,
  Spinner,
  Colors,
  Drawer,
  DrawerSize,
  Classes,
  // ButtonGroup, InputGroup, Input, Tag,H2, TextArea,Link
} from '@blueprintjs/core'
import { integrationsData } from '@/data/integrations'
import usePipelineManager from '@/hooks/usePipelineManager'
import Nav from '@/components/Nav'
import Sidebar from '@/components/Sidebar'
import AppCrumbs from '@/components/Breadcrumbs'
import Content from '@/components/Content'
import ContentLoader from '@/components/loaders/ContentLoader'
import StagePanel from '@/components/pipelines/StagePanel'
import TaskActivity from '@/components/pipelines/TaskActivity'
import { ReactComponent as GitlabProviderIcon } from '@/images/integrations/gitlab.svg'
import { ReactComponent as JenkinsProviderIcon } from '@/images/integrations/jenkins.svg'
import { ReactComponent as JiraProviderIcon } from '@/images/integrations/jira.svg'
import { ReactComponent as GitHubProviderIcon } from '@/images/integrations/github.svg'

import { ReactComponent as HelpIcon } from '@/images/help.svg'

const PipelineActivity = (props) => {
  const history = useHistory()
  const { pId } = useParams()
  const pollInterval = useRef()

  const [pipelineId, setPipelineId] = useState() // @todo REMOVE TEST RUN ID!
  const [activeProvider, setActiveProvider] = useState(integrationsData[0])
  const [pipelineName, setPipelineName] = useState()
  const [pollTimer, setPollTimer] = useState(5000)
  const [autoRefresh, setAutoRefresh] = useState(false)

  const [showInspector, setShowInspector] = useState(false)
  const [pipelineReady, setPipelineReady] = useState(false)

  const {
    runPipeline,
    cancelPipeline,
    fetchPipeline,
    activePipeline,
    // pipelineRun,
    isRunning,
    isFetching,
    errors: pipelineErrors,
    // setSettings: setPipelineSettings,
    lastRunId
  } = usePipelineManager(pipelineName)

  // useEffect(() => {
  //   setActiveProvider(providerId ? integrationsData.find(p => p.id === providerId) : integrationsData[0])
  // }, [providerId])

  useEffect(() => {
    setPipelineId(pId)
    // dayjs.extend(relativeTime)
    // dayjs.extend(updateLocale)
    // dayjs.updateLocale('en', {
    //   relativeTime: {
    //     future: 'in %s',
    //     past: '%s ago',
    //     s: '< 1min',
    //     m: 'a minute',
    //     mm: '%d minutes',
    //     h: 'an hour',
    //     hh: '%d hours',
    //     d: 'a day',
    //     dd: '%d days',
    //     M: 'a month',
    //     MM: '%d months',
    //     y: 'a year',
    //     yy: '%d years'
    //   }
    // })
    console.log('>>> REQUESTED PIPELINE ID ===', pId)
  }, [pId])

  useEffect(() => {
    if (pipelineId) {
      fetchPipeline(pipelineId)
      setAutoRefresh(activePipeline.status === 'TASK_RUNNING')
    }

    return () => {
      clearInterval(pollInterval.current)
    }
  }, [pipelineId, activePipeline.status, fetchPipeline])

  useEffect(() => {
    console.log('>>> TASKS KEY', activePipeline.tasks)
  }, [])

  useEffect(() => {
    setPipelineReady(activePipeline.ID !== null && !isFetching)
  }, [activePipeline.ID, isFetching])

  useEffect(() => {
    if (autoRefresh) {
      console.log('>> ACTIVITY POLLING ENABLED!')
      pollInterval.current = setInterval(() => {
        fetchPipeline(pipelineId)
      }, pollTimer)
    } else {
      console.log('>> ACTIVITY POLLING DISABLED!')
      clearInterval(pollInterval.current)
    }
  }, [autoRefresh, fetchPipeline, pipelineId, pollTimer])

  return (
    <>
      <div className='container container-pipeline-activity'>
        <Nav />
        <Sidebar />
        <Content>
          <main className='main'>
            <AppCrumbs
              items={[
                { href: '/', icon: false, text: 'Dashboard' },
                { href: '/pipelines/create', icon: false, text: 'Pipelines', disabled: true },
                { href: `/pipelines/activity/${pipelineId}`, icon: false, text: 'Pipeline Activity & Details', current: true },
              ]}
            />
            <div className='headlineContainer'>
              {/* <Link style={{ float: 'right', marginLeft: '10px', color: '#777777' }} to='/integrations'>
                <Icon icon='fast-backward' size={16} /> Go Back
              </Link> */}
              <div style={{ display: 'flex' }}>
                <div>
                  <span style={{ marginRight: '10px' }}>
                    <Icon icon='pulse' size={32} color={Colors.RED5} />
                  </span>
                </div>
                <div>
                  <h1 style={{ margin: 0 }}>
                    Pipeline Activity
                    {activePipeline?.ID && (
                      <>
                        <span style={{ paddingLeft: '10px' }}>&rarr;
                          <strong style={{ paddingLeft: '10px' }}>
                            #{pipelineId}
                          </strong>
                        </span>
                      </>
                    )}
                    <Popover
                      className='trigger-pipeline-activity-help'
                      popoverClassName='popover-help-pipeline-activity'
                      position={Position.RIGHT}
                      autoFocus={false}
                      enforceFocus={false}
                      usePortal={false}
                    >
                      <a href='#' rel='noreferrer'><HelpIcon width={19} height={19} style={{ marginLeft: '10px' }} /></a>
                      <>
                        <div style={{ textShadow: 'none', fontSize: '12px', padding: '12px', maxWidth: '300px' }}>
                          <div style={{ marginBottom: '10px', fontWeight: 700, fontSize: '14px', fontFamily: '"Montserrat", sans-serif' }}>
                            <Icon icon='help' size={16} /> Pipeline RUN Activity
                          </div>
                          <p>Need Help? &mdash; For better accuracy, ensure that all of your Data Integrations
                            successfully pass the <strong>Connection Test</strong>.
                          </p>
                        </div>
                      </>
                    </Popover>
                  </h1>
                  <p className='page-description mb-0'>View the collection stages for a Pipeline  Run.</p>
                  <p style={{ margin: '0 0 36px 0', padding: 0 }}>
                    You may <strong>Cancel</strong> a running pipeline before it completes.
                  </p>
                </div>
              </div>
            </div>
            {!autoRefresh && isFetching && (
              <ContentLoader title='Loading Pipeline Run ...' message='Please wait while pipeline activity is loaded.' />
            )}
            {activePipeline?.ID && (
              <>
                <StagePanel activePipeline={activePipeline} pipelineReady={pipelineReady} />
                <div style={{ marginBottom: '24px', width: '100%' }}>
                  <CSSTransition
                    in={pipelineReady}
                    timeout={300}
                    classNames='activity-panel'
                    // unmountOnExit
                  >
                    <Card
                      className='pipeline-activity-card'
                      elevation={Elevation.TWO}
                      style={{ padding: 0, width: '100%', display: 'flex', flexDirection: 'column' }}
                    >
                      <div
                        className='pipeline-activity' style={{
                          display: 'flex',
                          width: '100%',
                          justifyContent: 'space-between',
                          borderBottom: '1.0px solid #f0f0f0',
                          padding: '20px'
                        }}
                      >

                        <div className='pipeline-info' style={{ paddingRight: '12px' }}>
                          <h2 className='headline' style={{ marginTop: '0' }}>
                            <span
                              className='pipeline-name'
                              style={{
                                textOverflow: 'ellipsis',
                                overflow: 'hidden',
                                whiteSpace: 'nowrap',
                                display: 'block',
                                maxWidth: '430px',
                                color: activePipeline.status === 'TASK_FAILED' ? Colors.RED4 : ''
                              }}
                            >{activePipeline.name || 'Unamed Collection'}
                            </span>
                          </h2>
                          <div className='pipeline-timestamp'>
                            2021-12-08 08:00 AM (UTC)
                          </div>

                        </div>
                        <div className='pipeline-status' style={{ paddingRight: '12px' }}>
                          <label style={{ color: Colors.GRAY3 }}>Status</label>
                          <div style={{ fontSize: '15px', display: 'flex' }}>
                            <span style={{ marginRight: '4px', color: activePipeline.status === 'TASK_RUNNING' ? '#0066FF' : '' }}>
                              {activePipeline.status.replace('TASK_', '')}
                            </span>
                            {activePipeline.status === 'TASK_FAILED' && (
                              <Icon
                                icon='warning-sign' size={16}
                                color={Colors.RED5} style={{ alignSelf: 'flex-start', marginLeft: '5px', marginBottom: '2px' }}
                              />
                            )}
                            {activePipeline.status === 'TASK_COMPLETED' && (
                              <Icon
                                icon='tick' size={16}
                                color={Colors.GREEN5} style={{ alignSelf: 'flex-start', marginLeft: '5px', marginBottom: '2px' }}
                              />
                            )}
                            {activePipeline.status === 'TASK_RUNNING' && (
                              <Spinner
                                className='pipeline-status-spinner'
                                size={14}
                                intent={activePipeline.status === 'TASK_COMPLETED' ? 'success' : 'danger'}
                                value={activePipeline.status === 'TASK_COMPLETED' ? 1 : null}
                              />
                            )}
                          </div>
                        </div>
                        <div className='pipeline-duration' style={{ paddingRight: '12px' }}>
                          <label style={{ color: Colors.GRAY3 }}>Duration</label>
                          <div style={{ fontSize: '15px', whiteSpace: 'nowrap' }}>
                            {/* {activePipeline.spentSeconds >= 60 ? `${Number(activePipeline.spentSeconds / 60).toFixed(2)}mins` : `${activePipeline.spentSeconds}secs`} */}
                            {activePipeline.status === 'TASK_RUNNING' ? dayjs(activePipeline.CreatedAt).toNow(true) : dayjs(activePipeline.UpdatedAt).from(activePipeline.CreatedAt, true)}
                          </div>
                        </div>
                        <div className='pipeline-actions' style={{ display: 'flex', justifyContent: 'center', alignItems: 'center' }}>
                          {activePipeline.status === 'TASK_COMPLETED' && (
                            <a
                              className='bp3-button bp3-intent-primary'
                              href={GRAFANA_URL}
                              target='_blank'
                              rel='noreferrer'
                              style={{ backgroundColor: '#3bd477', color: '#ffffff' }}
                            >
                              <Icon icon='doughnut-chart' size={13} /> <span className='bp3-button-text'>Grafana</span>
                            </a>
                          )}
                          {activePipeline.status === 'TASK_RUNNING' && (
                            <Popover
                              key='popover-help-key-cancel-run'
                              className='trigger-pipeline-cancel'
                              popoverClassName='popover-pipeline-cancel'
                              position={Position.BOTTOM}
                              autoFocus={false}
                              enforceFocus={false}
                              usePortal={true}
                              disabled={activePipeline.status !== 'TASK_RUNNING'}
                            >
                              <Button
                                className={`btn-cancel-pipeline${activePipeline.status !== 'TASK_RUNNING' ? '-disabled' : ''}`}
                                icon='stop' text='CANCEL' intent={activePipeline.status !== 'TASK_RUNNING' ? '' : 'primary'}
                                disabled={activePipeline.status !== 'TASK_RUNNING'}
                              />
                              <>
                                <div style={{ fontSize: '12px', padding: '12px', maxWidth: '200px' }}>
                                  <p>Are you Sure you want to cancel this <strong>Run</strong>?</p>
                                  <div style={{ display: 'flex', width: '100%', justifyContent: 'flex-end' }}>
                                    <Button
                                      text='CANCEL' minimal
                                      small className={Classes.POPOVER_DISMISS}
                                      stlye={{ marginLeft: 'auto', marginRight: '3px' }}
                                    />
                                    <Button
                                      text='YES' icon='small-tick' intent={Intent.DANGER} small
                                      onClick={() => cancelPipeline(activePipeline.ID)}
                                    />
                                  </div>
                                </div>
                              </>
                            </Popover>
                          )}
                          {activePipeline.status === 'TASK_FAILED' && (
                            <Button
                              className='btn-restart-pipeline'
                              icon='reset'
                              text='RESTART'
                              intent='warning'
                              disabled={activePipeline.status !== 'TASK_FAILED'}
                            />
                          )}
                          <Button
                            icon='refresh'
                            style={{ marginLeft: '5px' }}
                            disabled={activePipeline.status === 'TASK_RUNNING' || activePipeline.status === 'TASK_FAILED'}
                            minimal
                          />
                        </div>

                      </div>
                      <TaskActivity activePipeline={activePipeline} />
                    </Card>
                  </CSSTransition>
                  <div style={{ display: 'flex', padding: '5px 3px', fontSize: '10px', color: '#777777', justifyContent: 'space-between' }}>
                    <div>
                      <Popover
                        className='trigger-pipeline-stage-help'
                        popoverClassName='popover-help-stage-activity'
                        position={Position.BOTTOM}
                        autoFocus={false}
                        enforceFocus={false}
                        usePortal={false}
                      >
                        <a href='#' rel='noreferrer'><Icon icon='help' color={Colors.GRAY3} size={12} style={{ marginTop: '0', marginRight: '5px' }} /></a>
                        <>
                          <div style={{ textShadow: 'none', fontSize: '12px', padding: '12px', maxWidth: '315px' }}>
                            <div style={{
                              marginBottom: '10px',
                              fontWeight: 700,
                              fontSize: '14px',
                              fontFamily: '"Montserrat", sans-serif'
                            }}
                            >
                              <Icon icon='help' size={16} /> Stages and Tasks
                            </div>
                            <p>
                              Monitor <strong>Duration</strong> and <strong>Progress</strong> completion for all tasks. <strong>Grafana</strong> access will be  enabled when the pipeline completes.

                            </p>

                          </div>
                        </>
                      </Popover>

                      {activePipeline.finishedTasks}/{activePipeline.totalTasks} Tasks Completed
                      <span style={{ padding: '0 2px' }}>
                        <Icon icon='dot' size={10} color={Colors.GRAY3} style={{ marginBottom: '3px' }} />
                      </span>
                      <strong>Created </strong> {activePipeline.CreatedAt}
                    </div>
                    <div>
                      {isFetching && (
                        <span style={{ color: Colors.GREEN5 }}>
                          <Icon icon='updated' size={11} style={{ marginBottom: '2px' }} /> Refreshing Activity...
                        </span>
                      )}
                    </div>
                  </div>
                </div>

                {/* <div className='stage-activity' style={{ alignSelf: 'flex-start', width: '100%' }}>
                </div> */}

                <div className='run-settings' style={{ alignSelf: 'flex-start', width: '100%' }}>
                  <div style={{ display: 'flex' }}>
                    <div style={{ display: 'flex', alignItems: 'flex-start', padding: '2px 8px 0 0' }}>
                      <Icon icon='cog' height={16} size={16} color='rgba(0,0,0,0.5)' />
                    </div>
                    <div>
                      <h2 className='headline' style={{ marginTop: 0 }}>
                        Run Settings
                      </h2>
                      <p>Data Provider settings configured for this pipeline execution.</p>
                    </div>
                  </div>

                  <div style={{ padding: '0 10px', display: 'flex', marginTop: '24px', justifyContent: 'space-between', width: '100%' }}>
                    <div className='jenkins-settings' style={{ display: 'flex' }}>
                      <div style={{ display: 'flex', padding: '2px 6px' }}>
                        <JenkinsProviderIcon width={24} height={24} />
                      </div>
                      <div>
                        <label style={{
                          lineHeight: '100%',
                          display: 'block',
                          fontSize: '14px',
                          marginTop: '0',
                          marginBottom: '10px'
                        }}
                        >
                          <strong style={{ fontSize: '11px', fontFamily: 'Montserrat', fontWeight: 800 }}>Jenkins</strong>
                          <br />Auto-configured
                        </label>
                        <span style={{ color: Colors.GRAY3 }}>(No Settings)</span>
                      </div>
                    </div>
                    <div className='jira-settings' style={{ display: 'flex' }}>
                      <div style={{ display: 'flex', padding: '2px 6px' }}>
                        <JiraProviderIcon width={24} height={24} />
                      </div>
                      <div>
                        <label style={{
                          lineHeight: '100%',
                          display: 'block',
                          fontSize: '14px',
                          marginTop: '0',
                          marginBottom: '10px'
                        }}
                        ><strong style={{ fontSize: '11px', fontFamily: 'Montserrat', fontWeight: 800 }}>JIRA</strong><br />Board IDs
                        </label>
                        {activePipeline.tasks.filter(t => t.plugin === 'jira').map((t, tIdx) => (
                          <div key={`board-id-key-${tIdx}`}>
                            <Icon icon='nest' size={12} color={Colors.GRAY4} style={{ marginRight: '6px' }} />
                            <span>
                              {t.options[Object.keys(t.options)[0]]} on Server #{t.options[Object.keys(t.options)[1]]}<br />
                            </span>
                          </div>
                        ))}
                      </div>
                    </div>
                    <div className='gitlab-settings' style={{ display: 'flex' }}>
                      <div style={{ display: 'flex', padding: '2px 6px' }}>
                        <GitlabProviderIcon width={24} height={24} />
                      </div>
                      <div>
                        <label style={{
                          lineHeight: '100%',
                          display: 'block',
                          fontSize: '14px',
                          marginTop: '0',
                          marginBottom: '10px'
                        }}
                        ><strong style={{ fontSize: '11px', fontFamily: 'Montserrat', fontWeight: 800 }}>GitLab</strong><br />Project IDs
                        </label>
                        {activePipeline.tasks.filter(t => t.plugin === 'gitlab').map((t, tIdx) => (
                          <div key={`project-id-key-${tIdx}`}>
                            <Icon icon='nest' size={12} color={Colors.GRAY4} style={{ marginRight: '6px' }} />
                            <span>
                              {t.options[Object.keys(t.options)[0]]}<br />
                            </span>
                          </div>
                        ))}
                      </div>
                    </div>
                    <div className='github-settings' style={{ display: 'flex' }}>
                      <div style={{ display: 'flex', padding: '2px 6px' }}>
                        <GitHubProviderIcon width={24} height={24} />
                      </div>
                      <div>
                        <label style={{ lineHeight: '100%', display: 'block', fontSize: '14px', marginTop: '0', marginBottom: '10px' }}>
                          <strong style={{ fontSize: '11px', fontFamily: 'Montserrat', fontWeight: 800 }}>GitHub</strong><br />Repositories
                        </label>
                        {activePipeline.tasks.filter(t => t.plugin === 'github').map((t, tIdx) => (
                          <div key={`repostitory-id-key-${tIdx}`}>
                            <Icon icon='nest' size={12} color={Colors.GRAY4} style={{ marginRight: '6px' }} />
                            <span>
                              <strong>{t.options[Object.keys(t.options)[0]]}</strong>
                              <span style={{ color: Colors.GRAY5, padding: '0 1px' }}>/</span>
                              <strong>{t.options[Object.keys(t.options)[1]]}</strong>
                            </span>
                          </div>
                        ))}
                      </div>
                    </div>
                  </div>
                </div>
              </>
            )}

            {!pipelineId && (
              <Card elevation={Elevation.TWO} style={{ display: 'flex', alignSelf: 'flex-start' }}>
                <div style={{ display: 'flex', alignSelf: 'flex-start', flexDirection: 'column' }}>
                  <h2 style={{ margin: '0 0 12px 0' }}>
                    <Icon
                      icon='warning-sign'
                      color={Colors.RED4} size={16} style={{ marginBottom: '4px' }}
                    /> Pipeline Run ID <strong>Missing</strong>...
                  </h2>
                  <p>Please provide a Pipeline ID to load Run activity and details.
                    <br /> Check the Address URL in your Browser and try again.
                  </p>
                </div>
              </Card>
            )}

            <div style={{
              marginTop: '40px',
              borderTop: '1px solid #f0f0f0',
              display: 'flex',
              width: '100%',
              justifyContent: 'flex-start'
            }}
            >

              <div style={{ padding: '8px', display: 'flex', width: '100%', justifyContent: 'space-between' }}>
                <span>See <strong style={{ textDecoration: 'underline' }}>All Jobs</strong> to monitor all pipeline activity.</span>
                <div>
                  <Button
                    onClick={() => setShowInspector((iS) => !iS)}
                    icon='code' text='Inspect JSON' small minimal
                    style={{ marginRight: '3px', color: Colors.GRAY3 }}
                  />
                  <Button
                    onClick={() => history.push('/pipelines/create')}
                    icon='add' text='Create New Run' small minimal
                  />
                </div>
              </div>
            </div>

          </main>

        </Content>

      </div>
      <Drawer
        className='drawer-json-inspector'
        icon='code'
        onClose={() => setShowInspector(false)}
        title={`RUN No. ${activePipeline.ID} JSON Payload`}
        position={Position.RIGHT}
        size={DrawerSize.SMALL}
        autoFocus
        canEscapeKeyClose
        canOutsideClickClose
        enforceFocus
        hasBackdrop
        isOpen={showInspector}
        usePortal
      >
        <div className={Classes.DRAWER_BODY}>
          <div className={Classes.DIALOG_BODY}>
            <h3 style={{ margin: 0, padding: '8px 0' }}>
              <span style={{ float: 'right', fontSize: '9px', color: '#aaaaaa' }}>application/json</span> JSON RESPONSE
            </h3>
            <p>
              If you are submitting a
              <strong>Bug-Report</strong> regarding a Pipeline Run, include the output below for better debugging.
            </p>
            <div className='formContainer'>
              <Card
                interactive={false}
                elevation={Elevation.ZERO}
                style={{ padding: '6px 12px', minWidth: '320px', width: '100%', maxWidth: '601px', marginBottom: '20px', overflow: 'auto' }}
              >

                <code>
                  <pre style={{ fontSize: '10px' }}>
                    {JSON.stringify(activePipeline, null, '  ')}
                  </pre>
                </code>
              </Card>
            </div>
          </div>
        </div>
      </Drawer>

    </>
  )
}

export default PipelineActivity
