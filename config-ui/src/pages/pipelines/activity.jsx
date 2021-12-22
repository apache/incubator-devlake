import React, { Fragment, useEffect, useCallback, useState, useRef } from 'react'
import { CSSTransition } from 'react-transition-group'
import { useHistory, useParams } from 'react-router-dom'
import { ToastNotification } from '@/components/Toast'
import { GRAFANA_URL } from '@/utils/config'
import request from '@/utils/request'
import * as dayjs from 'dayjs'
import * as relativeTime from 'dayjs/plugin/relativeTime'
import * as updateLocale from 'dayjs/plugin/updateLocale'
import {
  H2, Button, Icon, Intent,
  ButtonGroup, InputGroup, Input,
  Card, Elevation, Tag,
  Popover,
  Tooltip,
  Position,
  Spinner,
  Colors,
  Drawer,
  DrawerSize,
  TextArea,
  Link,
  Classes
} from '@blueprintjs/core'
import { integrationsData } from '@/data/integrations'
import usePipelineManager from '@/hooks/usePipelineManager'
import Nav from '@/components/Nav'
import Sidebar from '@/components/Sidebar'
import AppCrumbs from '@/components/Breadcrumbs'
import Content from '@/components/Content'
import ContentLoader from '@/components/loaders/ContentLoader'
import { ReactComponent as GitlabProviderIcon } from '@/images/integrations/gitlab.svg'
import { ReactComponent as JenkinsProviderIcon } from '@/images/integrations/jenkins.svg'
import { ReactComponent as JiraProviderIcon } from '@/images/integrations/jira.svg'
import { ReactComponent as GitHubProviderIcon } from '@/images/integrations/github.svg'
import { ReactComponent as PipelineRunningIcon } from '@/images/synchronize.svg'
import { ReactComponent as PipelineFailedIcon } from '@/images/no-synchronize.svg'
import { ReactComponent as PipelineCompleteIcon } from '@/images/check-circle.svg'

