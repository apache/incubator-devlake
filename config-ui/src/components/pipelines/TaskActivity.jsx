
import React, { useState, useEffect } from 'react'
// import { CSSTransition } from 'react-transition-group'
import { Providers } from '@/data/Providers'
import {
  Icon,
  Spinner,
  Colors,
  Tooltip,
  Position,
  Intent,
  Popover
} from '@blueprintjs/core'
import dayjs from '@/utils/time'
import StageLane from '@/components/pipelines/StageLane'

const TaskActivity = (props) => {
  const { activePipeline, stages = [] } = props

  const [progressDetail, setProgressDetail] = useState({
    totalSubTasks: 0,
    finishedSubTasks: 0,
    totalRecords: -1,
    finishedRecords: 0,
    subTaskName: null,
    subTaskNumber: 1
  })

  const getActiveTask = (tasks = []) => {
    return tasks.find(t => t.status === 'TASK_RUNNING')
  }

  useEffect(() => {
    setProgressDetail(initialDetail => getActiveTask(activePipeline.tasks)
      ? getActiveTask(activePipeline.tasks).progressDetail
      : initialDetail)
  }, [activePipeline])

  return (
    <>

      <div
        className='pipeline-task-activity' style={{
          // padding: '20px',
          padding: Object.keys(stages).length === 1 ? '10px' : 0,
          overflow: 'hidden',
          textOverflow: 'ellipsis'
        }}
      >
        {Object.keys(stages).length > 1 && (
          <div
            className='pipeline-multistage-activity'
          >
            {Object.keys(stages).map((sK, sIdx) => (
              <StageLane key={`stage-lane-key-${sIdx}`} stages={stages} sK={sK} sIdx={sIdx} />
            ))}
          </div>
        )}
        {Object.keys(stages).length === 1 && activePipeline?.ID && activePipeline.tasks && activePipeline.tasks.map((t, tIdx) => (
          <div
            className='pipeline-task-row'
            key={`pipeline-task-key-${tIdx}`}
            style={{
              display: 'flex',
              padding: '4px 0',
              justifyContent: 'space-between',
              fontSize: '12px',
              opacity: t.status === 'TASK_CREATED' ? 0.7 : 1
            }}
          >
            <div style={{ display: 'flex', justifyContent: 'center', paddingRight: '8px', width: '32px', minWidth: '32px' }}>
              {t.status === 'TASK_COMPLETED' && (
                <Tooltip content={`Task Complete [STAGE ${t.pipelineRow}]`} position={Position.TOP} intent={Intent.SUCCESS}>
                  <Icon icon='small-tick' size={18} color={Colors.GREEN5} style={{ marginLeft: '0' }} />
                </Tooltip>
              )}
              {t.status === 'TASK_FAILED' && (
                <Tooltip content={`Task Failed [STAGE ${t.pipelineRow}]`} position={Position.TOP} intent={Intent.PRIMARY}>
                  <Icon icon='warning-sign' size={14} color={Colors.RED5} style={{ marginLeft: '0', marginBottom: '3px' }} />
                </Tooltip>
              )}
              {t.status === 'TASK_RUNNING' && (
                <Tooltip content={`Task Running [STAGE ${t.pipelineRow}]`} position={Position.TOP}>
                  <Spinner
                    className='task-spinner'
                    size={14}
                    intent={t.status === 'TASK_COMPLETED' ? 'success' : 'warning'}
                    value={t.status === 'TASK_COMPLETED' ? 1 : t.progress}
                  />
                </Tooltip>
              )}
              {t.status === 'TASK_CREATED' && (
                <Tooltip content={`Task Created (Pending) [STAGE ${t.pipelineRow}]`} position={Position.TOP}>
                  <Icon icon='pause' size={14} color={Colors.GRAY3} style={{ marginLeft: '0', marginBottom: '3px' }} />
                </Tooltip>
              )}
            </div>
            <div
              className='pipeline-task-cell-name'
              style={{ padding: '0 8px', minWidth: '130px', display: 'flex', justifyContent: 'space-between' }}
            >
              <strong
                className='task-plugin-name'
                style={{ overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}
              >
                {t.plugin}
              </strong>

            </div>
            <div
              className='pipeline-task-cell-settings'
              style={{
                padding: '0 8px',
                display: 'flex',
                width: '25%',
                minWidth: '25%',
                // whiteSpace: 'nowrap',
                justifyContent: 'flex-start',
                // overflow: 'hidden',
                // textOverflow: 'ellipsis',
                flexGrow: 1
              }}
            >
              {t.plugin !== Providers.JENKINS && t.plugin !== 'refdiff' && (
                <div style={{ overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap', display: 'flex', justifyContent: 'flex-start' }}>
                  <span style={{ color: Colors.GRAY2 }}>
                    <Icon icon='link' size={8} style={{ marginBottom: '3px', alignSelf: 'flex-start' }} />{' '}
                    {t.options.projectId || t.options.boardId || t.options.owner || t.options.projectName}
                  </span>
                  {t.plugin === Providers.GITHUB && (
                    <span style={{ fontWeight: 600 }}>/{t.options.repositoryName || t.options.repo || '(Repository)'}</span>
                  )}
                  {t.plugin === Providers.GITEXTRACTOR && (
                    <div style={{ paddingLeft: '12px', display: 'inline-block' }}>
                      <span>{t.options.url}</span><br />
                      <strong>{t.options.repoId}</strong>
                    </div>
                  )}
                  {t.plugin === Providers.DBT && (
                    <span style={{ fontWeight: 600 }}>&nbsp;{t.options.projectPath}</span>
                  )}
                </div>
              )}
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
              <span style={{ whiteSpace: 'nowrap' }}>
                {(() => {
                  let statusRelativeTime = dayjs(t.CreatedAt).toNow(true)
                  switch (t.status) {
                    case 'TASK_COMPLETED':
                    case 'TASK_FAILED':
                      statusRelativeTime = t.finishedAt == null ? 'N/A' : dayjs(t.finishedAt).from(t.beganAt, true)
                      break
                    case 'TASK_RUNNING':
                    default:
                      statusRelativeTime = dayjs(t.beganAt).toNow(true)
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
              style={{ display: 'flex', flexGrow: 1, width: '60%' }}
            >
              {t.message && (
                <div style={{ width: '98%', whiteSpace: 'wrap', overflow: 'hidden', textOverflow: 'ellipsis', paddingLeft: '10px' }}>
                  <span style={{ color: t.status === 'TASK_FAILED' ? Colors.RED4 : Colors.GRAY3 }}>
                    {t.message.length > 255
                      ? (
                        <Popover><>`${t.message.slice(0, 255)}...`</>
                          <div style={{
                            maxWidth:
                            '300px',
                            maxHeight: '300px',
                            padding: '10px',
                            overflow: 'auto',
                            backgroundColor: '#f8f8f8'
                          }}
                          >
                            <h3 style={{ margin: '5px 0', color: Colors.RED5 }}>
                              <Icon
                                icon='warning-sign'
                                size={14}
                                color={Colors.RED5} style={{ float: 'left', margin: '2px 4px 0 0' }}
                              />
                              ERROR MESSAGE <small style={{ color: Colors.GRAY3 }}> (Extended) </small>
                            </h3>
                            {t.message}
                            <p style={{ margin: '10px 0', color: Colors.GRAY3 }}> &gt; Please check the console log for more details... </p>
                          </div>
                        </Popover>
                        )
                      : t.message}
                  </span>
                </div>
              )}
            </div>
          </div>
        ))}
        {(!activePipeline.tasks || activePipeline.tasks.length === 0) && (
          <>
            <div style={{ display: 'flex' }}>
              <Icon
                icon='warning-sign'
                size={12}
                color={Colors.ORANGE5} style={{ float: 'left', margin: '0 4px 0 0' }}
              />
              <p>
                <strong>Missing Configuration</strong>, this pipeline has no tasks.
                <br />Please create a new pipeline with a valid configuration.
              </p>
            </div>
          </>
        )}
      </div>
      {progressDetail && progressDetail.subTaskName !== null && (
        <div
          className='pipeline-progress-detail' style={{
            backgroundColor: 'rgb(235, 243, 255)',
            padding: '0',
            borderTop: '1px solid rgb(0, 102, 255)'
          }}
        >
          <h2 className='headline' style={{ margin: '10px 20px' }}>
            <span style={{ display: 'inline-block', margin: '0 5px', float: 'right' }}>
              <Spinner
                className='task-details-spinner'
                size={14}
                intent={Intent.NONE}
                value={null}
              />
            </span>
            TASK DETAILS
          </h2>
          <table className='bp3-html-table striped bordered' style={{ width: '100%' }}>
            <thead>
              <tr>
                <th style={{ paddingLeft: '20px' }}>Sub-Tasks</th>
                <th>Records</th>
                <th>Subtask ID & Name</th>
                <th />
              </tr>
            </thead>
            <tbody>
              <tr>
                <td style={{ paddingLeft: '20px' }}>
                  <span style={{ color: 'rgb(0, 102, 255)', fontWeight: 'bold' }}>
                    {progressDetail.finishedSubTasks}
                  </span> / {progressDetail.totalSubTasks}
                </td>
                <td>
                  <span style={{ color: 'rgb(0, 102, 255)', fontWeight: 'bold' }}>
                    {progressDetail.finishedRecords}
                  </span> / {progressDetail.totalRecords}
                  
                </td>
                <td>{progressDetail.subTaskNumber}: <span style={{ color: 'rgb(0, 102, 255)', fontWeight: 'bold' }}>{progressDetail.subTaskName}</span></td>
                <td />
              </tr>
            </tbody>
          </table>
        </div>
      )}
    </>
  )
}

export default TaskActivity
