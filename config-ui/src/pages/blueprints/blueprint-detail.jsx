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
import React, { useEffect, useState, useCallback, useRef } from 'react'
import { useParams, useHistory } from 'react-router-dom'
// import { CSSTransition } from 'react-transition-group'
import { DEVLAKE_ENDPOINT } from '@/utils/config'
import request from '@/utils/request'
import dayjs from '@/utils/time'
import { saveAs } from 'file-saver'
import {
  Button,
  Elevation,
  Intent,
  Switch,
  Card,
  Tooltip,
  Icon,
  Colors,
  Divider,
  Spinner,
  Classes,
  Position,
  Popover,
  Collapse,
  Dialog
} from '@blueprintjs/core'
import { NullBlueprint } from '@/data/NullBlueprint'
import { NullPipelineRun } from '@/data/NullPipelineRun'
// import { Providers, ProviderLabels, ProviderIcons } from '@/data/Providers'
import {
  StageStatus,
  TaskStatus,
  TaskStatusLabels,
  StatusColors,
  StatusBgColors
} from '@/data/Task'

// import TaskActivity from '@/components/pipelines/TaskActivity'
import CodeInspector from '@/components/pipelines/CodeInspector'
import StageLane from '@/components/pipelines/StageLane'
import { ToastNotification } from '@/components/Toast'
import BlueprintNavigationLinks from '@/components/blueprints/BlueprintNavigationLinks'

import useIntegrations from '@/hooks/useIntegrations'
import useBlueprintManager from '@/hooks/useBlueprintManager'
import usePipelineManager from '@/hooks/usePipelineManager'
import usePaginator from '@/hooks/usePaginator'

