
import React from 'react'
// import { CSSTransition } from 'react-transition-group'
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
  H4,
  // Alignment
} from '@blueprintjs/core'
import dayjs from '@/utils/time'

const StageTask = (props) => {
  const { stages = [], task, sK, sIdx } = props

  return (
    <>
      <Card
        elevation={task.status === 'TASK_RUNNING' ? Elevation.TWO : Elevation.ONE}
        className={`pipeline-task-module task-${task.status.split('_')[1].toLowerCase()}`}
        style={{
          // display: 'flex',
          // padding: '0',
          // border: task.status === 'TASK_COMPLETED' ? '2px solid rgba(61, 204, 145, 0.95)' : task.status === 'TASK_FAILED' ? '2px solid rgba(255, 115, 115, 0.9)' : '2px solid #ddd',
          // boxShadow: task.status === 'TASK_FAILED' ? '0 0 0 1px #e8471c, 0 0 0 3px rgba(232, 71, 28, 0.3), 0 1px 1px rgba(16, 22, 26, 0.2)' : 'initial',
          // borderRadius: '12px',
          // backgroundColor: task.status === 'TASK_COMPLETED' ? 'rgba(61, 204, 145, 0.15)' : '',
          // margin: '5px',
          // flex: 1,
          // width: '80%',
          // minWidth: '100px',
          // maxWidth: '210px',
          // minHeight: '32px',
          // height: 'auto',
          // justifyContent: 'flex-start',
          // alignItems: 'center'
        }}
      >
        <div className='task-module-status' style={{ display: 'flex', justifyContent: 'center', padding: '8px', width: '32px', minWidth: '32px' }}>
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
              {/* <Icon icon='stopwatch' size={14} color={Colors.GRAY3} style={{ marginLeft: '0', marginBottom: '3px' }} /> */}
            </Tooltip>
          )}
        </div>
        <div
          className='task-module-name'
          style={{
            flex: 1,
            // borderRight: '1px solid #dddddd',
          }}
        >
          <div style={{ padding: '4px 2px 4px 0' }}>
            {task.plugin}<br />
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
              {task.plugin === 'github' && (<>{task.options.owner}/{task.options.repositoryName}</>)}
            </span>
          </div>
        </div>
        <div
          className='task-module-progress'
          style={{
            // display: 'flex',
            // width: '50px',
            // overflow: 'hidden',
            // textOverflow: 'ellisis',
            // whiteSpace: 'nowrap',
            // padding: '8px',
            // justifyContent: 'flex-end',
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
