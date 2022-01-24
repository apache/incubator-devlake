import React, { useState, useEffect } from 'react'
// import { CSSTransition } from 'react-transition-group'
import {
  Providers,
  ProviderLabels
} from '@/data/Providers'
import {
  // Classes,
  Icon,
  Spinner,
  Colors,
  Tooltip,
  Position,
  Intent,
  Card,
  Elevation,
  Popover,
  TextArea,
  H4,
  // Alignment
} from '@blueprintjs/core'
import dayjs from '@/utils/time'
import StageTaskName from '@/components/pipelines/StageTaskName'
import StageTaskIndicator from '@/components/pipelines/StageTaskIndicator'

const StageTask = (props) => {
  const {
    stages = [],
    task,
    sK,
    sIdx,
  } = props

  const [taskModuleOpened, setTaskModuleOpened] = useState()

  return (
    <>
      <Card
        elevation={task.status === 'TASK_RUNNING' ? Elevation.TWO : Elevation.ONE}
        className={`pipeline-task-module task-${task.status.split('_')[1].toLowerCase()} ${task.ID === taskModuleOpened?.ID ? 'active' : ''}`}
        onClick={() => setTaskModuleOpened(task)}
        style={{

        }}
      >
        <StageTaskIndicator task={task} />
        {/* <div
          className='task-module-status'
          style={{
            display: 'flex',
            justifyContent: 'center',
            padding: '8px',
            width: '32px',
            minWidth: '32px'
          }}
        >
          {task.status === 'TASK_COMPLETED' && (
            <Tooltip content={`Task Complete [STAGE ${task.pipelineRow}]`} position={Position.TOP} intent={Intent.SUCCESS}>
              <Icon icon='small-tick' size={18} color={Colors.GREEN5} style={{ marginLeft: '0' }} />
            </Tooltip>
          )}
          {task.status === 'TASK_FAILED' && (
            <Tooltip content={`Task Failed [STAGE ${task.pipelineRow}]`} position={Position.TOP} intent={Intent.PRIMARY}>
              <Spinner
                className='task-module-spinner'
                size={14}
                intent={Intent.WARNING}
                value={task.progress}
              />
            </Tooltip>
          )}
          {task.status === 'TASK_RUNNING' && (
            <Tooltip content={`Task Running [STAGE ${task.pipelineRow}]`} position={Position.TOP}>
              <Spinner
                className='task-module-spinner'
                size={14}
                intent={task.status === 'TASK_COMPLETED' ? 'success' : 'warning'}
                value={task.status === 'TASK_COMPLETED' ? 1 : task.progress}
              />
            </Tooltip>
          )}
          {task.status === 'TASK_CREATED' && (
            <Tooltip content={`Task Created (Pending) [STAGE ${task.pipelineRow}]`} position={Position.TOP}>
              <Spinner
                className='task-module-spinner'
                size={14}
                value={0}
              />
              {/* <Icon icon='stopwatch' size={14} color={Colors.GRAY3} style={{ marginLeft: '0', marginBottom: '3px' }} /> *\\/}
            </Tooltip>
          )}
        </div> */}
        <div
          className='task-module-name'
          style={{
            flex: 1
          }}
        >
          <div style={{ padding: '4px 2px 4px 0' }}>
            <StageTaskName task={task} showDetails={taskModuleOpened} onClose={() => setTaskModuleOpened(null)} />
            {/* <Popover
              className='trigger-pipeline-activity-help'
              popoverClassName='popover-help-pipeline-activity'
              position={Position.RIGHT}
              autoFocus={false}
              enforceFocus={false}
              usePortal={true}
            >
              <span>{task.plugin}</span>
              <>
                <div style={{ textShadow: 'none', fontSize: '12px', padding: '12px', maxWidth: '320px' }}>
                  <div style={{ marginBottom: '10px', fontWeight: 700, fontSize: '14px', fontFamily: '"Montserrat", sans-serif' }}>
                    <Icon icon='layer' size={16} /> {ProviderLabels[task.plugin.toUpperCase()] || 'System Task'}<br />
                  </div>
                  <div style={{ fontSize: '10px' }}>
                    <div style={{ display: 'flex', flexDirection: 'row', justifyContent: 'space-between' }}>
                      <div>
                        <label style={{ color: Colors.GRAY2 }}>ID</label><br />
                        <span>{task.ID}</span>
                      </div>
                      <div style={{ marginLeft: '20px' }}>
                        <label style={{ color: Colors.GRAY2 }}>Created</label><br />
                        <span>{dayjs(task.CreatedAt).format('L LT')}</span>
                      </div>
                      {task.finishedAt && (
                        <div style={{ marginLeft: '20px' }}>
                          <label style={{ color: Colors.GRAY2 }}>Finished</label><br />
                          <span>{dayjs(task.finishedAt).format('L LT')}</span>
                        </div>
                      )}
                    </div>
                    <div style={{ marginTop: '6px' }}>
                      <label style={{ color: Colors.GRAY2 }}>Status</label><br />
                      <strong>{task.status}</strong> <strong className='bp3-tag bp3-intent-primary' style={{ minHeight: '16px', fontSize: '10px', padding: '2px 6px', borderRadius: '6px'}}>PROGRESS {task.progress * 100}%</strong>
                    </div>
                    <div style={{ marginTop: '6px' }}>
                      <label style={{ color: Colors.GRAY2 }}>Options</label><br />
                      <TextArea
                        readonly
                        fill value={JSON.stringify(task.options)} style={{ fontSize: '10px', backgroundColor: '#f8f8f8', resize: 'none' }}
                      />
                      {/* <span>
                        <pre style={{ margin: 0 }}>
                          <code>
                            {JSON.stringify(task.options)}
                          </code>
                        </pre>
                      </span> *\\/}
                    </div>
                    {task.message !== '' && (
                      <div style={{ marginTop: '6px' }}>
                        <label style={{ color: Colors.DARK_GRAY1 }}>Message</label><br />
                        <span style={{ color: task.status === 'TASK_FAILED' ? Colors.RED3 : Colors.BLACK }}>
                          {task.status === 'TASK_FAILED' && (
                            <Icon
                              icon='warning-sign'
                              color={Colors.RED5}
                              size={10}
                              style={{ marginRight: '3px' }}
                            />
                          )}
                          {task.message}
                        </span>
                      </div>
                    )}
                  </div>
                </div>
              </>
            </Popover> */}

            <span style={{
              opacity: 0.4,
              display: 'block',
              width: '90%',
              fontSize: '9px',
              overflow: 'hidden',
              whiteSpace: 'nowrap',
              textOverflow: 'ellipsis'
            }}
            >
              {task.plugin !== 'github' && (<>ID {task.options.projectId || task.options.boardId}</>)}
              {task.plugin === 'github' && (<>@{task.options.owner}/{task.options.repositoryName}</>)}
            </span>
          </div>
        </div>
        <div
          className='task-module-progress'
          style={{
            fontWeight: task.status === 'TASK_COMPLETED' ? 700 : 300
          }}
        >
          {Number(task.status === 'TASK_COMPLETED' ? 100 : (task.progress / 1) * 100).toFixed(0)}%
        </div>
      </Card>
    </>
  )
}

export default StageTask
