import React, { Fragment, useEffect, useState, useRef, useCallback } from 'react'
import { CSSTransition } from 'react-transition-group'
import { useHistory, useParams, Link } from 'react-router-dom'
import { GRAFANA_URL } from '@/utils/config'
import dayjs from '@/utils/time'
import {
  Button, Icon, Intent,
  Card, Elevation,
  Popover,
  Tooltip,
  Position,
  Spinner,
  Colors,
  Classes,
  // Drawer,
  // DrawerSize,
  ButtonGroup, InputGroup, Input, Tag, H2, TextArea
} from '@blueprintjs/core'
import { integrationsData } from '@/data/integrations'
import {
  Providers,
  ProviderLabels
} from '@/data/Providers'
import usePipelineManager from '@/hooks/usePipelineManager'
import Nav from '@/components/Nav'
import Sidebar from '@/components/Sidebar'
import AppCrumbs from '@/components/Breadcrumbs'
import Content from '@/components/Content'
import ContentLoader from '@/components/loaders/ContentLoader'
import PipelineIndicator from '@/components/widgets/PipelineIndicator'
import CodeInspector from '@/components/pipelines/CodeInspector'
import { ReactComponent as GitlabProviderIcon } from '@/images/integrations/gitlab.svg'
import { ReactComponent as JenkinsProviderIcon } from '@/images/integrations/jenkins.svg'
import { ReactComponent as JiraProviderIcon } from '@/images/integrations/jira.svg'
import { ReactComponent as GitHubProviderIcon } from '@/images/integrations/github.svg'
import { ReactComponent as BackArrowIcon } from '@/images/undo.svg'
import { ReactComponent as HelpIcon } from '@/images/help.svg'