const BlueprintDetail = (props) => {
  const { integrations: Integrations, ProviderLabels } = useIntegrations()

  const { bId } = useParams()

  const [blueprintId, setBlueprintId] = useState()
  const [activeBlueprint, setActiveBlueprint] = useState(NullBlueprint)
  // eslint-disable-next-line no-unused-vars
  const [blueprintConnections, setBlueprintConnections] = useState([])
  const [blueprintPipelines, setBlueprintPipelines] = useState([])
  const [inspectedPipeline, setInspectedPipeline] = useState(NullPipelineRun)
  const [currentRun, setCurrentRun] = useState()
  const [showCurrentRunTasks, setShowCurrentRunTasks] = useState(true)
  const [showInspector, setShowInspector] = useState(false)
  const [currentStages, setCurrentStages] = useState([])

  const pollTimer = 5000
  const pollInterval = useRef()
  const [autoRefresh, setAutoRefresh] = useState(false)

  const [expandRun, setExpandRun] = useState(null)
  const [isDownloading, setIsDownloading] = useState(false)

  const {
    blueprint,
    getNextRunDate,
    activateBlueprint,
    deactivateBlueprint,
    fetchBlueprint
  } = useBlueprintManager()

  const {
    activePipeline,
    pipelines,
    isFetchingAll: isFetchingAllPipelines,
    fetchPipeline,
    runPipeline,
    cancelPipeline,
    fetchAllPipelines,
    lastRunId,
    setSettings: setPipelineSettings,
    // eslint-disable-next-line no-unused-vars
    allowedProviders,
    // eslint-disable-next-line no-unused-vars
    detectPipelineProviders,
    logfile: pipelineLogFilename,
    getPipelineLogfile,
    rerunAllFailedTasks,
    rerunTask
  } = usePipelineManager()

  const {
    data: historicalRuns,
    pagedData: pagedHistoricalRuns,
    setData: setHistoricalRuns,
    renderControlsComponent: renderPagnationControls
  } = usePaginator()

  const buildPipelineStages = useCallback((tasks = []) => {
    let stages = {}
    console.log('>>>> RECEIVED PIPELINE TASKS FOR STAGE...', tasks)
    tasks?.forEach((tS) => {
      stages = {
        ...stages,
        [tS.pipelineRow]: tasks?.filter((t) => t.pipelineRow === tS.pipelineRow)
      }
    })
    console.log('>>> BUILDING PIPELINE STAGES...', stages)
    return stages
  }, [])

  const runBlueprint = useCallback(() => {
    if (activeBlueprint !== null) {
      runPipeline(activeBlueprint.id)
    }
  }, [activeBlueprint, runPipeline])

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

  const handlePipelineDialogClose = useCallback(() => {
    setExpandRun(null)
  }, [])

  const inspectRun = useCallback((pipelineRun) => {
    setInspectedPipeline(pipelineRun)
    setShowInspector(true)
  }, [])

  const viewPipelineRun = useCallback(
    (pipelineRun) => {
      const fetchPipelineTasks = async () => {
        const t = await request.get(
          `${DEVLAKE_ENDPOINT}/pipelines/${pipelineRun?.id}/tasks`
        )
        setExpandRun({
          ...pipelineRun,
          tasks: t.data?.tasks || []
        })
      }
      if (pipelineRun?.id !== null) {
        fetchPipelineTasks()
      }
    },
    [setExpandRun]
  )

  const handleInspectorClose = useCallback(() => {
    setInspectedPipeline(NullPipelineRun)
    setShowInspector(false)
  }, [])

  const getTaskStatusIcon = (status) => {
    let icon = null
    switch (status) {
      case TaskStatus.ACTIVE:
      case TaskStatus.RUNNING:
        icon = <Spinner size={14} intent={Intent.PRIMARY} />
        break
      case TaskStatus.COMPLETE:
        icon = <Icon icon='tick-circle' size={14} color={Colors.GREEN5} />
        break
      case TaskStatus.FAILED:
        icon = <Icon icon='delete' size={14} color={Colors.RED5} />
        break
      case TaskStatus.CANCELLED:
        icon = <Icon icon='undo' size={14} color={Colors.RED5} />
        break
      case TaskStatus.CREATED:
        icon = <Icon icon='stopwatch' size={14} color={Colors.GRAY3} />
        break
    }
    return icon
  }

  const downloadPipelineLog = useCallback(
    (pipeline) => {
      console.log(
        `>>> DOWNLOADING PIPELINE #${pipeline?.id}  LOG...`,
        getPipelineLogfile(pipeline?.id)
      )
      setIsDownloading(true)
      ToastNotification.clear()
      let downloadStatus = 404
      const checkStatusAndDownload = async (pipeline) => {
        const d = await request.get(getPipelineLogfile(pipeline?.id))
        downloadStatus = d?.status
        if (pipeline?.id && downloadStatus === 200) {
          saveAs(getPipelineLogfile(pipeline?.id), pipelineLogFilename)
          setIsDownloading(false)
        } else if (pipeline?.id && downloadStatus === 404) {
          ToastNotification.show({
            message: d?.message || 'Logfile not available',
            intent: 'danger',
            icon: 'error'
          })
          setIsDownloading(false)
        } else {
          ToastNotification.show({
            message: 'Pipeline Invalid or Missing',
            intent: 'danger',
            icon: 'error'
          })
          setIsDownloading(false)
        }
      }
      checkStatusAndDownload(pipeline)
    },
    [getPipelineLogfile, pipelineLogFilename]
  )

  useEffect(() => {
    setBlueprintId(bId)
    console.log('>>> REQUESTED BLUEPRINT ID ===', bId)
  }, [bId])

  useEffect(() => {
    if (blueprintId) {
      fetchBlueprint(blueprintId)
      fetchAllPipelines(blueprintId)
    }
  }, [lastRunId, autoRefresh, blueprintId, fetchBlueprint, fetchAllPipelines])

  useEffect(() => {
    console.log('>>>> SETTING ACTIVE BLUEPRINT...', blueprint)
    if (blueprint?.id) {
      setActiveBlueprint((b) => ({
        ...b,
        ...blueprint,
        id: blueprint.id,
        name: blueprint.name
      }))
      setBlueprintConnections(
        blueprint?.settings?.connections.map((connection, cIdx) => ({
          id: cIdx,
          provider: connection?.plugin,
          name: `${
            ProviderLabels[connection?.plugin.toUpperCase()]
          } Connection (ID #${connection?.connectionId})`,
          dataScope: connection?.scope
            .map((s) => [`${s.options?.owner}/${s?.options?.repo}`])
            .join(', '),
          dataDomains: []
        }))
      )
      setPipelineSettings({
        name: `${blueprint?.name} ${Date.now()}`,
        blueprintId: blueprint?.id,
        plan: blueprint?.plan
      })
    }
  }, [blueprint, setPipelineSettings, ProviderLabels])

  useEffect(() => {
    console.log('>>>> FETCHED ALL PIPELINES..', pipelines, activeBlueprint?.id)
    setBlueprintPipelines(
      pipelines.filter((p) => p.blueprintId === activeBlueprint?.id)
    )
  }, [pipelines, activeBlueprint])

  useEffect(() => {
    console.log('>>>> RELATED BLUEPRINT PIPELINES..', blueprintPipelines)
    fetchPipeline(blueprintPipelines[0]?.id)
    setHistoricalRuns(
      blueprintPipelines.map((p, pIdx) => ({
        id: p.id,
        status: p.status,
        statusLabel: TaskStatusLabels[p.status],
        statusIcon: getTaskStatusIcon(p.status),
        startedAt: p.beganAt ? dayjs(p.beganAt).format('L LTS') : '-',
        completedAt: p.finishedAt ? dayjs(p.updatedAt).format('L LTS') : ' - ',
        duration:
          p.beganAt && p.finishedAt
            ? dayjs(p.beganAt).from(p.finishedAt, true)
            : p.beganAt &&
              [TaskStatus.RUNNING, TaskStatus.CREATED].includes(p?.status)
            ? dayjs(p.beganAt).toNow(true)
            : ' - '
      }))
    )
  }, [blueprintPipelines, setHistoricalRuns])

  useEffect(() => {
    if (
      activePipeline?.id &&
      [
        TaskStatus.CREATED,
        TaskStatus.RUNNING,
        TaskStatus.COMPLETE,
        TaskStatus.FAILED
      ].includes(activePipeline.status)
    ) {
      setCurrentStages(buildPipelineStages(activePipeline.tasks))
      setAutoRefresh(
        [TaskStatus.RUNNING, TaskStatus.CREATED].includes(
          activePipeline?.status
        )
      )
      setCurrentRun((cR) => ({
        ...cR,
        id: activePipeline.id,
        status: activePipeline.status,
        statusLabel: TaskStatusLabels[activePipeline.status],
        icon: getTaskStatusIcon(activePipeline.status),
        startedAt: activePipeline.beganAt
          ? dayjs(activePipeline.beganAt).format('L LTS')
          : '-',
        duration: [TaskStatus.CREATED, TaskStatus.RUNNING].includes(
          activePipeline.status
        )
          ? dayjs(activePipeline.beganAt || activePipeline.createdAt).toNow(
              true
            )
          : dayjs(activePipeline.beganAt).from(
              activePipeline.finishedAt || activePipeline.updatedAt,
              true
            ),
        stage: `Stage ${activePipeline.stage}`,
        tasksFinished: Number(activePipeline.finishedTasks),
        tasksTotal: Number(activePipeline.totalTasks),
        error: activePipeline.message || null
      }))
    }
  }, [activePipeline])

  useEffect(() => {
    console.log('>> BUILDING CURRENT STAGES...', currentStages)
  }, [currentStages])

  useEffect(() => {
    if (activePipeline?.id) {
      if (autoRefresh) {
        console.log('>> ACTIVITY POLLING ENABLED!')
        pollInterval.current = setInterval(() => {
          fetchPipeline(activePipeline?.id)
          // setLastPipeline(activePipeline)
        }, pollTimer)
        return () => {
          clearInterval(pollInterval.current)
        }
      }
    }
  }, [autoRefresh, fetchPipeline, activePipeline?.id, pollTimer])

  // useEffect(() => {
  //   console.log('>> VIEW PIPELINE RUN....', expandRun)
  // }, [expandRun])

  return (
    <>
      <main className='main'>
        <div
          className='blueprint-header'
          style={{
            display: 'flex',
            width: '100%',
            justifyContent: 'space-between',
            marginBottom: '10px'
          }}
        >
          <div className='blueprint-name' style={{}}>
            <h2 style={{ fontWeight: 'bold' }}>{activeBlueprint?.name}</h2>
          </div>
          <div
            className='blueprint-info'
            style={{ display: 'flex', alignItems: 'center' }}
          >
            <div className='blueprint-schedule'>
              <span
                className='blueprint-schedule-interval'
                style={{ textTransform: 'capitalize', padding: '0 10px' }}
              >
                {activeBlueprint?.interval} (at{' '}
                {dayjs(getNextRunDate(activeBlueprint?.cronConfig)).format(
                  'hh:mm A'
                )}
                )
              </span>{' '}
              &nbsp;{' '}
              <span className='blueprint-schedule-nextrun'>
                {activeBlueprint?.isManual ? (
                  <strong>Manual Mode</strong>
                ) : (
                  <>
                    Next Run{' '}
                    {dayjs(
                      getNextRunDate(activeBlueprint?.cronConfig)
                    ).fromNow()}
                  </>
                )}
              </span>
            </div>
            <div className='blueprint-actions' style={{ padding: '0 10px' }}>
              <Button
                intent={Intent.PRIMARY}
                small
                text='Run Now'
                onClick={runBlueprint}
                disabled={
                  !activeBlueprint?.enable ||
                  [TaskStatus.CREATED, TaskStatus.RUNNING].includes(
                    currentRun?.status
                  )
                }
              />
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
                onChange={() => handleBlueprintActivation(activeBlueprint)}
                style={{
                  marginBottom: 0,
                  marginTop: 0,
                  color: !activeBlueprint?.enable ? Colors.GRAY3 : 'inherit'
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

        <BlueprintNavigationLinks blueprint={activeBlueprint} />

        <div
          className='blueprint-run'
          style={{
            width: '100%',
            alignSelf: 'flex-start',
            minWidth: '750px'
          }}
        >
          <h3>Current Run</h3>
          <Card
            className={`current-run status-${currentRun?.status.toLowerCase()}`}
            elevation={Elevation.TWO}
            style={{ padding: '12px', marginBottom: '8px' }}
          >
            {currentRun && (
              <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                <div>
                  <label style={{ color: '#94959F' }}>Status</label>
                  <div style={{ display: 'flex' }}>
                    <span style={{ marginRight: '6px', marginTop: '2px' }}>
                      {currentRun?.icon}
                    </span>
                    <h4
                      className={`status-${currentRun?.status.toLowerCase()}`}
                      style={{ fontSize: '15px', margin: 0, padding: 0 }}
                    >
                      {currentRun?.statusLabel}
                    </h4>
                  </div>
                </div>
                <div>
                  <label style={{ color: '#94959F' }}>Started at</label>
                  <h4 style={{ fontSize: '15px', margin: 0, padding: 0 }}>
                    {currentRun?.startedAt}
                  </h4>
                </div>
                <div>
                  <label style={{ color: '#94959F' }}>Duration</label>
                  <h4 style={{ fontSize: '15px', margin: 0, padding: 0 }}>
                    {currentRun?.duration}
                  </h4>
                </div>
                <div>
                  <label style={{ color: '#94959F' }}>Current Stage</label>
                  <h4 style={{ fontSize: '15px', margin: 0, padding: 0 }}>
                    {currentRun?.stage}
                  </h4>
                </div>
                <div>
                  <label style={{ color: '#94959F' }}>Tasks Completed</label>
                  <h4 style={{ fontSize: '15px', margin: 0, padding: 0 }}>
                    {currentRun?.tasksFinished} / {currentRun?.tasksTotal}
                  </h4>
                </div>
                <div
                  style={{
                    display: 'flex',
                    justifyContent: 'center',
                    alignItems: 'center'
                  }}
                >
                  {currentRun?.status === 'TASK_RUNNING' && (
                    <div style={{ display: 'block' }}>
                      {/* <Button intent={Intent.PRIMARY} outlined text='Cancel' onClick={cancelRun} /> */}
                      <Popover
                        key='popover-help-key-cancel-run'
                        className='trigger-pipeline-cancel'
                        popoverClassName='popover-pipeline-cancel'
                        position={Position.BOTTOM}
                        autoFocus={false}
                        enforceFocus={false}
                        usePortal={true}
                      >
                        <Button
                          // icon='stop'
                          text='Cancel'
                          intent={Intent.PRIMARY}
                          outlined
                        />
                        <>
                          <div
                            style={{
                              fontSize: '12px',
                              padding: '12px',
                              maxWidth: '200px'
                            }}
                          >
                            <p>
                              Are you Sure you want to cancel this{' '}
                              <strong>Pipeline Run</strong>?
                            </p>
                            <div
                              style={{
                                display: 'flex',
                                width: '100%',
                                justifyContent: 'flex-end'
                              }}
                            >
                              <Button
                                text='NO'
                                minimal
                                small
                                className={Classes.POPOVER_DISMISS}
                                style={{
                                  marginLeft: 'auto',
                                  marginRight: '3px'
                                }}
                              />
                              <Button
                                className={Classes.POPOVER_DISMISS}
                                text='YES'
                                icon='small-tick'
                                intent={Intent.DANGER}
                                small
                                onClick={() => cancelPipeline(currentRun?.id)}
                              />
                            </div>
                          </div>
                        </>
                      </Popover>
                    </div>
                  )}
                  {currentRun?.status === TaskStatus.COMPLETE && (
                    <Button
                      intent={Intent.PRIMARY}
                      onClick={rerunAllFailedTasks}
                    >
                      Run Failed Tasks
                    </Button>
                  )}
                </div>
              </div>
            )}
            {!currentRun && (
              <>
                <p style={{ margin: 0 }}>
                  There is no current run for this blueprint.
                </p>
              </>
            )}
            {currentRun?.error && (
              <div style={{ marginTop: '10px' }}>
                <p className='error-msg' style={{ color: '#E34040' }}>
                  {currentRun?.error}
                </p>
              </div>
            )}
          </Card>
          {currentRun && (
            <Card
              elevation={Elevation.TWO}
              style={{ padding: '12px', marginBottom: '8px' }}
            >
              <div
                className='blueprint-run-activity'
                style={{ display: 'flex', width: '100%' }}
              >
                <div
                  className='pipeline-task-activity'
                  style={{
                    flex: 1,
                    padding: Object.keys(currentStages).length === 1 ? '0' : 0,
                    overflow: 'hidden',
                    textOverflow: 'ellipsis'
                  }}
                >
                  {Object.keys(currentStages).length > 0 && (
                    <div className='pipeline-multistage-activity'>
                      {Object.keys(currentStages).map((sK, sIdx) => (
                        <StageLane
                          key={`stage-lane-key-${sIdx}`}
                          stages={currentStages}
                          sK={sK}
                          sIdx={sIdx}
                          showStageTasks={showCurrentRunTasks}
                          rerunTask={rerunTask}
                        />
                      ))}
                    </div>
                  )}
                </div>

                <Button
                  icon={showCurrentRunTasks ? 'chevron-down' : 'chevron-right'}
                  intent={Intent.NONE}
                  minimal
                  small
                  style={{
                    textAlign: 'center',
                    display: 'block',
                    float: 'right',
                    margin: '0 10px',
                    marginBottom: 'auto'
                  }}
                  onClick={() => setShowCurrentRunTasks((s) => !s)}
                />
              </div>
            </Card>
          )}
        </div>

        <div
          className='blueprint-historical-runs'
          style={{
            width: '100%',
            alignSelf: 'flex-start',
            minWidth: '750px'
          }}
        >
          <h3>Historical Runs</h3>
          <Card
            elevation={Elevation.TWO}
            style={{ padding: '0', marginBottom: '8px' }}
          >
            <table
              className='bp3-html-table bp3-html-table historical-runs-table'
              style={{ width: '100%' }}
            >
              <thead>
                <tr>
                  <th style={{ minWidth: '100px', whiteSpace: 'nowrap' }}>
                    Status
                  </th>
                  <th style={{ minWidth: '100px', whiteSpace: 'nowrap' }}>
                    Started at
                  </th>
                  <th style={{ minWidth: '100px', whiteSpace: 'nowrap' }}>
                    Completed at
                  </th>
                  <th style={{ minWidth: '100px', whiteSpace: 'nowrap' }}>
                    Duration
                  </th>
                  <th style={{ width: '100%', whiteSpace: 'nowrap' }} />
                </tr>
              </thead>
              <tbody>
                {pagedHistoricalRuns.map((run, runIdx) => (
                  <tr key={`historical-run-key-${runIdx}`}>
                    <td
                      style={{
                        width: '15%',
                        whiteSpace: 'nowrap',
                        borderBottom: '1px solid #f0f0f0'
                      }}
                    >
                      <span
                        style={{
                          display: 'inline-block',
                          float: 'left',
                          marginRight: '5px'
                        }}
                      >
                        {run.statusIcon}
                      </span>{' '}
                      {run.statusLabel}
                    </td>
                    <td
                      style={{
                        width: '25%',
                        whiteSpace: 'nowrap',
                        borderBottom: '1px solid #f0f0f0'
                      }}
                    >
                      {run.startedAt}
                    </td>
                    <td
                      style={{
                        width: '25%',
                        whiteSpace: 'nowrap',
                        borderBottom: '1px solid #f0f0f0'
                      }}
                    >
                      {run.completedAt}
                    </td>
                    <td
                      style={{
                        width: '15%',
                        whiteSpace: 'nowrap',
                        borderBottom: '1px solid #f0f0f0'
                      }}
                    >
                      {run.duration}
                    </td>
                    <td
                      style={{
                        textAlign: 'right',
                        borderBottom: '1px solid #f0f0f0',
                        whiteSpace: 'nowrap'
                      }}
                    >
                      <Tooltip intent={Intent.PRIMARY} content='View JSON'>
                        <Button
                          intent={Intent.PRIMARY}
                          minimal
                          small
                          icon='code'
                          onClick={() =>
                            inspectRun(
                              blueprintPipelines.find((p) => p.id === run.id)
                            )
                          }
                        />
                      </Tooltip>
                      <Tooltip
                        intent={Intent.PRIMARY}
                        content='Download Full Log'
                      >
                        <Button
                          intent={Intent.NONE}
                          loading={isDownloading}
                          minimal
                          small
                          icon='document'
                          style={{ marginLeft: '10px' }}
                          onClick={() =>
                            downloadPipelineLog(
                              blueprintPipelines.find((p) => p.id === run.id)
                            )
                          }
                        />
                      </Tooltip>
                      <Tooltip
                        intent={Intent.PRIMARY}
                        content='Show Run Activity'
                      >
                        <Button
                          intent={Intent.PRIMARY}
                          minimal
                          small
                          icon={
                            expandRun?.id === run.id
                              ? 'chevron-down'
                              : 'chevron-right'
                          }
                          style={{ marginLeft: '10px' }}
                          onClick={() =>
                            viewPipelineRun(
                              blueprintPipelines.find((p) => p.id === run.id)
                            )
                          }
                        />
                      </Tooltip>
                    </td>
                  </tr>
                ))}
                {historicalRuns.length === 0 && (
                  <tr>
                    <td colSpan={5}>
                      There are no historical runs associated with this
                      blueprint.
                    </td>
                  </tr>
                )}
              </tbody>
            </table>
          </Card>
        </div>
        {historicalRuns.length > 0 && (
          <div style={{ alignSelf: 'flex-end', padding: '10px' }}>
            {renderPagnationControls()}
          </div>
        )}
      </main>
      <CodeInspector
        isOpen={showInspector}
        activePipeline={inspectedPipeline}
        onClose={handleInspectorClose}
      />
      <Dialog
        className='dialog-view-pipeline'
        // icon=
        title={`Historical Run #${expandRun?.id}`}
        isOpen={expandRun !== null}
        onClose={handlePipelineDialogClose}
        onClosed={() => {}}
        canOutsideClickClose={true}
        style={{ backgroundColor: '#ffffff' }}
      >
        <div className={Classes.DIALOG_BODY}>
          {Object.keys(buildPipelineStages(expandRun?.tasks)).length > 0 && (
            <div className='pipeline-multistage-activity'>
              {Object.keys(buildPipelineStages(expandRun?.tasks)).map(
                (sK, sIdx) => (
                  <StageLane
                    key={`stage-lane-key-${sIdx}`}
                    stages={buildPipelineStages(expandRun?.tasks)}
                    sK={sK}
                    sIdx={sIdx}
                    showStageTasks={true}
                  />
                )
              )}
            </div>
          )}
        </div>
      </Dialog>
    </>
  )
}

export default BlueprintDetail
