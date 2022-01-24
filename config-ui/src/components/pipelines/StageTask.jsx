import React, { useState, useEffect } from 'react'
// import { CSSTransition } from 'react-transition-group'
import {
  Providers,
  ProviderLabels
} from '@/data/Providers'
import {
  Card,
  Elevation,
} from '@blueprintjs/core'
import dayjs from '@/utils/time'
import StageTaskName from '@/components/pipelines/StageTaskName'
import StageTaskIndicator from '@/components/pipelines/StageTaskIndicator'
import StageTaskCaption from '@/components/pipelines/StageTaskCaption'

const StageTask = (props) => {
  const {
    stages = [],
    task,
    sK,
    sIdx,
  } = props

  const [taskModuleOpened, setTaskModuleOpened] = useState()

  const generateStageTaskCssClasses = () => {
    return `pipeline-task-module task-${task.status.split('_')[1].toLowerCase()} ${task.ID === taskModuleOpened?.ID ? 'active' : ''}`
  }

  return (
    <>
      <Card
        elevation={task.status === 'TASK_RUNNING' ? Elevation.TWO : Elevation.ONE}
        className={generateStageTaskCssClasses()}
        onClick={() => setTaskModuleOpened(task)}
        style={{

        }}
      >
        <StageTaskIndicator task={task} />
        <div
          className='task-module-name'
          style={{
            flex: 1
          }}
        >
          <div style={{ padding: '4px 2px 4px 0' }}>
            <StageTaskName task={task} showDetails={taskModuleOpened} onClose={() => setTaskModuleOpened(null)} />
            <StageTaskCaption task={task} options={task.options} />
            {/* <span style={{
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
            </span> */}
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