const PipelineActivity = (props) => {
  const history = useHistory()
  const { pId } = useParams()
  const pollInterval = useRef()

  const [pipelineId, setPipelineId] = useState() // @todo REMOVE TEST RUN ID!
  const [activeProvider, setActiveProvider] = useState(integrationsData[0])
  const [pipelineName, setPipelineName] = useState()
  const [pollTimer, setPollTimer] = useState(5000)
  // const [autoRefresh, setAutoRefresh] = useState(true)

  const [showInspector, setShowInspector] = useState(false)

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
    dayjs.extend(relativeTime)
    dayjs.extend(updateLocale)
    dayjs.updateLocale('en', {
      relativeTime: {
        future: 'in %s',
        past: '%s ago',
        s: '< 1min',
        m: 'a minute',
        mm: '%d minutes',
        h: 'an hour',
        hh: '%d hours',
        d: 'a day',
        dd: '%d days',
        M: 'a month',
        MM: '%d months',
        y: 'a year',
        yy: '%d years'
      }
    })
    console.log('>>> REQUESTED PIPELINE ID ===', pId)
  }, [pId])

  useEffect(() => {
    if (pipelineId) {
      fetchPipeline(pipelineId)
      // @todo: ENABLE ACTIVITY POLLING
      // pollInterval.current = setInterval(() => {
      //   fetchPipeline(pipelineId)
      // }, pollTimer)
    }

    return () => {
      clearInterval(pollInterval.current)
    }
  }, [pipelineId, fetchPipeline])

  useEffect(() => {
    console.log('>>> TASKS KEY', activePipeline.tasks)
  }, [])

  // useEffect(() => {

  // }, [pipelineId, fetchPipeline])

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
                { href: `/pipelines/activity/${pipelineId}`, icon: false, text: 'Pipeline Activity', current: true },
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
                  </h1>
                  <p className='page-description mb-0'>View the collection stages for a Pipeline  Run.</p>
                  <p style={{ margin: '0 0 36px 0', padding: 0 }}>
                    You may <strong>Cancel</strong> a running pipeline before it completes.
                  </p>
                </div>
              </div>
            </div>
            {isFetching && (
              <ContentLoader title='Loading Pipeline Run ...' message='Please wait while pipeline activity is loaded.' />
            )}
            {!isFetching && activePipeline?.ID && (
              <>
                <div style={{ marginBottom: '24px', width: '100%' }}>
                  <Card
                    className='pipeline-activity-card'
                    elevation={Elevation.TWO}
                    style={{ width: '100%', display: 'flex', }}
                  >
                    <div
                      className='pipeline-activity' style={{
                        display: 'flex',
                        width: '100%',
                        justifyContent: 'space-between'
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
                        {/* <div>
                  <label>Run Date (UTC)</label>
                  <div>2021-12-08 08:00 AM</div>
                </div> */}
                      </div>
                      <div className='pipeline-status' style={{ paddingRight: '12px' }}>
                        <label style={{ color: Colors.GRAY3 }}>Status</label>
                        <div style={{ fontSize: '14px', display: 'flex' }}>
                          <span style={{ marginRight: '4px', color: activePipeline.status === 'TASK_RUNNING' ? '#0066FF' : '' }}>
                            {activePipeline.status}
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
                        <div style={{ fontSize: '14px', whiteSpace: 'nowrap' }}>
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
                            usePortal={false}
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
                  </Card>
                  <p style={{ padding: '5px 3px', fontSize: '10px' }}>
                    {activePipeline.finishedTasks}/{activePipeline.totalTasks} Tasks Completed
                    <span style={{ padding: '0 2px' }}><Icon icon='dot' size={10} color={Colors.GRAY4} /></span>
                    <strong>Created </strong> {activePipeline.CreatedAt}
                  </p>
                </div>

                <div className='stage-activity' style={{ alignSelf: 'flex-start', width: '100%' }}>
                  <h2 className='headline'>
                    <Icon icon='layers' height={16} size={16} color='rgba(0,0,0,0.5)' /> Stages and Tasks
                  </h2>
                  <p>Monitor <strong>Duration</strong> and <strong>Progress</strong> completion for all tasks. <strong>Grafana</strong> access will be  enabled when the pipeline completes.
                    Warning messages will also be displayed here when present.
                  </p>
                  <Card elevation={Elevation.TWO} style={{ display: 'flex', padding: 0 }}>

                    <ButtonGroup style={{ backgroundColor: 'transparent' }}>
                      <Button active>
                        <h3 style={{ margin: 0, fontSize: '20px' }}>
                          Stage 1
                        </h3>
                      </Button>
                      <Button disabled><h3 style={{ margin: 0, fontSize: '20px' }}>Stage 2</h3></Button>
                      <Button
                        minimal
                        style={{
                          marginLeft: '1px',
                          background: '#ffffff!!important',
                          width: 0,
                          height: 0,
                          borderTop: '16px solid transparent',
                          borderBottom: '16px solid transparent',
                          borderLeft: '16px solid rgba(206, 217, 224, 0.5)'
                        }}
                      />
                    </ButtonGroup>
                    <h3 style={{
                      textTransform: 'uppercase',
                      lineHeight: '33px',
                      margin: 0,
                      fontFamily: 'Montserrat',
                      fontWeight: 800,
                      letterSpacing: '2px'
                    }}
                    >Finished Tasks &middot; <span style={{ color: Colors.GREEN5 }}>{activePipeline.finishedTasks}</span>
                      <em style={{ color: '#dddddd', padding: '0 4px', textTransform: 'lowercase' }}>/</em>{activePipeline.totalTasks}
                    </h3>
                    <div style={{ display: 'flex', fontSize: '16px', fontWeight: 700, marginLeft: 'auto', lineHeight: '33px', padding: '0 10px' }}>
                      {(() => {
                        let statusIcon = null
                        switch (activePipeline.status) {
                          case 'TASK_COMPLETED':
                            statusIcon = (
                              <Icon
                                icon={<PipelineCompleteIcon width={24} height={24} style={{ margin: '0 6px 0 10px', display: 'flex', alignSelf: 'center' }} />}
                                size={24}
                              />
                            )
                            break
                          case 'TASK_FAILED':
                            statusIcon = (
                              <Icon
                                icon={<PipelineFailedIcon width={24} height={24} style={{ margin: '0 6px 0 10px', display: 'flex', alignSelf: 'center' }} />}
                                size={24}
                              />
                            )
                            break
                          case 'TASK_RUNNING':
                          default:
                            statusIcon = (
                              <Icon
                                icon={<PipelineRunningIcon width={24} height={24} style={{ margin: '0 6px 0 10px', display: 'flex', alignSelf: 'center' }} />}
                                size={24}
                              />
                            )
                            break
                        }
                        return statusIcon
                      })()}
                      {Number((activePipeline.finishedTasks / activePipeline.totalTasks) * 100).toFixed(1)}%
                    </div>
                  </Card>
                  <div style={{
                    paddingTop: '7px',
                    // borderTop: '1px solid #f5f5f5',
                    marginTop: '14px',
                    marginBottom: '36px'
                  }}
                  >
                    {activePipeline?.ID && activePipeline.tasks && activePipeline.tasks.map((t, tIdx) => (
                      <div
                        className='pipeline-task-row'
                        key={`pipeline-task-key-${tIdx}`}
                        style={{ display: 'flex', padding: '4px 6px', justifyContent: 'space-between', fontSize: '14px' }}
                      >
                        <div style={{ display: 'flex', justifyContent: 'center', paddingRight: '8px', width: '32px', minWidth: '32px' }}>
                          {t.status === 'TASK_COMPLETED' && (
                            <Icon icon='small-tick' size={18} color={Colors.GREEN5} style={{ marginLeft: '0' }} />
                          )}
                          {t.status === 'TASK_FAILED' && (
                            <Icon icon='warning-sign' size={14} color={Colors.RED5} style={{ marginLeft: '0', marginBottom: '3px' }} />
                          )}
                          {t.status === 'TASK_RUNNING' && (
                            <Spinner
                              className='task-spinner'
                              size={14}
                              intent={t.status === 'TASK_COMPLETED' ? 'success' : 'warning'}
                              value={t.status === 'TASK_COMPLETED' ? 1 : t.progress}
                            />
                          )}
                        </div>
                        <div
                          className='pipeline-task-cell-name'
                          style={{ padding: '0 8px', width: '20%', display: 'flex', justifyContent: 'space-between' }}
                        >
                          <strong
                            className='task-plugin-name'
                            style={{ overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}
                          >
                            {t.plugin}
                          </strong>
                          {/* <div style={{
                            width: '120px',
                            // backgroundColor: '#f0f0f0',
                            justifyContent: 'flex-end',
                            display: 'flex',
                            textAlign: 'right'
                          }}
                          >
                            <span style={{ fontSize: '10px', color: Colors.GRAY3, alignSelf: 'center' }}>
                              {t.plugin !== 'jenkins' && (
                                <span>
                                  {Object.keys(t.options)[0]}
                                </span>
                              )}
                              {t.plugin === 'github' && (
                                <span>
                                  <br />{Object.keys(t.options)[1]}
                                </span>
                              )}
                            </span>
                          </div>
                          <div style={{
                            width: '120px',
                            display: 'flex',
                            alignContent: 'flex-start',
                            // backgroundColor: '#f7f7f7',
                            marginRight: 'auto',
                            paddingLeft: '10px'
                          }}
                          >
                            <span style={{ fontSize: '10px', color: Colors.GRAY3, alignSelf: 'center', fontWeight: 600 }}>
                              {t.plugin !== 'jenkins' && (
                                <span>
                                  {t.options[Object.keys(t.options)[0]]}
                                </span>
                              )}
                              {t.plugin === 'github' && (
                                <span>
                                  <br />{t.options[Object.keys(t.options)[1]]}
                                </span>
                              )}
                            </span>
                          </div> */}
                          {/* {t.status === 'TASK_COMPLETED' && (
                            <Icon icon='small-tick' size={14} color={Colors.GREEN5} style={{ marginLeft: '5px' }} />
                          )}
                          {t.status === 'TASK_FAILED' && (
                            <Icon icon='warning-sign' size={11} color={Colors.RED5} style={{ marginLeft: '5px', marginBottom: '3px' }} />
                          )} */}
                        </div>
                        <div
                          className='pipeline-task-cell-duration'
                          style={{
                            padding: '0',
                            minWidth: '80px',
                            // whiteSpace: 'nowrap',
                            textAlign: 'right'
                          }}
                        >
                          {/* <span>{t.spentSeconds >= 60 ? `${Number(t.spentSeconds / 60).toFixed(2)}mins` : `${t.spentSeconds}secs`}</span> */}
                          <span style={{ whiteSpace: 'nowrap' }}>
                            {/* {t.status === 'TASK_COMPLETED' ? dayjs(t.UpdatedAt).from(t.CreatedAt, true) : dayjs(t.CreatedAt).toNow(true)} */}
                            {(() => {
                              let statusRelativeTime = dayjs(t.CreatedAt).toNow(true)
                              switch (t.status) {
                                case 'TASK_COMPLETED':
                                case 'TASK_FAILED':
                                  statusRelativeTime = dayjs(t.UpdatedAt).from(t.CreatedAt, true)
                                  break
                                case 'TASK_RUNNING':
                                default:
                                  statusRelativeTime = dayjs(t.CreatedAt).toNow(true)
                                  break
                              }
                              return statusRelativeTime
                            })()}
                          </span>
                        </div>
                        <div
                          className='pipeline-task-cell-progress'
                          style={{
                            padding: '0 8px',
                            minWidth: '100px',
                            textAlign: 'right'
                          }}
                        >
                          <span style={{ fontWeight: t.status === 'TASK_COMPLETED' ? 800 : 600 }}>
                            {Number(t.status === 'TASK_COMPLETED' ? 100 : (t.progress / 1) * 100).toFixed(2)}%
                          </span>
                        </div>
                        <div
                          className='pipeline-task-cell-message'
                          style={{ width: '70%', paddingLeft: '10px', fontSize: '12px' }}
                        >
                          {t.plugin !== 'jenkins' && (
                            <>
                              <span style={{ color: Colors.GRAY2 }}>
                                <Icon icon='link' size={8} style={{ marginBottom: '3px' }} /> {t.options[Object.keys(t.options)[0]]}
                              </span>
                              {t.plugin === 'github' && (
                                <span style={{ fontWeight: 60 }}>/{t.options[Object.keys(t.options)[1]]}</span>
                              )}
                            </>
                          )}
                          {t.message && (
                            <>
                              <span style={{ color: t.status === 'TASK_FAILED' ? Colors.RED4 : Colors.GRAY3, paddingLeft: '10px' }}>
                                {t.message}
                              </span>
                            </>
                          )}
                        </div>
                      </div>
                    ))}
                  </div>
                </div>

                <div className='run-settings' style={{ alignSelf: 'flex-start', width: '100%' }}>
                  <h2 className='headline'>
                    <Icon icon='cloud-upload' height={16} size={16} color='rgba(0,0,0,0.5)' /> Run Settings
                  </h2>
                  <p>Data Provider settings configured for this pipeline execution.</p>

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
                              <strong>{t.options[Object.keys(t.options)[0]]}</strong><span style={{ color: Colors.GRAY5, padding: '0 1px' }}>/</span>
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
                  <h2 style={{ margin: '0 0 12px 0' }}><Icon icon='warning-sign' color={Colors.RED4} size={16} style={{ marginBottom: '4px' }} /> Pipeline Run ID <strong>Missing</strong>...</h2>
                  <p>Please provide a Pipeline ID to load Run activity and details.<br /> Check the Address URL in your Browser and try again.</p>
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
            <p>If you are submitting a <strong>Bug-Report</strong> regarding a Pipeline Run, include the output below for better debugging.</p>
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