const Pipelines = (props) => {
  const history = useHistory()
  const { providerId } = useParams()
  const [activeProvider, setActiveProvider] = useState(integrationsData[0])

  const [isProcessing, setIsProcessing] = useState(false)
  const [refresh, setRefresh] = useState(false)
  const [activeStatus, setActiveStatus] = useState('all')
  const [latestPipeline, setLatestPipeline] = useState()
  const [showInspector, setShowInspector] = useState(false)
  const [inspectPipeline, setInspectPipeline] = useState(null)

  const {
    runPipeline,
    cancelPipeline,
    fetchPipeline,
    pipelines,
    pipelineCount,
    fetchAllPipelines,
    activePipeline,
    pipelineRun,
    isFetching,
    isFetchingAll,
    errors: pipelineErrors,
    setSettings: setPipelineSettings,
    lastRunId,
  } = usePipelineManager()

  // const [pipelines, setPipelines] = useState([
  //   {
  //     ID: 18092,
  //     CreatedAt: '2021-12-23T23:40:24.808Z',
  //     UpdatedAt: '2021-12-24T00:38:35.165Z',
  //     name: '#695 COLLECT 1640302757915',
  //     tasks: [
  //       [
  //         {
  //           plugin: 'github',
  //           options: {
  //             owner: 'e2corporation',
  //             repositoryName: 'getmdl-select'
  //           }
  //         },
  //         {
  //           plugin: 'jenkins',
  //           options: {}
  //         }
  //       ]
  //     ],
  //     totalTasks: 2,
  //     finishedTasks: 2,
  //     beganAt: '2021-12-23T23:40:24.86Z',
  //     finishedAt: '2021-12-24T00:38:35.163Z',
  //     status: 'TASK_FAILED',
  //     message: "Error 1364: Field 'origin_key' doesn't have a default value",
  //     spentSeconds: 3491
  //   },
  //   {
  //     ID: 455,
  //     CreatedAt: '2021-12-23T21:36:29.168Z',
  //     UpdatedAt: '2021-12-23T21:42:52.971Z',
  //     name: 'RETRY | config-ui trigger Thu Dec 23 2021 10:08:06 GMT-0500 (EST)',
  //     tasks: [
  //       [
  //         {
  //           plugin: 'gitlab',
  //           options: {
  //             projectId: 1967944
  //           }
  //         },
  //         {
  //           plugin: 'gitlab',
  //           options: {
  //             projectId: 4967944
  //           }
  //         },
  //         {
  //           plugin: 'gitlab',
  //           options: {
  //             projectId: 8967944
  //           }
  //         },
  //         {
  //           plugin: 'github',
  //           options: {
  //             owner: 'merico-dev',
  //             repositoryName: 'lake'
  //           }
  //         },
  //         {
  //           plugin: 'jenkins',
  //           options: {}
  //         },
  //         {
  //           plugin: 'jira',
  //           options: {
  //             boardId: 8,
  //             sourceId: 1
  //           }
  //         },
  //         {
  //           plugin: 'gitlab',
  //           options: {
  //             projectId: 8967944
  //           }
  //         }
  //       ]
  //     ],
  //     totalTasks: 7,
  //     finishedTasks: 4,
  //     beganAt: '2021-12-23T21:36:29.269Z',
  //     finishedAt: null,
  //     status: 'TASK_FAILED',
  //     message: '',
  //     spentSeconds: 0
  //   },
  //   {
  //     ID: 425,
  //     CreatedAt: '2021-12-21T18:03:21.365Z',
  //     UpdatedAt: '2021-12-21T18:03:21.729Z',
  //     name: 'COLLECT 1640109798011',
  //     tasks: [
  //       [
  //         {
  //           plugin: 'jenkins',
  //           options: {}
  //         }
  //       ]
  //     ],
  //     totalTasks: 1,
  //     finishedTasks: 1,
  //     beganAt: '2021-12-21T18:03:21.412Z',
  //     finishedAt: '2021-12-21T18:03:21.726Z',
  //     status: 'TASK_COMPLETED',
  //     message: '',
  //     spentSeconds: 0
  //   },
  //   {
  //     ID: 440,
  //     CreatedAt: '2021-12-23T15:08:06.025Z',
  //     UpdatedAt: '2021-12-23T20:09:53.786Z',
  //     name: 'config-ui trigger Thu Dec 23 2021 10:08:06 GMT-0500 (EST)',
  //     tasks: [
  //       [
  //         {
  //           plugin: 'gitlab',
  //           options: {
  //             projectId: 8967944
  //           }
  //         },
  //         {
  //           plugin: 'jira',
  //           options: {
  //             boardId: 8,
  //             sourceId: 1
  //           }
  //         },
  //         {
  //           plugin: 'jenkins',
  //           options: {}
  //         },
  //         {
  //           plugin: 'github',
  //           options: {
  //             owner: 'merico-dev',
  //             repositoryName: 'lake'
  //           }
  //         }
  //       ],
  //       [
  //         {
  //           plugin: 'gitlab',
  //           options: {
  //             projectId: 8967944
  //           }
  //         }
  //       ],
  //       [
  //         {
  //           plugin: 'gitlab',
  //           options: {
  //             projectId: 4967944
  //           }
  //         }
  //       ],
  //       [
  //         {
  //           plugin: 'gitlab',
  //           options: {
  //             projectId: 1967944
  //           }
  //         }
  //       ]
  //     ],
  //     totalTasks: 7,
  //     finishedTasks: 2,
  //     beganAt: '2021-12-23T15:08:06.12Z',
  //     finishedAt: null,
  //     status: 'TASK_FAILED',
  //     message: '',
  //     spentSeconds: 0
  //   },
  //   {
  //     ID: 18091,
  //     CreatedAt: '2021-12-23T23:37:58.588Z',
  //     UpdatedAt: '2021-12-23T23:38:59.083Z',
  //     name: 'COLLECT 1640302665464',
  //     tasks: [
  //       [
  //         {
  //           plugin: 'jenkins',
  //           options: {}
  //         },
  //         {
  //           plugin: 'github',
  //           options: {
  //             owner: 'e2corporation',
  //             repositoryName: 'getmdl-select'
  //           }
  //         }
  //       ]
  //     ],
  //     totalTasks: 2,
  //     finishedTasks: 2,
  //     beganAt: '2021-12-23T23:37:58.647Z',
  //     finishedAt: '2021-12-23T23:38:59.082Z',
  //     status: 'TASK_FAILED',
  //     message: "Error 1364: Field 'origin_key' doesn't have a default value",
  //     spentSeconds: 61
  //   },
  //   {
  //     ID: 18116,
  //     CreatedAt: '2022-01-05T18:04:53.903Z',
  //     UpdatedAt: '2022-01-05T18:04:54.648Z',
  //     name: 'COLLECT 1641405891157',
  //     tasks: [
  //       [
  //         {
  //           plugin: 'gitlab',
  //           options: {
  //             projectId: 8967944
  //           }
  //         },
  //         {
  //           plugin: 'github',
  //           options: {
  //             owner: 'merico-dev',
  //             repositoryName: 'lake'
  //           }
  //         },
  //         {
  //           plugin: 'jira',
  //           options: {
  //             boardId: 8,
  //             sourceId: 1
  //           }
  //         },
  //         {
  //           plugin: 'jenkins',
  //           options: {}
  //         }
  //       ]
  //     ],
  //     totalTasks: 4,
  //     finishedTasks: 2,
  //     beganAt: '2022-01-05T18:04:53.984Z',
  //     finishedAt: null,
  //     status: 'TASK_RUNNING',
  //     message: '',
  //     spentSeconds: 0
  //   }
  // ])

  const [filteredPipelines, setFilteredPipelines] = useState([])

  const filterPipelines = useCallback((status) => {
    setFilteredPipelines(status === 'all' ? pipelines : pipelines.filter((p) => p.status === status))
  }, [pipelines])

  const getPipelineCountByStatus = useCallback((status) => {
    return status === 'all' ? pipelines.length : pipelines.filter((p) => p.status === status).length
  }, [pipelines])

  useEffect(() => {
    fetchAllPipelines()
  }, [fetchAllPipelines])

  useEffect(() => {
    console.log('>>> Pipelines', pipelines)
    setFilteredPipelines(pipelines)
    if (pipelines.length > 0) {
      const latestPipelineRun = pipelines[0]
      fetchPipeline(latestPipelineRun.ID)
    }
  }, [pipelines, fetchPipeline])

  useEffect(() => {

  }, [pipelineCount])

  useEffect(() => {
    console.log('>> FILTER STATUS CHANGED ===> ', activeStatus)
    filterPipelines(activeStatus)
  }, [activeStatus])

  useEffect(() => {

  }, [refresh])

  useEffect(() => {
    console.log('>>> LATEST PIPELINE!', latestPipeline)
  }, [latestPipeline])

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
                { href: '/pipelines', icon: false, text: 'Manage Pipeline Runs', current: true },
              ]}
            />
            <div className='headlineContainer'>
              {/* <Link style={{ float: 'right', marginLeft: '10px', color: '#777777' }} to='/integrations'>
                <Icon icon='fast-backward' size={16} /> Go Back
              </Link> */}
              <div style={{ display: 'flex' }}>
                <div>
                  <span style={{ marginRight: '10px' }}>
                    <Icon icon='git-merge' size={32} />
                  </span>
                </div>
                <div>
                  <h1 style={{ margin: 0 }}>
                    Pipeline Run Logs
                    <Popover
                      className='trigger-manage-pipelines-help'
                      popoverClassName='popover-help-manage-pipelines'
                      position={Position.RIGHT}
                      autoFocus={false}
                      enforceFocus={false}
                      usePortal={false}
                    >
                      <a href='#' rel='noreferrer'><HelpIcon width={19} height={19} style={{ marginLeft: '10px' }} /></a>
                      <>
                        <div style={{ textShadow: 'none', fontSize: '12px', padding: '12px', maxWidth: '300px' }}>
                          <div style={{ marginBottom: '10px', fontWeight: 700, fontSize: '14px', fontFamily: '"Montserrat", sans-serif' }}>
                            <Icon icon='help' size={16} /> Manage Pipeline Runs
                          </div>
                          <p>Need Help? &mdash; Manage, Stop running and Restart failed pipelines.
                            Access <strong>Task Progress</strong> and Activity for all your pipelines.
                          </p>
                        </div>
                      </>
                    </Popover>
                  </h1>
                  <p className='page-description mb-0'>Manage Job Activity for all your pipeline runs.</p>
                  <p className=''>The most recent runs are show first, please select a time range.</p>
                </div>
                <div style={{ marginLeft: 'auto' }}>
                  <Button icon='add' intent={Intent.PRIMARY} text='Create Run' onClick={() => history.push('/pipelines/create')} />
                </div>
              </div>
            </div>

            {(isFetchingAll || !isFetchingAll) && (
              <>
                <div style={{ display: 'flex', marginTop: '30px', minHeight: '38px', width: '100%', justifyContent: 'space-between' }}>

                  <ButtonGroup className='filter-status-group' round='true' style={{ fontSize: '12px', zIndex: 0 }}>
                    <Button className='btn-pipeline-filter' intent={activeStatus === 'all' ? 'primary' : null} onClick={() => setActiveStatus('all')}>
                      <span style={{ marginRight: '10x', letterSpacing: '0', fontWeight: 900 }}>All&nbsp;</span>
                      <Tag className='tag-data-count'>{getPipelineCountByStatus('all')}</Tag>
                    </Button>
                    <Button className='btn-pipeline-filter' intent={activeStatus === 'TASK_RUNNING' ? 'primary' : null} onClick={() => setActiveStatus('TASK_RUNNING')}>
                      <span style={{ marginRight: '10x', letterSpacing: '0', fontWeight: 700 }}>Running&nbsp;</span>
                      <Tag className='tag-data-count'>{getPipelineCountByStatus('TASK_RUNNING')}</Tag>
                    </Button>
                    <Button className='btn-pipeline-filter' intent={activeStatus === 'TASK_COMPLETED' ? 'primary' : null} onClick={() => setActiveStatus('TASK_COMPLETED')}>
                      <span style={{ marginRight: '10x', letterSpacing: '0', fontWeight: 700 }}>Complete&nbsp;</span>
                      <Tag className='tag-data-count'>{getPipelineCountByStatus('TASK_COMPLETED')}</Tag>
                    </Button>
                    <Button className='btn-pipeline-filter' intent={activeStatus === 'TASK_FAILED' ? 'primary' : null} onClick={() => setActiveStatus('TASK_FAILED')}>
                      <span style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', marginRight: '10x', letterSpacing: '0', fontWeight: 700 }}>
                        <Icon icon='warning-sign' size={14} style={{ justifySelf: 'center', marginRight: '10px' }} />
                        Failed&nbsp;
                        <Tag className='tag-data-count'>{getPipelineCountByStatus('TASK_FAILED')}</Tag>

                      </span>
                    </Button>
                  </ButtonGroup>

                  <InputGroup
                    leftElement={<Icon icon='search' />}
                    placeholder='Search Pipelines'
                    rightElement={<Button text='GO' intent='primary' />}
                    round
                  />
                </div>

                <div style={{ display: 'flex', width: '100%', justifySelf: 'flex-start', marginTop: '8px' }}>
                  <Card interactive={false} elevation={Elevation.TWO} style={{ width: '100%', padding: '2px' }}>
                    <table className='bp3-html-table bp3-html-table-bordered connections-table' style={{ width: '100%' }}>
                      <thead>
                        <tr>
                          <th>ID</th>
                          <th>Pipeline Name</th>
                          <th>Duration</th>
                          <th>Status</th>
                          <th />
                        </tr>
                      </thead>
                      <tbody>
                        {!isFetchingAll && filteredPipelines.length > 0 && filteredPipelines.map((pipeline, pIdx) => (
                          <tr
                            key={`pipeline-row-${pIdx}`}
                            className={pipeline?.status === 'TASK_FAILED' ? 'pipeline-row pipeline-failed' : 'pipeline-row'}
                            style={{ verticalAlign: 'middle' }}
                          >
                            <td
                              style={{ cursor: 'pointer' }}
                              className='cell-id'
                            >
                              <Tooltip content='Pipeline Run ID' position={Position.TOP}>
                                <a href='#' onClick={() => history.push(`/pipelines/activity/${pipeline.ID}`)}>
                                  {pipeline.ID}
                                </a>
                              </Tooltip>
                            </td>

                            <td
                              onClick={(e) => history.push(`/pipelines/activity/${pipeline.ID}`)}
                              style={{ cursor: 'pointer' }}
                              className='cell-name'
                            >
                              {/* <Icon icon='power' color={Colors.GRAY4} size={10} style={{ float: 'right', marginLeft: '10px' }} /> */}

                              <span style={{ display: 'inline-block', float: 'right', color: '#999999', marginLeft: '15px' }}>{dayjs(pipeline.createdAt).format()}</span>

                              <strong style={{ lineHeight: '100%', fontSize: '12px', fontWeight: 800, textOverflow: 'ellipsis', overflow: 'hidden', display: 'block', whiteSpace: 'nowrap', maxWidth: '100%' }}>

                                {pipeline.name}
                                {pipeline.status === 'TASK_COMPLETED' && (<Icon icon='tick' size={10} color={Colors.GREEN5} style={{ margin: '0 10px', float: 'right', marginBottom: '2px' }} />)}
                              </strong>

                            </td>
                            <td
                              className='cell-duration'
                          // onClick={(e) => configureConnection(connection, e)}
                              style={{ cursor: 'pointer', whiteSpace: 'nowrap' }}
                            >

                              {/* {dayjs(pipeline.CreatedAt).toNow(pipeline.CreatedAt)} */}
                              {dayjs(pipeline.UpdatedAt).from(pipeline.CreatedAt, true)}
                            </td>
                            {/* <td
                          className='cell-status'
                          // onClick={(e) => configureConnection(connection, e)}
                          style={{ cursor: 'pointer' }}
                        >
                          {pipeline.status}
                        </td> */}
                            <td className='cell-status' style={{ textTransform: 'uppercase', whiteSpace: 'nowrap' }}>
                              <span style={{ display: 'inline-block', float: 'left', marginRight: '10px' }}>
                                <Tooltip content={`Progress ${pipeline.finishedTasks}/${pipeline.totalTasks} Tasks`}>
                                  {pipeline.status === 'TASK_RUNNING' &&
                                  (
                                    <Spinner
                                      style={{ margin: 0 }}
                                      className='mini-task-spinner' size={14} intent='warning'
                                      value={Number(pipeline.finishedTasks / pipeline.totalTasks).toFixed(1)}
                                    />
                                  )}
                                </Tooltip>
                                {pipeline.status === 'TASK_COMPLETED' &&
                                 (
                                   <Spinner
                                     style={{ margin: 0 }}
                                     className='mini-task-spinner' size={14} intent='success' value={1}
                                   />)}
                                {pipeline.status === 'TASK_FAILED' &&
                                 (
                                   <Spinner
                                     style={{ margin: 0 }} className='mini-task-spinner'
                                     size={14} intent='info' value={Number(pipeline.finishedTasks / pipeline.totalTasks).toFixed(1)}
                                   />)}
                              </span>
                              {pipeline.status === 'TASK_FAILED' && (
                                <strong style={{ color: Colors.RED5 }}>Failed</strong>
                              )}
                              {pipeline.status === 'TASK_COMPLETED' && (
                                <strong style={{ color: Colors.GREEN5 }}>Complete</strong>
                              )}
                              {pipeline.status === 'TASK_RUNNING' && (
                                <strong style={{ color: Colors.BLUE5 }}>Running</strong>
                              )}
                              {pipeline.status === 'TASK_CREATED' && (
                                <strong style={{ color: Colors.GRAY3 }}>
                                  <Icon icon='array' size={14} color={Colors.GRAY2} /> Pending...
                                </strong>
                              )}
                            </td>
                            <td className='cell-actions' style={{ padding: '0 10px', verticalAlign: 'middle' }}>
                              <div style={{
                                display: 'flex',
                                justifySelf: 'center',
                                gap: '3px',
                                alignSelf: 'center',
                                alignContent: 'center',
                                alignItems: 'center'
                              }}
                              >
                                <a
                                  href='#'
                                  onClick={() => history.push(`/pipelines/activity/${pipeline.ID}`)}
                                  data-provider={pipeline.id}
                                  className='bp3-button bp3-small bp3-minimal'
                                >
                                  <Icon icon='eye-open' size={16} />
                                </a>
                                {pipeline.status === 'TASK_FAILED' && (
                                  <a
                                    href='#'
                                    data-provider={pipeline.id}
                                    className='bp3-button bp3-small bp3-minimal'
                                  >
                                    <Icon icon='refresh' size={12} />
                                  </a>)}
                                {pipeline.status === 'TASK_RUNNING' && (

                                  <Popover
                                    key={`popover-help-key-cancel-run-${pipeline.ID}`}
                                    className='trigger-pipeline-cancel'
                                    popoverClassName='popover-pipeline-cancel'
                                    position={Position.BOTTOM}
                                    autoFocus={false}
                                    enforceFocus={false}
                                    usePortal={true}
                                    disabled={pipeline.status !== 'TASK_RUNNING'}
                                  >
                                    <a
                                      href='#'
                                      data-provider={pipeline.id}
                                      className='bp3-button bp3-small bp3-minimal'
                                    >
                                      <Icon icon='stop' size={16} style={{ color: Colors.RED5 }} />
                                    </a>
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
                                            className={Classes.POPOVER_DISMISS}
                                            text='YES' icon='small-tick' intent={Intent.DANGER} small
                                            onClick={() => cancelPipeline(activePipeline.ID)}
                                          />
                                        </div>
                                      </div>
                                    </>
                                  </Popover>

                                // <a
                                //   href='#'
                                //   data-provider={pipeline.id}
                                //   className='bp3-button bp3-small bp3-minimal'
                                // >
                                //   <Icon icon='stop' size={16} style={{ color: Colors.RED5 }} />
                                // </a>
                                )}
                                <a
                                  href='#'
                                  onClick={() => setInspectPipeline(pipeline) | setShowInspector(true)}
                                  data-provider={pipeline.id}
                                  className='bp3-button bp3-small bp3-minimal'
                                >
                                  <Icon icon='code' size={16} />
                                </a>
                              </div>
                              {/* {activeProvider?.multiSource && (
                            <DeleteAction
                              id={deleteId}
                              connection={connection}
                              text='Delete'
                              showConfirmation={() => setDeleteId(pipeline.ID)}
                              onConfirm={runDeletion}
                              onCancel={(e) => setDeleteId(false)}
                              isDisabled={isRunningDelete || isDeletingConnection}
                              isLoading={isRunningDelete || isDeletingConnection}
                            >
                              <DeleteConfirmationMessage title={`DELETE "${pipeline.name}"`} />
                            </DeleteAction>
                          )} */}

                            </td>
                          </tr>
                        ))}
                        {isFetchingAll && (
                          <tr>
                            <td className='loading-cell' colSpan='5' style={{ backgroundColor: '#f8f8f8' }}>
                              <ContentLoader
                                elevation={Elevation.ZERO}
                                cardStyle={{ border: '0 !important', boxShadow: 'none', backgroundColor: 'transparent' }}
                                title='Loading Pipeline Runs ...' message='Please wait while the pipeline run logs are loaded.'
                              />
                            </td>
                          </tr>
                        )}

                        {!isFetchingAll && filteredPipelines.length === 0 && (
                          <tr>
                            <td className='no-data-message-cell' colSpan='5' style={{ backgroundColor: '#fffcf0' }}>
                              <h3 style={{ fontWeight: 800, letterSpacing: '2px', textTransform: 'uppercase', margin: 0, fontFamily: '"Montserrat", sans-serif' }}>0 Pipeline Runs</h3>
                              <p style={{ margin: 0 }}>There are no pipeline logs for the current status <strong>{activeStatus}</strong>.</p>
                            </td>
                          </tr>
                        )}
                      </tbody>
                    </table>
                    {/* {maxConnectionsExceeded(sourceLimits[activeProvider.id], connections.length) && (
                  <p style={{ margin: 0, padding: '10px', backgroundColor: '#f0f0f0', borderTop: '1px solid #cccccc' }}>
                    <Icon icon='warning-sign' size='16' color={Colors.GRAY1} style={{ marginRight: '5px' }} />
                    You have reached the maximum number of allowed connections for this provider.
                  </p>
                )} */}
                  </Card>

                </div>
                <div style={{ marginTop: '10px', display: 'flex', width: '100%', justifySelf: 'flex-start' }}>

                  <div style={{ display: 'flex', width: '50%', fontSize: '11px', color: '#555555' }}>
                    <Icon icon='user' size={14} style={{ marginRight: '8px' }} />
                    <div>
                      <span>by {' '} <strong>Administrator</strong></span><br />
                      <span style={{ color: '#888888' }}>Displaying 6 of {pipelineCount} pipeline run log entries from API.</span>
                    </div>
                  </div>

                  <div style={{ display: 'flex', marginLeft: 'auto', marginRight: '20px' }}>

                    <Button small icon='add' style={{ marginRight: '5px' }} minimal onClick={() => history.push('/pipelines/create')} />
                    <Button small icon='refresh' minimal text='Refresh' onClick={() => fetchAllPipelines()} />
                  </div>
                  <div className='pagingation-controls' style={{ display: 'flex' }}>

                    <Button
                      className='pagination-btn btn-prev-page'
                      icon='step-backward' small text='PREV' style={{ marginRight: '5px' }} disabled
                    />
                    <Button className='pagination-btn btn-next-page' rightIcon='step-forward' small text='NEXT' />
                  </div>
                </div>
                <div style={{ height: '50px' }} />
              </>
            )}
            {/* <div style={{ marginTop: '100px', display: 'flex', width: '100%', justifyContent: 'flex-start' }}>
              <Button intent='secondary' icon='eye-open' text='VIEW' style={{ backgroundColor: '#eeeeee', color: '#888888' }} />
              <Button intent='primary' icon='doughnut-chart' text='View Graphs' style={{ backgroundColor: '#eeeeee', color: '#888888', marginLeft: '10px' }} />
              <Button intent='primary' icon='add' text='Add Provider' style={{ backgroundColor: '#eeeeee', color: '#888888', marginLeft: '10px' }} />
              <Button intent='success' icon='refresh' text='Running' style={{ backgroundColor: '#eeeeee', color: '#ffffff', marginLeft: '10px' }} />
            </div> */}
          </main>
        </Content>
      </div>
      {!isFetchingAll && activePipeline && (<PipelineIndicator pipeline={activePipeline} graphsUrl={GRAFANA_URL} onFetch={fetchPipeline} onCancel={cancelPipeline} />)}
      {!isFetchingAll && inspectPipeline && (<CodeInspector isOpen={showInspector} activePipeline={inspectPipeline} onClose={setShowInspector} />)}

    </>
  )
}

export default Pipelines
