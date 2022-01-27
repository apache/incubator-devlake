import React from 'react'
import {
  Icon,
  Spinner,
  Colors,
  Tooltip,
  Position,
  Intent,
} from '@blueprintjs/core'

const StageTaskIndicator = (props) => {
  const { task } = props

  return (
    <div
      className='task-module-status'
      style={{
        display: 'flex',
        justifyContent: 'center',
        padding: '8px',
        width: '32px',
        minWidth: '32px',
        outline: 'none'
      }}
    >
      {task.status === 'TASK_COMPLETED' && (
        <Tooltip content={`Task Complete [STAGE ${task.pipelineRow}]`} position={Position.TOP} intent={Intent.SUCCESS}>
          <Icon icon='small-tick' size={18} color={Colors.GREEN5} style={{ marginLeft: '0', outline: 'none' }} />
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
  )
}

export default StageTaskIndicator
